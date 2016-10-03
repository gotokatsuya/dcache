package dcache

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

type Cache struct {
	sync.RWMutex
	directory string
	maxSize   int
}

func NewCache(dir string, max int) (*Cache, error) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	if max == 0 {
		return nil, errors.New("max is zero")
	}

	return &Cache{
		directory: dir,
		maxSize:   max,
	}, nil
}

func (c *Cache) Get(key string, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	fileName := c.getFileNameByKey(key)
	file, err := os.Open(path.Join(c.directory, fileName))
	if err != nil {
		return err
	}
	if _, err := file.Read(data); err != nil {
		return err
	}
	return nil
}

func (c *Cache) Set(key string, data []byte) error {
	c.Lock()
	defer c.Unlock()
	if err := c.removeIfOverMaxSize(); err != nil {
		return err
	}
	fileName := c.getFileNameByKey(key)
	file, err := os.Create(path.Join(c.directory, fileName))
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return err
	}
	return nil
}

func (c *Cache) removeIfOverMaxSize() error {
	files, err := ioutil.ReadDir(c.directory)
	if err != nil {
		return err
	}
	if len(files) < c.maxSize {
		return nil
	}
	sortedFileInfos := SortFileInfosByModTimeAsc(files)
	return os.Remove(path.Join(c.directory, sortedFileInfos[0].Name()))
}

func (c *Cache) getFileNameByKey(key string) string {
	k := sha256.Sum256([]byte(key))
	return hex.EncodeToString(k[:])
}

func (c *Cache) Remove(key string) error {
	c.Lock()
	defer c.Unlock()
	fileName := c.getFileNameByKey(key)
	return os.Remove(path.Join(c.directory, fileName))
}

func (c *Cache) Clear() error {
	c.Lock()
	defer c.Unlock()
	return os.RemoveAll(c.directory)
}
