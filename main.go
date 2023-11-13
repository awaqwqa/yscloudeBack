package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"yscloudeBack/route"
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

	r := gin.Default()
	route.InitRoute(r, dbManager)
	err = r.Run(":24016")
	if err != nil {
		return
	}
}
