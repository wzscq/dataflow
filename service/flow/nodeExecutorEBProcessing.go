package flow

import (
    "time"
	"dataflow/common"
	"dataflow/data"
	"encoding/json"
	"log"
	"sort"
)
/*
余额处理过程涉及逻辑说明
1、参与核销的数据需要首先按照主核销字段从小到大排序
2、一次余额处理可以分为多个步骤，
	每个步骤包含以下配置：
		一个处理类型：同表处理（单表），跨表处理（双表） 目前暂时仅处理双表的情况
		一个或多个（目前仅处理2个的情况，后续补充完善多个的情况）modelID，只有一个modelID的情况就是同表处理（表内正负数核销）
			模型包含以下字段配置：
				方向：左表、右表
				是否调节表：是、否 （主要是决定了核销记录ID对应的字段）
				余额为0时是否继续核销：是、否，一般应该都是否，对于有调节表的情况有时是需要往业务表增加值的
		核销类型：根据结果表的核销类型填写对应取值，系统提供默认选择项
		一方为负数时核销方式：（补充：当两边同为负数或正数时核销绝对值小的数量）
			1、核销负数，正数数记录变大，
			2、核销正数，负数变小
			3、核销左侧，
			4、核销右侧
3、需要指定核销明细表modelID
4、需要指定比对分组表modelID
5、需要设置是否核销金额，如果同时指定了核销金额和数量，则以金额为主，数量按比例核销
6、需要设置是否核销数量
7、核销余额的处理，
	调节表不允许出现余额,如果调节表出现余额则报错，
	左右表出现余额时生成新的记录数据放入对应的待核销数据表中
8、备注一下单表处理策略
	可以先将单表的正负数转化为左右表，然后按照左右表核销，核销后再合并左右表核销结果组成单表
*/

const (
	WRITEOFF_SIDE_LEFT = "left"  //左表
	WRITEOFF_SIDE_RIGHT = "right" //右表
)

const (
	SOURCE_ORIGNAL="0"  //原始数据
	SOURCE_ENDINGBALANCE = "1"   //核销后余额
	SOURCE_CALLBACK = "2" //撤销反冲
)

const (
	CC_LEFT_ID = "left_id"  //左表ID字段
	CC_RIGHT_ID = "right_id"  //右表ID字段
	CC_LEFTAJUST_ID = "left_ajust_id"  //左表调节表ID字段
	CC_RIGHTAJUST_ID = "right_ajust_id"  //左表调节表ID字段
	CC_AMOUNT = "amount"  //核销金额字段
	CC_QUANTITY = "quantity"  //核销数量字段
	CC_GROUP_ID = "match_group"  //比对分组号字段
	CC_WRITEOFF_TYPE = "writeoff_type"  //核销类型字段
	CC_WRITEOFF_STATUS = "writeoff_status"  //核销状态
	CC_WRITEOFF_AMOUNT = "writeoff_amount"  //已核销金额
	CC_WRITEOFF_QUANTITY = "writeoff_quantity"  //已核销数量
	CC_AMOUNT_BALANCE = "amount_balance"  //核销后剩余金额
	CC_QUANTITY_BALANCE = "quantity_balance"  //核销后剩余数量
	CC_SOURCE = "source" //数据来源
	CC_SOURCE_ID = "source_id" //数据来源
)

const (
	WRITEOFF_STATUS_NOTYET = "0"  //未核销
	WRITEOFF_STATUS_ALREADY = "1"  //已核销
)

const (
	WRITEOFF_AMOUNT="1"    //核销金额
	WRITEOFF_QUANTITY="1"  //核销数量
	WRITEOFF_ZERO_STOP="1" //余额为0时停止核销
)

//只有一方为负数时的核销方式
const (
	MONO_NEGATIVE_WRITEOFF_NEGATIVE="0"  //消负数
	MONO_NEGATIVE_WRITEOFF_POSITIVE="1"  //消正数
	MONO_NEGATIVE_WRITEOFF_LEFT="2"     //消左表数量
	MONO_NEGATIVE_WRITEOFF_RIGHT="3"     //消右表数量
)

type bpModel struct  {
	ModelID string `json:"modelID"`
	Side string `json:"side"`
	Ajust string `json:"ajust"`
	ZeroStop string `json:"zeroStop"`
}

type bpStep struct {
	Models []bpModel `json:"models"`
	WriteoffType string `json:"writeoffType"`
	MonoNegativeWriteoffMethod string `json:"monoNegativeWriteoffMethod"`
}

type bpNodeConf struct {
	BPSteps []bpStep `json:"bpSteps"`
	WriteoffDetailModelID string  `json:"writeoffDetailModelID"`
	MatchGroupModelID string `json:"matchGroupModelID"`
	DealAmount string  `json:"dealAmount"`
	DealQuantity string `json:"dealQuantity"`
}

type nodeExecutorEBProcessing struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorEBProcessing)getNodeConf()(*bpNodeConf){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
    if err != nil {
        log.Println(err)
		return nil
    }
	log.Println(string(jsonStr))
	conf:=&bpNodeConf{}
    if err := json.Unmarshal(jsonStr, conf); err != nil {
        log.Println(err)
		return nil
    }

	return conf
}

func (nodeExecutor *nodeExecutorEBProcessing)createGroupDataItem(
	groupModelID string,
	groupID interface{})(*modelDataItem){
	return &modelDataItem{
		ModelID:&groupModelID,
		List:&[]map[string]interface{}{
			map[string]interface{}{
				data.CC_ID:groupID,
				data.SAVE_TYPE_COLUMN:data.SAVE_UPDATE,
				CC_WRITEOFF_STATUS:WRITEOFF_STATUS_ALREADY,
			},
		},
	}
}

func (nodeExecutor *nodeExecutorEBProcessing)getBPResult(
	nodeConf *bpNodeConf,
	bpModelMap map[string]*BPModelIterator,
	bpResult *modelDataItem)(*flowDataItem){
	log.Println("nodeExecutorEBProcessing getBPResult start")

	//如果没有任何核销记录，则不对数据做处理
	if bpResult.List==nil || len(*bpResult.List)==0 {
		return nil
	}

	//整合处理结果
	models:=[]modelDataItem{
		*bpResult,
	}
	var groupDataItem *modelDataItem
	for _,it:=range(bpModelMap){
		if groupDataItem==nil {
			groupDataItem=nodeExecutor.createGroupDataItem(nodeConf.MatchGroupModelID,it.getGroupID())
			models=append(models,*groupDataItem)
		}

		modelData:=it.getWriteoffResult()
		if modelData!=nil {
			models=append(models,*modelData)
		}
	}

	fDataItem:=flowDataItem{
		Models:models,
	}

	log.Println("nodeExecutorEBProcessing getBPResult end")
	return &fDataItem
}

func (nodeExecutor *nodeExecutorEBProcessing)getModelDataItem(
	modelID string,
	dataItem *flowDataItem)(*modelDataItem){
	
	for _,modelData:=range(dataItem.Models){
		if (*modelData.ModelID)==modelID {
			if modelData.List==nil || len(*modelData.List)==0 {
				return nil
			} 
			return &modelData
		}
	}

	return nil
}

func (nodeExecutor *nodeExecutorEBProcessing)getGroupID(dataItem *modelDataItem)(interface{},int){
	for _,row:=range(*dataItem.List){
		groupID,ok:=row[CC_GROUP_ID]
		if !ok {
			return nil,common.ResultNoGroupID
		}
		return groupID,common.ResultSuccess
	}
	return nil,common.ResultNoGroupID
}

func (nodeExecutor *nodeExecutorEBProcessing)sortItemRows(
	nodeConf *bpNodeConf,
	dataItem *modelDataItem){
	sortField:=CC_AMOUNT
	if nodeConf.DealAmount!=WRITEOFF_AMOUNT && 
	   nodeConf.DealQuantity==WRITEOFF_QUANTITY {
		sortField=CC_QUANTITY
	}

	compareAmount:=func(i,j int) bool {
		val1,_:=getFieldValue((*dataItem.List)[i],sortField)
		val2,_:=getFieldValue((*dataItem.List)[j],sortField)
		return val1<val2
	}
	sort.Slice(*dataItem.List, compareAmount)
}

func (nodeExecutor *nodeExecutorEBProcessing)loadModelIterator(
	nodeConf *bpNodeConf,
	step *bpStep,
	dataItem *flowDataItem,
	bpModelMap map[string]*BPModelIterator)(int){
	log.Println("nodeExecutorEBProcessing loadModelIterator start")
	for _,modelItem:=range(step.Models){
		_,ok:=bpModelMap[modelItem.ModelID]
		if !ok {
			//如果未获取到对应模型的数据则说明这个表没有数据需要核销，因此直接将结果表赋值未空
			//后续核销步骤中遇到有空的数据则直接跳过不处理
			modelDataItem:=nodeExecutor.getModelDataItem(modelItem.ModelID,dataItem)
			if modelDataItem!=nil {
				//对参与核销的数据做个排序，按从小到大顺序排序
				nodeExecutor.sortItemRows(nodeConf,modelDataItem)

				groupID,errorCode:=nodeExecutor.getGroupID(modelDataItem)
				if errorCode!=common.ResultSuccess {
					return errorCode
				}
				bpModelIterator:=createModelIterator(
					nodeConf.DealAmount,
					nodeConf.DealQuantity,
					&modelItem,
					modelDataItem,
					groupID)
				bpModelMap[modelItem.ModelID]=bpModelIterator
			}
		}
	}
	log.Println("nodeExecutorEBProcessing loadModelIterator end")
	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorEBProcessing)getWriteoffAmountByVal(
	step *bpStep,
	leftVal,rightVal float64)(float64,int){
	if leftVal==0 && rightVal==0 {
		//左右表的值同时等于0，结束核销
		return 0,common.ResultSuccess
	} else if leftVal<0 && rightVal <0 {
		//左右表的值同时小于等于0，取绝对值小的，也就是较大的值
		return ternaryOperatorFloat64(leftVal>rightVal,leftVal,rightVal),common.ResultSuccess
	} else if leftVal>0 && rightVal>0 {
		//左右表的值同时大于0，取绝对值小的，也就是较小的值
		return ternaryOperatorFloat64(leftVal<rightVal,leftVal,rightVal),common.ResultSuccess
	} else if leftVal==0 || rightVal == 0 {
		//当有一边为0值，则返回不为0的值
		return ternaryOperatorFloat64(leftVal!=0,leftVal,rightVal),common.ResultSuccess
	} else if leftVal<0 || rightVal < 0 {
		//当有一边值小于0，另一边大于0时，根据配置项返回相应的值
		switch step.MonoNegativeWriteoffMethod {
			case MONO_NEGATIVE_WRITEOFF_NEGATIVE:   //消负数
				return ternaryOperatorFloat64(leftVal<0,leftVal,rightVal),common.ResultSuccess
			case MONO_NEGATIVE_WRITEOFF_POSITIVE:   //消正数
				return ternaryOperatorFloat64(leftVal>0,leftVal,rightVal),common.ResultSuccess
			case MONO_NEGATIVE_WRITEOFF_LEFT:     //消左表数量
				return leftVal,common.ResultSuccess
			case MONO_NEGATIVE_WRITEOFF_RIGHT:     //消右表数量
				return rightVal,common.ResultSuccess
		}
	}

	return 0,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorEBProcessing)getWriteoffAmount(
	step *bpStep,
	bpModelMap map[string]*BPModelIterator)(float64,int){
	var leftVal,rightVal float64
	//注意这里的取数逻辑目前仅支持一个左表和一个右表的情况，其它情况暂时不支持，待后续优化完善
	for _,modelItem:=range(step.Models){
		it,ok:=bpModelMap[modelItem.ModelID]
		//取不到，说明对应的模型数据不存在，这时直接返回0，结束当前核销步骤
		if !ok {
			return 0,common.ResultSuccess
		}

		val,errorCode:=it.getBalance()
		if errorCode!=common.ResultSuccess {
			return 0,errorCode
		}

		//如果待核销数据为0，且配置中要求当前表余额等于0时停止核销，则直接返回0
		if val==0 && modelItem.ZeroStop == WRITEOFF_ZERO_STOP {
			return 0,common.ResultSuccess
		}

		if modelItem.Side == WRITEOFF_SIDE_LEFT {
			leftVal=val
		} else {
			rightVal=val
		}
	}

	return nodeExecutor.getWriteoffAmountByVal(step,leftVal,rightVal)
}

func (nodeExecutor *nodeExecutorEBProcessing)createWriteoffRecord(
	writeoffType string)(map[string]interface{}){
	return map[string]interface{}{
		data.SAVE_TYPE_COLUMN:data.SAVE_CREATE,
		CC_WRITEOFF_TYPE:writeoffType,
	}
}

func (nodeExecutor *nodeExecutorEBProcessing)writeoff(
	step *bpStep,
	bpModelMap map[string]*BPModelIterator,
	bpResult *modelDataItem,
	writeoffAmount float64){

	//创建核销中间结果表记录
	resultRow:=nodeExecutor.createWriteoffRecord(step.WriteoffType)
	//每个模型的余额扣除
	for _,modelItem:=range(step.Models){
		it,_:=bpModelMap[modelItem.ModelID]
		amount,quantity,id:=it.writeoff(writeoffAmount)
		if modelItem.Side == WRITEOFF_SIDE_LEFT {
			//组号，核销金额数量都以左侧表为准
			resultRow[CC_LEFT_ID]=id
			resultRow[CC_AMOUNT]=amount
			resultRow[CC_QUANTITY]=quantity
			resultRow[CC_GROUP_ID]=it.getGroupID()
		} else {
			resultRow[CC_RIGHT_ID]=id
		}
	}

	(*bpResult.List)=append((*bpResult.List),resultRow)
}

func (nodeExecutor *nodeExecutorEBProcessing)balanceProcessingStep(
	nodeConf *bpNodeConf,
	step *bpStep,
	dataItem *flowDataItem,
	bpModelMap map[string]*BPModelIterator,
	bpResult *modelDataItem)(int){
	log.Println("nodeExecutorEBProcessing balanceProcessingStep start")

	//将处理相关的Model数据加载到迭代器中
	errorCode:=nodeExecutor.loadModelIterator(nodeConf,step,dataItem,bpModelMap)
	if errorCode!=common.ResultSuccess {
		return errorCode
	}
	log.Println(bpModelMap)

	//开始循环迭代处理数据核销
	writeoffAmount,errorCode:=nodeExecutor.getWriteoffAmount(step,bpModelMap)
	log.Printf("nodeExecutorEBProcessing getWriteoffAmount writeoffAmount:%f,errorCode:%d",writeoffAmount,errorCode)
	if errorCode!=common.ResultSuccess {
		return errorCode
	}

	for writeoffAmount!=0 {
		//核销对应数据
		nodeExecutor.writeoff(step,bpModelMap,bpResult,writeoffAmount)

		writeoffAmount,errorCode=nodeExecutor.getWriteoffAmount(step,bpModelMap)
		log.Printf("nodeExecutorEBProcessing getWriteoffAmount writeoffAmount:%f,errorCode:%d",writeoffAmount,errorCode)
		if errorCode!=common.ResultSuccess {
			return errorCode
		}
	}

	log.Println("nodeExecutorEBProcessing balanceProcessingStep end")
	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorEBProcessing)balanceProcessing(
	nodeConf *bpNodeConf,
	dataItem *flowDataItem)(*flowDataItem,int){
	log.Println("nodeExecutorEBProcessing balanceProcessing start")
	//模型数据迭代器用于封装对模型数据的获取和核销逻辑
	//首先创建模型数据迭代器的容器，将所有模型的迭代器放入map中
	bpModelMap:=map[string]*BPModelIterator{}
	//存储核销中间结果
	bpResult:=&modelDataItem{
		ModelID:&nodeConf.WriteoffDetailModelID,
		List:&[]map[string]interface{}{},
	}
	//核销过程可以分多个步骤进行，这里每次处理一个步骤
	for _,step:=range(nodeConf.BPSteps){
		errorCode:=nodeExecutor.balanceProcessingStep(nodeConf,&step,dataItem,bpModelMap,bpResult)
		if errorCode!=common.ResultSuccess {
			return nil,errorCode
		}
	}

	//将模型迭代器数据转化未结果数据
	resultDataItem:=nodeExecutor.getBPResult(nodeConf,bpModelMap,bpResult)
	if resultDataItem!=nil {
		log.Println(*resultDataItem)
	} else {
		log.Println("getBPResult return with nil")
	}
	log.Println("nodeExecutorEBProcessing balanceProcessing end")
	return resultDataItem,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorEBProcessing)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	log.Println("nodeExecutorEBProcessing run start")

	req:=node.Input
	flowResult:=&flowReqRsp{
		FlowID:req.FlowID,
		FlowInstanceID:req.FlowInstanceID,
		Stage:req.Stage,
		DebugID:req.DebugID,
		UserRoles:req.UserRoles,
		GlobalFilterData:req.GlobalFilterData,
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
		SelectAll:req.SelectAll,
		GoOn:true,
	}

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_EB_PROCESSING,
	}
	//加载节点配置
	nodeConf:=nodeExecutor.getNodeConf()
	if nodeConf==nil {
		log.Printf("nodeExecutorEBProcessing run get node config error\n")
		return flowResult,common.CreateError(common.ResultNodeGroupConfigError,params)
	}

	if req.Data==nil || len(*req.Data)==0 {
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=req
		return req,nil
	}

	//装载结果数据
	resultData:=[]flowDataItem{}
	//对于每个分组数据进行单独核销处理，分组间没有关联
	//遍历每个数据分组
	for _,item:= range (*req.Data) {
		result,errorCode:=nodeExecutor.balanceProcessing(nodeConf,&item)
		if errorCode != common.ResultSuccess {
			return flowResult,common.CreateError(errorCode,params)
		}
		if result!=nil {
			resultData=append(resultData,*result)
		}
	}

	flowResult.Data=&resultData

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	log.Println("nodeExecutorEBProcessing run end")
	return flowResult,nil
}