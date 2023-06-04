package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"pn_server/cache"
	"pn_server/config"
	"pn_server/public/logs"
	"pn_server/record"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	db     record.RecordI = nil
	gcache cache.CacheI   = nil
)

func main() {
	// 测试 pprof
	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	// 初始化配置
	config.Init("./config-dev.yaml")

	// 初始化日志
	logs.Init(config.Global.Log)

	// 初始化 record
	r := record.NewLevelDBRecord(config.Global.Db.Leveldb.Path)
	db = r

	// 初始化缓存, 100K
	c := cache.New(100 * 1000)
	gcache = c

	// 初始化 gin
	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(gin.Logger())
	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	cfg.AddAllowHeaders("Session-Id")
	e.Use(cors.New(cfg))

	// 路由
	WithRouter(e)

	// 监听端口 7001
	e.Run(config.Global.Web.Addr)

}
