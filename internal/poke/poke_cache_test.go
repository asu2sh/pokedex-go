package internal

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val PokeMapResult
	}{
		{
			key: "https://example.com",
			val: PokeMapResult{
				Results: []PokeMapLocation{{Name: "example"}},
			},
		},
		{
			key: "https://example.com/path",
			val: PokeMapResult{
				Results: []PokeMapLocation{{Name: "example-path"}},
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewPokeCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if !reflect.DeepEqual(val, c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestStartCacheCleanup(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewPokeCache(baseTime)
	cache.Add("https://example.com", PokeMapResult{Results: []PokeMapLocation{{Name: "example"}}})

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
