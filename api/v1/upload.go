package v1

import (
	"github.com/gin-gonic/gin"
	"happy-dog/model"
	"happy-dog/utils/errmsg"
	"net/http"
)

func Upload(c *gin.Context) {
	file, fileHeader, _ := c.Request.FormFile("file")

	fileSize := fileHeader.Size
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": errmsg.ERROR,
			"url":  "",
			"msg":  errmsg.GetErrMsg(errmsg.ERROR),
		})
	}

	url, code := model.UploadFile(file, fileSize, int(user_id.(uint)))
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"url":  url,
	})
}
