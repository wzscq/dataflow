import { SplitPane } from "react-collapse-pane";
import { useDispatch, useSelector } from "react-redux";
import { useEffect } from "react";

import ControlPanel from './ControlPanel';
import FlowPanel from './FlowPanel';
import ParamPanel from './ParamPanel';
import JSONPanel from "./JSONPanel";
import { getFlowConfigAction } from "../../../api";


export default function Flow({appID}){
    const dispatch=useDispatch();
    const view = useSelector(state=>state.flow.view);
    const currentFlow=useSelector(state=>{
        console.log(state.flow);
        const flows=state.flow.openedflows;
        for(let flowIndex=0;flowIndex<flows.length;flowIndex++){
            if(flows[flowIndex].id===state.flow.currentFlow){
                return flows[flowIndex];
            }
        }
    });  

    useEffect(
        ()=>{
            //如果不是新建的流，则需要从后台加载配置
            if(currentFlow&&currentFlow.loaded===false&&currentFlow.isNew===false){
                dispatch(getFlowConfigAction({appID:appID,flowID:currentFlow.id}));
            }
        },
        [currentFlow,appID,dispatch]
    );

    const flowView=view==='flow'?
    (
        <SplitPane dir='ltr'initialSizes={[15,65,35]} split="vertical" collapse={false}>
            <ControlPanel />
            <FlowPanel key={currentFlow?.id}/>
            <ParamPanel />
        </SplitPane>
    ):(
        <JSONPanel/>
    );
   
    const flowContent=(currentFlow&&currentFlow.loaded===true)?flowView:null;

    return (
        <div className="flow-content">    
            {flowContent}
        </div>
    );
}