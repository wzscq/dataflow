package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
	"dataflow/common"
	"dataflow/flow"
	"dataflow/data"
    "dataflow/test"
    "log"
    "time"
    "runtime"
    "os"
)

func main() {
    confFile:="conf/conf.json"
    if len(os.Args)>1 {
        confFile=os.Args[1]
        log.Println(confFile)
    }
    //初始化配置
    conf:=common.InitConfig(confFile)
    //设置启动线程数量
    runtime.GOMAXPROCS(conf.Runtime.GoMaxProcs)

	//设置log打印文件名和行号
    log.SetFlags(log.Lshortfile | log.LstdFlags)

    //初始化时区
    var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone

	//初始化路由
	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowAllOrigins:true,
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"*"},
        AllowCredentials: true,
    }))

	dataRepo:=&data.DefatultDataRepository{}
    dataRepo.Connect(
        conf.Mysql.Server,
        conf.Mysql.User,
        conf.Mysql.Password,
        conf.Mysql.DBName,
        conf.Mysql.ConnMaxLifetime,
        conf.Mysql.MaxOpenConns,
        conf.Mysql.MaxIdleConns)

    flowExpired,_:=time.ParseDuration(conf.Redis.FlowInstanceExpired)
    flowInstanceRepository:=&flow.DefaultFlowInstanceRepository{}
    flowInstanceRepository.Init(conf.Redis.Server,conf.Redis.FlowInstanceDB,flowExpired,conf.Redis.Password)
	//初始化流控制器
	flowController:=&flow.FlowController{
		DataRepository:dataRepo,
        InstanceRepository:flowInstanceRepository,
        Mqtt:conf.Mqtt,
	}
    flowController.Bind(router)

    //测试数据生成
    testController:=&test.TestController{
		DataRepository:dataRepo,
	}
    testController.Bind(router)

	//启动监听服务
	router.Run(conf.Service.Port)
}