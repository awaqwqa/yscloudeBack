package controller

import (
	"fmt"
	"os"
	"path"
	"time"
	"yscloudeBack/source/app/middleware"
	"yscloudeBack/source/app/model"

	"github.com/gin-gonic/gin"
)

func (cm *ControllerMannager) LoadHandler() gin.HandlerFunc {
	db := cm.GetDbManager()
	client := cm.GetCluster()
	return func(ctx *gin.Context) {
		//绑定参数
		option := &model.BuildOption{}
		code, err := model.BindStruct(ctx, &option)
		if err != nil {
			model.BackError(ctx, code)
			return
		}

		// 获取 user
		value, isFind := ctx.Get(middleware.ContextName)
		if !isFind {
			model.BackError(ctx, model.CodeGetUserFalse)
			return
		}
		userName := value.(string)
		user, err := db.GetUserByUserName(userName)
		if err != nil {
			model.BackError(ctx, model.CodeGetUserFalse)
			return
		}
		//获取fbToken
		fbTokens, err := db.GetFbTokens()
		if err != nil {
			model.BackError(ctx, model.CodeGetFbTokenFalse)
			return
		}
		var fbToken string
		if len(fbTokens) > 0 {
			fbToken = fbTokens[0].Value
		}

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
			infoCopy, err := user.GetAllStructureInfoCopy()
			if err != nil {
				model.BackError(ctx, model.CodeUnknowError)
				return
			}
			for _, v := range infoCopy {
				if v.FileName == option.StructureName {
					fileSize = v.FileSize
					ok = true
				}
			}
			if !ok {
				model.BackError(ctx, model.CodeUnknowError)
				return
			}
		}
		option.StructureName = path.Join(user.GetDirPath(), option.StructureName)
		if args, err := compileBuildExecArgs(option, fbToken); err != nil {
			model.BackError(ctx, model.CodeUnknowError)
			return
		} else {
			resultChan := make(chan struct{ instanceID, err string }, 1)
			client.Run(option.TaskName, wrapperExec, args, func(instanceID string, err string) {
				resultChan <- struct{ instanceID, err string }{instanceID, err}
			})
			r := <-resultChan
			if r.err == "" {
				buildInstanceDetail := &model.BuildTaskInfo{
					Time:        time.Now().String(),
					InstanceID:  r.instanceID,
					StartArgs:   args,
					BuildOption: option,
					FileSize:    fileSize,
				}
				//m.AppendBuildInstanceDetail(a, buildInstanceDetail)
				//c.JSON(200, gin.H{"instance_id": r.instanceID})
				model.BackSuccess(ctx, buildInstanceDetail.InstanceID)
			} else {
				model.BackErrorByString(ctx, "不能创建导入任务")
				return
			}
		}
		//compileBuildExecArgs := func(option *model.BuildOption, fbToken string) (args []string, err error) {
		//	if fbToken == "" {
		//		return nil, fmt.Errorf("fbtoken not provided")
		//	}
		//	if option.RentalServerCode == "" {
		//		return nil, fmt.Errorf("rental server code not provided")
		//	}
		//	args = []string{
		//		"--convert-dir", convertDir,
		//		"--file", option.StructureName,
		//		"--user-token", fbToken,
		//		"--server", option.RentalServerCode,
		//		"--pos", fmt.Sprintf("[%v,%v,%v]", option.PosX, option.PosY, option.PosZ),
		//	}
		//	if option.RentalServerPassword != "" {
		//		args = append(args, "--server-password", option.RentalServerPassword)
		//	}
		//	return args, nil
		//}
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
