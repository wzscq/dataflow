import ControlItem from './ControlItem';

import './index.css';

const controls=[
    {type: 'start',label: 'Start'},
    {type: 'end',label: 'End'},
    {type: 'query',label: 'Query Data'},
    {type: 'CRVRequest',label: 'CRV Request'},
    {type: 'requestQuery',label:'Request Query'},
    {type: 'relatedQuery',label:'Related Query'},
    {type: 'esi',label: 'ESI'},
	{type: 'fieldGroup',label: 'Group By Fields'},
    {type: 'numericGroup',label: 'Group By Numeric Field'},
    {type: 'splitExtraQuantity',label: 'Split extra quantity'},
	{type: 'match',label: 'Match By Field'},
    {type: 'dataTransfer',label:'Data Transfer'},
    {type: 'dataTransform',label:'Data Transform'},
    {type: 'groupTransform',label:'Group Transform'},
    {type: 'ebProcessing',label: 'EB Processing'},
	{type: 'verifyMatch',label: 'Match Verify'},
    {type: 'verifyValue',label: 'Value Verify'},
    {type: 'numberCompare',label: 'Number Compare'},
	{type: 'filter',label: 'Filter'},
    {type: 'saveMatched',label: 'Save Matched Group'},
    {type: 'createMatchResult',label: 'Create Match Result'},
	{type: 'saveNotMatched',label: 'Update Not Matched Reason'},
    {type: 'save',label: 'Save'},
    {type: 'flow',label: 'Call Flow'},
    {type: 'delay',label: 'Delay'},
    {type: 'log',label: 'log'},
    {type: 'debug',label: 'Debug'}
];

export default function FlowPanel(){
    const constrols1=controls.map(item=><ControlItem key={item.type} {...item}/>);

    return (
        <div className='control-wrapper'>
            {constrols1}
        </div>
    );
}