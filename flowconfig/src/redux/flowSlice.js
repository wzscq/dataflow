import { createSlice } from '@reduxjs/toolkit';
import { message } from 'antd';
import {
    flowListAction,
    getFlowConfigAction,
    saveFlowConfigAction,
    addFlowAction,
    deleteFlowAction
} from '../api';

var g_newFlowIndex=0;
// Define the initial state using that type
const initialState = {
    view:'flow',
    currentNode:"",
    currentFlow:null,
    openedflows:[],
    flows:[],
    loaded:false,
    pending:false
}

const getNewFlowID=()=>{
    return 'new flow '+g_newFlowIndex++;
}

const closeFlowFunc=(state,flowID)=>{
    let closedIndex=0;
    for(let i=0;i<state.openedflows.length;i++){
        if(state.openedflows[i].id===flowID){
            closedIndex=i;
        }
    }
    delete state.openedflows[closedIndex];
    state.openedflows=state.openedflows.filter(item=>item);
    if(state.openedflows.length===0){
        state.currentFlow="";
    } else {
        if (closedIndex>=state.openedflows.length){
            closedIndex=state.openedflows.length-1;
        }
        state.currentFlow=state.openedflows[closedIndex].id;
    }
    
    state.currentNode="";
}

export const flowSlice = createSlice({
    name: 'flow',
    initialState,
    reducers: {
        addNewFlow: (state,action) => {
            const flowID=getNewFlowID();
            state.openedflows.push({
                loaded:true,
                isNew:true,
                isModified:true,
                id:flowID,
                nodes:[],
                edges:[],
            });
            state.currentFlow=flowID
            state.currentNode="";
        },
        setNodes:(state,action) => {
            state.openedflows.forEach((item)=>{
                if(item.id===state.currentFlow){
                    item.nodes=action.payload;
                    item.isModified=true;
                }
            });
        },
        setEdges:(state,action) => {
            state.openedflows.forEach((item)=>{
                if(item.id===state.currentFlow){
                    item.edges=action.payload;
                    item.isModified=true;
                }
            });
        },
        setView:(state,action) =>{
            state.view=action.payload;
        },
        setCurrentNode:(state,action)=>{
            state.currentNode=action.payload;
        },
        setCurrentFlow:(state,action)=>{
            state.currentFlow=action.payload;
            state.currentNode="";
        },
        openFlow:(state,action)=>{
            //一次允许打开多个flow
            action.payload.forEach(flowID=>{
                state.openedflows.push({
                    loaded:false,
                    isNew:false,
                    isModified:false,
                    id:flowID,
                    nodes:[],
                    edges:[],
                });
            });
            state.currentFlow=action.payload[0];
            state.currentNode="";
        },
        closeFlow:(state,action)=>{
            closeFlowFunc(state,action.payload);
        },
        updateNodeData:(state,action)=>{
            state.openedflows.forEach((item)=>{
                if(item.id===state.currentFlow){
                    item.nodes.forEach(element => {
                        if(element.id===state.currentNode){
                            element.data=action.payload;
                        }
                    });
                    item.isModified=true;
                }
            });
        },
        updateFlowDescription:(state,action)=>{
            const {flowID,description}=action.payload;
            state.openedflows.forEach((item)=>{
                if(item.id===flowID){
                    item.isModified=true;
                    item.description=description;
                }
            });
        },
    },
    extraReducers: (builder) => {
        builder.addCase(flowListAction.pending, (state, action) => {
            state.pending=true;
            state.loaded=true;
        });
        builder.addCase(flowListAction.fulfilled, (state, action) => {
            console.log("flowListAction fulfilled:",action);
            state.pending=false;
            console.log(action);
            if(action.payload.error){
                message.error(action.payload.message);
            } else {
                state.flows=action.payload.result;
            }
        });
        builder.addCase(flowListAction.rejected , (state, action) => {
            console.log("queryFlowsAction return error:",action);
            state.pending=false;
            if(action.error&&action.error.message){
                message.error(action.error.message);
            } else {
                message.error("位置错误！");
            }
        });
        builder.addCase(getFlowConfigAction.pending, (state, action) => {
            state.pending=true;
            state.openedflows.forEach((item)=>{
                if(item.id===state.currentFlow){
                    item.loaded=true;
                }
            });
        });
        builder.addCase(getFlowConfigAction.fulfilled, (state, action) => {
            console.log("getFlowConfigAction fulfilled:",action);
            state.pending=false;
            console.log(action);
            if(action.payload.error){
                message.error(action.payload.message);
            } else {
                state.openedflows.forEach((item)=>{
                    if(item.id===state.currentFlow){
                        item.nodes=action.payload.result.nodes;
                        item.edges=action.payload.result.edges;
                        item.description=action.payload.result.description;
                    }
                });
            }
        });
        builder.addCase(getFlowConfigAction.rejected , (state, action) => {
            console.log("getFlowConfigAction return error:",action);
            state.pending=false;
            if(action.error&&action.error.message){
                message.error(action.error.message);
            } else {
                message.error("位置错误！");
            }
        });
        builder.addCase(saveFlowConfigAction.pending, (state, action) => {
            state.pending=true;
        });
        builder.addCase(saveFlowConfigAction.fulfilled, (state, action) => {
            console.log("saveFlowConfigAction fulfilled:",action);
            state.pending=false;
            console.log(action);
            if(action.payload.error){
                message.error(action.payload.message);
            } else {
                state.openedflows.forEach((item)=>{
                    if(item.id===state.currentFlow){
                        item.isModified=false;
                        item.isNew=false;
                    }
                });
            }
        });
        builder.addCase(saveFlowConfigAction.rejected , (state, action) => {
            console.log("saveFlowConfigAction return error:",action);
            state.pending=false;
            if(action.error&&action.error.message){
                message.error(action.error.message);
            } else {
                message.error("位置错误！");
            }
        });

        builder.addCase(addFlowAction.pending, (state, action) => {
            state.pending=true;
        });
        builder.addCase(addFlowAction.fulfilled, (state, action) => {
            console.log("addFlowAction fulfilled:",action);
            state.pending=false;
            if(action.payload.error){
                message.error(action.payload.message,5);
            } else {
                state.openedflows.forEach((item)=>{
                    if(item.id===state.currentFlow){
                        item.isModified=false;
                        item.isNew=false;
                        item.id=action.payload.result;
                        state.currentFlow=item.id;
                        state.loaded=false;
                    }
                });
            }
        });
        builder.addCase(addFlowAction.rejected , (state, action) => {
            console.log("addFlowAction return error:",action);
            state.pending=false;
            if(action.error&&action.error.message){
                message.error(action.error.message);
            } else {
                message.error("位置错误！");
            }
        });

        builder.addCase(deleteFlowAction.pending, (state, action) => {
            state.pending=true;
        });
        builder.addCase(deleteFlowAction.fulfilled, (state, action) => {
            console.log("deleteFlowAction fulfilled:",action);
            state.pending=false;
            if(action.payload.error){
                message.error(action.payload.message,5);
            } else {
                closeFlowFunc(state,state.currentFlow);
                state.loaded=false;
            }
        });
        builder.addCase(deleteFlowAction.rejected , (state, action) => {
            console.log("deleteFlowAction return error:",action);
            state.pending=false;
            if(action.error&&action.error.message){
                message.error(action.error.message);
            } else {
                message.error("未知错误！");
            }
        });
    } 
});

// Action creators are generated for each case reducer function
export const { 
    addNewFlow,
    setNodes,
    setEdges,
    setView,
    setCurrentNode,
    updateNodeData,
    updateFlowDescription,
    setCurrentFlow,
    openFlow,
    closeFlow
} = flowSlice.actions

export default flowSlice.reducer