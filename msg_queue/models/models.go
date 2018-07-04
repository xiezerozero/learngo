package models

import (
	"config"
	"database/sql"
	"github.com/go-redis/redis"
)

var dbConn *sql.DB
var redisConn *redis.Client

func init() {
	dbConn = config.InitMysqlConnection("db.json")
	redisConn = config.InitRedisClient("redis.json")
}

func DefaultDBConn() *sql.DB {
	return dbConn
}

func DefaultRedisConn() *redis.Client {
	return redisConn
}
