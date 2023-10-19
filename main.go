package main

import (
	"github.com/gin-gonic/gin"
	"yscloudeBack/route"
)

func main() {
	r := gin.Default()
	route.InitRoute(r)
	err := r.Run(":24016")
	if err != nil {
		return
	}
}
