package tree

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type MyNode struct {
	Id       uint64
	Name     string
	Pid      uint64
	Level    uint32
	Children []*MyNode
}

func TestBuild(t *testing.T) {

	nodes := []*MyNode{
		{Id: 1, Name: "p1", Pid: 0},
		{Id: 2, Name: "p2", Pid: 0},
		{Id: 3, Name: "p3", Pid: 0},
		{Id: 4, Name: "s1", Pid: 1},
		{Id: 5, Name: "s2", Pid: 2},
		{Id: 6, Name: "s3", Pid: 3},
		{Id: 7, Name: "s4", Pid: 8},
		{Id: 8, Name: "s4", Pid: 0},
		{Id: 9, Name: "s6", Pid: 8},
		{Id: 10, Name: "16", Pid: 6},
	}

	tree := NewTree[MyNode]()
	rs, _ := tree.Build(nodes)

	for _, node := range rs {
		printNode(node, 0)
	}

}
func printNode(node *MyNode, level int) {
	prefix := strings.Repeat("  ", level)
	fmt.Printf("%sId: %d, Name: %s, Pid: %d Lvl %d \n", prefix, node.Id, node.Name, node.Pid, node.Level)

	for _, child := range node.Children {
		printNode(child, level+1)
	}
}

func TestReflect(t *testing.T) {
	type Node struct {
		Id       uint64
		Name     string
		Pid      uint64
		Children []Node
	}

	node := Node{
		Id: 1, Name: "p1", Pid: 0,
	}
	v := reflect.ValueOf(node)
	value := v.FieldByName("Id")
	kind := value.Kind()
	b := kind >= reflect.Int && kind <= reflect.Float64
	fmt.Println(b)
	fmt.Println(value)
	var pidValue int = 0
	pv := reflect.ValueOf(pidValue)
	fmt.Println(pv.Kind())

}
