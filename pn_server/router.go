package main

import (
	"dc3/public/utils"
	"encoding/json"
	"pn_server/config"
	"pn_server/public/errs"
	"pn_server/public/logs"
	"pn_server/public/response"
	"pn_server/record"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

func WithRouter(e *gin.Engine) {
	api := e.Group("/api/pn_server")

	// 系统信息
	api.POST("/system/info", SystemInfo)

	// 用户相关
	userApi := api.Group("/user")
	{
		userApi.POST("/register", UserRegister)
		userApi.POST("/login", UserLogin)
	}

	// 笔记相关
	noteApi := api.Group("/note")
	noteApi.Use(Auth())
	{
		noteApi.POST("/sync", Sync)
	}

	// 测试接口
	api.POST("/test/user_all", UserAll)
}

type C2S_NoteSync struct {
	Cmds []NoteCmd `json:"cmds"`
}

type S2C_NoteSync struct {
	Notes []Note `json:"notes"`
}

type NoteCmd struct {
	Tp               string `json:"tp"` // add, update, delete
	NoteID           string `json:"note_id"`
	EncryptedContent string `json:"encrypted_content"`
	CreateTime       string `json:"create_time"`
	UpdateTime       string `json:"update_time"`
	CmdTime          string `json:"cmd_time"`
}

type Note struct {
	NoteID           string `json:"note_id"`
	EncryptedContent string `json:"encrypted_content"`
	CreateTime       string `json:"create_time"`
	UpdateTime       string `json:"update_time"`
}

func Sync(c *gin.Context) {
	var req C2S_NoteSync
	var err error

	if err = c.ShouldBindJSON(&req); err != nil {
		logs.Errorf("err: %v", err)
		response.Err(c, err)
		return
	}

	// 获取 session
	session, ok := c.Keys["session"].(Session)
	if !ok {
		logs.Errorf("session is empty")
		response.Err(c, ErrUnauthorized)
		return
	}

	// 获取用户
	user, err := db.GetUser(session.UserID)
	if err != nil {
		logs.Errorf("err: %v", err)
		response.Err(c, err)
		return
	}

	// 排序命令
	sort.Slice(req.Cmds, func(i, j int) bool {
		// 创建时间逆序
		return req.Cmds[i].CmdTime > req.Cmds[j].CmdTime
	})

	// 执行命令
	for _, cmd := range req.Cmds {
		Process(user, cmd)
	}

	// 保存用户
	if err = db.SetUser(user.ID, user); err != nil {
		logs.Errorf("err: %v", err)
		response.Err(c, err)
		return
	}

	// 返回笔记
	notes := make([]Note, 0, len(user.Notes))
	for _, n := range user.Notes {
		notes = append(notes, Note{
			NoteID:           n.ID,
			EncryptedContent: n.EncryptedContent,
			CreateTime:       Int64ToStrTime(n.CreateTime),
			UpdateTime:       Int64ToStrTime(n.UpdateTime),
		})
	}

	var resp S2C_NoteSync
	resp.Notes = notes
	response.Ok(c, &resp)
}

func Process(u *record.User, cmd NoteCmd) {
	switch cmd.Tp {
	case "add":
		n := &record.Note{
			Base: record.Base{
				ID:         cmd.NoteID,
				CreateTime: StrTimeToInt64(cmd.CreateTime),
				UpdateTime: StrTimeToInt64(cmd.UpdateTime),
			},
			EncryptedContent: cmd.EncryptedContent,
		}
		u.AddNote(n)
	case "update":
		u.RemoveNote(cmd.NoteID)
		n := &record.Note{
			Base: record.Base{
				ID:         cmd.NoteID,
				CreateTime: StrTimeToInt64(cmd.CreateTime),
				UpdateTime: StrTimeToInt64(cmd.UpdateTime),
			},
			EncryptedContent: cmd.EncryptedContent,
		}
		u.AddNote(n)
	case "delete":
		u.RemoveNote(cmd.NoteID)
	default:
		logs.Errorf("unknown cmd tp: %v", cmd.Tp)
	}
}

func StrTimeToInt64(t string) int64 {
	i, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		logs.Errorf("[StrTimeToInt64]err: %v", err)
	}
	return i
}

func Int64ToStrTime(i int64) string {
	return strconv.FormatInt(i, 10)
}

func UserAll(c *gin.Context) {
	users, err := db.GetAllUsers()
	if err != nil {
		logs.Errorf("err: %v", err)
		logs.Errorf("err: %v", err)
		response.Err(c, err)
		return
	}
	response.Ok(c, users)
}

// 鉴权中间件, 从请求头中获取 account 和 password_hash
var ErrUnauthorized = errs.New(401, errs.WithMsg("unauthorized"))

type Session struct {
	UserID string `json:"user_id"`
}

func MarshalSession(session Session) (string, error) {
	bs, err := json.Marshal(session)
	return string(bs), err
}

func UnmarshalSession(session string) (Session, error) {
	var s Session
	err := json.Unmarshal([]byte(session), &s)
	return s, err
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 验证和刷新 session
		sessionID := c.GetHeader("Session-Id")
		if sessionID == "" {
			logs.Errorf("session_id is empty")
			response.Err(c, ErrUnauthorized)
			c.Abort()
			return
		}

		// 从缓存中获取 value
		value, ok := gcache.Get(sessionID)
		if !ok {
			logs.Errorf("session not found")
			response.Err(c, ErrUnauthorized)
			c.Abort()
			return
		}

		// 解析 session
		s, err := UnmarshalSession(value)
		if err != nil {
			logs.Errorf("unmarshal session err: %v", err)
			response.Err(c, ErrUnauthorized)
			c.Abort()
			return
		}

		if len(c.Keys) == 0 {
			c.Keys = make(map[string]interface{})
		}

		c.Keys["session"] = s

		// 刷新 session
		gcache.Set(sessionID, value, config.Global.SessionExpireDuation)

		c.Next()
	}
}

type C2S_SystemInfo struct{}

type S2C_SystemInfo struct {
	Version        string `json:"version"`
	NeedPrivateKey bool   `json:"need_private_key"`
	AdminContact   string `json:"admin_contact"`
}

func SystemInfo(c *gin.Context) {
	// 获取系统信息
	version := config.Global.Version
	needPrivateKey := config.Global.RegisterPrivateKey != ""

	// 返回系统信息
	var resp S2C_SystemInfo
	resp.Version = version
	resp.NeedPrivateKey = needPrivateKey
	resp.AdminContact = config.Global.AdminContact
	response.Ok(c, &resp)
}

type C2S_UserRegister struct {
	Account      string `json:"account"`
	PasswordHash string `json:"password_hash"`
	PrivateKey   string `json:"private_key"`
}

type S2C_UserRegister struct {
}

var ErrUserAlreadyExists = errs.New(1001, errs.WithMsg("user already exists"))
var ErrInvalidParam = errs.New(1002, errs.WithMsg("invalid param"))
var ErrInvalidPrivateKey = errs.New(1003, errs.WithMsg("invalid private key"))

func UserRegister(c *gin.Context) {
	var req C2S_UserRegister
	var err error

	if err = c.ShouldBindJSON(&req); err != nil {
		logs.Errorf("err: %v", err)
		logs.Errorf("err: %v", err)
		response.Err(c, err)
		return
	}

	// 参数错误
	if req.Account == "" || req.PasswordHash == "" {
		logs.Errorf("err: %v", err)
		response.Err(c, ErrInvalidParam)
		return
	}

	// 用户已经存在
	uid, err := db.GetUserIDByAccount(req.Account)
	if err == nil || uid != "" {
		logs.Errorf("err: %v", err)
		response.Err(c, ErrUserAlreadyExists)
		return
	}

	// 注册密钥错误
	pk := config.Global.RegisterPrivateKey
	if pk != "" && req.PrivateKey != pk {
		logs.Errorf("err: %v", err)
		response.Err(c, ErrInvalidPrivateKey)
		return
	}

	// 创建用户
	u := record.NewUser(req.Account, req.PasswordHash)
	if err = db.SetUser(u.ID, u); err != nil {
		logs.Errorf("err: %v", err)
		response.Err(c, err)
		return
	}

	// 返回用户ID
	var resp S2C_UserRegister
	response.Ok(c, &resp)
}

type C2S_UserLogin struct {
	Account      string `json:"account"`
	PasswordHash string `json:"password_hash"`
}

type S2C_UserLogin struct {
	SessionID string `json:"session_id"`
}

var ErrUserNotExists = errs.New(2001, errs.WithMsg("user not exists"))
var ErrInvalidPassword = errs.New(2002, errs.WithMsg("invalid password"))

func UserLogin(c *gin.Context) {
	var req C2S_UserLogin
	var err error

	if err = c.ShouldBindJSON(&req); err != nil {
		logs.Errorf("err: %v", err)
		logs.Errorf("err: %v", err)
		response.Err(c, err)
		return
	}

	// 参数错误
	if req.Account == "" || req.PasswordHash == "" {
		logs.Errorf("err: %v", err)
		response.Err(c, ErrInvalidParam)
		return
	}

	// 用户不存在
	uid, err := db.GetUserIDByAccount(req.Account)
	if err != nil {
		logs.Errorf("err: %v", err)
		response.Err(c, err)
		return
	}

	// 获取用户信息
	u, err := db.GetUser(uid)
	if err != nil {
		logs.Errorf("err: %v", err)
		response.Err(c, err)
		return
	}

	// 密码错误
	if u.PasswordHash != req.PasswordHash {
		logs.Errorf("err: %v", err)
		response.Err(c, ErrInvalidPassword)
		return
	}

	// 创建 Session
	sessionID := utils.RandString(32)
	session := Session{
		UserID: u.ID,
	}
	sessionStr, err := MarshalSession(session)
	if err != nil {
		logs.Errorf("err: %v", err)
		response.Err(c, err)
		return
	}
	gcache.Set(sessionID, sessionStr, config.Global.SessionExpireDuation)

	// 返回 SessionID
	var resp S2C_UserLogin
	resp.SessionID = sessionID

	response.Ok(c, &resp)
}
