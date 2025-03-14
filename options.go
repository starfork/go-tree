package tree

import (
	"fmt"
	"reflect"
)

// Options 配置项结构
type Options[Node any] struct {
	pidValue int
	Id       string
	parentId string
	children string
	lvl      string

	addChild func(parent, child *Node) error
}

// Option 配置项函数类型
type Option[Node any] func(o *Options[Node])

// WithGetID 配置 getID 的实现
func WithID[Node any](f string) Option[Node] {
	return func(o *Options[Node]) {
		o.Id = f
	}
}
func WithLevl[Node any](f string) Option[Node] {
	return func(o *Options[Node]) {
		o.lvl = f
	}
}
func WithPidValue[Node any](f int) Option[Node] {
	return func(o *Options[Node]) {
		o.pidValue = f
	}
}

// WithGetParentID 配置 getParentID 的实现
func WithParentID[Node any](f string) Option[Node] {
	return func(o *Options[Node]) {
		o.parentId = f
	}
}
func WithChildren[Node any](f string) Option[Node] {
	return func(o *Options[Node]) {
		o.children = f
	}
}

// WithAddChild 配置 addChild 的实现
func WithAddChild[Node any](f func(parent, child *Node) error) Option[Node] {
	return func(o *Options[Node]) {
		o.addChild = f
	}
}

// DefaultOptions 默认配置
func DefaultOptions[Node any]() *Options[Node] {
	var children = "Children"
	return &Options[Node]{
		pidValue: 0,
		Id:       "Id",
		parentId: "Pid",
		lvl:      "Level",
		children: children,

		// addChild: func(parent, child *Node) error {
		// 	v := reflect.ValueOf(child).Elem()
		// 	cf := v.FieldByName(children)
		// 	if !cf.IsValid() || cf.Kind() != reflect.Slice || !cf.CanSet() {
		// 		return fmt.Errorf("field Children is not a settable slice")
		// 	}
		// 	cf.Set(reflect.Append(cf, reflect.ValueOf(parent)))
		// 	return nil
		// },
		addChild: func(parent, child *Node) error {
			v := reflect.ValueOf(parent).Elem()
			cf := v.FieldByName(children)
			if !cf.IsValid() || cf.Kind() != reflect.Slice || !cf.CanSet() {
				return fmt.Errorf("field %s is not a settable slice", children)
			}

			// 追加子节点
			newChildren := reflect.Append(cf, reflect.ValueOf(child))
			cf.Set(newChildren)
			return nil
		},
	}

}

func isInteger(v any) bool {
	kind := reflect.TypeOf(v).Kind()
	return kind >= reflect.Int && kind <= reflect.Uint64
}

// 比较字段vlaue
func compareField(v, target any) (bool, error) {
	if !isInteger(v) || !isInteger(target) {
		return false, fmt.Errorf("v is not integer")
	}
	if !isInteger(target) {
		return false, fmt.Errorf("target is not integer")
	}

	val1 := reflect.ValueOf(v)
	val2 := reflect.ValueOf(target)

	num1 := val1.Convert(reflect.TypeOf(int64(0))).Int()
	num2 := val2.Convert(reflect.TypeOf(int64(0))).Int()
	return num1 == num2, nil
}
