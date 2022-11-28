import { useState } from "react";
import { useSelector,useDispatch } from "react-redux";
import {Space,Button,Input} from 'antd';
import AceEditor from "react-ace";
import 'ace-builds/src-noconflict/mode-json';
import 'ace-builds/src-noconflict/ext-language_tools';
import 'ace-builds/src-noconflict/ext-searchbox';
import 'ace-builds/src-noconflict/theme-monokai';

import {closeDialog} from '../../../../../redux/dialogSlice';
import { updateFlowDescription } from '../../../../../redux/flowSlice';

import './index.css';

export default function DialogFlowDescription({dialogIndex}){
    const dispatch = useDispatch();
    const dialogItem=useSelector(state=>state.dialog.dialogs[dialogIndex]);
    const {flowID,description}=dialogItem.param;
    const [flowDesc,setFlowDesc]=useState(description);

    const onCancel=()=>{
        dispatch(closeDialog());
    }

    const onOk=()=>{
        dispatch(updateFlowDescription({flowID:flowID,description:flowDesc}));
        dispatch(closeDialog());
    }

    const onChange=(newValue)=>{
        setFlowDesc(newValue);
    }

    return (
        <div className="dialog-func-script">
            <div style={{maxHeight:350,overflowY:'auto'}}>
                <AceEditor
                    placeholder="Placeholder Text"
                    mode="json"
                    theme="monokai"
                    name="funcScript"
                    fontSize={12}
                    showPrintMargin={true}
                    showGutter={true}
                    highlightActiveLine={true}
                    onChange={onChange}
                    value={flowDesc}
                    setOptions={{
                    enableBasicAutocompletion: true,
                    enableLiveAutocompletion: true,
                    enableSnippets: false,
                    showLineNumbers: true,
                    tabSize: 2,
                }}/>
            </div>
            <div className="dialog-bottom-bar">
                <Space style={{float:'right'}}>
                    <Button style={{minWidth:100}} size="small" onClick={onOk}>Ok</Button>
                    <Button style={{minWidth:100}} size="small" onClick={onCancel}>Close</Button>
                </Space>
            </div>
        </div>
    );
}