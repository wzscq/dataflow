package flow

import (
	"log"
)

type matchExecutorMany2One struct {}

func (matchExecutor *matchExecutorMany2One)sortValues(matchValues,sortedMatchValues,valueIndex *[]int){
	for index,val:=range *matchValues {
		if index==0 {
			(*sortedMatchValues)[0]=val
			(*valueIndex)[0]=0
		} else {
			for i:=index;i>0;i-- {
				if val>=(*sortedMatchValues)[i-1] {
					(*sortedMatchValues)[i]=val
						(*valueIndex)[i]=index
					break
				} else {
					(*sortedMatchValues)[i]=(*sortedMatchValues)[i-1]
						(*valueIndex)[i]=(*valueIndex)[i-1]
					if i==1 {
						(*sortedMatchValues)[0]=val
						(*valueIndex)[0]=index
					}
				}
			}
		}
	}
}

func (matchExecutor *matchExecutorMany2One)getMatchGroupByResult(
	rightRow int,
	leftRows,matchResult,valueIndex *[]int)(*matchGroup){

	leftIndex:=make([]int,len(*matchResult))
	count:=0
	for _,idx:=range (*matchResult) {
		//获取代比对数组数据索引
		leftIdx:=(*valueIndex)[idx]  
		leftIndex[count]=(*leftRows)[leftIdx]
		count++
	}

	return &matchGroup{
		LeftRows:leftIndex,
		RightRows:[]int{
			rightRow,
		},
	}
}

//一对多的逻辑，每次从左表拿一个数字，从右表拿所有数字参与比对
func (matchExecutor *matchExecutorMany2One)getMatchGroup(
	matchValue *matchValue,
	forMatch *matchGroup,
	tolerance matchTolerance)(*matchGroup){
	
	log.Println("start matchExecutorMany2One getMatchGroup")

	valueIndex:=make([]int,len(forMatch.LeftRows))  //保存排序前后的索引对照
	sortedMatchValues:=make([]int,len(forMatch.LeftRows))
	matchValues:=make([]int,len(forMatch.LeftRows))	

	//获取左表所有待匹配数值
	for index,row:=range forMatch.LeftRows {
		matchValues[index]=matchValue.LeftValues[row]
	}
	log.Println("getMatchGroup matchValues ",matchValues)
	//按照数值大小对数值进行排序
	matchExecutor.sortValues(&matchValues,&sortedMatchValues,&valueIndex)
	//每次从右表取一个数字
	for _,row:=range forMatch.RightRows {
		goal:=matchValue.RightValues[row]
		//对排序后的数据执行匹配逻辑
		result:=combinationGoal(sortedMatchValues,goal,tolerance)
		if result!=nil&&result.MatchCount==1 {
			//根据比对结果生成匹配分组
			return matchExecutor.getMatchGroupByResult(row,&forMatch.LeftRows,&result.Output,&valueIndex)
		}
	}
	
	log.Println("end matchExecutorMany2One getMatchGroup")
	return nil
}