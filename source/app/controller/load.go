package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"
	"yscloudeBack/source/app/cluster"
	"yscloudeBack/source/app/model"
	"yscloudeBack/source/app/utils"

	"github.com/gin-gonic/gin"
)

const (
	TYPE_LOAD_KEY     = "key_load"
	TYPE_LOAD_BALANCE = "balance_load"
)

type StructureForm struct {
	FileGroupName string `json:"file_group_name"`
	StructureName string `json:"structure_name"`
	ServerCode    string `json:"server_code"`
	ServerPasswd  string `json:"server_passwd"`
	PosX          int    `json:"pox_x"`
	PoxY          int    `json:"pos_y"`
	PosZ          int    `json:"pos_z"`
	Type          string `json:"type"`
	Key           string `json:"key"`
}

func CheckLoadStrucForm(form *StructureForm) error {
	if form.FileGroupName == "" {
		return fmt.Errorf("file group name cant be nil")
	}
	if form.StructureName == "" {
		return fmt.Errorf("struct name cant be nil")
	}
	if form.ServerCode == "" {
		return fmt.Errorf("server code cant be nil")
	}
	if form.Type == "" {
		return fmt.Errorf("type cant be nil")
	}
	if form.Type == TYPE_LOAD_KEY && form.Key == "" {
		return fmt.Errorf("when type is key_load ,the key is not alllowed be nil ")
	}
	return nil
}
func (cm *ControllerMannager) GetInstanceStatus() gin.HandlerFunc {
	archiveManager := cm.archiveManager
	client := cm.GetCluster()
	return func(ctx *gin.Context) {
		var form struct {
			InstanceId string `json:"instance_id"`
		}
		code, err := model.BindStruct(ctx, &form)
		if err != nil {
			model.BackError(ctx, code)
			return
		}
		instanceID := form.InstanceId
		res, err := archiveManager.GetArchive(fmt.Sprintf("instance.%v.detail", instanceID))
		if err != nil {
			utils.Error(err.Error())
			model.BackErrorByString(ctx, fmt.Sprintf("cant find instace.%v.detail file", instanceID))
			return
		}

		detail := cluster.InstanceDetail{}
		err = json.Unmarshal(res, &detail)
		errS := ""
		if err != nil {
			errS = err.Error()
		}
		model.BackSuccess(ctx, gin.H{"status": detail.Status, "name": detail.Name, "detail": detail.StatusDetail, "error": errS})
		//c.JSON(200, gin.H{"status": detail.Status, "name": detail.Name, "detail": detail.StatusDetail, "error": errS})
		return
		//choker := make(chan struct{}, 1)
		//client.Status(instanceID, func(status, detail, name string, err string) {
		//	model.BackSuccess(ctx, gin.H{"status": status, "name": name, "detail": detail, "error": err})
		//	//c.JSON(200, gin.H{"status": status, "name": name, "detail": detail, "error": err})
		//	close(choker)
		//})
		//<-choker

	}
}
func (cm *ControllerMannager) GetStreamOutput() gin.HandlerFunc {
	//db := cm.GetDbManager()
	return func(ctx *gin.Context) {
		var form struct {
			InsId string `json:"instance_id"`
		}
		code, err2 := model.BindStruct(ctx, &form)
		if err2 != nil {
			model.BackError(ctx, code)
			return
		}
		instanceID := form.InsId
		utils.Info(instanceID)
		choker := make(chan struct{}, 1)
		stop, err := cm.streamController.AttachListener(instanceID, func(msg string) error {
			_, e := ctx.Writer.WriteString(msg)
			if e != nil {
				return e
			}
			choker <- struct{}{}
			return nil

		})
		if err != nil {
			close(choker)
			model.BackErrorByString(ctx, err.Error())
			return
		}
		ctx.Status(200)
		<-choker
		stop()
	}
}
func (cm *ControllerMannager) LoadHandler() gin.HandlerFunc {
	db := cm.GetDbManager()
	client := cm.GetCluster()
	return func(ctx *gin.Context) {
		//绑定参数
		var form StructureForm
		option := &model.BuildOption{}
		code, err := model.BindStruct(ctx, &form)
		if err != nil {
			model.BackError(ctx, code)
			return
		}
		// 检查form是否正常
		err = CheckLoadStrucForm(&form)
		if err != nil {
			model.BackErrorByString(ctx, err.Error())
			return
		}

		user, err := cm.GetUserFromCtx(ctx)
		if err != nil {
			model.BackErrorByString(ctx, err.Error())
			return
		}
		tokens, err := db.GetFbTokens()
		if err != nil {
			utils.Error(err.Error())
			model.BackErrorByString(ctx, fmt.Sprintf("cant get fb tokens of struct_loader"))
			return
		}
		if len(tokens) == 0 {
			model.BackErrorByString(ctx, fmt.Sprintf("fbtoken num is 0"))
			return
		}
		fbToken := tokens[0]
		func() {
			//set taskName
			//
			//utils.Info("fbToken is : %v", fbToken)
			option.TaskName = form.StructureName + "::" + user.UserName
			option.Auth.Token = fbToken.Value
			option.PosX = form.PosX
			option.PosZ = form.PosZ
			option.PosY = form.PoxY
			option.StructureName = form.StructureName
			option.RentalServerCode = form.ServerCode
			option.RentalServerPassword = form.ServerPasswd
		}()
		//文件地址
		workDir, _ := os.Getwd()
		wrapperExec := path.Join(workDir, "builder_wrapper")
		// 转化目录名字
		convertDir := path.Join(workDir, "converted")
		compileBuildExecArgs := func(option *model.BuildOption, fbToken string) (args []string, err error) {
			if fbToken == "" {
				return nil, fmt.Errorf("fbtoken not provided")
			}
			if option.RentalServerCode == "" {
				return nil, fmt.Errorf("rental server code not provided")
			}
			args = []string{
				"--convert-dir", convertDir,
				"--file", option.StructureName,
				"--user-token", fbToken,
				"--server", option.RentalServerCode,
				"--pos", fmt.Sprintf("[%v,%v,%v]", option.PosX, option.PosY, option.PosZ),
			}
			if option.RentalServerPassword != "" {
				args = append(args, "--server-password", option.RentalServerPassword)
			}
			return args, nil
		}
		// TODO: 接受建筑路径作为建筑名字
		// TODO: client来跑
		fileSize := int64(0)
		{
			ok := false
			infoCopy, err2 := cm.filer.GetFileGroupHeadData(user.UserName, form.FileGroupName)
			if err2 != nil {
				utils.Error(err2.Error())
				model.BackErrorByString(ctx, fmt.Sprintf("cant get file group structs infos"))
				return
			}
			for _, v := range infoCopy {
				if v.Name() == option.StructureName {
					fileSize = v.Size()
					ok = true
				}
			}
			if !ok {
				model.BackErrorByString(ctx, fmt.Sprintf("cant find the %v file from fileGroup", form.StructureName))
				return
			}
		}
		// 这里structureName
		user_path := cm.filer.GetUserPath(user.UserName)
		filepath.Join(user_path, form.FileGroupName)
		file_path := filepath.Join(user_path, form.StructureName)
		option.StructureName = file_path
		args, err := compileBuildExecArgs(option, fbToken.Value)
		if err != nil {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}

		resultChan := make(chan struct{ instanceID, err string }, 1)
		//utils.Info("args:%v", args)
		//utils.Info("taskName:%v,wrapperExec:%v,args:%v", option.TaskName, wrapperExec, args)
		client.Run(option.TaskName, wrapperExec, args, func(instanceID string, err string) {
			resultChan <- struct{ instanceID, err string }{instanceID, err}
		})

		r := <-resultChan
		if r.err != "" {
			model.BackErrorByString(ctx, "不能创建导入任务")
			utils.Error("不能创建导入任务")
			return
		}
		buildInstanceDetail := &model.BuildTaskInfo{
			Time:        time.Now().String(),
			InstanceID:  r.instanceID,
			StartArgs:   args,
			BuildOption: option,
			FileSize:    fileSize,
		}
		//m.AppendBuildInstanceDetail(a, buildInstanceDetail)
		//c.JSON(200, gin.H{"instance_id": r.instanceID})
		//model.BackSuccess(ctx, buildInstanceDetail.InstanceID)
		utils.Info("%v 任务id", buildInstanceDetail.InstanceID)
		model.BackSuccess(ctx, fmt.Sprintf("%v", buildInstanceDetail.InstanceID))
		return

	}
}

// compileBuildExecArgs = func(option *BuildOption, fbToken string) (args []string, err error) {
// 	if fbToken == "" {
// 		return nil, fmt.Errorf("fbtoken not provided")
// 	}
// 	if option.RentalServerCode == "" {
// 		return nil, fmt.Errorf("rental server code not provided")
// 	}
// 	args = []string{
// 		"--convert-dir", convertDir,
// 		"--file", option.StructureName,
// 		"--user-token", fbToken,
// 		"--server", option.RentalServerCode,
// 		"--pos", fmt.Sprintf("[%v,%v,%v]", option.PosX, option.PosY, option.PosZ),
// 	}
// 	if option.RentalServerPassword != "" {
// 		args = append(args, "--server-password", option.RentalServerPassword)
// 	}
// 	return args, nil
// }
