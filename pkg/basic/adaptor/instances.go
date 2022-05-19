package adaptor

var inst Instances

// Instances manage the instances of adaptors
type Instances struct {
	rawAPIOper          RawAPIOper // rawAPIOper is the instance for query raw api
	kmsOper             KMSOper
	serviceOper         ServiceOper
	namespaceOper       NamespaceOper
	fileServerOper      FileServerOper
	polyOper            PolyAPIOper
	evalerOper          EvalerOper
	rawPolyOper         RawPolyOper
	apiStatOper         APIStatOper
	appCenterServerOper AppCenterServerOper
}

// GetInst return the adaptor instance
func getInst() *Instances {
	return &inst
}

// Operation is enum of operation, create/update/delete/query/request/polyapi/rawapi
type Operation uint

// RequireScriptReady check if an operation need script ready
func (op Operation) RequireScriptReady() bool {
	if op == OpRequest || op == OpPublish {
		return true
	}
	return false
}

// operations
const (
	OpCreate Operation = iota + 1
	OpUpdate
	OpEdit // edit arrange
	OpDelete
	OpQuery
	OpBuild
	OpAddRawAPI
	OpAddPolyAPI
	OpAddService
	OpAddSub
	OpRequest
	OpPublish
)
