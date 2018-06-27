package config

import (
	"os"
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/go-redis/redis"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	Addr string
	Password string
	Db int
}

func InitRedisClient(fileName string) *redis.Client {
	file, e := os.Open(fileName)

	if e != nil {
		log.Println("redis config file error: ", e.Error())
		return nil
	}
	defer file.Close()
	b, _ := ioutil.ReadAll(file)
	var config config
	json.Unmarshal(b, &config)

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Addr,
		Password:config.Password,
		DB: config.Db,
	})
	return redisClient
}

type dbConfig struct {
	DriverName string `json:"driver_name"`
	User, Password string
	Addr string
	DbName string
}

func InitMysqlConnection(fileName string) *sql.DB {
	configFile, e := os.Open(fileName)

	defer configFile.Close()

	var config dbConfig
	b, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(b, &config)

	db ,e := sql.Open(config.DriverName, fmt.Sprintf("%s:%s@tcp(%s)/%s", config.User, config.Password, config.Addr, config.DbName))
	if e != nil {
		panic("连不上数据库")
	}
	e = db.Ping()
	if e != nil {
		panic("连不上数据库")
	}
	return db
}
