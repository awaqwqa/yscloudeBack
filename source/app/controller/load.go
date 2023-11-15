package controller

import (
	"fmt"
	"os"
	"path"
	"yscloudeBack/source/app/cluster"
	"yscloudeBack/source/app/db"
	"yscloudeBack/source/app/middleware"
	"yscloudeBack/source/app/model"

	"github.com/gin-gonic/gin"
)

func LoadHandler(db *db.DbManager, client *cluster.ClusterRequester) gin.HandlerFunc {
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
		option.OptionalFBToken = fbToken
		//文件地址
		workDir, _ := os.Getwd()
		wrapperExec := path.Join(workDir, "builder_wrapper")
		convertDir := path.Join(workDir, "converted")
		// 文件
		build_option := &model.BuildOption{}
		code, err2 := model.BindStruct(ctx, &build_option)
		if err2 != nil {
			model.BackError(ctx, code)
			return
		}
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

		var lf *model.LoadForm
		code, err := model.BindStruct(ctx, &lf)
		if err != nil {
			model.BackError(ctx, code)
			return
		}
		model.BackError(ctx, model.CodeUnknowError)
		return
		//TODO: 一些导入的处理
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
