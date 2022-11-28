package flow

type matchExecutorOne2One struct {

}

func (matchExecutor *matchExecutorOne2One)getMatchGroup(
	matchValue *matchValue,
	forMatch *matchGroup,
	tolerance matchTolerance)(*matchGroup){
	
	for _,leftIndex:=range (forMatch.LeftRows) {
		leftVal:=matchValue.LeftValues[leftIndex]
		for _,rightIndex:=range (forMatch.RightRows) {
			rightVal:=matchValue.RightValues[rightIndex]	
			diff:=leftVal - rightVal
			if diff >= tolerance.Left && diff <= tolerance.Right {
				//找到匹配的数据
				return &matchGroup{
					LeftRows:[]int{
						leftIndex,
					},
					RightRows:[]int{
						rightIndex,
					},
				}
			}
		}
	}

	return nil
}