import { useDispatch, useSelector } from "react-redux";
import { Button, Space } from "antd";
import AceEditor from "react-ace";
import 'ace-builds/src-noconflict/mode-json';
import 'ace-builds/src-noconflict/ext-language_tools';
import 'ace-builds/src-noconflict/ext-searchbox';
import 'ace-builds/src-noconflict/theme-monokai';

import {closeDialog} from '../../../../../redux/dialogSlice';
import './index.css';

export default function DialogDebugInfoDetail({dialogIndex}){
    const dispatch = useDispatch();
    const dialogItem=useSelector(state=>state.dialog.dialogs[dialogIndex]);
    const {nodeID,index,flowID}=dialogItem.param;
    const debugInfo=useSelector(state=>state.debug.debugMessages.filter(item=>item.id===nodeID&&item.flowID===flowID)[index]);

    const onCancel=()=>{
        dispatch(closeDialog());
    }

    return (
        <div>
            <AceEditor
                style={{height:400,width:"100%"}}
                placeholder=""
                mode="json"
                theme="monokai"
                name="debugResult"
                fontSize={12}
                showPrintMargin={true}
                showGutter={true}
                highlightActiveLine={true}
                value={JSON.stringify(debugInfo,null,'\t')}
                setOptions={{
                    enableBasicAutocompletion: true,
                    enableLiveAutocompletion: true,
                    enableSnippets: false,
                    showLineNumbers: true,
                    tabSize: 2,
                }}/>
            <div className="dialog-bottom-bar">
                <Space style={{float:'right'}}>
                    <Button style={{minWidth:100}} size="small" onClick={onCancel}>Close</Button>
                </Space>
            </div>
        </div>
    );
}