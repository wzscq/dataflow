package flow

import (
	"log"
)

type matchExecutorOne2Many struct {}

func (matchExecutor *matchExecutorOne2Many)sortValues(matchValues,sortedMatchValues,valueIndex *[]int){
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

func (matchExecutor *matchExecutorOne2Many)getMatchGroupByResult(
	leftRow int,
	rightRows,matchResult,valueIndex *[]int)(*matchGroup){
		
	rightIndex:=make([]int,len(*matchResult))
	count:=0
	for _,idx:=range (*matchResult) {
			rightIdx:=(*valueIndex)[idx]  //第一是左表数据，索引值减一后为右表数据的实际索引
			rightIndex[count]=(*rightRows)[rightIdx]
			count++
	}

	return &matchGroup{
		LeftRows:[]int{
			leftRow,
		},
		RightRows:rightIndex,
	}
}

//一对多的逻辑，每次从左表拿一个数字，从右表拿所有数字参与比对
func (matchExecutor *matchExecutorOne2Many)getMatchGroup(
	matchValue *matchValue,
	forMatch *matchGroup,
	tolerance matchTolerance)(*matchGroup){
	
	log.Println("start matchExecutorOne2Many getMatchGroup")

	valueIndex:=make([]int,len(forMatch.RightRows))  //保存排序前后的索引对照
	sortedMatchValues:=make([]int,len(forMatch.RightRows))
	matchValues:=make([]int,len(forMatch.RightRows))	

	//获取右表所有待匹配数值
	for index,row:=range forMatch.RightRows {
		matchValues[index]=matchValue.RightValues[row]
	}
	//按照数值大小对数值进行排序
	matchExecutor.sortValues(&matchValues,&sortedMatchValues,&valueIndex)
	//log.Println("getMatchGroup matchValues ",matchValues)
	//每次从左表取一个数字
	for _,row:=range forMatch.LeftRows {
		goal:=matchValue.LeftValues[row]
		result:=combinationGoal(sortedMatchValues,goal,tolerance)

		if result!=nil&&result.MatchCount==1 {
			//根据比对结果生成匹配分组
			return matchExecutor.getMatchGroupByResult(row,&forMatch.RightRows,&result.Output,&valueIndex)
		}
	}
	
	log.Println("end matchExecutorOne2Many getMatchGroup")
	return nil
}