package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var mu sync.Mutex
var wg sync.WaitGroup
var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type MapWithExpiration struct {
	sync.RWMutex
	items map[string]interface{}
}

var checkTimeMap = make(map[string]time.Time)

func main() {
	m := MapWithExpiration{
		items: make(map[string]interface{}),
	}
	go m.DeleteTimeCheck()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			key1 := randStr(5)
			key2 := randStr(4)
			m.Set(key1, "value1", i)
			m.Set(key2, "value2", i)
			fmt.Println(m.Get(key1))
			fmt.Println(m.Get(key2))
			fmt.Println("-----------")
			time.Sleep(time.Second * 2)
			fmt.Println(m.Get(key1))
			fmt.Println(m.Get(key2))
			fmt.Println("-----------")
			time.Sleep(time.Second * 2)
			fmt.Println(m.Get(key1))
			fmt.Println(m.Get(key2))
			fmt.Println("-----------")
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func (m *MapWithExpiration) Set(key string, value interface{}, t int) {
	mu.Lock()
	m.items[key] = value
	checkTimeMap[key] = time.Now().Add(time.Second * time.Duration(t))
	mu.Unlock()
}

func (m *MapWithExpiration) Get(key string) (interface{}, bool) {
	m.RLock()
	defer m.RUnlock()
	value, found := m.items[key]
	return value, found
}

func (m *MapWithExpiration) DeleteTimeCheck() {
	for {
		mu.Lock()
		for key, expirationTime := range checkTimeMap {
			if time.Now().After(expirationTime) {
				println("Удален:", key)
				delete(m.items, key)
				delete(checkTimeMap, key)
			}
		}
		mu.Unlock()
		time.Sleep(time.Second)
	}
}

func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
