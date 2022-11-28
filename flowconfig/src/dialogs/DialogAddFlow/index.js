import { useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Button, Space,Input} from "antd";
import {closeDialog} from '../../redux/dialogSlice';
import {addFlowAction} from '../../api';
import './index.css';

export default function DialogAddFlow({dialogIndex}){
    const dispatch = useDispatch();
    const [flowID,setFlowID]=useState("");
    const dialogItem=useSelector(state=>state.dialog.dialogs[dialogIndex]);
    const currentFlow=useSelector(state=>{
        console.log(state.flow);
        const flows=state.flow.openedflows;
        for(let flowIndex=0;flowIndex<flows.length;flowIndex++){
            if(flows[flowIndex].id===state.flow.currentFlow){
                return flows[flowIndex];
            }
        }
    }); 

    const onCancel=()=>{
        dispatch(closeDialog());
    }

    const onOk=()=>{
        //open flows
        const {nodes,edges}=currentFlow;
        const actionParams={
            appID:dialogItem.param.appID,
            flowID:flowID,
            flowConf:{nodes,edges}
        }
        dispatch(addFlowAction(actionParams));
        dispatch(closeDialog());
    }

    return (
        <div>
            <Input value={flowID} onChange={(e)=>setFlowID(e.target.value)} placeholder="Please input new flow ID"/>
            <div className="dialog-bottom-bar">
                <Space style={{float:'right'}}>
                    <Button style={{minWidth:100}} size="small" onClick={onCancel}>Cancel</Button>
                    <Button disabled={(flowID.length===0)} style={{minWidth:100}} size="small" onClick={onOk}>Ok</Button>
                </Space>
            </div>
        </div>
    );
}