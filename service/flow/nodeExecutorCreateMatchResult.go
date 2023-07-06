package flow

import (
    "time"
	"dataflow/common"
	"dataflow/data"
	"log"
	"encoding/json"
	"strconv"
	"fmt"
	"github.com/rs/xid"
)

const (
	AGGREGATION_FIRST = "first"
	AGGREGATION_SUM = "sum"
)

const (
	CC_MATCH_GROUP = "match_group"
	CC_MATCH_STATUS = "match_status"
	CC_MATCH_MESSAGE = "match_failure_reason"
	CC_MATCH_RESULT = "match_result"
)

const (
	MATCH_STATUS_NONE = "0"  //未比对
	MATCH_STATUS_SUCCESS = "1"  //比对成功
	MATCH_STATUS_FAILURE = "2" //比对失败
	MATCH_STATUS_REVOKED = "3" //比对撤销
)

type groupModelField struct {
	Field string `json:"field"` 
	SourceModel string `json:"sourceModel"` 
	SourceField string `json:"sourceField"` 
	Aggregation string `json:"aggregation"` 
}

type groupModelConfig struct {
	ModelID string `json:"modelID"` 
	Fields []groupModelField `json:"fields"` 
}

type createMatchResultConfig struct {
	MatchResult string  `json:"matchResult"`  
	GroupModel groupModelConfig `json:"groupModel"`
}

type nodeExecutorCreateMatchResult struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorCreateMatchResult)getNodeConfig()(*createMatchResultConfig){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
    if err != nil {
        log.Println(err)
		return nil
    }
	conf:=&createMatchResultConfig{}
    if err := json.Unmarshal(jsonStr, conf); err != nil {
        log.Println(err)
		return nil
    }

	return conf
}

func (nodeExecutor *nodeExecutorCreateMatchResult)aggregationFirst(
	groupRow,modelRow map[string]interface{},
	groupField *groupModelField,
	index int){

	if index == 0 {
		value,ok:=modelRow[groupField.SourceField]
		if ok {
			groupRow[groupField.Field]=value
		}
	}
}

func (nodeExecutor *nodeExecutorCreateMatchResult)aggregationSum(
	groupRow,modelRow map[string]interface{},
	groupField *groupModelField,
	index int){
	
	//获取之前的数据
	preVal,ok:=groupRow[groupField.Field]
	if !ok {
		log.Println("nodeExecutorSaveMatched aggragationSum source field "+groupField.Field+" not found")
		preVal="0"
	}
	spreVal,_:=preVal.(string)
	preValFloat64,err:=strconv.ParseFloat(spreVal, 64)
	if err !=nil {
		log.Println("nodeExecutorSaveMatched aggragationSum can not convert aggregeted value to float64")
		return
	}
	//获取当前数据
	curVal,ok:=modelRow[groupField.SourceField]
	if !ok {
		log.Println("nodeExecutorSaveMatched aggragationSum no source field:"+groupField.SourceField)
		return
	}
	scurVal,_:=curVal.(string)
	curValFloat64,err:=strconv.ParseFloat(scurVal, 64)
	if err !=nil {
		log.Println("nodeExecutorSaveMatched aggragationSum can not convert source value to float64")
		return
	}

	newVal:=preValFloat64+curValFloat64
	groupRow[groupField.Field]=fmt.Sprint(newVal)
}

func (nodeExecutor *nodeExecutorCreateMatchResult)aggregationCount(
	groupRow,modelRow map[string]interface{},
	groupField *groupModelField,
	index int){
	
	//获取之前的数据
	preVal,ok:=groupRow[groupField.Field]
	if !ok {
		log.Println("nodeExecutorSaveMatched aggregationCount source field "+groupField.Field+" not found")
		preVal="0"
	}
	spreVal,_:=preVal.(string)
	preValInt64,err:=strconv.ParseInt(spreVal,10,64)
	if err !=nil {
		log.Println("nodeExecutorSaveMatched aggregationCount can not convert aggregeted value to int64")
		return
	}
	
	newVal:=preValInt64+1
	groupRow[groupField.Field]=fmt.Sprint(newVal)
}

func (nodeExecutor *nodeExecutorCreateMatchResult)aggragation(
	groupRow,modelRow map[string]interface{},
	groupField *groupModelField,
	index int){
	
	if groupField.Aggregation == AGGREGATION_FIRST {
		nodeExecutor.aggregationFirst(groupRow,modelRow,groupField,index)
	} else if groupField.Aggregation == AGGREGATION_SUM {
		nodeExecutor.aggregationSum(groupRow,modelRow,groupField,index)
	} else if groupField.Aggregation == AGGREGATION_COUNT {
		nodeExecutor.aggregationCount(groupRow,modelRow,groupField,index)
	}
}

func (nodeExecutor *nodeExecutorCreateMatchResult)createGroup(
	dataItem *flowDataItem,
	createMatchResultConf *createMatchResultConfig,
	instance *flowInstance,
	groupIndex int,
	batchNo string)(string,*modelDataItem){

	groupID:=fmt.Sprintf("%s%05d",batchNo,groupIndex)
	
	groupRow:=map[string]interface{}{}
	//根据配置生成分组记录数据
	//迭代每个model
	for _,modelData:=range dataItem.Models {
		if modelData.List !=nil && len(*modelData.List)>0 {
			for index,modelRow:=range  (*modelData.List) {
				for _,groupField:=range createMatchResultConf.GroupModel.Fields {
					if groupField.SourceModel == *modelData.ModelID {
						nodeExecutor.aggragation(groupRow,modelRow,&groupField,index)
					}
				}
			}
		}
	}

	groupRow[CC_MATCH_RESULT]=createMatchResultConf.MatchResult
	groupRow[data.SAVE_TYPE_COLUMN]=data.SAVE_CREATE
	groupRow[data.CC_ID]=groupID

	modelID:=createMatchResultConf.GroupModel.ModelID
	groupData:=&modelDataItem{
		ModelID:&modelID,
		List:&[]map[string]interface{}{
			groupRow,
		},
	}
	
	log.Println("nodeExecutorCreateMatchResult createGroup ")
	return groupID,groupData
}

func (nodeExecutor *nodeExecutorCreateMatchResult)updateModels(
	dataItem *flowDataItem,
	groupID string,
	instance *flowInstance){

	for modelIdx,modelData:=range (dataItem.Models) {
		if modelData.List!=nil && len(*modelData.List)>0 {
			rows:=make([]map[string]interface{},len(*modelData.List))
			for index,row:=range (*modelData.List) {
				rows[index]=map[string]interface{}{}
				rows[index][data.CC_ID]=row[data.CC_ID]
				rows[index][data.CC_VERSION]=row[data.CC_VERSION]
				rows[index][data.SAVE_TYPE_COLUMN]=data.SAVE_UPDATE
				rows[index][CC_MATCH_GROUP]=groupID
				rows[index][CC_MATCH_STATUS]=MATCH_STATUS_SUCCESS
				rows[index][CC_MATCH_MESSAGE]=""
			}

			modelData.List=&rows
			dataItem.Models[modelIdx]=modelData
		}
	}
}

func (nodeExecutor *nodeExecutorCreateMatchResult)createMatchResult(
	dataItem *flowDataItem,
	createMatchResultConf *createMatchResultConfig,
	instance *flowInstance,
	groupIndex int,
	batchNo string)(*common.CommonError){
	
	log.Println("start nodeExecutorCreateMatchResult createMatchResult")
	
	//创建分组记录数据
	groupID,groupData:=nodeExecutor.createGroup(dataItem,createMatchResultConf,instance,groupIndex,batchNo)
	
	//将分组号更新到左右表，同时更新左右表数据的匹配状态
	nodeExecutor.updateModels(dataItem,groupID,instance)
	
	dataItem.Models=append(dataItem.Models,*groupData)
	
	log.Println("end nodeExecutorCreateMatchResult createMatchResult")
	return nil
}

func (nodeExecutor *nodeExecutorCreateMatchResult)getBatchNumber()(string){
	return xid.New().String()
}

func (nodeExecutor *nodeExecutorCreateMatchResult)copyDataItem(
	item,newItem *flowDataItem){
	newItem.VerifyResult=append(newItem.VerifyResult,item.VerifyResult...)
	newItem.Models=append(newItem.Models,item.Models...)
}

func (nodeExecutor *nodeExecutorCreateMatchResult)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	log.Println("nodeExecutorCreateMatchResult run start")

	jsonStr, _ := json.MarshalIndent(node, "", "    ")
	log.Println(string(jsonStr))

	req:=node.Input
	flowResult:=&flowReqRsp{
		FlowID:req.FlowID,
		FlowInstanceID:req.FlowInstanceID,
		Stage:req.Stage,
		DebugID:req.DebugID,
		UserRoles:req.UserRoles,
		UserID:req.UserID,
		AppDB:req.AppDB,
		FlowConf:req.FlowConf,
		ModelID:req.ModelID,
		ViewID:req.ViewID,
		FilterData:req.FilterData,
		Filter:req.Filter,
		List:req.List,
		Total:req.Total,
		SelectedRowKeys:req.SelectedRowKeys,
		Pagination:req.Pagination,
		Operation:req.Operation,
		GoOn:true,
	}

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_CREATE_MATCH_RESULT,
	}
	//加载节点配置
	createMatchResultConf:=nodeExecutor.getNodeConfig()
	if createMatchResultConf==nil {
		log.Printf("nodeExecutorCreateMatchResult run get node config error\n")
		return flowResult,common.CreateError(common.ResultNodeConfigError,params)
	}

	if req.Data==nil || len(*req.Data)==0 {
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=req
		log.Printf("nodeExecutorCreateMatchResult run end without data\n")
		return req,nil
	}

	batchNo:=nodeExecutor.getBatchNumber()

	resultData:=make([]flowDataItem,len(*req.Data))
	//这里的分组操作是在数据已经分组的基础上再次分组，分组数据不能跨原来的分组
	//遍历每个数据分组
	for index,item:= range (*req.Data) {
		nodeExecutor.copyDataItem(&item,&(resultData[index]))
		err:=nodeExecutor.createMatchResult(&(resultData[index]),createMatchResultConf,instance,index,batchNo)
		if err !=nil {
			err.Params["nodeID"]=node.ID
			err.Params["nodeType"]=NODE_CREATE_MATCH_RESULT
			return node.Input,err
		}
	}
	
	log.Println("resultData:",resultData)
	flowResult.Data=&resultData
	
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	log.Println("nodeExecutorCreateMatchResult run end")
	return flowResult,nil
}