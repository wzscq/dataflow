package flow

import (
	"time"
	"sync"
)

var g_instance_id int64
var g_instance_id_mutex sync.Mutex

func GetInstanceID()(string){
	g_instance_id_mutex.Lock()
	nowNumber:=time.Now().Unix()
	if nowNumber>g_instance_id {
		g_instance_id=nowNumber
	} else {
		g_instance_id+=1
	}
	t:=time.Unix(g_instance_id,0)
	g_instance_id_mutex.Unlock()
	return t.Format("20060102150405")
}