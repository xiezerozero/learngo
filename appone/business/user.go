package business

import (
	"net/http"
	"appone/models"
	"strings"
	"strconv"
	"github.com/astaxie/beego/validation"
	"errors"
)

type User struct {
	Base
}

// 单个用户
func (this *User) UserBusiness(writer http.ResponseWriter, r *http.Request)  {
	u, e := doUserBusiness(r)
	if e != nil {
		this.status = 0
		this.msg = e.Error()
		this.data = nil
	} else {
		this.status = 1
		this.msg = ""
		this.data = u
	}
	this.json(writer)
}

// 用户列表
func (this *User) UserList(writer http.ResponseWriter, r *http.Request)  {
	u, e := doUserList(r)
	if e != nil {
		this.status = 0
		this.msg = e.Error()
		this.data = nil
	} else {
		this.status = 1
		this.msg = ""
		this.data = u
	}
	this.json(writer)
}

func (this *User) UserTx(writer http.ResponseWriter, r *http.Request) {

	data := models.GetUserList()
	this.status = 1
	this.msg = ""
	this.data = data
	this.json(writer)
}

func doUserList(r *http.Request) ([]models.User, error) {
	options := map[string]string{}
	agestring := r.URL.Query().Get("age")
	utypestring := r.URL.Query().Get("type")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")
	_, e := strconv.Atoi(agestring)
	if e == nil {
		options["age"] = agestring
	}
	_, e = strconv.Atoi(utypestring)
	if e == nil {
		options["type"] = utypestring
	}
	if _, e = strconv.Atoi(page); e == nil {
		options["page"] = page
	}
	if _, e = strconv.Atoi(pageSize); e == nil {
		options["pageSize"] = pageSize
	}
	users, e := models.GetUsers(options)
	return users, e
}


// 处理业务
func doUserBusiness(r *http.Request) (models.User, error) {
	var u models.User
	if strings.ToLower(r.Method) == "get" {	//select
		idstring := r.URL.Query().Get("id")
		userId, e := strconv.Atoi(idstring)
		if e != nil {
			return u, errors.New("id参数不正确")
		} else {
			u, e = models.GetUserById(userId)
			if e != nil {
				return u, e
			} else {
				return u, nil
			}
		}
	} else {	//post:  insert/update
		r.ParseForm()
		valid := validation.Validation{}
		name := r.Form.Get("name")
		ageString := r.Form.Get("age")
		valid.Required(name, "name")
		valid.Numeric(ageString, "age")
		if valid.HasErrors() {
			return u, errors.New(valid.Errors[0].Message)
		} else {
			age, _ := strconv.Atoi(ageString)
			u.Name = name
			u.Age = age
			u, e := models.AddNewUser(u)
			if e != nil {
				return u, e
			} else {
				return u, nil
			}
		}
	}
}

