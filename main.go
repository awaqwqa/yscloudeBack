package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
	"yscloudeBack/route"
	"yscloudeBack/source/app/archiveManager"
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

	// 这里由controllerManager负责调控 所以defer controllerManager.close()会关闭所有的线程
	client := cluster.NewClusterRequester()
	err = client.Init("ws://localhost:3002/")
	if err != nil {
		//utils.Error(err.Error())
		panic(err.Error())
	}
	// 设置streamController
	stc := cluster.NewStreamController(client)

	acm := archiveManager.NewArchiveManager("./archive")
	go ArchiveListen(client, acm)

	cm := controller.NewControllerManager()
	err = cm.Init(stc, dbManager, client, newFiler, acm)
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
func ArchiveListen(client *cluster.ClusterRequester, archiveManager *archiveManager.ArchiveManager) {
	for {
		//看起来像是实列转化为建筑
		instanceToArchive := make([]cluster.InstanceDetail, 0)
		choker := make(chan struct{}, 1)
		utils.Info("details:")
		client.List(func(details []cluster.InstanceDetail) {
			for _, d := range details {
				fmt.Printf("%v,", d)
				if d.Status == "Finished" {
					instanceToArchive = append(instanceToArchive, d)
				}
			}
			close(choker)
		})
		<-choker
		utils.Info("task:")
		for _, _instance := range instanceToArchive {
			instance := _instance
			fmt.Printf("%v,", instance)
			d, _ := json.Marshal(instance)
			archiveManager.ArchiveNew(fmt.Sprintf("instance.%v.detail", instance.InstanceID), d)
			// 说明一下journal干什么的
			client.Journal(instance.InstanceID, 0, func(journal string, err string) {
				archiveManager.ArchiveNew(fmt.Sprintf("instance.%v.journal", instance.InstanceID), []byte(journal))
				client.Rm(instance.InstanceID, func(err string) {
					fmt.Printf("instance %v archived\n", instance.InstanceID)
				})
			})
		}
		time.Sleep(60 * time.Second)
	}
}
