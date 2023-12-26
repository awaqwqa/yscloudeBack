package main

import (
	"context"
	"fmt"
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
	filer := filer.NewFiler("./user_file")
	//clusterRequester
	// 这里由controllerManager负责调控 所以defer contollerManager.close()会关闭所有的线程
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	client := cluster.NewClusterRequester()
	err = client.Init("ws://localhost:3002/")
	if err != nil {
		utils.Error(err.Error())
	}
	fmt.Println(ctx)
	cm := controller.NewControllerManager()
	err = cm.SetDbManager(dbManager)
	if err != nil {
		panic(err)
	}
	err = cm.SetCluster(client)
	if err != nil {
		utils.Error(err.Error())
		return
	}

	err = cm.SetFiler(filer)
	if err != nil {
		utils.Error(err.Error())
		return
	}
	//go func() {
	//	err := client.InitReadLoop(ctx)
	//	if err != nil {
	//		cancelFn()
	//		utils.Error(err.Error())
	//	}
	//}()

	//cmdController := utils.NewCmdController()
	//cmdController.Init()
	//cmdController.Listen()

	//loger
	//utils.NewLoggerManager("./log")

	r := gin.Default()

	route.InitRoute(r, cm)
	err = r.Run(":24016")
	if err != nil {
		return
	}
}
