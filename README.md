# defergroup

一个「可以撤销的」defer 工具。

## 1. 缘起

在业务代码中常有这样的场景：在一个函数里要申请多个资源并返回。如果所有资源都申请成功，则返回它们全部；如果其中任何一个资源申请失败，要释放**已经申请的资源**。

在 golang 里，通常是在申请成功后立刻 defer 释放资源，但这种方式无法满足「全部成功时不释放」的需求。

本包为这种场景提供了一个方便的解决方案。

## 2. 用法

```go
package main

import "github.com/seedjyh/defergroup"

type Resource struct {
	a *A
	b *B
}

func NewResource() (*Resource, error) {
	// 在最开头创建一个 DeferGroup 并 defer 其 Do
	dg := new(DeferGroup)
	defer dg.Do()

	res := new(Resource)

	if a, err := OpenA(); err != nil {
		return nil, err
	} else {
		res.a = a
		dg.Register(func() { CloseA(a) }) // 每申请成功一个「子资源」就将其关闭函数注册到 dg
	}

	if b, err := OpenB(); err != nil {
		return nil, err
	} else {
		res.b = b
		dg.Register(func() { CloseB(b) })
	}

	// 执行这个 UnregisterAll 可以去掉所有已注册的 Close 函数，从而避免 defer 影响
	dg.UnregisterAll()
	return res, nil
}

```
