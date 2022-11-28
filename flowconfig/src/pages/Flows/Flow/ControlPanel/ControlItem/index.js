import { useCallback, useRef } from 'react';
import './index.css';

export default function ControlItem({label,type}){
    const refControl=useRef(null);

    const onDragStart = useCallback((event) => {
        const controlBounds = refControl.current.getBoundingClientRect();
        const offsetX=event.clientX - controlBounds.left;
        const offsetY=event.clientY - controlBounds.top;
        event.dataTransfer.setData('application/crvflowconfig/offsetX', offsetX);
        event.dataTransfer.setData('application/crvflowconfig/offsetY', offsetY);
        event.dataTransfer.setData('application/crvflowconfig/nodeType', type);
        event.dataTransfer.setData('application/crvflowconfig/nodeLabel', label);
        event.dataTransfer.effectAllowed = 'move';
    },[type,label,refControl]);

    return (<div ref={refControl} draggable className={"control-item control-item-"+type} onDragStart={onDragStart}>{label}</div>);
}