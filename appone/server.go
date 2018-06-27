package main

import (
	"net/http"
	"appone/business"
	"log"
)

func main() {

	user := business.User{}
	http.HandleFunc("/user", user.UserBusiness)
	http.HandleFunc("/users", user.UserList)
	log.Println("add all handle func....")
	log.Fatal(http.ListenAndServe(":1234", nil))

}

