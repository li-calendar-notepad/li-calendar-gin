package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// 参考：https://blog.csdn.net/u014459543/article/details/108429469

// var cacheAdapter *cache.Cache

type GoCacheStruct struct {
	gocahce *cache.Cache
}

// 创建一个goCache结构体
// cache.New(5*time.Minute, 60*time.Second)
func NewGoCache(defaultExpiration time.Duration, cleanupInterval time.Duration) *GoCacheStruct {
	cacheAdapter := cache.New(defaultExpiration, cleanupInterval)
	return &GoCacheStruct{
		gocahce: cacheAdapter,
	}
}

func (c *GoCacheStruct) Set(k string, x interface{}, d time.Duration) {
	c.gocahce.Set(k, x, d)
}

func (c *GoCacheStruct) Get(k string) (interface{}, bool) {
	return c.gocahce.Get(k)
}

//设置cache 无时间参数
func (c *GoCacheStruct) SetDefault(k string, x interface{}) {
	c.gocahce.SetDefault(k, x)
}

//删除 cache
func (c *GoCacheStruct) Delete(k string) {
	c.gocahce.Delete(k)
}

// Add() 加入缓存
func (c *GoCacheStruct) Add(k string, x interface{}, d time.Duration) {
	c.gocahce.Add(k, x, d)
}

// IncrementInt() 对已存在的key 值自增n
func (c *GoCacheStruct) IncrementInt(k string, n int) (num int, err error) {
	return c.gocahce.IncrementInt(k, n)
}

//ItemCount 获取已存在key的数量
func (c *GoCacheStruct) ItemCount() int {
	return c.gocahce.ItemCount()
}

//Flush 删除当前已存在的所有key
func (c *GoCacheStruct) Flush() {
	c.gocahce.Flush()
}
