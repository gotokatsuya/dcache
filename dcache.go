package dcache

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	directory string
	maxSize   int

	logging bool
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

func (c *Cache) Logging(flag bool) {
	c.logging = flag
}

func (c *Cache) log(v interface{}) {
	if !c.logging {
		return
	}
	log.Println(v)
}

func (c *Cache) getFileNameByKey(key string) string {
	k := sha256.Sum256([]byte(key))
	return path.Join(c.directory, hex.EncodeToString(k[:]))
}

func (c *Cache) Get(key string, data []byte) bool {
	c.RLock()
	defer c.RUnlock()
	fileName := c.getFileNameByKey(key)
	file, err := os.Open(fileName)
	if err != nil {
		c.log(err)
		return false
	}
	if _, err := file.Read(data); err != nil {
		c.log(err)
		return false
	}
	if err := c.hit(fileName); err != nil {
		c.log(err)
	}
	return true
}

func (c *Cache) hit(fileName string) error {
	now := time.Now()
	return os.Chtimes(fileName, now, now)
}

func (c *Cache) Set(key string, data []byte) bool {
	c.Lock()
	defer c.Unlock()
	if err := c.removeIfOverMaxSize(); err != nil {
		c.log(err)
		return false
	}
	fileName := c.getFileNameByKey(key)
	file, err := os.Create(fileName)
	if err != nil {
		c.log(err)
		return false
	}
	if _, err := file.Write(data); err != nil {
		c.log(err)
		return false
	}
	return true
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

func (c *Cache) Remove(key string) error {
	c.Lock()
	defer c.Unlock()
	fileName := c.getFileNameByKey(key)
	return os.Remove(fileName)
}

func (c *Cache) Clear() error {
	c.Lock()
	defer c.Unlock()
	return os.RemoveAll(c.directory)
}
