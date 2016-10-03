package dcache

import (
	"fmt"
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

	var (
		key  = "key_1"
		data []byte
	)

	if cache.Get(key, data) {
		t.Fatal("cache already exists")
	}
	if !cache.Set(key, []byte("hello world")) {
		t.Fatal()
	}
	if !cache.Get(key, data) {
		t.Fatal()
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

	var (
		key  = "key_2"
		data []byte
	)
	if !cache.Set(key, []byte("hello world")) {
		t.Fatal()
	}
	if !cache.Get(key, data) {
		t.Fatal()
	}
	if err := cache.Remove(key); err != nil {
		t.Fatal(err)
	}
	if cache.Get(key, data) {
		t.Fatal("cache should be deleted")
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

	for i := 3; i <= 7; i++ {
		key := fmt.Sprintf("key_%d", i)
		if !cache.Set(key, []byte(fmt.Sprintf("hello world %d", i))) {
			t.Fatal()
		}
		time.Sleep(1 * time.Second)
	}

	var (
		key  = "key_3"
		data []byte
	)
	if cache.Get(key, data) {
		t.Fatal("key_3 should be removed")
	}

	cache.Clear()
}

func TestHit(t *testing.T) {
	const (
		dir  = "test"
		size = 3
	)
	cache, err := NewCache(dir, size)
	if err != nil {
		t.Fatal(err)
	}

	for i := 3; i <= 5; i++ {
		key := fmt.Sprintf("key_%d", i)
		if !cache.Set(key, []byte(fmt.Sprintf("hello world %d", i))) {
			t.Fatal()
		}
		time.Sleep(1 * time.Second)
	}

	var (
		key  = "key_3"
		data []byte
	)
	if !cache.Get(key, data) {
		t.Fatal()
	}

	for i := 6; i <= 7; i++ {
		key := fmt.Sprintf("key_%d", i)
		if !cache.Set(key, []byte(fmt.Sprintf("hello world %d", i))) {
			t.Fatal()
		}
		time.Sleep(1 * time.Second)
	}

	key = "key_3"
	if !cache.Get(key, data) {
		t.Fatal()
	}
	key = "key_6"
	if !cache.Get(key, data) {
		t.Fatal()
	}
	key = "key_7"
	if !cache.Get(key, data) {
		t.Fatal()
	}

	cache.Clear()
}
