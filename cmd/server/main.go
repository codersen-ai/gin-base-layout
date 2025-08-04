package main

import (
	"flag"
	"fmt"

	"github.com/q1mi/gin-base-layout/internal/conf"
	"github.com/q1mi/gin-base-layout/internal/dao"
	"github.com/q1mi/gin-base-layout/internal/server"
	"github.com/q1mi/gin-base-layout/pkg/jwt"
	"github.com/q1mi/gin-base-layout/pkg/logging"
	"github.com/q1mi/gin-base-layout/pkg/snowflake"
)

var confPath = flag.String("conf", "./config/config.yaml", "配置文件路径")

func main() {
	// 加载配置
	flag.Parse()
	cfg := conf.Load(*confPath)

	// 初始化日志
	logger, err := logging.NewLogger(cfg)
	if err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer logger.Sync()

	dao.MustInitMySQL(cfg)  // 初始化 MySQL 连接
	dao.MustInitRedis(cfg)  // 初始化 Redis
	jwt.MustInit(cfg)       // 初始化 jwt
	snowflake.MustInit(cfg) // 初始化 snowflake

	// 初始化路由
	r := server.SetupRoutes(cfg)
	// 启动服务
	err = r.Run(fmt.Sprintf(":%d", cfg.GetInt("server.port")))
	if err != nil {
		panic(err)
	}
}
