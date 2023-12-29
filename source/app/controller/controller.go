package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"yscloudeBack/source/app/cluster"
	"yscloudeBack/source/app/db"
	"yscloudeBack/source/app/filer"
	"yscloudeBack/source/app/middleware"
	"yscloudeBack/source/app/model"
	"yscloudeBack/source/app/utils"
)

func NewControllerManager() *ControllerMannager {
	return &ControllerMannager{
		isDbWork:      false,
		isClusterWork: false,
	}
}

type ControllerMannager struct {
	dbManager        *db.DbManager
	cluster          *cluster.ClusterRequester
	streamController *cluster.StreamController
	filer            *filer.Filer
	isDbWork         bool
	isClusterWork    bool
}

func (cm *ControllerMannager) GetUserFromCtx(ctx *gin.Context) (*model.User, error) {
	name, err := middleware.GetContextName(ctx)
	if err != nil {
		return &model.User{}, err
	}
	user, err := cm.GetDbManager().GetUserByUserName(name)
	if err != nil {
		utils.Error(err.Error())
		return &model.User{}, fmt.Errorf("cant get user struct by name ")
	}
	return user, nil

}
func (cm *ControllerMannager) Init(sc *cluster.StreamController, db *db.DbManager, cluster *cluster.ClusterRequester, filer *filer.Filer) error {
	err := cm.SetStreamController(sc)
	if err != nil {
		return err
	}
	err = cm.SetDbManager(db)
	if err != nil {
		return err
	}
	err = cm.SetCluster(cluster)
	if err != nil {
		return err
	}

	err = cm.SetFiler(filer)
	if err != nil {
		return err
	}
	return nil
}
func (cm *ControllerMannager) SetFiler(filer *filer.Filer) error {
	if cm.filer != nil {
		return fmt.Errorf("filer is exited ,you cant set filer again")
	}
	cm.filer = filer
	return nil
}
func (cm *ControllerMannager) GetFiler() *filer.Filer {
	return cm.filer
}
func (cm *ControllerMannager) GetDbManager() *db.DbManager {
	return cm.dbManager
}
func (cm *ControllerMannager) GetCluster() *cluster.ClusterRequester {
	return cm.cluster
}
func (cm *ControllerMannager) SetDbManager(db *db.DbManager) error {
	if cm.isDbWork || cm.dbManager != nil {
		return fmt.Errorf("dbWork still work")
	}

	cm.dbManager = db
	cm.isDbWork = true
	return nil
}
func (cm *ControllerMannager) SetStreamController(sc *cluster.StreamController) error {
	cm.streamController = sc
	return nil
}
func (cm *ControllerMannager) SetCluster(cluster *cluster.ClusterRequester) error {
	if cm.isClusterWork {
		return fmt.Errorf("cluster still work")
	}
	cm.cluster = cluster
	cm.isClusterWork = true
	ctx, cancelFn := context.WithCancel(context.Background())
	go func() {
		err := cm.cluster.InitReadLoop(ctx)
		if err != nil {
			cancelFn()
			utils.Error(err.Error())
		}
	}()
	return nil
}
