package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
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
	db := cm.GetDbManager()
	return func(c *gin.Context) {
		if db == nil {
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request: %s", err.Error())
			return
		}

		// 检查文件大小
		if file.Size <= 0 || file.Size > 5<<20 { // 5<<20 是 5MB
			c.String(http.StatusBadRequest, "File size should be between 0 and 5 MB.")
			return
		}

		// 检查文件后缀名
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".schematic" && ext != ".bdx" {
			c.String(http.StatusBadRequest, "Invalid file extension. Only .schematic and .bdx are allowed.")
			return
		}

		// 处理文件（例如保存到磁盘）
		// ...

		// 返回成功响应
		c.String(http.StatusOK, "File %s uploaded successfully with size of %d.", file.Filename, file.Size)
	}
}
