package cache

import (
	"strconv"
	"testing"
)

func TestCacheSetAndGet(t *testing.T) {
	cache := New[string, int](10)

	for i := 0; i < 100; i++ {
		key := strconv.Itoa(i)
		err := cache.Set(key, i)
		if err != nil {
			t.Errorf("unexpected error during Set: %v", err)
		}
	}

	for i := 0; i < 90; i++ {
		key := strconv.Itoa(i)
		if cache.Exists(key) {
			t.Errorf("expected key %v to be evicted", key)
		}
	}

	for i := 90; i < 100; i++ {
		key := strconv.Itoa(i)
		if !cache.Exists(key) {
			t.Errorf("expected key %v to still exist", key)
		}
	}
}

func TestCacheSetOverride(t *testing.T) {
	cache := New[string, int](1)

	for i := 0; i < 1; i++ {
		err := cache.Set(strconv.Itoa(i), i)
		if err != nil {
			t.Log("an error occurred when setting entry", err)
		}
	}

	err := cache.Set("1", 100)
	if err != nil {
		t.Log("an error occurred when setting entry", err)
	}

	val, ok := cache.Get("1")
	if cache.order.Len() == 0 && (!ok || val != 100) {
		t.Error("1 should of been evicted from cache")
	}
}

func TestCacheExists(t *testing.T) {
	cache := New[string, int](2)
	err := cache.Set("one", 1)
	if err != nil {
		t.Error("an error occurred when setting entry", err)
	}

	err = cache.Set("two", 2)
	if err != nil {
		t.Error("an error occurred when setting entry", err)
	}

	ok := cache.Exists("one")
	if !ok {
		t.Error("one should be cached")
	}

	ok = cache.Exists("two")
	if !ok {
		t.Error("two should be cached")
	}

	ok = cache.Exists("three")
	if ok {
		t.Error("three isn't a entry")
	}

	err = cache.Set("three", 3)
	if err != nil {
		t.Error("an error occurred when setting entry", err)
	}

	ok = cache.Exists("three")
	if !ok {
		t.Error("three should be cached")
	}

	ok = cache.Exists("one")
	if ok {
		t.Error("one should of been evicted from cache")
	}
}

func TestCacheRemove(t *testing.T) {
	cache := New[string, int](2)
	for i := 0; i < 4; i++ {
		err := cache.Set(strconv.Itoa(i), i)
		if err != nil {
			t.Errorf("unexpected error during Set: %v", err)
		}
	}

	err := cache.Remove("1")
	if err == nil {
		t.Error("1 should not exist in cache")
	}

	err = cache.Remove("3")
	if err != nil {
		t.Error(err)
	}

	ok := cache.Exists("3")
	if ok {
		t.Error("3 should not exist in cache")
	}
}

func TestCacheSize(t *testing.T) {
	cache := New[string, int](2)
	for i := 0; i < 10; i++ {
		err := cache.Set(strconv.Itoa(i), i)
		if err != nil {
			t.Errorf("unexpected error during Set: %v", err)
		}
	}

	size := cache.Size()
	if size != 2 {
		t.Error("cache size should be 2")
	}
}

type user struct {
	name string
}

func TestCacheStruct(t *testing.T) {
	cache := New[string, *user](1)

	err := cache.Set("foo", &user{
		name: "foo",
	})
	if err != nil {
		t.Error(err)
	}

	err = cache.Set("bar", &user{
		name: "bar",
	})
	if err != nil {
		t.Error(err)
	}

	_, ok := cache.Get("foo")
	if ok {
		t.Error("foo should be evicted from cache")
	}

	val, ok := cache.Get("bar")
	if !ok {
		t.Error("bar should be cached")
	}

	if val.name != "bar" {
		t.Error("bar should be bar")
	}
}
