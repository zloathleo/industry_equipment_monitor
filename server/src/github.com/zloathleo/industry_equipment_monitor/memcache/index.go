//缓存实时数据
package memcache

import (
	"github.com/coocood/freecache"
	"github.com/zloathleo/industry_equipment_monitor/common/logger"
	"github.com/zloathleo/industry_equipment_monitor/dstruct"
	"github.com/zloathleo/industry_equipment_monitor/utils"
	"runtime/debug"
)

var GlobalMemCache *PointValueMemCache

type PointValueMemCache struct {
	cache             *freecache.Cache
	CacheValueMap     *dstruct.ValueMap //给每秒历史存储使用的,每秒历史存储执行后清空
	LastCacheValueMap *dstruct.ValueMap //给每小时最后时刻值迁移使用的,每小时执行后清空
}

func InitMemCahce() {
	GlobalMemCache = newPointValueCacheSystem()
	logger.Warnln("memcache init ok.")
}

func newPointValueCacheSystem() *PointValueMemCache {
	caches := new(PointValueMemCache)
	cacheSize := 10 * 1024 * 1024
	caches.cache = freecache.NewCache(cacheSize)
	caches.CacheValueMap = dstruct.NewValueMap()
	caches.LastCacheValueMap = dstruct.NewValueMap()
	debug.SetGCPercent(10)
	return caches
}

//获取缓存中的当前值
func (caches *PointValueMemCache) GetCurrentValue(pn string) (bool, float64) {
	valueBytes, _ := caches.cache.Get([]byte(pn))
	if valueBytes != nil {
		return true, utils.Float64FromBytes(valueBytes)
	} else {
		return false, 0
	}
}

//获取缓存中的当前值列表
func (caches *PointValueMemCache) GetCurrentValueList() map[string]float64 {
	cacheCopy := caches.LastCacheValueMap.SafeCopyAndClear()
	result := make(map[string]float64)
	for key := range cacheCopy {
		exist, value := caches.GetCurrentValue(key)
		if exist {
			result[key] = value
		}
	}
	return result
}

//保存值到缓存
func (caches *PointValueMemCache) SaveCurrentValue(pn string, value float64) {
	//60s 超时
	caches.cache.Set([]byte(pn), utils.Float64ToBytes(value), 60)
	caches.CacheValueMap.Set(pn, value)
	caches.LastCacheValueMap.Set(pn, value)
}
