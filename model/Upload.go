package model

import (
	"context"
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
	url := Img + "/" + strconv.Itoa(user_id) + "_avatar.jpg"

	code := uploadAvtar(DB, user_id, url)
	if code == errmsg.ERROR {
		return "", code
	}
	
	return url, errmsg.SUCCESS
}

//http://so0hvjygo.hn-bkt.clouddn.com/customer/w.jpg
//http://so0hvjygo.hn-bkt.clouddn.com/customer/customer/w.jpg
