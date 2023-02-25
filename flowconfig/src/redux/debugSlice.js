import { createSlice } from '@reduxjs/toolkit';
import { message } from 'antd';

import { startFlowAction,getMqttServer } from '../api';


const initialState = {
    pending:false,
    debugMessages:[],
    mqttConfLoaded:false,
    operation:null,
    mqttConf:{
        broker:"121.36.192.249",
        wsPort:9101,
        user:"mosquitto",
        password:"123456"
    }
}

const downloadFile=({data,fileName})=>{
    let blob=[data];
    var a = document.createElement('a');
    var url = window.URL.createObjectURL(new Blob(blob));
    a.href = url;
    a.download = fileName;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
}

export const debugSlice = createSlice({
    name: 'debug',
    initialState,
    reducers: {
        addDebugInfo:(state,action)=>{
            state.debugMessages.push(action.payload);
        },
        clearDebugInfo:(state,action)=>{
            state.debugMessages=[];
        },
        clearOperation:(state,action)=>{
            state.operation=null;
        }
    },
    extraReducers: (builder) => {
        builder.addCase(startFlowAction.pending, (state, action) => {
            state.pending=true;

        });
        builder.addCase(startFlowAction.fulfilled, (state, action) => {
            console.log("startFlowAction fulfilled:",action);
            state.pending=false;
            console.log(action);
            if(action.payload.error){
                if(action.payload.params){
                    message.error(action.payload.message+"\n"+JSON.stringify(action.payload.params));
                } else {
                    message.error(action.payload.message);
                }
            } else {
                if(action.payload.download===true){
                    downloadFile(action.payload);
                } else if(action.payload.result?.operation){
                    //如果流返回了operation,则需要通过crvframe中的mainframe执行这个operation
                    state.operation=action.payload.result.operation;
                }
            }
        });
        builder.addCase(startFlowAction.rejected , (state, action) => {
            console.log("startFlowAction return error:",action);
            state.pending=false;
            if(action.error&&action.error.message){
                message.error(action.error.message);
            } else {
                message.error("未知错误！");
            }
        });
        //获取MQTT配置信息
        builder.addCase(getMqttServer.pending, (state, action) => {
            state.mqttConfLoaded=true;
        });
        builder.addCase(getMqttServer.fulfilled, (state, action) => {
            console.log("getMqttServer fulfilled:",action);
            console.log(action);
            state.mqttConf=action.payload.result;
            console.log(action.payload.result);
        });
        builder.addCase(getMqttServer.rejected , (state, action) => {
            console.log("getMqttServer return error:",action);
            if(action.error&&action.error.message){
                message.error(action.error.message);
            } else {
                message.error("获取MQTT服务配置失败");
            }
        });
    }
});

export const { 
    addDebugInfo,
    clearDebugInfo,
    clearOperation
} = debugSlice.actions

export default debugSlice.reducer