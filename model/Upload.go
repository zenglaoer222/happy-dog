package model

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"happy-dog/utils"
	"happy-dog/utils/errmsg"
	"mime/multipart"
	"strconv"
)

var AccessKey = utils.AccessKey
var SecretKey = utils.SecretKey
var Bucket = utils.Bucket
var Img = utils.QiniuServer

func UploadFile(file multipart.File, fileSize int64, user_id int) (string, int) {
	putPolicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AccessKey, SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false, // 是否使用https域名
	}

	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	// 上传到customer文件夹中
	key := "customer/" + strconv.Itoa(user_id) + "_avatar.jpg"
	//err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	err = formUploader.Put(context.Background(), &ret, upToken, key, file, fileSize, &putExtra)
	if err != nil {
		return "", errmsg.ERROR
	}
	url := Img + "customer" + "/" + strconv.Itoa(user_id) + "_avatar.jpg"

	code := uploadAvtar(DB, user_id, url)
	if code == errmsg.ERROR {
		return "", code
	}

	return url, errmsg.SUCCESS
}

func UploadFileForShop(file multipart.File, fileSize int64, user_id int) (string, int) {
	// 上传到customer文件夹中
	key := "shop/" + strconv.Itoa(user_id) + "_front.jpg"
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", Bucket, key), // 仅允许覆盖当前 key
	}
	mac := qbox.NewMac(AccessKey, SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false, // 是否使用https域名
	}

	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	//err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	err := formUploader.Put(context.Background(), &ret, upToken, key, file, fileSize, &putExtra)
	if err != nil {
		fmt.Println("七牛云上传失败: %v", err) // 使用日志库记录错误
		return "", errmsg.ERROR
	}
	url := Img + "shop" + "/" + strconv.Itoa(user_id) + "_front.jpg"

	code := uploadFrontUrl(DB, user_id, url)
	if code == errmsg.ERROR {
		return "", code
	}

	return url, errmsg.SUCCESS
}

func UploadFileForProduct(file multipart.File, fileSize int64, product_id int) (string, int) {
	// 上传到product文件夹中
	key := "product/" + strconv.Itoa(product_id) + "_.jpg"
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", Bucket, key), // 仅允许覆盖当前 key
	}
	mac := qbox.NewMac(AccessKey, SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false, // 是否使用https域名
	}

	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	//err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	err := formUploader.Put(context.Background(), &ret, upToken, key, file, fileSize, &putExtra)
	if err != nil {
		fmt.Println("七牛云上传失败: %v", err) // 使用日志库记录错误
		return "", errmsg.ERROR
	}
	url := Img + "product" + "/" + strconv.Itoa(product_id) + "_.jpg"

	code := uploadProductUrl(DB, product_id, url)
	if code == errmsg.ERROR {
		return "", code
	}

	return url, errmsg.SUCCESS
}

//http://so0hvjygo.hn-bkt.clouddn.com/customer/w.jpg
//http://so0hvjygo.hn-bkt.clouddn.com/customer/customer/w.jpg
