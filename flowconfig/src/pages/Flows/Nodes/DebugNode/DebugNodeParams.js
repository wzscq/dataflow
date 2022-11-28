import {Row,Col,Button,Input} from 'antd';
import { useDispatch, useSelector } from 'react-redux';
import { PlusSquareOutlined,MinusSquareOutlined,AlignCenterOutlined } from '@ant-design/icons';
import { updateNodeData } from '../../../../redux/flowSlice';
import {openDialog} from '../../../../redux/dialogSlice';

export default function DebugNodeParams({node}){
    const dispatch=useDispatch();
    const debugInfo=useSelector(state=>state.debug.debugMessages.filter(item=>item.id===node.id));

    const onNodeDataChange=(data)=>{
        dispatch(updateNodeData(data));
    }

    const onDetail=(index)=>{
        dispatch(openDialog({type:'debugInfoDetail',title:'Debug Info',param:{nodeID:node.id,index:index}}));
    }

    const showDebugInfoList=node.data.__showDebugInfoList===false?false:true;

    const debugInfoList=debugInfo.map((item,index)=>{
        return (
            <Row className="param-panel-row" style={{display:showDebugInfoList?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-3" span={10}>{index}</Col>
                <Col className="param-panel-row-inputwithbutton" span={14}>
                    <Input disabled={true} value={item.startTime}/>
                    <Button className="button"  onClick={()=>{onDetail(index)}} size='small' icon={<AlignCenterOutlined />} />
                </Col>
            </Row>
        );
    });

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" span={10}>
            <div className='button' onClick={(e)=>onNodeDataChange({...node.data,__showDebugInfoList:!showDebugInfoList})}>
              {showDebugInfoList?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>DebugInfo</span>
          </Col>
          <Col className="param-panel-row-input" span={14}>
          </Col>
        </Row>
        {debugInfoList}
      </>
    );
}

