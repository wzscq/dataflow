package flow

import (
    "time"
	"log"
)

type instanceNode struct {
	ID string `json:"id"`
	Completed bool `json:"completed"`
	StartTime string `json:"startTime"`
	EndTime *string `json:"endTime,omitempty"`
	UserID string `json:"userID,omitempty"`
	NodeType string `json:"nodeType,omitempty"`
	Input *flowReqRsp `json:"input,omitempty"`
	Output *flowReqRsp `json:"output,omitempty"`
	Priority int64 `json:"priority,omitempty"`
}

func createInstanceNode(id,nodeType string,priority int64,input* flowReqRsp)(*instanceNode){
	log.Printf("createInstanceNode id %s，type: %s \n",id,nodeType)
	return &instanceNode{
		ID:id,
		NodeType:nodeType,
		Completed:false,
		StartTime:time.Now().Format("2006-01-02 15:04:05"),
		Input:input,
		Priority:priority,
	}
}

