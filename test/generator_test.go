package test

import (
	"go-id-generator/hander"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestIdGenerator(t *testing.T) {
	// 假设节点ID为1
	workerID := int64(1)
	snowflake, err := hander.NewSnowflake(workerID)
	if err != nil {
		t.Errorf(err.Error())
	}
	snowflake.Generate()
}

func TestBatch(t *testing.T) {
	//十台机器码，随机写1W，查看是否有重复
	mp := make(map[int64]int8)
	var wg sync.WaitGroup
	var mu sync.Mutex
	rand.Seed(time.Now().UnixNano())
	num := 10000
	for i := 1; i <= num; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			snowflake, err := hander.NewSnowflake(rand.Int63n(1023))
			if err != nil {
				panic(err)
			}
			id := snowflake.Generate()
			mu.Lock()
			mp[id] = 1
			mu.Unlock()
		}()
	}
	wg.Wait()
	if len(mp) != num {
		t.Errorf("Datacenter Expect ID Num: %d, Actual ID Num: %d\n", num, len(mp))
	}
}
