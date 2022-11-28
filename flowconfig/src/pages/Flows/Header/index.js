import {useSelector,useDispatch} from 'react-redux';
import {Button,Space,Popconfirm} from 'antd';
import mqtt from 'mqtt';
import { 
    SwapOutlined,
    FolderOpenOutlined,
    SaveOutlined,
    CaretRightOutlined,
    ClearOutlined } from '@ant-design/icons';
import './index.css';
import { setView } from '../../../redux/flowSlice';
import {openDialog} from '../../../redux/dialogSlice';
import {addDebugInfo,clearDebugInfo} from '../../../redux/debugSlice';
import { saveFlowConfigAction,deleteFlowAction, startFlowAction,getMqttServer } from '../../../api';
import { useEffect } from 'react';

var g_MQTTClient=null;

//const server="ws://121.36.192.249:9101";
/*const options={
  username:"mosquitto",
  password:"123456",
};*/
const topic="flowdebug/";

export default function Header({appID}){
    const dispatch=useDispatch();
    const view=useSelector(state=>state.flow.view);
    const {pending:isRunning,mqttConf,mqttConfLoaded}=useSelector(state=>state.debug);
    const currentFlow=useSelector(state=>{
        console.log(state.flow);
        const flows=state.flow.openedflows;
        for(let flowIndex=0;flowIndex<flows.length;flowIndex++){
            if(flows[flowIndex].id===state.flow.currentFlow){
                return flows[flowIndex];
            }
        }
    }); 

    const setViewType=()=>{
        dispatch(setView((view==='flow'?'json':'flow')));
    };

    const openFlow=()=>{
        dispatch(openDialog({type:'openFlow',title:'Open Flow',param:{appID:appID}}));
    }

    const saveFlow=()=>{
        const {id,nodes,edges,isNew,description}=currentFlow;
        if(isNew===true){
            dispatch(openDialog({type:'addFlow',title:'Add New Flow',param:{appID:appID}}));
        } else {
            dispatch(saveFlowConfigAction({appID:appID,flowID:id,flowConf:{nodes,edges,description}}));
        }
    }

    const deleteFlow=()=>{
        const {id}=currentFlow;
        dispatch(deleteFlowAction({appID:appID,flowID:id}));
    }

    const startFlow=(debugID)=>{
        const {id,nodes,edges}=currentFlow;
        dispatch(startFlowAction({appID:appID,flowID:id,flowConf:{nodes,edges},userID:'admin',debugID:debugID}));
    }

    const UUIDGeneratorBrowser = () =>([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
        (c ^ (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (c / 4)))).toString(16)
    );

    const getDebugID=()=>{
        return UUIDGeneratorBrowser();
    }

    const connectMqtt=()=>{
        console.log("connectMqtt ... ");
        if(g_MQTTClient!==null){
            g_MQTTClient.end();
            g_MQTTClient=null;
        }

        const debugID=getDebugID();
        const server='ws://'+mqttConf.broker+':'+mqttConf.wsPort;
        const options={
            username:mqttConf.user,
            password:mqttConf.password,
        }
        console.log("connect to mqtt server ... "+server+" with options:",options);
        g_MQTTClient  = mqtt.connect(server,options);
        g_MQTTClient.on('connect', () => {
            console.log("connected to mqtt server "+server+".");
            console.log("subscribe topics ...");
            g_MQTTClient.subscribe(topic+debugID, (err) => {
                if(!err){
                    console.log("subscribe topics success.");
                    console.log("topic:",topic+debugID);
                    //发送流执行请求
                    startFlow(debugID);
                } else {
                    console.log("subscribe topics error :"+err.toString());
                }
            });
        });
        g_MQTTClient.on('message', (topic, payload, packet) => {
            console.log("receive message topic :"+topic+" content :"+payload.toString());
            dispatch(addDebugInfo(JSON.parse(payload.toString())));
        });
        g_MQTTClient.on('close', () => {
            console.log("mqtt client is closed.");
        });
    }

    const runFlow=()=>{
        console.log('runFlow',isRunning);
        if(isRunning===false){
            connectMqtt();
        }
    }

    const clearDebugInfoFunc=()=>{
        dispatch(clearDebugInfo());
    }

    useEffect(()=>{
        if(mqttConfLoaded===false){
            dispatch(getMqttServer());
        }
    },[mqttConfLoaded,dispatch]);

    useEffect(
        ()=>{
            if(isRunning===false&&g_MQTTClient!==null){
                setTimeout(() => {
                    if(g_MQTTClient!==null){
                        g_MQTTClient.end();
                        g_MQTTClient=null;
                    }
                }, 5000);
            }
        },
        [isRunning]
    );

    const viewLable=(view==='flow'?'Json':'Flow');

    return (
        <div className="header-bar">
            <Space >
                <Button size='small' icon={<FolderOpenOutlined />} onClick={openFlow}>Open</Button>
                <Button disabled={(!currentFlow||currentFlow.isModified===false)} size='small' icon={<SaveOutlined />} onClick={saveFlow}>Save</Button>
                <Popconfirm title={"Are you sure to delete the flow with id:"+currentFlow?.id} onConfirm={deleteFlow}>
                    <Button disabled={(!currentFlow||currentFlow.isNew===true)} size='small' icon={<SaveOutlined />} >Delete</Button>
                </Popconfirm>
                <Button disabled={(!currentFlow||isRunning===true)} size='small' icon={<CaretRightOutlined />} onClick={runFlow}>Run</Button>
                <Button disabled={(!currentFlow||isRunning===true)} size='small' icon={<ClearOutlined />} onClick={clearDebugInfoFunc}>Clear Debug Info</Button>
                <Button size='small' icon={<SwapOutlined />} onClick={setViewType}>{viewLable}</Button>
            </Space>
        </div>
    );
}