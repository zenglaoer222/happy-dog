package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"happy-dog/model"
	"happy-dog/utils/errmsg"
	"net/http"
	"strconv"
)

func CreateProduct(c *gin.Context) {
	var code int
	sid, ok := c.Get("user_id")
	shop_id, _ := sid.(uint)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status": errmsg.ERROR,
			"msg":    errmsg.GetErrMsg(errmsg.ERROR),
		})
		return
	}

	// 1. 接收文件
	file, fileHeader, err := c.Request.FormFile("file_image")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": errmsg.ERROR,
			"msg":    errmsg.GetErrMsg(errmsg.ERROR),
		})
		return
	}

	fileSize := fileHeader.Size

	// 3. 获取其他表单参数
	name := c.PostForm("name")
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
	stockNum, _ := strconv.Atoi(c.PostForm("stock_num"))

	// 4. 创建商品对象
	product := model.Product{
		Name:     name,
		Price:    price,
		ShopId:   int(shop_id),
		StockNum: stockNum,
		Picture:  "", // 保存文件路径
	}

	// 5. 创建数据库记录
	code = model.CreateProduct(model.DB, &product)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}

	//修改url
	code, lastId := model.GetLastId(model.DB)
	if code == errmsg.ERROR {
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}

	var url string
	url, code = model.UploadFileForProduct(file, fileSize, int(lastId))

	product.Picture = url

	msg := errmsg.GetErrMsg(code)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"msg":    msg,
		"data":   product,
	})
}

func GetProductList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	code, productList, total := model.GetProductList(model.DB, uint(id), pageNum)
	if code == errmsg.ERROR {
		c.JSON(200, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
		"data":    productList,
		"total":   total,
	})
}

func DeleteProduct(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Param("pid"))

	code, _, sid := model.ExistProduct(model.DB, uint(pid))

	shop_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(200, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
		})
		return
	}
	fmt.Println(shop_id)
	fmt.Println(sid)

	if uint(sid) != shop_id.(uint) {
		c.JSON(200, gin.H{
			"status":  errmsg.ERROR,
			"message": "权限不匹配",
		})
		return
	}
	code = model.DeleteProduct(model.DB, pid)
	if code == errmsg.ERROR {
		c.JSON(200, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})

}
