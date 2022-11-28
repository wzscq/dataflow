import { Modal} from 'antd';
import { useSelector } from 'react-redux';

import DialogOpenFlow from './DialogOpenFlow';
import DialogAddFlow from './DialogAddFlow';
import DialogDebugInfoDetail from '../pages/Flows/Nodes/DebugNode/DialogDebugInfoDetail';
import DialogFuncScript from '../pages/Flows/Nodes/DataTransformNode/DialogFuncScript';
import DialogFilter from '../pages/Flows/Nodes/QueryNode/DialogFilter';
import DialogFlowDescription from '../pages/Flows/Flow/ParamPanel/DialogFlowDescription';

import './index.css';

const dialogRepository={
    openFlow:DialogOpenFlow,
    addFlow:DialogAddFlow,
    debugInfoDetail:DialogDebugInfoDetail,
    funcScript:DialogFuncScript,
    queryFilter:DialogFilter,
    flowDescription:DialogFlowDescription
};

export default function Dialog(){
    const dialogs=useSelector(state=>state.dialog.dialogs);
    if (dialogs.length===0) {
        return null;
    }

    const dialogControls=dialogs.map((item,index)=>{
        const DialogComponent=dialogRepository[item.type];
        if(DialogComponent){
            return (
                <Modal 
                    className={"dialog-"+item.type}
                    title={item.title}
                    closable={false}
                    zIndex={100+index} 
                    visible={true} 
                    centered={true}
                    footer={null}>
                    <DialogComponent dialogIndex={index}/>
                </Modal>
            );
        }
        return null;
    });

    return (<>{dialogControls}</>);
}