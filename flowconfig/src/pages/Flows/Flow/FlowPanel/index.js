import {useCallback,useRef,useState} from 'react';
import {useSelector,useDispatch} from 'react-redux';
import ReactFlow, { applyEdgeChanges,applyNodeChanges,addEdge,MiniMap, Controls } from 'react-flow-renderer';

import {setNodes,setEdges,setCurrentNode} from '../../../../redux/flowSlice';
import { nodeTypes } from '../../Nodes';

import './index.css';

export default function FlowPanel(){
    const dispatch=useDispatch();
    const reactFlowWrapper = useRef(null);
    const [reactFlowInstance, setReactFlowInstance] = useState(null);

    const {nodes,edges}=useSelector(state=>{
        const flows=state.flow.openedflows;
        for(let flowIndex=0;flowIndex<flows.length;flowIndex++){
            if(flows[flowIndex].id===state.flow.currentFlow){
                return flows[flowIndex];
            }
        }
    });    

    const onNodesChange = useCallback(
        (changes) => {
            if(changes.length===nodes.length){
                
            } else {
                console.log('onNodesChange',changes,nodes);
                dispatch(setNodes(applyNodeChanges(changes, nodes)));
            }
        },
        [nodes,dispatch]
    );
      
    const onEdgesChange = useCallback(
        (changes) => {
            console.log('onEdgesChange',changes);
            dispatch(setEdges(applyEdgeChanges(changes, edges)))
        },
        [edges,dispatch]
      );

    const onConnect = useCallback(
        (connection) => dispatch(setEdges(addEdge(connection, edges))),
        [edges,dispatch]
    );

    const onDragOver = useCallback((event) => {
        event.preventDefault();
        event.dataTransfer.dropEffect = 'move';
    }, []);

    const getId=()=>{
        let idx=nodes.length;
        while(true){
            const nodeID="dataflownode_"+idx;
            const node=nodes.find(item=>item.id===nodeID);
            if(!node){
                return nodeID;
            }
            idx++;
        }
    }
    
    const onDrop = (event) => {
        event.preventDefault();

        const reactFlowBounds = reactFlowWrapper.current.getBoundingClientRect();
        const type = event.dataTransfer.getData('application/crvflowconfig/nodeType');

        // check if the dropped element is valid
        if (typeof type === 'undefined' || !type) {
            return;
        }

        const offsetX=event.dataTransfer.getData('application/crvflowconfig/offsetX');
        const offsetY=event.dataTransfer.getData('application/crvflowconfig/offsetY');
        const position = reactFlowInstance.project({
            x: event.clientX - reactFlowBounds.left - offsetX,
            y: event.clientY - reactFlowBounds.top - offsetY,
        });
        
        const label = event.dataTransfer.getData('application/crvflowconfig/nodeLabel');
        const newNode = {
            id: getId(),
            type,
            position,
            data: { label: `${label}` },
        };

        console.log("setNodes",nodes.concat(newNode));

        dispatch(setNodes(nodes.concat(newNode)));
    }

    const onSelectionChange=useCallback(
        ({ nodes, edges })=>{
            if(nodes.length===1){
                dispatch(setCurrentNode(nodes[0].id));
            }
        },
        [dispatch]
    )

    const onPaneClick=()=>{
        dispatch(setCurrentNode(""));
    }

    const onNodeDoubleClick=(event, node)=>{
        const newNode = {
            ...node,
            id: getId(),
            position:{
                x:node.position.x+50,
                y:node.position.y+20,
            },
            selected:false
        };
        dispatch(setNodes(nodes.concat(newNode)));
    }

    return (
        <div className='flow-wrapper' ref={reactFlowWrapper}>
            <ReactFlow
                nodes={nodes}
                edges={edges}
                onNodesChange={onNodesChange}
                onEdgesChange={onEdgesChange}
                onConnect={onConnect}
                onDragOver={onDragOver}
                onDrop={onDrop}
                onInit={setReactFlowInstance}
                nodeTypes={nodeTypes}
                onSelectionChange={onSelectionChange}
                deleteKeyCode="Delete"
                onPaneClick={onPaneClick}
                onNodeDoubleClick={onNodeDoubleClick}
                >
                <MiniMap />
                <Controls />
            </ReactFlow>
        </div>
    );
}