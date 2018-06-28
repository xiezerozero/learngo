package models

import (
	"database/sql"
	"errors"
	"github.com/go-redis/redis"
	"appone/config"
	"log"
	"strconv"
	"fmt"
)


var dbConn *sql.DB
var redisConn *redis.Client
var GAMEINFO_REDIS_KEY string = "game_info"

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
	Data interface{} `json:"data"`
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
	pageSize := 10
	page := 1

	if userPage, ok := options["page"]; ok {
		page, _ = strconv.Atoi(userPage)
	}
	if userPageSize, ok := options["pageSize"]; ok {
		pageSize, _ = strconv.Atoi(userPageSize)
	}
	selectSql += " limit %d, %d "
	selectSql = fmt.Sprintf(selectSql, pageSize * (page - 1), pageSize)
	log.Println(selectSql)
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
	if e = rows.Err(); e != nil {
		// 获取失败
	}

	return u, nil
}

// 测试transaction使用
func UserInsertAndUpdate() (e error) {
	tx, e := dbConn.Begin()

	if e != nil {
		return
	}

	insertSql := "INSERT INTO `user`(`name`, `passwd`,`age`,`type`) VALUES('ASD','asdsad',12,1);"
	_, e = tx.Exec(insertSql)
	if e != nil {
		tx.Rollback()
		return
	}
	updateSql := "UPDATE `user` SET `passwd`='rwe' WHERE `id`=18"
	_, e = tx.Exec(updateSql)
	if e != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return nil
}

// 测试不确定多少字段的返回
func GetUserList() (data []map[string]string) {
	data = []map[string]string{}
	userSql := "select * from `user`"
	rows, e := dbConn.Query(userSql)
	if e != nil {
		return
	}
	// 获取字段名的[]string
	cols, e := rows.Columns()
	if e != nil {
		return
	}
	// 定义数据存在bytes字段
	values := make([]sql.RawBytes, len(cols))
	// scan参数需要[]interface{}类型分散传入
	vals := make([]interface{}, len(cols))
	for i, _ := range cols {
		vals[i] = &values[i]
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(vals...)
		m := map[string]string{}
		for i, v := range values {
			var value string
			if v == nil {
				value = "null"
			} else {
				value = string(v)
			}
			m[cols[i]] = value
		}
		data = append(data, m)
	}
	return
}

func PutGameInfoInRedis(gameInfo map[string]int) {

	for name, minute := range gameInfo {
		// redis.hincrby 不管是否存在该field
		_, e := redisConn.HIncrBy(GAMEINFO_REDIS_KEY, name, int64(minute)).Result()
		if e != nil {
			log.Println("field " + name + " 写入失败:", e.Error())
		}
	}

}
