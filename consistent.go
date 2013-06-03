package consistent

import (
    "errors"
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

type uints []uint32

func (x uints) Len() int           { return len(x) }
func (x uints) Less(i, j int) bool { return x[i] < x[j] }
func (x uints) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type Consistent struct {
	circle      map[uint32]string
	members     map[string]bool
	sortedList  uints
	virtualNode int
	sync.RWMutex
}

var ErrEmptyCircle = errors.New("empty consistent circle")

func (c *Consistent) hashKey(key string) uint32 {
	if len(key) < 64 {
		scratch := make([]byte, 64)
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func New() *Consistent {
	return &Consistent{
		circle:      make(map[uint32]string),
		members:     make(map[string]bool),
		sortedList:  make(uints, 0),
		virtualNode: 20,
	}
}

func (c *Consistent) KeyNum(key string, num int) string {
	return key + ":" + fmt.Sprintf("%d", num)
}

func (c *Consistent) Add(key string) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.members[key]; ok {
		return
	}

	for n := 0; n < c.virtualNode; n++ {
		c.circle[c.hashKey(c.KeyNum(key, n))] = key
	}
	c.members[key] = true
	c.updateList()
}

func (c *Consistent) search(key uint32) int {
	i := sort.Search(len(c.sortedList), func(x int) bool { return c.sortedList[x] > key })
	if i >= len(c.sortedList) {
		i = 0
	}
	return i
}

func (c *Consistent) Get(name string) (string, error) {
	c.Lock()
	defer c.Unlock()

	if len(c.circle) == 0 {
		return "", ErrEmptyCircle
	}
	key := c.hashKey(name)
	i := c.search(key)
	return c.circle[c.sortedList[i]], nil
}

func (c *Consistent) Remove(key string) {
	c.Lock()
	defer c.Unlock()

	for n := 0; n < c.virtualNode; n++ {
		delete(c.circle, c.hashKey(c.KeyNum(key, n)))
	}
	delete(c.members, key)
	c.updateList()
}

func (c *Consistent) updateList() {
	list := c.sortedList[:0]
	for k, _ := range c.circle {
		list = append(list, k)
	}
	sort.Sort(list)
	c.sortedList = list
}
