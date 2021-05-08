package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
# tree flow scratchpad
startNode
	gte10
		- true:
			containsAnyOdd
				-true:
					gte10_with_odd
				-false:
					gte10_without_odd
		- false:
			containsA5
				- true:
					lt10_with_five
				-false:
					lt10_without_five
*/

type mockedContext struct {
	intValue int
	intSlice []int
}

func buildTree() tree {
	gte10Condition := func(treeCtx interface{}) bool {
		ctx, _ := treeCtx.(mockedContext)
		return ctx.intValue >= 10
	}

	containsAnyOddCondition := func(treeCtx interface{}) bool {
		ctx, _ := treeCtx.(mockedContext)
		for _, nb := range ctx.intSlice {
			if nb%2 != 0 {
				return true
			}
		}
		return false
	}

	gte10WithOddUseCase := NewUseCase("gte10_with_odd")
	gte10WithoutOddUseCase := NewUseCase("gte10_without_odd")

	containsAnyOddNode := NewNode(containsAnyOddCondition)
	containsAnyOddNode.Bind(true, gte10WithOddUseCase)
	containsAnyOddNode.Bind(false, gte10WithoutOddUseCase)

	containsA5Condition := func(treeCtx interface{}) bool {
		ctx, _ := treeCtx.(mockedContext)
		for _, nb := range ctx.intSlice {
			if nb == 5 {
				return true
			}
		}
		return false
	}

	lt10With5 := NewUseCase("lt10_with_five")
	lt10Without5 := NewUseCase("lt10_without_five")

	containsA5Node := NewNode(containsA5Condition)
	containsA5Node.Bind(true, lt10With5)
	containsA5Node.Bind(false, lt10Without5)

	gte10Node := NewNode(gte10Condition)
	gte10Node.Bind(true, containsAnyOddNode)
	gte10Node.Bind(false, containsA5Node)

	return New(gte10Node)
}

func Test_tree_Resolve(t *testing.T) {
	tests := []struct {
		name    string
		tree    tree
		treeCtx interface{}
		want    Result
	}{
		{
			name: "gte10_with_odd",
			tree: buildTree(),
			treeCtx: mockedContext{
				intValue: 10,
				intSlice: []int{10, 10, 10, 9},
			},
			want: Result{
				StepInfo: StepInfo{
					StepType: USE_CASE,
				},
				Value: UseCaseName("gte10_with_odd"),
			},
		},
		{
			name: "gte10_without_odd",
			tree: buildTree(),
			treeCtx: mockedContext{
				intValue: 10,
				intSlice: []int{10, 10, 10, 8},
			},
			want: Result{
				StepInfo: StepInfo{
					StepType: USE_CASE,
				},
				Value: UseCaseName("gte10_without_odd"),
			},
		},
		{
			name: "lt10_with_five",
			tree: buildTree(),
			treeCtx: mockedContext{
				intValue: 9,
				intSlice: []int{10, 10, 10, 5},
			},
			want: Result{
				StepInfo: StepInfo{
					StepType: USE_CASE,
				},
				Value: UseCaseName("lt10_with_five"),
			},
		},
		{
			name: "lt10_without_five",
			tree: buildTree(),
			treeCtx: mockedContext{
				intValue: 9,
				intSlice: []int{10, 10, 10, 4},
			},
			want: Result{
				StepInfo: StepInfo{
					StepType: USE_CASE,
				},
				Value: UseCaseName("lt10_without_five"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := &tree{
				StartNode: tt.tree.StartNode,
			}

			if got := tree.Resolve(tt.treeCtx); !assert.ObjectsAreEqualValues(tt.want, got) {
				t.Fatalf("Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}
