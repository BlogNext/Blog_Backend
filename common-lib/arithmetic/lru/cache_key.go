package lru

import (
	"bytes"
	"fmt"
)

//生成缓存key
func GenerateCacheKey(parameters ...interface{}) string {
	cacheKey := new(bytes.Buffer)

	for _, parameter := range parameters {
		cacheKey.WriteString(fmt.Sprint(parameter))
	}

	return cacheKey.String()
}
