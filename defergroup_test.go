package defergroup

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

// doTest 一个测试用例。
//
//	allocResult 表示一系列「子资源」的申请结果
//	closeCalled 表示这些「子资源」是否被关闭（被调用了关闭函数）
func doTest(allocResult []bool) (closeCalled []bool) {
	closeCalled = make([]bool, len(allocResult))
	type Sub struct{}

	NewSub := func(index int) (*Sub, error) {
		if allocResult[index] {
			return new(Sub), nil
		} else {
			return nil, errors.New("fail")
		}
	}
	CloseSub := func(index int) {
		closeCalled[index] = true
	}

	type Resource struct {
		subs []*Sub
	}

	NewResource := func() (*Resource, error) {
		gc := new(DeferGroup)
		defer gc.Do()

		res := new(Resource)

		for i := range allocResult {
			i := i
			if sub, err := NewSub(i); err != nil {
				return nil, err
			} else {
				res.subs = append(res.subs, sub)
				gc.Register(func() {
					CloseSub(i)
				})
			}
		}

		gc.UnregisterAll()
		return res, nil
	}

	_, _ = NewResource()
	return
}

func TestDeferGroup_Do(t *testing.T) {
	assert.Equal(t, []bool{false, false}, doTest([]bool{true, true}))   // 两个资源都申请成功时，都不释放
	assert.Equal(t, []bool{true, false}, doTest([]bool{true, false}))   // 第一个成功，第二个失败时，释放第一个
	assert.Equal(t, []bool{false, false}, doTest([]bool{false, false})) // 第一个就失败时，都不释放（因为任何子资源都没有申请成功）
}
