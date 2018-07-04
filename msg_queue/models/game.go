package models

import (
	"appone/models"
	"encoding/json"
	"log"
	"time"
)

func QueueGameInDB(res string) (e error) {

	var game models.Game
	e = json.Unmarshal([]byte(res), &game)
	if e != nil {
		log.Println(e.Error())
		return
	}
	stm, e := dbConn.Prepare("INSERT INTO `game_upload`(`game_name`, `game_duration`,`create_time`,`dtdate`) VALUES(?,?,?,?)")
	if e != nil {
		log.Println(e.Error())
		return
	}
	t := time.Now()
	timeStamp := t.Unix()
	dtDate := t.Format("20060102")
	log.Println("haomiao:", t.UnixNano()/1e6)
	log.Println("time: ", timeStamp, ", ", dtDate)
	_, e = stm.Exec(game.Name, game.Minute, timeStamp, dtDate)
	if e != nil {
		log.Println(e.Error())
		return
	}
	return nil
}
