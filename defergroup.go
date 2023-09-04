// Package defergroup 用于解决依次申请多个资源，途中失败时要关闭已申请时的所有资源的情况。
//
//	例：要申请 f1, f2, f3 三个资源，全部成功时返回三个资源，但其中任何一个失败，则所有资源都要释放。
//	如果成功申请 f1 后，申请 f2 失败，此时要先释放 f1 再返回。
//	对于一次要申请很多资源的场合，每个 if-failure 分支都要释放已申请的资源，很容易错。
//	这个包就是用于解决这种问题。
package defergroup

type Func func()

// DeferGroup 保存所有可能需要调用的资源释放函数，批量释放。
//
// 用法
//
//	func foo() error {
//		dg := new(DeferGroup)
//		defer dg.Do()
//		...
//		if err := allocSomething(); err != nil {
//		    return err
//		} else {
//		    dg.Register(closeFunc1())
//		}
//		...
//		dg.UnregisterAll()
//		return nil
//	 }
type DeferGroup struct {
	funcs []Func
}

// Register 注册一个函数。
func (c *DeferGroup) Register(f Func) {
	c.funcs = append(c.funcs, f)
}

// Do 依次调用所有注册过的函数。推荐用 defer 调用。
func (c *DeferGroup) Do() {
	for _, f := range c.funcs {
		f()
	}
}

// UnregisterAll 取消关闭。在外部函数返回「成功」之前调用，可以避免 Do 生效。
func (c *DeferGroup) UnregisterAll() {
	c.funcs = nil
}
