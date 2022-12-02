package flow

import (
	"dataflow/common"
	"dataflow/data"
	"log"
)

type BPModelIterator struct {
	dealAmount string 
	dealQuantity string 
	modelData *modelDataItem
	groupID interface{}
	balance float64
	amountBalance float64
	quantityBalance float64
	writeoffQuantity float64
	currentRow int
	currentRowID interface{}
}

func createModelIterator(
	dealAmount string,
	dealQuantity string, 
	bpModel *bpModel,
	modelData *modelDataItem,
	groupID interface{})(*BPModelIterator){
	
	return &BPModelIterator{
		dealAmount:dealAmount,
		dealQuantity:dealQuantity,
		modelData:modelData,
		groupID:groupID,
		currentRow:-1,
		amountBalance:0,
		quantityBalance:0,
		balance:0,
		writeoffQuantity:0,
	}
}

func (it *BPModelIterator)getWriteoffBanlance(row map[string]interface{})(int){
	var errorCode int
	if it.dealAmount == WRITEOFF_AMOUNT {
		it.amountBalance,errorCode=getFieldValue(row,CC_AMOUNT)
		if errorCode!=common.ResultSuccess {
			return errorCode
		}
		it.balance=it.amountBalance
	}
	if it.dealQuantity == WRITEOFF_QUANTITY {
		it.quantityBalance,errorCode=getFieldValue(row,CC_QUANTITY)
		if errorCode!=common.ResultSuccess {
			return errorCode
		}
		//如果仅指定了核销数量，则将默认balance设置未数量的balance
		if it.dealAmount != WRITEOFF_AMOUNT {
			it.balance=it.quantityBalance
		}
	}
	//重新获取了balance的时候将核销的数量置为0
	it.writeoffQuantity=0
	return common.ResultSuccess
}

func (it *BPModelIterator)getRowID(row map[string]interface{})(interface{}){
	rowID,_:=row[data.CC_ID]
	return rowID
}

func (it *BPModelIterator)getBalance()(float64,int){
	log.Printf("%s BPModelIterator getBalance start with it.quantityBalance:%f,it.amountBalance:%f,it.balance:%f,it.writeoffQuantity:%f",
	*it.modelData.ModelID,it.quantityBalance,it.amountBalance,it.balance,it.writeoffQuantity)
	//banlance不等于0，说明之前获取的记录的balance还没有核销完，直接返回这个余额
	if it.balance!=0 {
		return it.balance,common.ResultSuccess
	}

	//如果之前的余额已经核销完成，则从原始数据中的取一个新的不等于0的余额进行核销
	var errorCode int
	for (it.currentRow+1)<len(*it.modelData.List)&&it.balance==0 {
		it.currentRow=it.currentRow+1
		log.Printf("BPModelIterator getBalance currentRow %d",it.currentRow)
		row:=(*it.modelData.List)[it.currentRow]
		errorCode=it.getWriteoffBanlance(row)
		if errorCode!=common.ResultSuccess {
			return 0,errorCode
		}
		it.currentRowID=it.getRowID(row)
	}
	log.Printf("%s BPModelIterator getBalance end with it.quantityBalance:%f,it.amountBalance:%f,it.balance:%f,it.writeoffQuantity:%f",
	*it.modelData.ModelID,it.quantityBalance,it.amountBalance,it.balance,it.writeoffQuantity)
	return it.balance,common.ResultSuccess
}

func (it *BPModelIterator)writeoff(amount float64)(float64,float64,interface{}){
	log.Printf("%s BPModelIterator writeoff start with it.quantityBalance:%f,it.amountBalance:%f,it.balance:%f,amount:%f,it.writeoffQuantity:%f",
	*it.modelData.ModelID,it.quantityBalance,it.amountBalance,it.balance,amount,it.writeoffQuantity)
				
	it.balance=it.balance-amount
	var writeoffAmount,writeoffQuantity float64
	if it.dealAmount == WRITEOFF_AMOUNT {
		writeoffAmount=amount
	}

	if it.dealQuantity == WRITEOFF_QUANTITY {
		if it.dealAmount != WRITEOFF_AMOUNT {
			writeoffQuantity=amount
		} else {
			//如果指定了金额核销，则数量的核销将按照金额的比例进行核销
			if it.balance==0 {
				//全部金额消完，则将数量也同时消完，用整体减去已经核销的部分
				writeoffQuantity=it.quantityBalance-it.writeoffQuantity
				it.writeoffQuantity=it.quantityBalance
			} else {
				//按比例核销数量
				writeoffQuantity=amount*it.quantityBalance/it.amountBalance
				//将已经核销的数量记录下来，金额消完的时候用整体减去已经核销的部分
				it.writeoffQuantity=it.writeoffQuantity+writeoffQuantity
			}
		}
	}

	log.Printf("%s BPModelIterator writeoff end with it.quantityBalance:%f,it.amountBalance:%f,it.balance:%f,amount:%f,it.writeoffQuantity:%f",
	*it.modelData.ModelID,it.quantityBalance,it.amountBalance,it.balance,amount,it.writeoffQuantity)

	return writeoffAmount,writeoffQuantity,it.currentRowID
}

func (it *BPModelIterator)getGroupID()(interface{}){
	return it.groupID
}

func (it *BPModelIterator)copyMap(in map[string]interface{})(map[string]interface{}){
	out:=map[string]interface{}{}
	for key,val:=range(in){
		out[key]=val
	}
	return out
}

func (it *BPModelIterator)getUpateWholeWriteoff(orgRow map[string]interface{})(map[string]interface{}){
	updateRow:=map[string]interface{}{}
	updateRow[data.CC_ID]=orgRow[data.CC_ID]
	updateRow[data.CC_VERSION]=orgRow[data.CC_VERSION]
	updateRow[data.SAVE_TYPE_COLUMN]=data.SAVE_UPDATE
	updateRow[CC_WRITEOFF_STATUS]=WRITEOFF_STATUS_ALREADY
	if it.dealAmount == WRITEOFF_AMOUNT {
		updateRow[CC_WRITEOFF_AMOUNT]=orgRow[CC_AMOUNT]
		updateRow[CC_AMOUNT_BALANCE]=0.0
	}
	if it.dealQuantity == WRITEOFF_QUANTITY {
		updateRow[CC_WRITEOFF_QUANTITY]=orgRow[CC_QUANTITY]
		updateRow[CC_QUANTITY_BALANCE]=0.0
	}
	return updateRow
}

func (it *BPModelIterator)getUpatePartialWriteoff(orgRow map[string]interface{})(map[string]interface{}){
	updateRow:=map[string]interface{}{}
	updateRow[data.CC_ID]=orgRow[data.CC_ID]
	updateRow[data.CC_VERSION]=orgRow[data.CC_VERSION]
	updateRow[data.SAVE_TYPE_COLUMN]=data.SAVE_UPDATE
	updateRow[CC_WRITEOFF_STATUS]=WRITEOFF_STATUS_ALREADY
	if it.dealAmount == WRITEOFF_AMOUNT {
		updateRow[CC_WRITEOFF_AMOUNT]=it.amountBalance-it.balance
		updateRow[CC_AMOUNT_BALANCE]=it.balance
	}
	if it.dealQuantity == WRITEOFF_QUANTITY {
		if it.dealAmount == WRITEOFF_AMOUNT {
			updateRow[CC_WRITEOFF_QUANTITY]=it.writeoffQuantity
			updateRow[CC_QUANTITY_BALANCE]=it.quantityBalance - it.writeoffQuantity
		} else {
			updateRow[CC_WRITEOFF_QUANTITY]=it.amountBalance-it.balance
			updateRow[CC_QUANTITY_BALANCE]=it.balance
		}
	}
	
	return updateRow
}

func (it *BPModelIterator)getCreatePartialBalance(orgRow map[string]interface{})(map[string]interface{}){
	orgRow[data.SAVE_TYPE_COLUMN]=data.SAVE_CREATE
	orgRow[CC_WRITEOFF_STATUS]=WRITEOFF_STATUS_NOTYET
	if it.dealAmount == WRITEOFF_AMOUNT {
		orgRow[CC_AMOUNT]=it.balance
	}
	if it.dealQuantity == WRITEOFF_QUANTITY {
		if it.dealAmount == WRITEOFF_AMOUNT {
			orgRow[CC_QUANTITY]=it.quantityBalance - it.writeoffQuantity
		} else {
			orgRow[CC_QUANTITY]=it.balance
		}
	}
	orgRow[CC_SOURCE]=SOURCE_ENDINGBALANCE
	orgRow[CC_SOURCE_ID]=orgRow[data.CC_ID]
	delete(orgRow,data.CC_ID)
	delete(orgRow,data.CC_VERSION)
	return orgRow
}

func (it *BPModelIterator)getUpateNoWriteoff(orgRow map[string]interface{})(map[string]interface{}){
	updateRow:=map[string]interface{}{}
	updateRow[data.CC_ID]=orgRow[data.CC_ID]
	updateRow[data.CC_VERSION]=orgRow[data.CC_VERSION]
	updateRow[data.SAVE_TYPE_COLUMN]=data.SAVE_UPDATE
	updateRow[CC_WRITEOFF_STATUS]=WRITEOFF_STATUS_ALREADY
	if it.dealAmount == WRITEOFF_AMOUNT {
		updateRow[CC_WRITEOFF_AMOUNT]=0.0
		updateRow[CC_AMOUNT_BALANCE]=orgRow[CC_AMOUNT]
	}
	if it.dealQuantity == WRITEOFF_QUANTITY {
		updateRow[CC_WRITEOFF_QUANTITY]=0.0
		updateRow[CC_QUANTITY_BALANCE]=orgRow[CC_QUANTITY]
	}
	return updateRow
}

func (it *BPModelIterator)getUpateWholeBalance(orgRow map[string]interface{})(map[string]interface{}){
	orgRow[data.SAVE_TYPE_COLUMN]=data.SAVE_CREATE
	orgRow[CC_WRITEOFF_STATUS]=WRITEOFF_STATUS_NOTYET
	orgRow[CC_SOURCE]=SOURCE_ENDINGBALANCE
	orgRow[CC_SOURCE_ID]=orgRow[data.CC_ID]
	delete(orgRow,data.CC_ID)
	delete(orgRow,data.CC_VERSION)
	return orgRow
}

func (it *BPModelIterator)getWriteoffResult()(*modelDataItem){
	//因为核销时顺序处理的，所有currentRow之前的数据都是已经核销完成的
	//currentRow如果balance==0，则也是核销完成的，如果balance!=0则说名有余额存在
	//对于有所有未核销和余额，原记录仍然标记未已核销状态，同时生成新对应的待核销记录
	rows:=[]map[string]interface{}{}
	for index,row:=range(*it.modelData.List){
		if index < it.currentRow {
			rows=append(rows,it.getUpateWholeWriteoff(row))
		} else if (index == it.currentRow){
			if it.balance==0{
				rows=append(rows,it.getUpateWholeWriteoff(row))
			} else {
				rows=append(rows,it.getUpatePartialWriteoff(row))
				rows=append(rows,it.getCreatePartialBalance(it.copyMap(row)))
			}
		} else {
			rows=append(rows,it.getUpateNoWriteoff(row))
			rows=append(rows,it.getUpateWholeBalance(it.copyMap(row)))
		}
	}

	modelID:=*it.modelData.ModelID
	modelData:=modelDataItem{
		ModelID:&modelID,
		List:&rows,
	}
	return &modelData
}

