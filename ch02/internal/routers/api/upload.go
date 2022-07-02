package api

import (
	"demo/ch02/global"
	"demo/ch02/internal/service"
	"demo/ch02/pkg/app"
	"demo/ch02/pkg/convert"
	"demo/ch02/pkg/errcode"
	"demo/ch02/pkg/upload"
	"github.com/gin-gonic/gin"
)

type Upload struct{}

func NewUpload() Upload {
	return Upload{}
}

func (u Upload) UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	// 参数获取与检测
	// 通过 c.Request.FormFile() 读取入参 file 字段的上传文件信息
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}
	// 使用入参 type 字段作为上传文件类型的确认依据
	fileType := convert.StrTo(c.PostForm("type")).MustInt()
	if fileHeader == nil || fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	// 调用 service 方法完成文件上传、文件保存，并返回文件展示地址
	svc := service.New(c.Request.Context())
	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		global.Logger.Errorf(c, "svc.UploadFile err: %v", err)
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}

	// 返回文件展示地址
	response.ToResponse(gin.H{
		"file_access_url": fileInfo.AccessUrl,
	})
}
