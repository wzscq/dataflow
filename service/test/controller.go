package test

import (
	"log"
	"github.com/gin-gonic/gin"
	"buoyancyinfo.com/matchflow/common"
	"buoyancyinfo.com/matchflow/data"
	"net/http"
	"math/rand"
	"fmt"
	"time"
)

type testReq struct {
	UserID string `json:"userID"`
	AppDB string `json:"appDB"`
}

type TestController struct {
	DataRepository data.DataRepository
}

func (controller *TestController) createRandomDeliveryCS(appDB,userID string,r *rand.Rand,id string)(int){
	row:=map[string]interface{}{}
	//根据配置生成分组记录数据
	row["id"]=id
	row["purchaser"]=fmt.Sprintf("p_%d",r.Intn(20))
	row["seller"]=fmt.Sprintf("s_%d",r.Intn(20))
	row["po"]=fmt.Sprintf("po_%d",r.Intn(20))
	row["so"]=fmt.Sprintf("so_%d",r.Intn(20))
	row["delivery_no"]=fmt.Sprintf("dn_%d",r.Intn(20))
	row["material"]=fmt.Sprintf("m_%d",r.Intn(20))
	price:=r.Intn(100)
	quantity:=r.Intn(100)
	amount:=price*quantity
	row["price"]=fmt.Sprintf("%d",price)
	row["quantity"]=fmt.Sprintf("%d",quantity)
	row["amount"]=fmt.Sprintf("%d",amount)
	row["match_status"]="0"

	row[data.SAVE_TYPE_COLUMN]=data.SAVE_CREATE

	saver:=data.Save{
		List:&[]map[string]interface{}{
			row,
		},
		AppDB:appDB,
		UserID:userID,
		ModelID:"match_deliveries_cs",
	}

	_,errorCode:=saver.Execute(controller.DataRepository)
	return errorCode
}

func (controller *TestController) createRandomDelivery(appDB,userID string,r *rand.Rand)(int){
	row:=map[string]interface{}{}
	//根据配置生成分组记录数据
	row["purchaser"]=fmt.Sprintf("p_%d",r.Intn(20))
	row["seller"]=fmt.Sprintf("s_%d",r.Intn(20))
	row["po"]=fmt.Sprintf("po_%d",r.Intn(20))
	row["so"]=fmt.Sprintf("so_%d",r.Intn(20))
	row["delivery_no"]=fmt.Sprintf("dn_%d",r.Intn(20))
	row["material"]=fmt.Sprintf("m_%d",r.Intn(20))
	price:=r.Intn(100)
	quantity:=r.Intn(100)
	amount:=price*quantity
	row["price"]=fmt.Sprintf("%d",price)
	row["quantity"]=fmt.Sprintf("%d",quantity)
	row["amount"]=fmt.Sprintf("%d",amount)
	row["match_status"]="0"

	row[data.SAVE_TYPE_COLUMN]=data.SAVE_CREATE

	saver:=data.Save{
		List:&[]map[string]interface{}{
			row,
		},
		AppDB:appDB,
		UserID:userID,
		ModelID:"match_deliveries",
	}

	_,errorCode:=saver.Execute(controller.DataRepository)
	return errorCode
}

func (controller *TestController) createRandomReceipt(appDB,userID string,r *rand.Rand)(int){
	row:=map[string]interface{}{}
	//根据配置生成分组记录数据
	row["purchaser"]=fmt.Sprintf("p_%d",r.Intn(20))
	row["seller"]=fmt.Sprintf("s_%d",r.Intn(20))
	row["po"]=fmt.Sprintf("po_%d",r.Intn(20))
	row["delivery_no"]=fmt.Sprintf("dn_%d",r.Intn(20))
	row["receipt_no"]=fmt.Sprintf("rc_%d",r.Intn(20))
	row["material"]=fmt.Sprintf("m_%d",r.Intn(20))
	price:=r.Intn(100)
	quantity:=r.Intn(100)
	amount:=price*quantity
	row["price"]=fmt.Sprintf("%d",price)
	row["quantity"]=fmt.Sprintf("%d",quantity)
	row["amount"]=fmt.Sprintf("%d",amount)
	row["match_status"]="0"

	row[data.SAVE_TYPE_COLUMN]=data.SAVE_CREATE

	saver:=data.Save{
		List:&[]map[string]interface{}{
			row,
		},
		AppDB:appDB,
		UserID:userID,
		ModelID:"match_receipts",
	}

	_,errorCode:=saver.Execute(controller.DataRepository)
	return errorCode
}

func (controller *TestController) generateData(c *gin.Context){
	log.Println("TestController generateData start")

	/*count:=0
	for i := 0; i < 20000000000; i++ {
		count+=1
	}
	log.Printf("TestController generateData end count %d \n",count)
	return;*/

	var req testReq
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestController generateData end")
		return
	}

	r := rand.New(rand.NewSource(999))
	for i := 0; i < 20000; i++ {
		//生成收货单数据
		/*errorCode:=controller.createRandomReceipt(req.AppDB,req.UserID,r)
		if errorCode!=common.ResultSuccess {
			rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			log.Println("TestController generateData end")
			return
		}
		errorCode=controller.createRandomDelivery(req.AppDB,req.UserID,r)
		if errorCode!=common.ResultSuccess {
			rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			log.Println("TestController generateData end")
			return
		}*/

		id:= time.Now().Format("20060102150405")
		id=id+fmt.Sprintf("_%d",i)
		errorCode:=controller.createRandomDeliveryCS(req.AppDB,req.UserID,r,id)
		if errorCode!=common.ResultSuccess {
			rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			log.Println("TestController generateData end")
			return
		}
	} 
	
	rsp:=common.CreateResponse(nil,nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("TestController generateData end")
}

func (controller *TestController) Bind(router *gin.Engine) {
	log.Println("Bind TestController")
	router.POST("/test/generateData", controller.generateData)
}