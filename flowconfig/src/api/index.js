import { createAsyncThunk } from '@reduxjs/toolkit';
import axios from 'axios';

export const getHost=()=>{
    const rootElement=document.getElementById('root');
    const host=rootElement?.getAttribute("host");
    console.log("host:"+host);
    return host;
}
  
const host=getHost()+process.env.REACT_APP_SERVICE_API_PREFIX; //'/frameservice';

export const flowListAction = createAsyncThunk(
    'flowList',
    async ({appID},_) => {
      const config={
        url:host+'/flow/list',
        method:'post',
        data:{
          appDB:appID
        },
        headers:{
          appDB:appID
        }
      }
      const response =await axios(config);
      return response.data;
    }
);

export const getFlowConfigAction = createAsyncThunk(
  'getFlowConfig',
  async ({appID,flowID},_) => {
    const config={
      url:host+'/flow/getConfig',
      method:'post',
      data:{
        appDB:appID,
        flowID
      },
      headers:{
        appDB:appID
      }
    }
    const response =await axios(config);
    return response.data;
  }
);

export const saveFlowConfigAction = createAsyncThunk(
  'saveFlowConfig',
  async ({appID,flowID,flowConf},_) => {
    const config={
      url:host+'/flow/saveConfig',
      method:'post',
      data:{
        appDB:appID,
        flowID:flowID,
        flowConf:flowConf
      },
      headers:{
        appDB:appID
      }
    }
    const response =await axios(config);
    return response.data;
  }
);

export const addFlowAction = createAsyncThunk(
  'addFlow',
  async ({appID,flowID,flowConf},_) => {
    const config={
      url:host+'/flow/addFlow',
      method:'post',
      data:{
        appDB:appID,
        flowID:flowID,
        flowConf:flowConf
      },
      headers:{
        appDB:appID
      }
    }
    const response =await axios(config);
    return response.data;
  }
);

export const deleteFlowAction = createAsyncThunk(
  'deleteFlow',
  async ({appID,flowID},_) => {
    const config={
      url:host+'/flow/deleteFlow',
      method:'post',
      data:{
        appDB:appID,
        flowID:flowID
      },
      headers:{
        appDB:appID
      }
    }
    const response =await axios(config);
    return response.data;
  }
);

export const startFlowAction = createAsyncThunk(
  'startFlow',
  async ({appID,flowID,userID,flowConf,debugID},_) => {
    const config={
      url:host+'/flow/start',
      method:'post',
      data:{
        appDB:appID,
        flowID:flowID,
        userID:userID,
        flowConf:flowConf,
        debugID
      },
      headers:{
        appDB:appID,
        userID:userID
      },
      responseType:'blob'
    }
    const response =await axios(config);
    if(response.data.type==='application/octet-stream'){
      let fileName=response.headers['content-disposition'];
      if(fileName){
        fileName=fileName.substring("attachment; filename=".length);
        fileName=decodeURI(fileName);
      } else {
        fileName="no_file_name";
      }
      return {data:response.data,download:true,fileName:fileName}
    }

    const p=new Promise((resolve) => {
      var reader = new FileReader()
      reader.onload = e => resolve(JSON.parse(e.target.result))
      reader.readAsText(response.data)
    });

    const jsonData=await p;
  
    return jsonData;
  }
);

export const getMqttServer = createAsyncThunk(
  'getMqttServer',
  async () => {
    const config={
      url:host+'/flow/getMqttServer',
      method:'post'
    }
    const response =await axios(config);
    return response.data;
  }
);