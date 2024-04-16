package common

type CommonRsp struct {
	ErrorCode int `json:"errorCode"`
	Message string `json:"message"`
	Error bool `json:"error"`
	Result interface{} `json:"result"`
	Params map[string]interface{} `json:"params"`
}

type CommonError struct {
	ErrorCode int `json:"errorCode"`
	Message string `json:"message"`
	Params map[string]interface{} `json:"params"`
}

const (
	ResultSuccess = 10000000
	ResultJsonDecodeError = 10000007
	ResultSQLError=10000009
	ResultQueryFieldNotFound=10000010
	ResultQueryWrongPagination=10000011
	ResultNotSupported=10000012
	ResultQueryWrongFilter=10000013
	ResultNotSupportedSaveType=10000015
	ResultNotSupportedValueType=10000016
	ResultNoIDWhenUpdate=10000017
	ResultNoVersionWhenUpdate=10000018
	ResultDuplicatePrimaryKey=10000019
	ResultWrongDataVersion=10000020
	ResultNoRelatedModel=10000022
	ResultNoRelatedField=10000023
	ResultNotSupportedFieldType=10000024
	ResultBase64DecodeError=10000028
	ResultNotDeleteData=10000036
	ResultLoadExcelFileError=10000040
	ResultESIFileAlreadyImported=10000041
	ResultExcelSheetNotExist = 10000042
	ResultJsonMarshalError = 10000043 

	ResultWrongRequest = 10100001
	ResultOpenFileError = 10100002
	ResultNoNodeOfGivenID = 10100004
	ResultNoExecutorForNodeType = 10100005
	ResultNodeConfigError = 10100006
	ResultNoDataForGroup = 10100007
	ResultNoKeyFieldForGroup = 10100008
	ResultNotSupportedMatchType = 10100009
	ResultNoLeftMatchData = 10100010
	ResultNoRightMatchData = 10100011
	ResultNoMatchField = 10100012
	ResultMatchFieldTypeError = 10100013
	ResultMatchValueToFloat64Error = 10100014
	ResultFlowIDAlreadyExist = 10100015
	ResultDeleteFlowFileError = 10100016
	ResultNoDataForSaveMatched = 10100017
	ResultNodeFilterConfigError = 10100018
	ResultNodeGroupConfigError = 10100019
	ResultRelatedQueryNoRelatedField = 10100020
	ResultNoBPModelData = 10100021
	ResultNoGroupID = 10100022
	ResultNoWriteoffField = 10100023
	ResultWriteoffValueTypeError = 10100024
	ResultExecuteTransformFunctionError = 10100025
	ResultCreateTransformFunctionError = 10100026
	ResultNoModelField = 10100027
	ResultFieldTypeError = 10100028
	ResultNoDataForExport = 10100029
	ResultNotSupportedNode = 10100030
	ResultCallExternalAPIError = 10100031
)

var errMsg = map[int]CommonRsp{
	ResultExcelSheetNotExist:CommonRsp{
		ErrorCode:ResultExcelSheetNotExist,
		Message:"读取Excel文件时指定的sheet不存在，请确认您选择的Excel文件是否正确",
		Error:true,
	},
	ResultCallExternalAPIError:CommonRsp{
		ErrorCode:ResultCallExternalAPIError,
		Message:"调用外部API时发生错误，请与管理员联系处理",
		Error:true,
	},
	ResultNoDataForExport:CommonRsp{
		ErrorCode:ResultNoDataForExport,
		Message:"没有符合条件的数据",
		Error:true,
	},
	ResultNotSupportedNode:CommonRsp{
		ErrorCode:ResultNotSupportedNode,
		Message:"执行流程时遇到不支持的节点类型，请与管理员联系处理",
		Error:true,
	},
	ResultESIFileAlreadyImported:CommonRsp{
		ErrorCode:ResultESIFileAlreadyImported,
		Message:"您选择的文件和已经导入的文件名称相同，不能重复导入相同的文件",
		Error:true,
	},
	ResultSuccess:CommonRsp{
		ErrorCode:ResultSuccess,
		Message:"操作成功",
		Error:false,
	},
	ResultFieldTypeError:CommonRsp{
		ErrorCode:ResultFieldTypeError,
		Message:"获取模型字段值时，字段类型错误，请与管理员联系处理",
		Error:true,
	},
	ResultNoModelField:CommonRsp{
		ErrorCode:ResultNoModelField,
		Message:"未取到指定模型字段的值，请与管理员联系处理",
		Error:true,
	},
	ResultCreateTransformFunctionError:CommonRsp{
		ErrorCode:ResultCreateTransformFunctionError,
		Message:"创建数据转换函数实例时发生错误，请与管理员联系处理",
		Error:true,
	},
	ResultExecuteTransformFunctionError:CommonRsp{
		ErrorCode:ResultExecuteTransformFunctionError,
		Message:"执行数据转换函数时发生错误，请与管理员联系处理",
		Error:true,
	},
	ResultWrongDataVersion:CommonRsp{
		ErrorCode:ResultWrongDataVersion,
		Message:"您没有修改数据的权限或数据已被其他用户修改，请刷新页面数据后重新尝试",
		Error:true,
	},
	ResultNotSupportedSaveType:CommonRsp{
		ErrorCode:ResultNotSupportedSaveType,
		Message:"保存数据请求中提供的保存操作类型不正确，请与管理员联系处理",
		Error:true,
	},
	ResultNoVersionWhenUpdate:CommonRsp{
		ErrorCode:ResultNoVersionWhenUpdate,
		Message:"更新数据请求中缺少Version字段，请与管理员联系处理",
		Error:true,
	},
	ResultSQLError:CommonRsp{
		ErrorCode:ResultSQLError,
		Message:"执行查询语句时发生错误，请与管理员联系处理",
		Error:true,
	},
	ResultBase64DecodeError:CommonRsp{
		ErrorCode:ResultBase64DecodeError,
		Message:"保存文件时文件内容Base64解码失败，请与管理员联系处理",
		Error:true,
	},
	ResultNoIDWhenUpdate:CommonRsp{
		ErrorCode:ResultNoIDWhenUpdate,
		Message:"更新或删除数据请求中缺少ID字段，请与管理员联系处理",
		Error:true,
	},
	ResultNotDeleteData:CommonRsp{
		ErrorCode:ResultNotDeleteData,
		Message:"删除数据失败，数据不存在或您没有权限删除相应数据",
		Error:true,
	},
	ResultQueryFieldNotFound:CommonRsp{
		ErrorCode:ResultQueryFieldNotFound,
		Message:"执行查询请求中没有提供查询字段信息，请与管理员联系处理",
		Error:true,
	},
	ResultQueryWrongPagination:CommonRsp{
		ErrorCode:ResultQueryWrongPagination,
		Message:"执行查询请求中提供的分页信息不正确，请与管理员联系处理",
		Error:true,
	},
	ResultQueryWrongFilter:CommonRsp{
		ErrorCode:ResultQueryWrongFilter,
		Message:"执行查询请求中提供的过滤信息不正确，请与管理员联系处理",
		Error:true,
	},
	ResultNotSupportedValueType:CommonRsp{
		ErrorCode:ResultNotSupportedValueType,
		Message:"保存数据请求中提供的字段值类型不支持，请与管理员联系处理",
		Error:true,
	},
	ResultNotSupported:CommonRsp{
		ErrorCode:ResultNotSupported,
		Message:"遇到不支持的过滤条件格式，请与管理员联系处理",
		Error:true,
	},
	ResultNotSupportedFieldType:CommonRsp{
		ErrorCode:ResultNotSupportedFieldType,
		Message:"遇到不支持的字段类型，请与管理员联系处理",
		Error:true,
	},
	ResultNoRelatedModel:CommonRsp{
		ErrorCode:ResultNoRelatedModel,
		Message:"关联字段中没有配置对应的关联数据模型，请与管理员联系处理",
		Error:true,
	},
	ResultNoRelatedField:CommonRsp{
		ErrorCode:ResultNoRelatedField,
		Message:"一对多关联字段中没有配置对应的关联字段，请与管理员联系处理",
		Error:true,
	},
	ResultLoadExcelFileError:CommonRsp{
		ErrorCode:ResultLoadExcelFileError,
		Message:"加载Excel文件失败，您选择的Excel文件格式不正确或文件损坏，请选择正确文件并重新尝试",
		Error:true,
	},
	ResultWrongRequest:CommonRsp{
		ErrorCode:ResultWrongRequest,
		Message:"请求参数错误，请检查参数是否完整，参数格式是否正确",
		Error:true,
	},
	ResultOpenFileError:CommonRsp{
		ErrorCode:ResultOpenFileError,
		Message:"打开配置文件时发生错误，请与管理员联系处理",
		Error:true,
	},
	ResultJsonDecodeError:CommonRsp{
		ErrorCode:ResultJsonDecodeError,
		Message:"解析JSON文件时发生错误，请与管理员联系处理",
		Error:true,
	},
	ResultJsonMarshalError:CommonRsp{
		ErrorCode:ResultJsonMarshalError,
		Message:"将对象转换为JSON文本时发生错误，请与管理员联系处理",
		Error:true,
	},
	ResultNoNodeOfGivenID:CommonRsp{
		ErrorCode:ResultNoNodeOfGivenID,
		Message:"执行流时找不到对应ID的节点，请与管理员联系处理",
		Error:true,
	},
	ResultNoExecutorForNodeType:CommonRsp{
		ErrorCode:ResultNoExecutorForNodeType,
		Message:"执行流时遇到不支持的节点类型，请与管理员联系处理",
		Error:true,
	},
	ResultNodeConfigError:CommonRsp{
		ErrorCode:ResultNodeConfigError,
		Message:"节点配置格式不正确，请与管理员联系处理",
		Error:true,
	},
	ResultNoDataForGroup:CommonRsp{
		ErrorCode:ResultNoDataForGroup,
		Message:"执行分组节点时待分组数据为空，请与管理员联系处理",
		Error:true,
	},
	ResultNoDataForSaveMatched:CommonRsp{
		ErrorCode:ResultNoDataForGroup,
		Message:"执行保存已匹配节点时数据为空，请与管理员联系处理",
		Error:true,
	},
	ResultNoKeyFieldForGroup:CommonRsp{
		ErrorCode:ResultNoKeyFieldForGroup,
		Message:"执行分组节点时待分组数据中缺少分组字段，请与管理员联系处理",
		Error:true,
	},
	ResultNotSupportedMatchType:CommonRsp{
		ErrorCode:ResultNotSupportedMatchType,
		Message:"执行匹配分组节点时遇到不支持的匹配类型，请与管理员联系处理",
		Error:true,
	},
	ResultNoLeftMatchData:CommonRsp{
		ErrorCode:ResultNoLeftMatchData,
		Message:"执行匹配分组节点时左表数据不存在，请与管理员联系处理",
		Error:true,
	},
	ResultNoRightMatchData:CommonRsp{
		ErrorCode:ResultNoRightMatchData,
		Message:"执行匹配分组节点时右表数据不存在，请与管理员联系处理",
		Error:true,
	},
	ResultNoMatchField:CommonRsp{
		ErrorCode:ResultNoMatchField,
		Message:"执行匹配分组节点时匹配字段不存在，请与管理员联系处理",
		Error:true,
	},
	ResultDuplicatePrimaryKey:CommonRsp{
		ErrorCode:ResultDuplicatePrimaryKey,
		Message:"创建数据时发现关键字重复，数据库不能创建新的记录",
		Error:true,
	},
	ResultMatchFieldTypeError:CommonRsp{
		ErrorCode:ResultMatchFieldTypeError,
		Message:"执行匹配分组节点时匹配字段数据类型不支持，请与管理员联系处理",
		Error:true,
	},
	ResultMatchValueToFloat64Error:CommonRsp{
		ErrorCode:ResultMatchValueToFloat64Error,
		Message:"执行匹配分组节点时匹配字段值无法转换为数值类型，请与管理员联系处理",
		Error:true,
	},
	ResultFlowIDAlreadyExist:CommonRsp{
		ErrorCode:ResultFlowIDAlreadyExist,
		Message:"保存流失败，已经存在相同ID的流",
		Error:true,
	},
	ResultDeleteFlowFileError:CommonRsp{
		ErrorCode:ResultDeleteFlowFileError,
		Message:"删除流配置文件时发生错误，请刷新后重新尝试",
		Error:true,
	},
	ResultNodeFilterConfigError:CommonRsp{
		ErrorCode:ResultDeleteFlowFileError,
		Message:"过滤节点的配置错误，请与管理员联系处理",
		Error:true,
	},
	ResultNodeGroupConfigError:CommonRsp{
		ErrorCode:ResultNodeGroupConfigError,
		Message:"字段分组节点的配置错误，请与管理员联系处理",
		Error:true,
	},
	ResultRelatedQueryNoRelatedField:CommonRsp{
		ErrorCode:ResultRelatedQueryNoRelatedField,
		Message:"执行关联查询时关联表中的关联字段不存在，请与管理员联系处理",
		Error:true,
	},
	ResultNoBPModelData:CommonRsp{
		ErrorCode:ResultNoBPModelData,
		Message:"执行分组数据处理时未找到对应模型数据，请与管理员联系处理",
		Error:true,
	},
	ResultNoGroupID:CommonRsp{
		ErrorCode:ResultNoGroupID,
		Message:"执行分组数据处理时待处理数据中缺少分组字段，请与管理员联系处理",
		Error:true,
	},
	ResultNoWriteoffField:CommonRsp{
		ErrorCode:ResultNoWriteoffField,
		Message:"执行分组数据处理时待处理数据中缺少核销字段，请与管理员联系处理",
		Error:true,
	},
	ResultWriteoffValueTypeError:CommonRsp{
		ErrorCode:ResultWriteoffValueTypeError,
		Message:"执行分组数据处理时待处理数据核销字段数据类型错误，请与管理员联系处理",
		Error:true,
	},
}

func CreateResponse(err *CommonError,result interface{})(*CommonRsp){
	if err==nil {
		commonRsp:=errMsg[ResultSuccess]
		commonRsp.Result=result
		return &commonRsp
	}

	commonRsp,ok:=errMsg[err.ErrorCode]
	if !ok {
		commonRsp=CommonRsp{
			ErrorCode:err.ErrorCode,
			Message:err.Message,
			Error:true,
		}
	}
	commonRsp.Result=result
	commonRsp.Params=err.Params
	return &commonRsp
}

func CreateError(errorCode int,params map[string]interface{})(*CommonError){
	return &CommonError{
		ErrorCode:errorCode,
		Params:params,
	}
}