package tree

import (
	"fmt"
	"reflect"
)

type Tree[Node any] struct {
	opts Options[Node]
}

func NewTree[Node any](opt ...Option[Node]) *Tree[Node] {
	opts := DefaultOptions[Node]()
	for _, o := range opt {
		o(opts)
	}
	return &Tree[Node]{
		opts: *opts,
	}
}

// 大多数场景，父节点的id都是小于子节点的。如果不是，则下面这个循环会遗漏
func (e *Tree[Node]) Build(data []*Node) ([]*Node, error) {
	idMap := make(map[interface{}]*Node)
	childMap := make(map[interface{}][]*Node)

	// Step 1: 构建 idMap 和 childMap
	for _, item := range data {
		id := reflect.ValueOf(item).Elem().FieldByName(e.opts.Id)
		pid := reflect.ValueOf(item).Elem().FieldByName(e.opts.parentId)

		if !id.IsValid() || !id.CanInterface() || !pid.IsValid() || !pid.CanInterface() {
			return nil, fmt.Errorf("field %s or %s not valid or not accessible", e.opts.Id, e.opts.parentId)
		}

		idValue := id.Interface()
		pidValue := pid.Interface()

		idMap[idValue] = item
		childMap[pidValue] = append(childMap[pidValue], item)
	}

	// Step 2: 将子节点挂载到对应的父节点，并设置层级
	var setLevel func(node *Node, level int)
	setLevel = func(node *Node, level int) {
		if e.opts.lvl != "" {
			lvlField := reflect.ValueOf(node).Elem().FieldByName(e.opts.lvl)
			if lvlField.IsValid() && lvlField.CanSet() {
				if lvlField.Kind() == reflect.Int || lvlField.Kind() == reflect.Int64 {
					lvlField.SetInt(int64(level))
				} else {
					return // lvl字段类型不匹配，跳过
				}
			}
		}

		if children, ok := childMap[reflect.ValueOf(node).Elem().FieldByName(e.opts.Id).Interface()]; ok {
			for _, child := range children {
				if err := e.opts.addChild(node, child); err != nil {
					continue
				}
				setLevel(child, level+1)
			}
		}
	}

	// Step 3: 找出根节点
	var roots []*Node
	for _, node := range data {
		pid := reflect.ValueOf(node).Elem().FieldByName(e.opts.parentId).Interface()

		isRoot, err := compareField(pid, e.opts.pidValue)
		if err != nil {
			return nil, err
		}

		if isRoot {
			setLevel(node, 1) // 根节点设置为第1层
			roots = append(roots, node)
		}
	}

	return roots, nil
}
