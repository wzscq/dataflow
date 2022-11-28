import {Row,Col,Button,Input} from 'antd';
import { AlignCenterOutlined } from '@ant-design/icons';
import { useDispatch } from "react-redux";

import {openDialog} from '../../../../redux/dialogSlice';

export default function FlowParam({flow}){
    const dispatch = useDispatch();

    const onEditDescription=()=>{
        dispatch(openDialog({type:'flowDescription',title:'Flow Description',param:{flowID:flow.id,description:flow.description}}));
    }

    return (
        <>
            <Row className="param-panel-row" gutter={24}>
                <Col className="param-panel-row-title" span={24}>Flow: {flow.id} </Col>
            </Row>
            <Row className="param-panel-row" gutter={24}>
                <Col className="param-panel-row-label" span={10}>Description</Col>
                <Col className="param-panel-row-inputwithbutton" span={14}>
                    <Input disabled={true} value={flow.description}/>
                    <Button className="button"  onClick={onEditDescription} size='small' icon={<AlignCenterOutlined />} />
                </Col>
            </Row>
        </>
    );
}