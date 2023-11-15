package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"yscloudeBack/route"
	"yscloudeBack/source/app/cluster"
	"yscloudeBack/source/app/db"
	"yscloudeBack/source/app/utils"
)

func main() {

	Db, err := gorm.Open(sqlite.Open("yscloudBack.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dbManager := db.NewDbManager(Db)
	err = dbManager.Init()
	if err != nil {
		utils.Error(err.Error())
		return
	}

	cmdController := utils.NewCmdController()
	cmdController.Init()
	cmdController.Listen()

	//loger
	utils.NewLoggerManager("./log")
	//clusterRequester
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3002/", http.Header{})
	if err != nil {
		panic(err)
	}
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	client := cluster.NewClusterRequester(conn)
	go func() {
		err := client.InitReadLoop(ctx)
		if err != nil {
			cancelFn()
			panic(err)
		}
	}()

	r := gin.Default()

	route.InitRoute(r, dbManager, client)
	err = r.Run(":24016")
	if err != nil {
		return
	}
}
