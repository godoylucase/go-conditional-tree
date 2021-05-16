package tree

const (
	USE_CASE StepType = (iota + 1) * 100
	NODE
	TREE
)

type UseCaseName string
type StepType uint
type Condition func(treeCtx interface{}) bool

func (st StepType) String() string {
	m := map[StepType]string{
		USE_CASE: "use_case",
		NODE:     "node",
		TREE:     "tree",
	}

	stepType, ok := m[st]
	if !ok {
		panic("step type not supported")
	}

	return stepType
}

type Step interface {
	Resolve(treeCtx interface{}) Result
	GetType() StepType
}

type Bindable interface {
	Bind(branch bool, step Step)
}

type Result struct {
	Value interface{}
	Err   error
}
