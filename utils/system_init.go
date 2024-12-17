// @Author TrandLiu
// @Date 2024/12/12 22:24:00
// @Desc
package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("config mysql", viper.Get("mysql"))
}
func InitMysql() {
	//自定义模板打印sql语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger})
	fmt.Println("mysql inited///")
	//user := models.UserBasic{}
	//DB.Find(&user)
	//fmt.Println(user)
}

// 初始化redis
func InitRedis() {
	Red := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addrr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	fmt.Println(Red)
}

const PublishKey = "websocket"

// Publish发布消息到Redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("publish....", msg)
	err = Red.Publish(ctx, channel, msg).Err()
	return err
}

// Subscribe订阅redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel)
	fmt.Println("Subcribe", ctx)
	message, err := sub.ReceiveMessage(ctx)
	if err != nil {
		return "", err
	}
	fmt.Println("Subscribe.......", message.Payload)
	return message.Payload, err
}
