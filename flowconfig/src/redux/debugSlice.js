import { createSlice } from '@reduxjs/toolkit';
import { message } from 'antd';

import { startFlowAction,getMqttServer } from '../api';


const initialState = {
    pending:false,
    debugMessages:[],
    mqttConfLoaded:false,
    mqttConf:{
        broker:"121.36.192.249",
        wsPort:9101,
        user:"mosquitto",
        password:"123456"
    }
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
    clearDebugInfo
} = debugSlice.actions

export default debugSlice.reducer