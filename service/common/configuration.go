package common

import (
	"log"
	"os"
	"encoding/json"
)

type RedisConf struct {
	Server string `json:"server"`
	TokenDB int `json:"tokenDB"`
	TokenExpired string `json:"tokenExpired"`
	Password string `json:"password"`
	FlowInstanceDB int `json:"flowInstanceDB"`
	FlowInstanceExpired string `json:"flowInstanceExpired"`
}

type mysqlConf struct {
	Server string `json:"server"`
	Password string `json:"password"`
	User string `json:"user"`
	DBName string `json:"dbName"`
	ConnMaxLifetime int `json:"connMaxLifetime"` 
  MaxOpenConns int `json:"maxOpenConns"`
  MaxIdleConns int `json:"maxIdleConns"`
}

type serviceConf struct {
	Port string `json:"port"`
}

type fileConf struct {
	Root string `json:"root"`
}

type runtimeConf struct {
	GoMaxProcs int `json:"goMaxProcs"`
}

type MqttConf struct {
	Broker string `json:"broker"`
	Port int `json:"port"`
	WSPort int `json:"wsPort"`
	Password string `json:"password"`
	User string `json:"user"`
	ClientID string `json:"clientID"`
	StartFlowTopic string `json:"startFlowTopic"`
	TaskNotificationTopic string `json:"taskNotificationTopic"`
}

type Config struct {
	Redis  RedisConf  `json:"redis"`
	Mysql  mysqlConf  `json:"mysql"`
	Service serviceConf `json:"service"`
	Runtime runtimeConf `json:"runtime"`
	Mqtt MqttConf `json:"mqtt"`
}

var gConfig Config

func InitConfig(confFile string)(*Config){
	log.Println("init configuation start ...")
	fileName := confFile
	filePtr, err := os.Open(fileName)
	if err != nil {
        log.Fatal("Open file failed [Err:%s]", err.Error())
    }
    defer filePtr.Close()

	// 创建json解码器
    decoder := json.NewDecoder(filePtr)
    err = decoder.Decode(&gConfig)
	if err != nil {
		log.Println("json file decode failed [Err:%s]", err.Error())
	}
	log.Println("init configuation end")
	return &gConfig
}

func GetConfig()(*Config){
	return &gConfig
}