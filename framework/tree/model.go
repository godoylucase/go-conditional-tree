package tree

type leafs struct {
	False Step
	True  Step
}

type tree struct {
	StartNode *node
}

func New(startNode *node) *tree {
	return &tree{StartNode: startNode}
}

func (t *tree) Resolve(treeCtx interface{}) *Result {
	return t.StartNode.Resolve(treeCtx)
}

func (t *tree) GetType() StepType {
	return TREE
}

type node struct {
	Condition Condition
	Leafs     leafs
}

func NewNode(c Condition, l leafs) *node {
	return &node{
		Condition: c,
		Leafs:     l,
	}
}

func (n *node) Resolve(treeCtx interface{}) *Result {
	if n.Condition(treeCtx) {
		return n.Leafs.True.Resolve(treeCtx)
	}

	return n.Leafs.False.Resolve(treeCtx)
}

func (n *node) GetType() StepType {
	return NODE
}

func (n *node) Bind(branch bool, step Step) {
	if branch {
		n.Leafs.True = step
	} else {
		n.Leafs.False = step
	}
}

type useCase struct {
	Name UseCaseName
}

func NewUseCase(ucn UseCaseName) *useCase {
	return &useCase{Name: ucn}
}

func (uc *useCase) GetType() StepType {
	return USE_CASE
}

func (uc *useCase) Resolve(treeCtx interface{}) *Result {
	return &Result{
		StepInfo: &StepInfo{
			StepType: USE_CASE,
		},
		Value:    uc.Name,
	}
}