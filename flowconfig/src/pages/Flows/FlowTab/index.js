import {useDispatch, useSelector} from 'react-redux';
import {Tabs} from 'antd';

import {addNewFlow,setCurrentFlow,closeFlow} from '../../../redux/flowSlice';

import './index.css';

const { TabPane } = Tabs;

export default function FlowTab(){
    const dispatch=useDispatch();
    const {openedflows,currentFlow}=useSelector(state=>state.flow);
    
    const flowTabs=openedflows.map(item=>{
        console.log(item);
        const tabTitle=item.isModified===true?(
                <><span style={{color:'red'}}>*</span><span>{item.id}</span></>
            ):(
                <span>{item.id}</span>
            );
        return (<TabPane tab={tabTitle} key={item.id} />);
    });

    const onTabChange=(newActiveKey)=>{
        console.log('onTabChange',newActiveKey)
        dispatch(setCurrentFlow(newActiveKey));
    }

    const onTabEdit=(targetKey,action)=>{
        if (action === 'add') {
            dispatch(addNewFlow());
        } else {
            dispatch(closeFlow(targetKey));
        }
    }

    return (
        <div className="flow-tab-bar">
            <Tabs type="editable-card" size='small' onChange={onTabChange} activeKey={currentFlow} onEdit={onTabEdit}>
                {flowTabs}
            </Tabs>
        </div>
    );
}