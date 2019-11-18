/*
 * Copyright 2012-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package SpringCore

import (
	"github.com/go-spring/go-spring-parent/spring-const"
)

//
// 定义 Condition 接口
//
type Condition interface {
	Matches(ctx SpringContext) bool
}

//
// 基于方法的 Condition 实现
//
type FunctionCondition struct {
	fn func(ctx SpringContext) bool
}

//
// 工厂函数
//
func NewFunctionCondition(fn func(ctx SpringContext) bool) *FunctionCondition {
	if fn == nil {
		panic("fn can't be null")
	}
	return &FunctionCondition{fn}
}

func (c *FunctionCondition) Matches(ctx SpringContext) bool {
	return c.fn(ctx)
}

//
// 基于 Property 的 Condition 实现
//
type PropertyCondition struct {
	name string
}

//
// 工厂函数
//
func NewPropertyCondition(name string) *PropertyCondition {
	return &PropertyCondition{name}
}

func (c *PropertyCondition) Matches(ctx SpringContext) bool {
	_, ok := ctx.GetDefaultProperty(c.name, "")
	return ok
}

//
// 基于 Bean 的 Condition 实现
//
type BeanCondition struct {
	beanId string
}

//
// 工厂函数
//
func NewBeanCondition(beanId string) *BeanCondition {
	return &BeanCondition{beanId}
}

func (c *BeanCondition) Matches(ctx SpringContext) bool {
	_, ok := ctx.FindBeanByName(c.beanId)
	return ok
}

//
// 基于 Missing Bean 的 Condition 实现
//
type MissingBeanCondition struct {
	beanId string
}

//
// 工厂函数
//
func NewMissingBeanCondition(beanId string) *MissingBeanCondition {
	return &MissingBeanCondition{beanId}
}

func (c *MissingBeanCondition) Matches(ctx SpringContext) bool {
	_, ok := ctx.FindBeanByName(c.beanId)
	return !ok
}

type OpMode int

const (
	OpMode_None OpMode = 0
	OpMode_Or   OpMode = 1
	OpMode_And  OpMode = 2
)

//
// 定义 Condition 链条节点
//
type ConditionNode struct {
	next *ConditionNode // 下一个节点
	op   OpMode         // 计算方式
	cond Condition      // 条件
}

//
// 工厂函数
//
func NewConditionNode() *ConditionNode {
	return &ConditionNode{
		op: OpMode_None,
	}
}

func (c *ConditionNode) Matches(ctx SpringContext) bool {

	if c.next != nil && c.next.cond == nil {
		panic("last op need a cond triggered")
	}

	if c.cond == nil && c.op == OpMode_None {
		return true
	}

	if r := c.cond.Matches(ctx); c.next != nil {

		switch c.op {
		case OpMode_Or: // or
			if r {
				return r
			} else {
				return c.next.Matches(ctx)
			}
		case OpMode_And: // and
			if r {
				return c.next.Matches(ctx)
			} else {
				return false
			}
		default:
			panic("error op mode")
		}

	} else {
		return r
	}
}

//
// 定义 Condition 服务
//
type Conditional struct {
	node *ConditionNode
	curr *ConditionNode
}

//
// 工厂函数
//
func NewConditional() *Conditional {
	node := NewConditionNode()
	return &Conditional{
		node: node,
		curr: node,
	}
}

func (c *Conditional) Matches(ctx SpringContext) bool {
	return c.node.Matches(ctx)
}

//
// c=a||b
//
func (c *Conditional) Or() *Conditional {
	node := NewConditionNode()
	c.curr.op = OpMode_Or
	c.curr.next = node
	c.curr = node
	return c
}

//
// c=a&&b
//
func (c *Conditional) And() *Conditional {
	node := NewConditionNode()
	c.curr.op = OpMode_And
	c.curr.next = node
	c.curr = node
	return c
}

func (c *Conditional) Group(cond *Conditional) *Conditional {
	c.curr.cond = cond.node
	return c
}

func (c *Conditional) ConditionOnProperty(name string) *Conditional {
	c.curr.cond = NewPropertyCondition(name)
	return c
}

func (c *Conditional) ConditionalOnBean(beanId string) *Conditional {
	c.curr.cond = NewBeanCondition(beanId)
	return c
}

func (c *Conditional) ConditionalOnMissingBean(beanId string) *Conditional {
	c.curr.cond = NewMissingBeanCondition(beanId)
	return c
}

func (c *Conditional) ConditionalOnExpression(expression string) *Conditional {
	panic(SpringConst.UNIMPLEMENTED_METHOD)
}

func (c *Conditional) ConditionOnMatches(fn func(ctx SpringContext) bool) *Conditional {
	c.curr.cond = NewFunctionCondition(fn)
	return c
}
