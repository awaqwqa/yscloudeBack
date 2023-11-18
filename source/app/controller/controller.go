package controller

import (
	"fmt"
	"yscloudeBack/source/app/cluster"
	"yscloudeBack/source/app/db"
)

func NewControllerManager() *ControllerMannager {
	return &ControllerMannager{
		isDbWork:      false,
		isClusterWork: false,
	}
}

type ControllerMannager struct {
	dbManager     *db.DbManager
	cluster       *cluster.ClusterRequester
	isDbWork      bool
	isClusterWork bool
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
func (cm *ControllerMannager) SetCluster(cluster *cluster.ClusterRequester) error {
	if cm.isClusterWork {
		return fmt.Errorf("cluster still work")
	}
	cm.cluster = cluster
	cm.isClusterWork = true
	return nil
}
