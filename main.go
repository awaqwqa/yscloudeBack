package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"yscloudeBack/route"
	"yscloudeBack/source/app/model"
)

func main() {

	Db, err := gorm.Open(sqlite.Open("yscloudBack.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dbManager := model.NewDbManager(Db)
	r := gin.Default()
	route.InitRoute(r, dbManager)
	err = r.Run(":24016")
	if err != nil {
		return
	}
}
