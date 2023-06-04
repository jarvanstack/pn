package record

import (
	"dc3/public/errs"
	"dc3/public/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

type RecordI interface {
	GetUserIDByAccount(account string) (string, error)
	GetUser(id string) (*User, error)
	SetUser(id string, user *User) error
	GetAllUsers() (map[string]*User, error)
}

type Base struct {
	ID         string `json:"id"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

type User struct {
	Base
	Account      string           `json:"account"`
	PasswordHash string           `json:"password_hash"`
	Notes        map[string]*Note `json:"notes"`
}

func NewID() string {
	return utils.RandString(32)
}

func Now() int64 {
	return time.Now().Unix()
}

func NewUser(account, passwordHash string) *User {
	return &User{
		Base: Base{
			ID:         NewID(),
			CreateTime: Now(),
			UpdateTime: Now(),
		},
		Account:      account,
		PasswordHash: passwordHash,
	}
}

func (u *User) AddNote(n *Note) {
	if u.Notes == nil {
		u.Notes = make(map[string]*Note)
	}
	u.Notes[n.ID] = n
}

func (u *User) RemoveNote(id string) {
	delete(u.Notes, id)
}

type Note struct {
	Base
	EncryptedContent string `json:"encrypted_content"`
}

var _ RecordI = &LevelDBRecord{}

type LevelDBRecord struct {
	db *leveldb.DB
}

func NewLevelDBRecord(path string) *LevelDBRecord {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Fatalf("open db error: %v, path:%v", err, path)
	}
	return &LevelDBRecord{
		db: db,
	}
}

var ErrUserNotFound = errs.New(1000, errs.WithMsg("user not found"))

func (r *LevelDBRecord) GetUser(id string) (*User, error) {
	has, err := r.db.Has(userKey(id), nil)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrUserNotFound
	}

	value, err := r.db.Get(userKey(id), nil)
	if err != nil {
		return nil, err
	}
	var u User
	err = json.Unmarshal(value, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func userKey(userID string) []byte {
	return []byte("user:" + userID)
}

func (r *LevelDBRecord) SetUser(id string, user *User) error {
	value, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// 更新 account -> userID 的映射
	olduser, err := r.GetUser(id)
	if err != nil {
		fmt.Printf("get user error: %v\n", err)
	} else {
		err = r.db.Delete(accountKey(olduser.Account), nil)
		if err != nil {
			return err
		}
	}
	err = r.db.Put(accountKey(user.Account), []byte(id), nil)
	if err != nil {
		return err
	}

	// 更新 user
	return r.db.Put(userKey(id), value, nil)
}

func (r *LevelDBRecord) GetAllUsers() (map[string]*User, error) {
	iter := r.db.NewIterator(nil, nil)
	defer iter.Release()

	users := make(map[string]*User)
	for iter.Next() {
		var u User
		err := json.Unmarshal(iter.Value(), &u)
		if err != nil {
			return nil, err
		}
		users[u.Account] = &u
	}
	return users, nil
}

func accountKey(userID string) []byte {
	return []byte("account_userid:" + userID)
}

func (r *LevelDBRecord) GetUserIDByAccount(account string) (string, error) {
	value, err := r.db.Get(accountKey(account), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return string(value), nil
}
