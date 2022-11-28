import { useState } from "react";
import { useSelector,useDispatch } from "react-redux";
import {Space,Button,Input} from 'antd';
import AceEditor from "react-ace";
import 'ace-builds/src-noconflict/mode-javascript';
import 'ace-builds/src-noconflict/ext-language_tools';
import 'ace-builds/src-noconflict/ext-searchbox';
import 'ace-builds/src-noconflict/theme-monokai';

import {closeDialog} from '../../../../../redux/dialogSlice';
import { updateNodeData } from '../../../../../redux/flowSlice';

import './index.css';

export default function DialogFuncScript({dialogIndex}){
    const dispatch = useDispatch();
    const dialogItem=useSelector(state=>state.dialog.dialogs[dialogIndex]);
    const {node,modelIndex,fieldIndex,initContent}=dialogItem.param;
    console.log(dialogItem);
    let initScript="";
    //注意这里的Script对话框被两个节点引用，dataTransform，groupTransform
    //groupTransform中没有modelIndex和fieldIndex
    if(modelIndex>=0&&fieldIndex>=0){
        initScript=node.data?.models[modelIndex].fields[fieldIndex].funcScript;
    } else {
        initScript=node.data.funcScript;
    }

    if(!initScript){
        initScript={content:initContent,name:""};
    }

    console.log("initScript:",JSON.stringify(initScript),modelIndex,fieldIndex);
    
    const [scriptObj,setScriptObj]=useState(initScript);

    const onCancel=()=>{
        dispatch(closeDialog());
    }

    const onOk=()=>{
        if(modelIndex>=0&&fieldIndex>=0){
            const models=[...node.data?.models];
            const modelFields=[...(models[modelIndex].fields)];
            modelFields[fieldIndex]={...modelFields[fieldIndex],funcScript:scriptObj};
            models[modelIndex]={...models[modelIndex],fields:modelFields};
            dispatch(updateNodeData({...node.data,models:models}));
        } else {
            dispatch(updateNodeData({...node.data,funcScript:scriptObj}));
        }
        dispatch(closeDialog());
    }

    const onChange=(newValue)=>{
        setScriptObj({...scriptObj,content:newValue});
    }

    const onFuncNameChange=(e)=>{
        setScriptObj({...scriptObj,name:e.target.value});
    }

    return (
        <div className="dialog-func-script">
            <div>
                <span>Function Name:</span>
                <Input size="small" value={scriptObj?.name} onChange={onFuncNameChange}/>
            </div>
            <div>
                <span>Function Content:</span>
                    <AceEditor
                        style={{height:350,width:"100%"}}
                        placeholder="Placeholder Text"
                        mode="javascript"
                        theme="monokai"
                        name="funcScript"
                        fontSize={12}
                        showPrintMargin={true}
                        showGutter={true}
                        highlightActiveLine={true}
                        onChange={onChange}
                        value={scriptObj?.content}
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