package flow

import (
	"log"
	"time"
)

type matchResult struct {
	Input    []int
	Output   []int
	Cost     float64
	Goal     int
	LeftTolerance int
	RightTolerance int
	MatchCount int
	FindCount int
	CalcCount int
	OutPutLength int
}

func combinationGoal(values []int,goal int,tolerance matchTolerance)(*matchResult){
	output := make([]int, len(values))
	
	result := matchResult{
		Input:values,
		Output:output,
		Cost:0,
		Goal:goal,
		LeftTolerance:tolerance.Left,
		RightTolerance:tolerance.Right,
		MatchCount:0,
		FindCount:1,
		CalcCount:0,
		OutPutLength:0,
	}

	result.doReconciliation()
	return &result
}

func (result *matchResult)doReconciliation() {
	tmpVal :=0
  start :=0
  end :=len(result.Input)
	count :=0
	
  for ;start<len(result.Input)&&tmpVal+result.Input[start]<(result.Goal+result.LeftTolerance);start++ {
		result.Output[start]=start  //result.Input[start]
    tmpVal=tmpVal+result.Input[start]
	}
    
	//lastMatchCount:=0
  startTime:=time.Now()

	for ;count<=start&&result.MatchCount<result.FindCount;{
			//log.Println("begin match, times:",count,"start:",start," end:",end, " tmpVal:",tmpVal,result.Output[:start-count])
			if count==0 {
					result.calc(start,end,tmpVal,start)
			} else {
					result.calc(start-count+1,end,tmpVal,start-count)
			}
			//log.Println("end match matchCount:",result.MatchCount-lastMatchCount,"\n")
			count++
			if start>=count {
					tmpVal=tmpVal-result.Input[start-count]
			}
			//lastMatchCount=result.MatchCount;
	}

	result.Output=result.Output[:result.OutPutLength]
  endTime:=time.Now()
	diff := endTime.Sub(startTime)
	//log.Println("duration:",diff.Seconds())
	result.Cost = diff.Seconds()
}

func (result *matchResult)calc(start int,end int,valPre int,level int){
    for i:=start;i<end&&result.MatchCount<result.FindCount;i++ {
        result.CalcCount=result.CalcCount+1
        tmpVal:=valPre+result.Input[i]
        if tmpVal < result.Goal+result.LeftTolerance {
            result.Output[level]=i //result.Input[i]
            result.calc(i+1,end,tmpVal,level+1)
				} else {
            if tmpVal<=result.Goal+result.RightTolerance {
								result.MatchCount++
                log.Println("start:",start," i:",i," preVal:",valPre," tmpVal:",tmpVal," calcCount:",result.CalcCount," matchCount:",result.MatchCount)
								result.Output[level]=i //result.Input[i]
								result.OutPutLength=level+1
                log.Println(result.Output[:level+1])
            }
            break
        }
    }
}

