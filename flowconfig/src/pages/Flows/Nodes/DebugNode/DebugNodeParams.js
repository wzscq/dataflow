import {Row,Col,Button,Input} from 'antd';
import { useDispatch, useSelector } from 'react-redux';
import { PlusSquareOutlined,MinusSquareOutlined,AlignCenterOutlined } from '@ant-design/icons';
import { updateNodeData } from '../../../../redux/flowSlice';
import {openDialog} from '../../../../redux/dialogSlice';

export default function DebugNodeParams({node,labelWidth}){
    const dispatch=useDispatch();
    const flowID=useSelector(state=>state.flow.currentFlow);
    console.log(flowID);
    const debugInfo=useSelector(state=>state.debug.debugMessages.filter(item=>item.id===node.id&&item.flowID===flowID));
    console.log(debugInfo);
    const onNodeDataChange=(data)=>{
        dispatch(updateNodeData(data));
    }

    const onDetail=(index)=>{
        dispatch(openDialog({type:'debugInfoDetail',title:'Debug Info',param:{nodeID:node.id,index:index,flowID:flowID}}));
    }

    const showDebugInfoList=node.data.__showDebugInfoList===false?false:true;

    const debugInfoList=debugInfo.map((item,index)=>{
        return (
            <Row className="param-panel-row" style={{display:showDebugInfoList?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>{index}</Col>
                <Col className="param-panel-row-inputwithbutton" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                    <Input disabled={true} value={item.startTime}/>
                    <Button className="button"  onClick={()=>{onDetail(index)}} size='small' icon={<AlignCenterOutlined />} />
                </Col>
            </Row>
        );
    });

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <div className='button' onClick={(e)=>onNodeDataChange({...node.data,__showDebugInfoList:!showDebugInfoList})}>
              {showDebugInfoList?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>DebugInfo</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
          </Col>
        </Row>
        {debugInfoList}
      </>
    );
}

