package flow

const (
	MATCHTYPE_MANY2MANY = "many2many"
	MATCHTYPE_MANY2ONE = "many2one"
	MATCHTYPE_ONE2MANY = "one2many"
	MATCHTYPE_ONE2ONE = "one2one"
)

type matchExecutor interface {
	getMatchGroup(matchValue *matchValue,forMatch *matchGroup,tolerance matchTolerance)(*matchGroup)
}

func getMatchExecutor(matchType string)(matchExecutor){
	switch matchType {
	case MATCHTYPE_ONE2ONE:
		return &matchExecutorOne2One{}
	case MATCHTYPE_ONE2MANY:
		return &matchExecutorOne2Many{}
	case MATCHTYPE_MANY2ONE:
		return &matchExecutorMany2One{}
	case MATCHTYPE_MANY2MANY:
		return &matchExecutorMany2Many{}
	}
	return nil
} 