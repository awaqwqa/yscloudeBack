package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"yscloudeBack/route"
	"yscloudeBack/source/app/model"
	"yscloudeBack/source/app/utils"
)

func main() {

	Db, err := gorm.Open(sqlite.Open("yscloudBack.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dbManager := model.NewDbManager(Db)
	err = dbManager.Init()
	if err != nil {
		return
	}

	cmdControler := utils.NewCmdController()
	cmdControler.Init()
	cmdControler.Listen()
	r := gin.Default()
	route.InitRoute(r, dbManager)
	err = r.Run(":24016")
	if err != nil {
		return
	}
}
