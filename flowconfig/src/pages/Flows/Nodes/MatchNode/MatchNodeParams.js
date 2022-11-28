import {Row,Col,Button,Input,Select} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusSquareOutlined,MinusOutlined } from '@ant-design/icons';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';

const { Option } = Select;

export default function MatchNodeParams({node}){
    const dispatch=useDispatch();

    const setShowSteps=()=>{
          const showSteps=node.data?.__showSteps===false?true:false;
          dispatch(updateNodeData({...node.data,__showSteps:showSteps}));
        }

    const setShowModels=()=>{
        const showModels=node.data?.__showModels===false?true:false;
        dispatch(updateNodeData({...node.data,__showModels:showModels}));
      }

    const onAddMatchStep=()=>{
      const steps=node.data?.steps?[...node.data?.steps]:[];
      steps.push({matchType:"",tolerance:{left:"0",right:"0"},__showStep:true,__showTolerance:true});
      dispatch(updateNodeData({...node.data,steps:steps}));
    }

    const onAddMatchModel=()=>{
      const models=node.data?.models?[...node.data?.models]:[];
      models.push({side:"left",modelID:"",field:"",__showModel:true});
      dispatch(updateNodeData({...node.data,models:models}));
    }

    const onDelStep=(index)=>{
        const steps=[...node.data?.steps];
        delete steps[index];
        dispatch(updateNodeData({...node.data,steps:steps.filter(item=>item)}));
    }

    const onDelModel=(index)=>{
      const models=[...node.data?.models];
      delete models[index];
      dispatch(updateNodeData({...node.data,models:models.filter(item=>item)}));
  }
    
    const onStepChange=(index,value)=>{
        const steps=[...node.data?.steps];
        steps[index]=value;
        dispatch(updateNodeData({...node.data,steps:steps}))
    }

    const onModelChange=(index,value)=>{
        const models=[...node.data?.models];
        models[index]=value;
        dispatch(updateNodeData({...node.data,models:models}))
    }

    const showSteps=node.data?.__showSteps===false?false:true;
    const showModels=node.data?.__showModels===false?false:true;

    const models=node.data.models?.map((item,index)=>{
      return (<>
        <Row className="param-panel-row" style={{display:showModels?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-1" span={10}>
                <div className='button' onClick={()=>onModelChange(index,{...item,__showModel:!item.__showModel})}>
                    {item.__showModel?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                Model {index}
            </Col>
            <Col className="param-panel-row-input" span={14}>
                <Button className="button"  onClick={()=>onDelModel(index)} size='small' icon={<MinusOutlined />} />
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showModels&&item.__showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" span={10}>ModelID</Col>
            <Col className="param-panel-row-input" span={14}>
                <Input value={item.modelID} onChange={(e)=>onModelChange(index,{...item,modelID:e.target.value})}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showModels&&item.__showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" span={10}>Field</Col>
            <Col className="param-panel-row-input" span={14}>
                <Input value={item.field} onChange={(e)=>onModelChange(index,{...item,field:e.target.value})}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showModels&&item.__showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" span={10}>Side</Col>
            <Col className="param-panel-row-input" span={14}>
                <Select value={item.side} size='small' onChange={(value)=>onModelChange(index,{...item,side:value})}>
                    <Option key='left'>Left</Option>
                    <Option key='right'>Right</Option>
                </Select>
            </Col>
        </Row>
      </>);
    })
   
    const steps=node.data.steps?.map((item,index)=>{
        return (<>
            <Row className="param-panel-row" style={{display:showSteps?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-1" span={10}>
                    <div className='button' onClick={()=>onStepChange(index,{...item,__showStep:!item.__showStep})}>
                        {item.__showStep?<MinusSquareOutlined />:<PlusSquareOutlined />}
                    </div>
                    Step {index}
                </Col>
                <Col className="param-panel-row-input" span={14}>
                    <Button className="button"  onClick={()=>onDelStep(index)} size='small' icon={<MinusOutlined />} />
                </Col>
            </Row>
            <Row className="param-panel-row" style={{display:showSteps&&item.__showStep?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-2" span={10}>Match Type</Col>
                <Col className="param-panel-row-input" span={14}>
                    <Select value={item.matchType} size='small' onChange={(value)=>onStepChange(index,{...item,matchType:value})}>
                        <Option key='one2one'>one2one</Option>
                        <Option key='one2many'>one2many</Option>
                        <Option key='many2one'>many2one</Option>
                        <Option key='many2many'>many2many</Option>
                    </Select>
                </Col>
            </Row>
            <Row className="param-panel-row" style={{display:showSteps&&item.__showStep?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-2" span={10}>
                    <div className='button' onClick={()=>onStepChange(index,{...item,__showTolerance:!item.__showTolerance})}>
                        {item.__showTolerance?<MinusSquareOutlined />:<PlusSquareOutlined />}
                    </div>
                    Tolerance {index}
                </Col>
                <Col className="param-panel-row-input" span={14}>
                    <sapn>{"["+item.tolerance.left+","+item.tolerance.right+"]"}</sapn>
                </Col>
            </Row>
            <Row className="param-panel-row" style={{display:showSteps&&item.__showStep&&item.__showTolerance?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-3" span={10}>Left</Col>
                <Col className="param-panel-row-input" span={14}>
                    <Input value={item.tolerance.left} onChange={(e)=>onStepChange(index,{...item,tolerance:{...item.tolerance,left:e.target.value}})}/>
                </Col>
            </Row>
            <Row className="param-panel-row" style={{display:showSteps&&item.__showStep&&item.__showTolerance?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-3" span={10}>Right</Col>
                <Col className="param-panel-row-input" span={14}>
                    <Input value={item.tolerance.right} onChange={(e)=>onStepChange(index,{...item,tolerance:{...item.tolerance,right:e.target.value}})}/>
                </Col>
            </Row>
        </>);
    });

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" span={10}>
            <div className='button' onClick={setShowModels}>
              {showModels?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>Models</span>
          </Col>
          <Col className="param-panel-row-input" span={14}>
            <Button onClick={onAddMatchModel} className='button' size='small' icon={<PlusOutlined />} />
          </Col>
        </Row>
        {models}
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" span={10}>
            <div className='button' onClick={setShowSteps}>
              {showSteps?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>Steps</span>
          </Col>
          <Col className="param-panel-row-input" span={14}>
            <Button onClick={onAddMatchStep} className='button' size='small' icon={<PlusOutlined />} />
          </Col>
        </Row>
        {steps}
      </>
    );
}

