import { useRef, useState } from "react";
import {useSelector } from "react-redux";

import CommonParam from "./CommonParams";
import FlowParams from './FlowParams';
import { nodeParams } from "../../Nodes";
import './index.css';

var isMouseDown=false;
var mouseDownLeft=150;

export default function ParamPanel(){
    const [splitLeft,setSplitLeft]=useState(150);
    const refSplitBar=useRef();

    const onSplitMouseDown=(e)=>{
        //开始拉动，记录鼠标的起始位置
        isMouseDown=true;
        mouseDownLeft=e.clientX
        console.log(mouseDownLeft)
    }

    const onSplitonMouseMove=(e)=>{
        //鼠标移动过程中，改变分割条的位置
        if(refSplitBar.current&&isMouseDown===true){
            console.log(e)
            const diff=e.clientX-mouseDownLeft;
            refSplitBar.current.style.left=splitLeft+diff+'px';
        }
    }

    const onSplitMouseUp=(e)=>{
        //释放鼠标后移动完成，更新控件的位置
        if(refSplitBar.current&&isMouseDown===true){
            console.log(e)
            isMouseDown=false;
            const diff=e.clientX-mouseDownLeft;
            refSplitBar.current.style.left=splitLeft+diff+'px';
            setSplitLeft(splitLeft+diff);
        }
    }

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

    let paramControl=null;
    if(currentNode===null){
        if(currentFlow===null){
            return null;
        }
        paramControl= (<FlowParams flow={currentFlow} labelWidth={splitLeft+12} />);
    } else {
        const NodeParams=nodeParams[currentNode.type];
        paramControl=(<>
            <CommonParam type="label" node={currentNode} labelWidth={splitLeft+12}/>
            {NodeParams?<NodeParams node={currentNode} labelWidth={splitLeft+12}/>:null}
        </>);
    }

    return (
        <div className="param-panel" 
            onMouseMove={onSplitonMouseMove}
            onMouseUp={onSplitMouseUp}
        >
            <div className="param-panel-content">
                {paramControl}
                <div 
                    ref={refSplitBar}
                    className="param-panel-split"  
                    style={{left:splitLeft}}
                    onMouseDown={onSplitMouseDown}
                />
            </div>
        </div>);
}