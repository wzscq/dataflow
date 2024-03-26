package common

import (
	"log"
	"github.com/spf13/viper"
)

type RedisConf struct {
	Server string `json:"server" mapstructure:"server"`
	TokenDB int `json:"tokenDB" mapstructure:"tokenDB"`
	TokenExpired string `json:"tokenExpired" mapstructure:"tokenExpired"`
	Password string `json:"password" mapstructure:"password"`
	FlowInstanceDB int `json:"flowInstanceDB" mapstructure:"flowInstanceDB"`
	FlowInstanceExpired string `json:"flowInstanceExpired" mapstructure:"flowInstanceExpired"`
	TLS			    string   `json:"tls" mapstructure:"tls"`
}

type mysqlConf struct {
	Server string `json:"server" mapstructure:"server"`
	Password string `json:"password" mapstructure:"password"`
	User string `json:"user" mapstructure:"user"`
	DBName string `json:"dbName" mapstructure:"dbName"`
	ConnMaxLifetime int `json:"connMaxLifetime" mapstructure:"connMaxLifetime"` 
  MaxOpenConns int `json:"maxOpenConns" mapstructure:"maxOpenConns"`
  MaxIdleConns int `json:"maxIdleConns" mapstructure:"maxIdleConns"`
  TLS			    string   `json:"tls" mapstructure:"tls"`
}

type serviceConf struct {
	Port string `json:"port" mapstructure:"port"`
}

type fileConf struct {
	Root string `json:"root" mapstructure:"root"`
}

type runtimeConf struct {
	GoMaxProcs int `json:"goMaxProcs" mapstructure:"goMaxProcs"`
}

type MqttConf struct {
	Broker string `json:"broker" mapstructure:"broker"`
	Port int `json:"port" mapstructure:"port"`
	WSPort int `json:"wsPort" mapstructure:"wsPort"`
	Password string `json:"password" mapstructure:"password"`
	User string `json:"user" mapstructure:"user"`
	ClientID string `json:"clientID" mapstructure:"clientID"`
	StartFlowTopic string `json:"startFlowTopic" mapstructure:"startFlowTopic"`
	TaskNotificationTopic string `json:"taskNotificationTopic" mapstructure:"taskNotificationTopic"`
}

type Config struct {
	Redis  RedisConf  `json:"redis" mapstructure:"redis"`
	Mysql  mysqlConf  `json:"mysql" mapstructure:"mysql"`
	Service serviceConf `json:"service" mapstructure:"service"`
	Runtime runtimeConf `json:"runtime" mapstructure:"runtime"`
	Mqtt MqttConf `json:"mqtt" mapstructure:"mqtt"`
}

var gConfig Config

func InitConfig(confFile string) *Config {
	log.Println("init configuation start ...")

	viper.SetDefault("mysql.tls", "false")
	viper.SetDefault("redis.tls", "false")

	viper.BindEnv("redis.server", "DATAFLOW_REDIS_SERVER")
	viper.BindEnv("redis.password", "DATAFLOW_REDIS_PASSWORD")
	viper.BindEnv("redis.tls", "DATAFLOW_REDIS_TLS")
	viper.BindEnv("mysql.server", "DATAFLOW_MYSQL_SERVER")
	viper.BindEnv("mysql.password", "DATAFLOW_MYSQL_PASSWORD")
	viper.BindEnv("mysql.user", "DATAFLOW_MYSQL_USER")
	viper.BindEnv("mysql.tls", "DATAFLOW_MYSQL_TLS")
	viper.BindEnv("mqtt.broker", "DATAFLOW_MQTT_BROKER")
	viper.BindEnv("mqtt.port", "DATAFLOW_MQTT_PORT")
	viper.BindEnv("mqtt.wsPort", "DATAFLOW_MQTT_WSPORT")
	viper.BindEnv("mqtt.password", "DATAFLOW_MQTT_PASSWORD")
	viper.BindEnv("mqtt.user", "DATAFLOW_MQTT_USER")

	viper.SetConfigFile(confFile)

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("ReadInConfig failed [Err:" + err.Error() + "]")
		return nil
	}

	err = viper.Unmarshal(&gConfig)
	if err != nil {
		log.Println("Unmarshal failed [Err:" + err.Error() + "]")
		return nil
	}
	log.Println("init configuation end")
	return &gConfig
}

/*func InitConfig(confFile string)(*Config){
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
}*/

func GetConfig()(*Config){
	return &gConfig
}