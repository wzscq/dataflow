package mqtt

import (
	"log"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type eventHandler interface {
	StartFlow(reqPayload []byte)
}

type MQTTClient struct {
	Broker string
	Port int 
	User string
	Password string
	ClientID string
	StartFlowTopic string
	Client mqtt.Client
	Handler eventHandler
}

func (mqc *MQTTClient) getClient()(mqtt.Client){
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d",mqc.Broker,mqc.Port))
	opts.SetClientID(mqc.ClientID)
	opts.SetUsername(mqc.User)
	opts.SetPassword(mqc.Password)
	opts.SetAutoReconnect(true)

	opts.SetDefaultPublishHandler(mqc.messagePublishHandler)
	opts.OnConnect = mqc.connectHandler
	opts.OnConnectionLost = mqc.connectLostHandler
	opts.OnReconnecting = mqc.reconnectingHandler

	client:=mqtt.NewClient(opts)
	if token:=client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error)
		return nil
	}
	return client
}

func (mqc *MQTTClient) connectHandler(client mqtt.Client){
	log.Println("MQTTClient connectHandler connect status: ",client.IsConnected())
	if client.IsConnected() {
		mqc.Client=client
		client.Subscribe(mqc.StartFlowTopic,0,mqc.onStartFlow)
	}
}

func (mqc *MQTTClient) connectLostHandler(client mqtt.Client, err error){
	log.Println("MQTTClient connectLostHandler connect status: ",client.IsConnected(),err)
}

func (mqc *MQTTClient) messagePublishHandler(client mqtt.Client, msg mqtt.Message){
	log.Println("MQTTClient messagePublishHandler topic: ",msg.Topic())
}

func (mqc *MQTTClient) reconnectingHandler(Client mqtt.Client,opts *mqtt.ClientOptions){
	log.Println("MQTTClient reconnectingHandler ")
}

func (mqc *MQTTClient)onStartFlow(Client mqtt.Client, msg mqtt.Message){
	log.Println("MQTTClient onStartFlow ",msg.Topic())
	mqc.Handler.StartFlow(msg.Payload())		
}

func (mqc *MQTTClient) Init(){
	mqc.Client=mqc.getClient()
}