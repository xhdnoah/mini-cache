package singleflight

import "sync"

// 代表正在进行中，或已经结束的请求
type call struct {
	wg  sync.WaitGroup // 避免重入
	val interface{}
	err error
}

// 主数据结构，管理不同 key 的请求(call
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// 针对相同的 key, 无论 Do 被调用多少次，fn 只会被调用一次，返回 fn 返回值和错误
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		// Lazy Initialization
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		// 如果请求正在进行，则等待
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	// 发请求前加锁
	c.wg.Add(1)
	// 添加到 g.m，表明 key 已有对应的请求在处理
	g.m[key] = c
	g.mu.Unlock()

	// 发起请求
	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	// 更新 g.m
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
