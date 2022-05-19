package errcode

import (
	"github.com/quanxiang-cloud/polyapi/pkg/lib/errdefiner"
)

// error code defines
var (
	//*************************ADD NEW ONES HERE...

	ErrAPIPath                 = r.MustReg(140019990043, "无效的API Path：%v")
	ErrInitAppPath             = r.MustReg(140019990042, "应用初始化失败：%v")
	ErrUseDifferentApp         = r.MustReg(140019990041, "不能跨APP引用：%v")
	ErrServiceAuth             = r.MustReg(140019990040, "鉴权配置错误：%s")
	ErrExistPoly               = r.MustReg(140019990039, "api（%v）已被编排")
	ErrExceedingMaximumLimit   = r.MustReg(140019990038, "数据长度超过 %v")
	ErrUniqueCustomerRootField = r.MustReg(140019990037, "body带有'%v'，不能再定义其他字段")
	ErrUnrecognizedAPIType     = r.MustReg(140019990036, "无法识别的API类型'%v'")
	ErrTooLong                 = r.MustReg(140019990035, "标识%v长度超过%v")
	ErrParameterIn             = r.MustReg(140019990034, "无效的in参数:%s")
	ErrDuplicateName           = r.MustReg(140019990033, "参数名字重复:%s")
	ErrNoNamespacePermit       = r.MustReg(140019990032, "没有该分组访问权限")
	ErrIsolateAuthType         = r.MustReg(140019990031, "组外API仅支持无密钥鉴权方法(none,sytem)")
	ErrAuthContent             = r.MustReg(140019990030, "鉴权方法配置不合法：%v")
	ErrCharacterSet            = r.MustReg(140019990029, "不支持的字符：%v")
	ErrCharacterTooLong        = r.MustReg(140019990028, "标识（%v）长度超过 %v")
	ErrInvalidHost             = r.MustReg(140019990027, "主机地址格式不正确：%s")
	ErrSignatureCMDError       = r.MustReg(140019990026, "鉴权方法配置错误：%s")
	ErrImportDuplicateID       = r.MustReg(140019990025, "存在重复ID")
	ErrAPIInvalid              = r.MustReg(140019990024, "相关API已失效，不能进行该操作")
	ErrMissingPathArg          = r.MustReg(140019990023, "节点%s：URL [%s] 缺少路径参数 %s")
	ErrAPIKeyServiceMismatch   = r.MustReg(140019990022, "API私钥分组设置不匹配")
	ErrServiceWithoutKeys      = r.MustReg(140019990021, "分组设置下尚未上传API私钥")
	ErrPolyNotBuild            = r.MustReg(140019990020, "该编排API尚未准备就绪")
	ErrServiceWithRaw          = r.MustReg(140019990019, "该分组设置包含API")
	ErrServiceWithKeys         = r.MustReg(140019990018, "该分组设置包含私钥")
	ErrNSWithService           = r.MustReg(140019990017, "该分组包含分组设置")
	ErrNSWithPoly              = r.MustReg(140019990016, "该分组包含聚合API")
	ErrNSWithRaw               = r.MustReg(140019990015, "该分组包含API")
	ErrNSWithSubs              = r.MustReg(140019990014, "该分组包含子分组")
	ErrCreateExistsRaw         = r.MustReg(140019990013, "要注册的API已存在")
	ErrCreateExistsPoly        = r.MustReg(140019990012, "要创建的编排API已存在")
	ErrCreateExistsService     = r.MustReg(140019990011, "要创建的分组设置已存在")
	ErrCreateExistsNS          = r.MustReg(140019990010, "要创建的分组已存在：%v")
	ErrActiveEnabled           = r.MustReg(140019990009, "对象启用中，不能进行该操作")
	ErrActiveDisabled          = r.MustReg(140019990008, "对象未启用，不能进行该操作")
	ErrInvalidAppPathType      = r.MustReg(140019990007, "不支持的AppPathType [%v], 仅支持：%v")
	ErrInvalidDocType          = r.MustReg(140019990006, "不支持的DocType [%v], 仅支持：%v")
	ErrInvalidAuthType         = r.MustReg(140019990005, "不支持的AuthType [%v], 仅支持：%v")
	ErrAPINotOpen              = r.MustReg(140019990004, "该API暂未开放")
	ErrNameInvalid             = r.MustReg(140019990003, "名字标识（%v）不合法，只支持：数字/英文字母/下划线/中划线")
	ErrNameTooLong             = r.MustReg(140019990002, "名字标识(%v)长度超过%v")
	ErrNotFound                = r.MustReg(140019990001, "数据未找到")

	//--------------------------------------------------------------------------

	ErrBuildNodeCnt            = r.MustReg(140010030001, "编排节点数必须多于2")
	ErrBuildMissOutput         = r.MustReg(140010030002, "缺少输出节点 %s(%s)")
	ErrBuildIsolateNode        = r.MustReg(140010030003, "编排节点%s(%s)不可达")
	ErrBuildNoNameNode         = r.MustReg(140010030004, "不能编排无名节点%s")
	ErrBuildMissStart          = r.MustReg(140010030005, "编排必须从输入节点%s(%s)开始")
	ErrBuildDuplicateNode      = r.MustReg(140010030006, "重复引用节点%s(%s)")
	ErrBuildDuplicateNodeName  = r.MustReg(140010030007, "节点重名%s(%s)-%s(%s)")
	ErrBuildDuplicateNodeAlias = r.MustReg(140010030008, "节点别名重复%s(%s)-%s(%s)")
	ErrBuildMultiInput         = r.MustReg(140010030009, "编排多于1个输入节点%s-%s")
	ErrBuildOutputName         = r.MustReg(140010030010, "输出节点必须命名为%s, 而不是%s")
	ErrBuildOutputWithNext     = r.MustReg(140010030011, "输出节点%s不能有接续节点%v")
	ErrBuildMultiOutput        = r.MustReg(140010030012, "编排多于1个输出节点%s-%s")
	ErrBuildIfWithNext         = r.MustReg(140010030013, "条件节点%s不能有接续节点%v")
	ErrBuildRequestWithoutNext = r.MustReg(140010030014, "请求节点%s缺少接续节点")
	ErrBuildUnknownNodeType    = r.MustReg(140010030015, "不支持的节点类型%s(%s)")
	ErrBuildUnknownNextNode    = r.MustReg(140010030016, "从节点%s(%s)引用的接续节点%s不存在")
	ErrBuildNodeFail           = r.MustReg(140010030017, "节点(%s)编排失败：%s")
	//ErrBuildDifferentAPP             = r.MustReg(140010030018, "不能编排不同APP下的API: %s")
	ErrVMExecuteFail                 = r.MustReg(140010030019, "编排API执行失败,请检查编排重新调试上线")
	ErrVMEvalFail                    = r.MustReg(140010030020, "参数计算失败:%s")
	ErrBuildNodeName                 = r.MustReg(140010030021, "节点标识(%s)不符合规范")
	ErrBuildInvalidOperator          = r.MustReg(140010030022, "%v不是有效操作符")
	ErrBuildIfNodeMissingBothYN      = r.MustReg(140010030023, "条件节点%v不能同时缺失是和否两个分支")
	ErrBuildRequestNodeMissingRawAPI = r.MustReg(140010030024, "请求节点未配置API路径")
	ErrBuildInputRequiredButMissing  = r.MustReg(140010030025, "必填数据%v未配置")

	//--------------------------------------------------------------------------
	// Signature

	ErrInputArgValidateMismatch = r.MustReg(140010020004, "参数验证失败：%v.%v")
	ErrInputValueExpired        = r.MustReg(140010020003, "参数已过期：%v.%v")
	ErrInputValueInvalid        = r.MustReg(140010020002, "参数值不合法：%v.%v")
	ErrInputMissingArg          = r.MustReg(140010020001, "缺少参数：%v.%v")

	//--------------------------------------------------------------------------

	ErrTimestampFormat   = r.MustReg(140010010007, "时间格式'%v'错误, 仅支持%v")
	ErrParameterError    = r.MustReg(140010010006, "输入参数错误：%v")
	ErrDataFormatInvalid = r.MustReg(140010010005, "格式不合法：%v.%v(%v)")
	ErrInternal          = r.MustReg(140010010004, "内部错误:%v")
	ErrSystemBusy        = r.MustReg(140010010003, "系统繁忙，请稍候重试")
	ErrGateBlockedAPI    = r.MustReg(140010010002, "API服务繁忙，请稍候重试")
	ErrGateBlockedIP     = r.MustReg(140010010001, "访问受限，请联系管理员")

	//--------------------------------------------------------------------------
)

//------------------------------------------------------------------------------

// exports
const (
	ErrParams = errdefiner.ErrParams
	Internal  = errdefiner.Internal
	Unknown   = errdefiner.Unknown
	Success   = errdefiner.Success
)

// func exports
var (
	Errorf             = errdefiner.Errorf
	NewErrorWithString = errdefiner.NewErrorWithString
	r                  = errdefiner.NewErrorDefiner()
)
