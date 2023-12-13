package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
	"yscloudeBack/source/app/filer"
	"yscloudeBack/source/app/middleware"
	"yscloudeBack/source/app/model"
)

func (cm *ControllerMannager) GetStructs() gin.HandlerFunc {
	db := cm.GetDbManager()
	return func(ctx *gin.Context) {
		userName, isok := ctx.Get(middleware.ContextName)
		if !isok {
			model.BackError(ctx, model.CodeInvalidToken)
			return
		}
		value, isok := userName.(string)
		if !isok {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		user, err := db.GetUserByUserName(value)
		if err != nil {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		infoCopy, err := user.GetAllStructureInfoCopy()
		if err != nil {
			model.BackError(ctx, model.CodeUnknowError)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code":    model.CodeSuccess,
			"msg":     model.CodeSuccess.Msg(),
			"structs": infoCopy,
		})
	}
}
func (cm *ControllerMannager) UploadFile() gin.HandlerFunc {
	//db := cm.GetDbManager()
	return func(c *gin.Context) {
		file_group := c.PostForm("file_group")
		//utils.Info("upload file : file_group >> %v", file_group)
		name, err := middleware.GetContextName(c)
		if err != nil {
			model.BackErrorByString(c, err.Error())
			return
		}
		//user, err := db.GetUserByUserName(name)
		//if err != nil {
		//	model.BackError(c, model.CodeGetUserFalse)
		//	return
		//}
		file, err := c.FormFile("file")
		if err != nil {
			model.BackErrorByString(c, fmt.Sprintf("Bad request: %s", err.Error()))
			return
		}

		// 检查文件大小
		if file.Size <= 0 || file.Size > 5<<20 { // 5<<20 是 5MB
			model.BackErrorByString(c, fmt.Sprintf("File size should be between 0 and 5 MB."))
			return
		}

		// 检查文件后缀名
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".schematic" && ext != ".bdx" {
			model.BackErrorByString(c, "Invalid file extension. Only .schematic and .bdx are allowed.")
			return
		}

		// 处理文件（例如保存到磁盘）
		fileData := make([]byte, file.Size)
		src, err := file.Open()
		if err != nil {
			model.BackErrorByString(c, "cant get file ")
			return
		}
		defer src.Close()
		_, err = src.Read(fileData)
		if err != nil {
			model.BackErrorByString(c, err.Error())
			return
		}
		user_path := cm.filer.GetUserPath(name)
		err = filer.CheckUserPath(user_path)
		if err != nil {
			err = nil
			err := filer.CreateFolder(user_path)
			if err != nil {
				model.BackErrorByString(c, err.Error())
				return
			}
		}
		file_group_path := filepath.Join(user_path, file_group)
		err = filer.CheckFileGroupPath(file_group_path)
		if err != nil {
			model.BackErrorByString(c, err.Error())
			return
		}
		if !filer.IsPathValid(file.Filename) {
			model.BackErrorByString(c, "this file path is not allowed")
			return
		}
		err = cm.filer.NewTask(name, file_group, file.Filename, fileData)
		if err != nil {
			model.BackErrorByString(c, err.Error())
			return
		}
		//structure, err := user.NewUserStructure(file.Filename, fileData, file_group)
		//if err != nil {
		//	model.BackErrorByString(c, "cant upload file into dir")
		//	return
		//}
		//err = db.AddStructure(structure)
		//if err != nil {
		//	model.BackErrorByString(c, "cant upload structure into db")
		//	return
		//}
		//dbStructure, err := db.GetStructureByHash(structure.FileHash)
		//if err != nil {
		//	model.BackErrorByString(c, "bind structure false")
		//	return
		//}
		//err = db.AssociateStuctureWithUser(user.ID, dbStructure.ID)
		//if err != nil {
		//	model.BackErrorByString(c, "bind structure false")
		//	return
		//}
		model.BackSuccess(c, fmt.Sprintf("File %s uploaded successfully with size of %d.", file.Filename, file.Size))
		// 返回成功响应
	}
}
