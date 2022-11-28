import {useSelector } from "react-redux";

import CommonParam from "./CommonParams";
import FlowParams from './FlowParams';
import { nodeParams } from "../../Nodes";
import './index.css';

export default function ParamPanel(){
    const currentFlow=useSelector(state=>{
        let flow=null;
        state.flow.openedflows.forEach(flowItem=>{
            if(flowItem.id===state.flow.currentFlow){
                flow=flowItem;
            }
        });
        return flow;
    });

    const currentNode=useSelector(state=>{
        let node=null;
        if(currentFlow!==null){
            currentFlow.nodes.forEach(nodeItem=>{
                if(nodeItem.id===state.flow.currentNode){
                    node=nodeItem
                }
            });
        }
        return node;
    });

    if(currentNode===null){
        if(currentFlow===null){
            return null;
        }
        return (<FlowParams flow={currentFlow} />);
    }

    const NodeParams=nodeParams[currentNode.type];

    return (
        <div className="param-panel">
            <CommonParam type="label" node={currentNode}/>
            {NodeParams?<NodeParams node={currentNode}/>:null}
        </div>);
}