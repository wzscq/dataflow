package flow

import (
	"github.com/go-redis/redis/v8"
	"time"
	"log"
	"encoding/json"
)

type FlowInstanceRepository interface {
	Init(url string,db int,expire time.Duration,password string)
	saveInstance(instance *flowInstance)(error)
	loadInstance(instance *flowInstance)(error)
}

type flowInstanceCahe struct {
	//CompletedNodes []*instanceNode `json:"completedNodes"`
	DebugID *string `json:"debugID,omitempty"`
	WaitingNodes []*instanceNode `json:"waitingNodes"`
}

type DefaultFlowInstanceRepository struct {
	client *redis.Client
	expire time.Duration
}

func (repo *DefaultFlowInstanceRepository)Init(url string,db int,expire time.Duration,password string){
	repo.client=redis.NewClient(&redis.Options{
        Addr:     url,
        Password: password, // no password set
        DB:       db,  // use default DB
    })
	repo.expire=expire
}

func (repo *DefaultFlowInstanceRepository)saveInstance(instance *flowInstance)(error){
	log.Println("saveInstance")
	// Create JSON from the instance data.
	instanceCahe:=&flowInstanceCahe{
		WaitingNodes:instance.WaitingNodes,
		DebugID:instance.DebugID,
	}
  bytes, err := json.Marshal(instanceCahe)
	if err!=nil {
		log.Println("save flow instance error:",err.Error())
		return err
	}
  // Convert bytes to string.
  jsonStr := string(bytes)
	return repo.client.Set(repo.client.Context(), instance.InstanceID, jsonStr, repo.expire).Err()
}

func (repo *DefaultFlowInstanceRepository)loadInstance(instance *flowInstance)(error){
	jsonStr,err:=repo.client.Get(repo.client.Context(), instance.InstanceID).Result()
	if err!=nil {
		log.Println("get flow instance error:",err.Error())
		return err
	}
	// Get byte slice from string.
  bytes := []byte(jsonStr)
	instanceCahe:=&flowInstanceCahe{}
	err = json.Unmarshal(bytes, instanceCahe)
	if err!=nil {
		log.Println("get flow instance error:",err.Error())
		return err
	}
	//instance.CompletedNodes=instanceCahe.CompletedNodes
	instance.DebugID=instanceCahe.DebugID
	instance.WaitingNodes=instanceCahe.WaitingNodes
	return nil
}

