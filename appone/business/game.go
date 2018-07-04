package business

import (
	"appone/models"
	"net/http"
	"strconv"
	"strings"
)

type Game struct {
	Base
}

func (this *Game) Upload(writer http.ResponseWriter, r *http.Request) {

	uploadGame(r)
	this.status = 1
	this.msg = "ok"
	this.json(writer)
}

func uploadGame(r *http.Request) {
	r.ParseForm()
	infoString := r.PostForm.Get("infos")
	infoSlice := strings.Split(infoString, ";")
	// name => minute
	gameInfo := map[string]int{}
	for _, info := range infoSlice {
		infos := strings.Split(info, ",")
		if len(infos) != 2 {
			continue
		}
		if infos[0] == "" {
			continue
		}
		minute, e := strconv.Atoi(infos[1])
		if e != nil {
			continue
		}
		gameInfo[infos[0]] = minute
	}
	models.PutGameInfoInRedis(gameInfo)
}
