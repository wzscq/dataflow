import {useSelector} from 'react-redux';
import JSONInput from 'react-json-editor-ajrm';
import locale    from 'react-json-editor-ajrm/locale/en';

export default function JSONPanel(){
    const currentFlow=useSelector(state=>{
        const flows=state.flow.openedflows;
        for(let flowIndex=0;flowIndex<flows.length;flowIndex++){
            if(flows[flowIndex].id===state.flow.currentFlow){
                return flows[flowIndex];
            }
        }
    });  

    return (
        <JSONInput
            id          = 'a_unique_id'
            locale      = { locale }
            height      = '100%'
            width='100%'
            placeholder={currentFlow}/>);
}