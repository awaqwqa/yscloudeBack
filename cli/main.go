package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"yscloudeBack/route"
)

func main() {
	r := gin.Default()
	route.InitRoute(r)
	db, err := gorm.Open(sqlite.Open("yscloudBack.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = r.Run(":24016")
	if err != nil {
		return
	}
}
