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

	output := []*Node{}

	for _, item := range data {
		pid := reflect.ValueOf(item).Elem().FieldByName(e.opts.parentId)
		if !pid.IsValid() || !pid.CanInterface() {
			return nil, fmt.Errorf("field %s not valid or not accessible", e.opts.parentId)
		}

		b, err := compareField(pid.Interface(), e.opts.pidValue)
		if err != nil {
			return nil, err
		}

		if b {
			output = append(output, item)
		} else {
			for _, oitem := range output {
				id := reflect.ValueOf(oitem).Elem().FieldByName(e.opts.Id)
				if !id.IsValid() || !id.CanInterface() {
					continue
				}

				b1, err := compareField(pid.Interface(), id.Interface())
				if err != nil || !b1 {
					continue
				}

				if err := e.opts.addChild(item, oitem); err != nil {
					return nil, err
				}
			}
		}
	}

	return output, nil
}
