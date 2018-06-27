package models

import (
	"database/sql"
	"errors"
	"github.com/go-redis/redis"
	"appone/config"
	"log"
)


var dbConn *sql.DB
var redisConn *redis.Client

func init() {
	dbConn = config.InitMysqlConnection("db.json")
	redisConn = config.InitRedisClient("redis.json")
}

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Passwd string `json:"passwd,omitempty"`
	Age int `json:"age"`
}

type Feedback struct {
	Id int `json:"id"`
	NetbarId int `json:"netbar_id"`
	NetbarName string `json:"netbar_name"`
	CreateTime int `json:"create_time"`
	Content string `json:"content"`
}

type ApiResponse struct {
	Status int `json:"status"`
	Msg string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func GetUserById(userId int) (User, error) {
	var u User
	stm, e := dbConn.Prepare("SELECT id,name,age FROM `user` where id = ?")
	if e != nil {
		return u, e
	}
	defer stm.Close()
	var id, age int
	var name string
	e = stm.QueryRow(userId).Scan(&id, &name,&age)
	if e == sql.ErrNoRows {
		return u, errors.New("找不到记录")
	}
	if e != nil {
		return u, e
	}
	u.Id = id
	u.Name = name
	u.Age = age
	return u, nil
}

func AddNewUser(u User) (User, error) {
	stm, e := dbConn.Prepare("INSERT INTO `user`(`name`,`age`) VALUES(?,?)")

	if e != nil {
		return u, e
	}
	defer stm.Close()
	result , e := stm.Exec(u.Name, u.Age)
	if e != nil {
		return u, e
	}
	lastId, e := result.LastInsertId()
	if e != nil {
		return u, e
	}
	u.Id = int(lastId)
	return u, nil
}

func GetUsers(options map[string]string) ([]User, error) {
	u := []User{}
	selectSql := "select id, name, age from `user` where 1 "
	if age, ok := options["age"]; ok {
		selectSql += " and age > " + age
	}
	if utype, ok := options["type"]; ok {
		selectSql += " and type = " + utype
	}
	stm, e := dbConn.Prepare(selectSql)
	if e != nil {
		return u, e
	}
	defer stm.Close()
	rows, e := stm.Query()
	if e != nil {
		return u, e
	}
	defer rows.Close()
	var id ,age int
	var name string
	for rows.Next() {
		e = rows.Scan(&id, &name, &age)
		if e != nil {
			log.Fatal(e)
		}
		u = append(u, User{
			Id: id,
			Name:name,
			Age:age,
		})
	}

	return u, nil
}
