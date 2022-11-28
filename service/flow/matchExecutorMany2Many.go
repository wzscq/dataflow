package flow


type matchExecutorMany2Many struct {}

func (matchExecutor *matchExecutorMany2Many)sortValues(matchValues,sortedMatchValues,valueIndex *[]int){
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

func (matchExecutor *matchExecutorMany2Many)getMatchGroupByResult(
	forMatch *matchGroup,
	matchResult,valueIndex *[]int)(*matchGroup){
		
	leftIndex:=make([]int,len(*matchResult))
	rightIndex:=make([]int,len(*matchResult))
	
	leftRowCount:=len(forMatch.LeftRows)
	leftCount:=0
	rightCount:=0
	orgIdx:=0
	for _,idx:=range (*matchResult) {
		//获取代比对数组数据索引,
		//valueIndex[idx]==0表示代比对的第一个数据，这里就是左表的数据，直接忽略
		orgIdx=(*valueIndex)[idx]
		if orgIdx < leftRowCount {
			leftIndex[leftCount]=(forMatch.LeftRows)[orgIdx]
			leftCount++
		} else {
			orgIdx=orgIdx-leftRowCount  //索引值减去左表数量后为右表数据的实际索引
			rightIndex[rightCount]=(forMatch.RightRows)[orgIdx]
			rightCount++
		}
	}

	leftIndex=leftIndex[:leftCount]
	rightIndex=rightIndex[:rightCount]
	
	return &matchGroup{
		LeftRows:leftIndex,
		RightRows:rightIndex,
	}
}

//一对多的逻辑，每次从左表拿一个数字，从右表拿所有数字参与比对
func (matchExecutor *matchExecutorMany2Many)getMatchGroup(
	matchValue *matchValue,
	forMatch *matchGroup,
	tolerance matchTolerance)(*matchGroup){
	
	//log.Println("start matchExecutorMany2Many getMatchGroup")

	leftCount:=len(forMatch.LeftRows)
	rightCount:=len(forMatch.RightRows)
	valueIndex:=make([]int,leftCount+rightCount)  //保存排序前后的索引对照
	sortedMatchValues:=make([]int,leftCount+rightCount)
	matchValues:=make([]int,leftCount+rightCount)	

	//获取左表所有待匹配数值,左表数据取反
	for index,row:=range forMatch.LeftRows {
		matchValues[index]=matchValue.LeftValues[row]*-1
	}
	//获取右表所有待匹配数值
	for index,row:=range forMatch.RightRows {
		matchValues[index+leftCount]=matchValue.RightValues[row]
	}
	//log.Println("getMatchGroup matchValues ",matchValues)

	//按照数值大小对数值进行排序
	matchExecutor.sortValues(&matchValues,&sortedMatchValues,&valueIndex)
	//对排序后的数据执行匹配逻辑
	result:=combinationGoal(sortedMatchValues,0,tolerance)
	if result!=nil&&result.MatchCount==1 {
		//根据比对结果生成匹配分组
		return matchExecutor.getMatchGroupByResult(forMatch,&result.Output,&valueIndex)
	}

	//log.Println("end matchExecutorMany2Many getMatchGroup")
	return nil
}