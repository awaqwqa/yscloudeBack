package controller

import (
	"yscloudeBack/source/app/db"
	"yscloudeBack/source/app/model"

	"github.com/gin-gonic/gin"
)

func LoadHandler(db *db.DbManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
