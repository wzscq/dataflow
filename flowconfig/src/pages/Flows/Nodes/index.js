import StartNode from "./StartNode";
import EndNode from "./EndNode";
import {QueryNode,QueryNodeParams} from "./QueryNode";
import {RequestQueryNode,RequestQueryNodeParams} from "./RequestQueryNode";
import {RelatedQueryNode,RelatedQueryNodeParams} from "./RelatedQueryNode";
import {GroupNode,GroupNodeParams} from "./GroupNode";
import {NumericGroupNode,NumericGroupNodeParams} from "./NumericGroupNode";
import {MatchNode,MatchNodeParams}  from "./MatchNode";
import {VerifyMatchNode,VerifyMatchNodeParams} from "./VerifyMatchNode";
import {VerifyValueNode,VerifyValueNodeParams} from "./VerifyValueNode";
import {NumberCompareNode,NumberCompareNodeParams} from "./NumberCompareNode";
import {FilterNode,FilterNodeParams} from "./FilterNode";
import {DelayNode,DelayNodeParams} from "./DelayNode";
import {SaveMatchedNode,SaveMatchedNodeParams} from "./SaveMatchedNode";
import {CreateMatchResultNode,CreateMatchResultParams} from "./CreateMatchResultNode";
import SaveNotMatchedNode from "./SaveNotMatchedNode";
import SaveNode from "./SaveNode";
import InsertNode from "./InsertNode";
import DeleteNode from "./DeleteNode";
import LogNode from "./LogNode";
import {DebugNode,DebugNodeParams} from "./DebugNode";
import {EBProcessingNode,EBProcessingNodeParams} from "./EBProcessingNode";
import {DataTransferNode,DataTransferParams} from './DataTransferNode';
import {DataTransformNode,DataTransformParams} from './DataTransformNode';
import {GroupTransformNode,GroupTransformParams} from './GroupTransformNode';
import {SplitExtraQuantityNode,SplitExtraQuantityParams} from './SplitExtraQuantityNode';
import {ESINode,ESINodeParams} from './ESINode';
import {CRVRequestNode,CRVRequestNodeParams} from './CRVRequestNode';
import {FlowNode,FlowNodeParams} from './FlowNode';
import {FlowAsyncNode,FlowAsyncNodeParams} from './FlowAsyncNode';
import {ReturnCRVResultNode,ReturnCRVResultNodeParams} from './ReturnCRVResultNode';
import {ReturnCRVErrorNode,ReturnCRVErrorNodeParams} from './ReturnCRVErrorNode';
import {CRVFormNode,CRVFormNodeParams} from './CRVFormNode';
import {TaskInfoNode,TaskInfoParams} from './TaskInfoNode';
import {ExportExcelNode,ExportExcelParams} from './ExportExcelNode';
import {CallExternalAPINodeParams,CallExternalAPINode} from './CallExternalAPINode';

export const nodeTypes={
  start:StartNode,
  end:EndNode,
  query:QueryNode,
	requestQuery:RequestQueryNode,
	relatedQuery:RelatedQueryNode,
	fieldGroup:GroupNode,
	numericGroup:NumericGroupNode,
	match:MatchNode,
	verifyMatch:VerifyMatchNode,
	verifyValue:VerifyValueNode,
	numberCompare:NumberCompareNode,
	filter:FilterNode,
	delay:DelayNode,
	saveMatched:SaveMatchedNode,
	saveNotMatched:SaveNotMatchedNode,
	createMatchResult:CreateMatchResultNode,
	log:LogNode,
	debug:DebugNode,
	ebProcessing:EBProcessingNode,
	save:SaveNode,
	insert:InsertNode,
	delete:DeleteNode,
	dataTransfer:DataTransferNode,
	dataTransform:DataTransformNode,
	groupTransform:GroupTransformNode,
	splitExtraQuantity:SplitExtraQuantityNode,
	esi:ESINode,
	CRVRequest:CRVRequestNode,
	callFlow:FlowNode,
	callFlowAsync:FlowAsyncNode,
	returnCRVResult:ReturnCRVResultNode,
	returnCRVError:ReturnCRVErrorNode,
	CRVForm:CRVFormNode,
	taskInfo:TaskInfoNode,
	exportExcel:ExportExcelNode,
	callExternalAPI:CallExternalAPINode
}

export const nodeParams={
  start:null,
  end:null,
  query:QueryNodeParams,
	requestQuery:RequestQueryNodeParams,
	relatedQuery:RelatedQueryNodeParams,
	fieldGroup:GroupNodeParams,
	numericGroup:NumericGroupNodeParams,
	match:MatchNodeParams,
	verifyMatch:VerifyMatchNodeParams,
	verifyValue:VerifyValueNodeParams,
	numberCompare:NumberCompareNodeParams,
	filter:FilterNodeParams,
	delay:DelayNodeParams,
	saveMatched:SaveMatchedNodeParams,
	createMatchResult:CreateMatchResultParams,
	saveNotMatched:null,
	log:null,
	save:null,
	insert:null,
	delete:null,
	debug:DebugNodeParams,
	ebProcessing:EBProcessingNodeParams,
	dataTransfer:DataTransferParams,
	dataTransform:DataTransformParams,
	groupTransform:GroupTransformParams,
	splitExtraQuantity:SplitExtraQuantityParams,
	esi:ESINodeParams,
	CRVRequest:CRVRequestNodeParams,
	callFlow:FlowNodeParams,
	callFlowAsync:FlowAsyncNodeParams,
	returnCRVResult:ReturnCRVResultNodeParams,
	returnCRVError:ReturnCRVErrorNodeParams,
	CRVForm:CRVFormNodeParams,
	taskInfo:TaskInfoParams,
	exportExcel:ExportExcelParams,
	callExternalAPI:CallExternalAPINodeParams
}