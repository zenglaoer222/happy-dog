package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	AppMode  string
	HttpPort string

	JwtKey string

	Db          string
	Dbhost      string
	Dbport      string
	Dbuser      string
	Dbpass      string
	Dbname      string
	AccessKey   string
	SecretKey   string
	Bucket      string
	QiniuServer string

	RedisAddr string
	RedisPwd  string
	RedisDb   int
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径是否正确", err)
	}
	fmt.Println("init")
	LoadServer(file)
	LoadData(file)
	LoadQiniu(file)
	LoadRedis(file)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
	JwtKey = file.Section("server").Key("JwtKey").MustString("jfje92mf")
}

func LoadData(file *ini.File) {
	Db = file.Section("database").Key("Db").MustString("mysql")
	Dbhost = file.Section("database").Key("Dbhost").MustString("127.0.0.1")
	Dbport = file.Section("database").Key("Dbport").MustString("3306")
	Dbuser = file.Section("database").Key("Dbuser").MustString("ginblog")
	Dbpass = file.Section("database").Key("Dbpass").MustString("admin123")
	Dbname = file.Section("database").Key("Dbname").MustString("ginblog")
}

func LoadQiniu(file *ini.File) {
	AccessKey = file.Section("qiniu").Key("AccessKey").String()
	SecretKey = file.Section("qiniu").Key("SecretKey").String()
	Bucket = file.Section("qiniu").Key("Bucket").String()
	QiniuServer = file.Section("qiniu").Key("QiniuServer").String()
}

func LoadRedis(file *ini.File) {
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPwd = file.Section("redis").Key("RedisPwd").String()
	RedisDb, _ = file.Section("redis").Key("RedisDb").Int()
}
