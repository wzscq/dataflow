import {Row,Col,Input,Select} from 'antd';
import { PlusSquareOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';

const { Option } = Select;

export default function CreateMatchResultParams({node,labelWidth}){
    const dispatch=useDispatch();

    const onNodeDataChange=(data)=>{
        dispatch(updateNodeData(data));
    }

    const showTaskInfo=node.data.__showTaskInfo===false?false:true;
    const showTaskStepInfo=node.data.__showTaskStepInfo===false?false:true;

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <div className='button' onClick={()=>onNodeDataChange({...node.data,__showTaskInfo:!showTaskInfo})}>
              {showTaskInfo?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>Task Info</span>
          </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24} style={{display:showTaskInfo?"flex":"none"}}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>Name</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data.task?.name} onChange={(e)=>onNodeDataChange({...node.data,task:{...node.data.task,name:e.target.value}})}/>
            </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24} style={{display:showTaskInfo?"flex":"none"}}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>Status</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={node.data.task?.executionStatus} size='small' onChange={(value)=>onNodeDataChange({...node.data,task:{...node.data.task,executionStatus:value}})}>
                    <Option key='0'>待执行</Option>
                    <Option key='1'>执行中</Option>
                    <Option key='2'>执行完成</Option>
                </Select>
            </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24} style={{display:showTaskInfo?"flex":"none"}}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>Progress</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data.task?.executionProgress} onChange={(e)=>onNodeDataChange({...node.data,task:{...node.data.task,executionProgress:e.target.value}})}/>
            </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24} style={{display:showTaskInfo?"flex":"none"}}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>Result</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={node.data.task?.resultStatus} size='small' onChange={(value)=>onNodeDataChange({...node.data,task:{...node.data.task,resultStatus:value}})}>
                    <Option key='0'>未执行</Option>
                    <Option key='1'>执行成功</Option>
                    <Option key='2'>执行错误</Option>
                </Select>
            </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24} style={{display:showTaskInfo?"flex":"none"}}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>ErrorCode</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data.task?.errorCode} onChange={(e)=>onNodeDataChange({...node.data,task:{...node.data.task,errorCode:e.target.value}})}/>
            </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24} style={{display:showTaskInfo?"flex":"none"}}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>Message</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data.task?.message} onChange={(e)=>onNodeDataChange({...node.data,task:{...node.data.task,message:e.target.value}})}/>
            </Col>
        </Row>

        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <div className='button' onClick={()=>onNodeDataChange({...node.data,__showTaskStepInfo:!showTaskStepInfo})}>
              {showTaskStepInfo?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>Step Info</span>
          </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24} style={{display:showTaskStepInfo?"flex":"none"}}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>Name</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data.step?.name} onChange={(e)=>onNodeDataChange({...node.data,step:{...node.data.step,name:e.target.value}})}/>
            </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24} style={{display:showTaskStepInfo?"flex":"none"}}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>Status</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={node.data.step?.executionStatus} size='small' onChange={(value)=>onNodeDataChange({...node.data,step:{...node.data.step,executionStatus:value}})}>
                    <Option key='0'>待执行</Option>
                    <Option key='1'>执行中</Option>
                    <Option key='2'>执行完成</Option>
                </Select>
            </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24} style={{display:showTaskStepInfo?"flex":"none"}}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>Progress</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data.step?.executionProgress} onChange={(e)=>onNodeDataChange({...node.data,step:{...node.data.step,executionProgress:e.target.value}})}/>
            </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24} style={{display:showTaskStepInfo?"flex":"none"}}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>Result</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={node.data.step?.resultStatus} size='small' onChange={(value)=>onNodeDataChange({...node.data,step:{...node.data.step,resultStatus:value}})}>
                    <Option key='0'>未执行</Option>
                    <Option key='1'>执行成功</Option>
                    <Option key='2'>执行错误</Option>
                </Select>
            </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24} style={{display:showTaskStepInfo?"flex":"none"}}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>ErrorCode</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data.step?.errorCode} onChange={(e)=>onNodeDataChange({...node.data,step:{...node.data.step,errorCode:e.target.value}})}/>
            </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24} style={{display:showTaskStepInfo?"flex":"none"}}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>Message</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data.step?.message} onChange={(e)=>onNodeDataChange({...node.data,step:{...node.data.step,message:e.target.value}})}/>
            </Col>
        </Row>
      </>
    );
}

