package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"yscloudeBack/route"
	"yscloudeBack/source/app/cluster"
	"yscloudeBack/source/app/controller"
	"yscloudeBack/source/app/db"
	"yscloudeBack/source/app/filer"
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
		panic(err)
	}
	//这里塞一个相对地址
	newFiler := filer.NewFiler("./user_file")
	//clusterRequester
	// 这里由controllerManager负责调控 所以defer controllerManager.close()会关闭所有的线程
	client := cluster.NewClusterRequester()
	err = client.Init("ws://localhost:3002/")
	if err != nil {
		utils.Error(err.Error())
	}
	// 设置streamController
	stc := cluster.NewStreamController(client)
	cm := controller.NewControllerManager()
	err = cm.Init(stc, dbManager, client, newFiler)
	if err != nil {
		utils.Error(err.Error())
		return
	}

	r := gin.Default()

	route.InitRoute(r, cm)
	err = r.Run(":24016")
	if err != nil {
		return
	}
}
