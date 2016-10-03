package dcache

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	const (
		dir  = "test"
		size = 10000
	)
	cache, err := NewCache(dir, size)
	if err != nil {
		t.Fatal(err)
	}
	if cache.directory != dir {
		t.Fatal("directory is wrong")
	}
	if cache.maxSize != size {
		t.Fatal("maxSize is wrong")
	}

	cache.Clear()
}

func TestGetSet(t *testing.T) {
	const (
		dir  = "test"
		size = 10000
	)
	cache, err := NewCache(dir, size)
	if err != nil {
		t.Fatal(err)
	}

	const key = "key_1"
	var data []byte
	err = cache.Get(key, data)
	if err == nil {
		t.Fatal("cache already exists")
	}

	err = cache.Set(key, []byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}

	err = cache.Get(key, data)
	if err != nil {
		t.Fatal(err)
	}

	cache.Clear()
}

func TestRemove(t *testing.T) {
	const (
		dir  = "test"
		size = 10000
	)
	cache, err := NewCache(dir, size)
	if err != nil {
		t.Fatal(err)
	}

	const key = "key_2"

	err = cache.Set(key, []byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}

	var data []byte
	err = cache.Get(key, data)
	if err != nil {
		t.Fatal(err)
	}

	err = cache.Remove(key)
	if err != nil {
		t.Fatal(err)
	}

	err = cache.Get(key, data)
	if err == nil {
		t.Fatal("cache already exists")
	}

	cache.Clear()
}

func TestRemoveIfOverMaxSize(t *testing.T) {
	const (
		dir  = "test"
		size = 3
	)
	cache, err := NewCache(dir, size)
	if err != nil {
		t.Fatal(err)
	}

	var key = "key_3"
	err = cache.Set(key, []byte("hello world3"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	key = "key_4"
	err = cache.Set(key, []byte("hello world4"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	key = "key_5"
	err = cache.Set(key, []byte("hello world5"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	key = "key_6"
	err = cache.Set(key, []byte("hello world6"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	key = "key_7"
	err = cache.Set(key, []byte("hello world7"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	key = "key_3"
	var data []byte
	err = cache.Get(key, data)
	if err == nil {
		t.Fatal("key_3 should be removed")
	}

	cache.Clear()
}
