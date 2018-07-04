package main

import (
	"appone/business"
	"log"
	"net/http"
)

func main() {
	routes := routes()
	for url, fun := range routes {
		http.HandleFunc(url, fun)
	}

	log.Println("add all handle func....")
	log.Fatal(http.ListenAndServe(":1234", nil))
}

func routes() map[string]func(http.ResponseWriter, *http.Request) {
	routesMap := make(map[string]func(http.ResponseWriter, *http.Request))
	user := business.User{}
	routesMap["/user"] = user.UserBusiness
	routesMap["/users"] = user.UserList
	routesMap["/usertx"] = user.UserTx

	game := business.Game{}
	routesMap["/game/load"] = game.Upload

	return routesMap
}
