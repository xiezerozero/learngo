package main

import (
	models2 "appone/models"
	"github.com/go-redis/redis"
	"log"
	"msg_queue/models"
	"time"
)

func main() {

	gameUploadBusiness()

}

func gameUploadBusiness() {
	redisConn := models.DefaultRedisConn()

	for {
		result, e := redisConn.RPop(models2.GAME_UPLOAD_REDIS_KEY).Result()
		if e != nil {
			if e == redis.Nil {
				log.Println("redis game load list is empty")
				time.Sleep(time.Second)
			} else {
				log.Println(e.Error())
			}
			continue
		}

		go models.QueueGameInDB(result)
	}
}
