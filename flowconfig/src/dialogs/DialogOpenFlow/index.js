import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Button, Space,List,Checkbox } from "antd";
import { RedoOutlined } from '@ant-design/icons';
import {closeDialog} from '../../redux/dialogSlice';
import {openFlow} from '../../redux/flowSlice';
import {flowListAction} from '../../api';
import './index.css';

export default function DialogOpenFlow({dialogIndex}){
    const dispatch = useDispatch();
    const dialogItem=useSelector(state=>state.dialog.dialogs[dialogIndex]);
    const {flows,openedflows,loaded,pending}=useSelector(state=>state.flow);
    const [selectedItems,setSelectedItems]=useState({});

    const data=flows.filter(flowID=>!openedflows.find(item=>item.id===flowID));

    const onCancel=()=>{
        dispatch(closeDialog());
    }

    const onOk=()=>{
        //open flows
        const selectedFlows=data.filter((flowID,index)=>(selectedItems[index]===index));
        dispatch(openFlow(selectedFlows));
        dispatch(closeDialog());
    }

    const onItemClick=(index)=>{
        if(selectedItems[index]===index){
            delete selectedItems[index]
        } else {
            selectedItems[index]=index
        }
        setSelectedItems({...selectedItems});
    }

    const reloadFlows=()=>{
        dispatch(flowListAction({appID:dialogItem.param.appID}));
    }

    useEffect(
        ()=>{
            if(pending===false&&loaded===false){
                dispatch(flowListAction({appID:dialogItem.param.appID}));
            }
        },
        [pending,loaded,dispatch,dialogItem]
    );

    return (
        <div>
            <div className="dialog-header-bar">
                <Space style={{float:'right'}}>
                    <Button loading={pending} size="small" onClick={reloadFlows} icon={<RedoOutlined />} />
                </Space>
            </div>
            <div className="flow-list">
                <List
                    size="small"
                    bordered
                    dataSource={data}
                    renderItem={(item,index)=>{
                        return (
                            <List.Item size='small'>
                                <Checkbox onChange={()=>onItemClick(index)}>{item}</Checkbox>
                            </List.Item>
                        );
                    }}
                />
            </div>
            <div className="dialog-bottom-bar">
                <Space style={{float:'right'}}>
                    <Button style={{minWidth:100}} size="small" onClick={onCancel}>Cancel</Button>
                    <Button disabled={(Object.keys(selectedItems).length===0)} style={{minWidth:100}} size="small" onClick={onOk}>Ok</Button>
                </Space>
            </div>
        </div>
    );
}