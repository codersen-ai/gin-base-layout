package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

func MustInitMySQL(cfg *viper.Viper) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.GetString("mysql.user"),
		cfg.GetString("mysql.password"),
		cfg.GetString("mysql.host"),
		cfg.GetString("mysql.port"),
		cfg.GetString("mysql.dbname"),
	)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Errorf("connect db fail: %w", err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("connect db fail: %w", err))
	}
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.GetInt("mysql.max_idle_conns"))
	sqlDB.SetMaxOpenConns(cfg.GetInt("mysql.max_open_conns"))
	sqlDB.SetConnMaxLifetime(cfg.GetDuration("mysql.max_lifetime"))
	// query.SetDefault(db)
}

func MustInitRedis(conf *viper.Viper) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GetString("data.redis.addr"),
		Password: conf.GetString("data.redis.password"),
		DB:       conf.GetInt("data.redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Errorf("init redis failed, err:%w", err))
	}
	RedisClient = rdb
}
