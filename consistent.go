//modify by gonghh

package consistent

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

type entry struct {
	Idx  uint16
	Key  string
	Data interface{}
}

const maxObjSize = 0x10000

type entrys []*entry

func (x entrys) Len() int           { return len(x) }
func (x entrys) Less(i, j int) bool { return x[i].Idx < x[j].Idx }
func (x entrys) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

var default_Consistent *Consistent = nil

func init() {
	default_Consistent = New()
}

type Consistent struct {
	circle       map[uint16]*entry
	members      map[string]int
	members_lock *sync.Mutex

	objs        []*entry
	objs_lock   *sync.RWMutex
	virtualNode int
}

var ErrEmptyCircle = errors.New("empty consistent circle")

func New() *Consistent {
	return &Consistent{
		circle:       make(map[uint16]*entry),
		members:      make(map[string]int),
		members_lock: &sync.Mutex{},

		objs:        nil,
		objs_lock:   &sync.RWMutex{},
		virtualNode: 20,
	}
}

func keyNum(key string, num int) string {
	return key + ":" + fmt.Sprintf("%d", num)
}

func hashKey(key string) uint16 {
	return uint16(crc32.ChecksumIEEE([]byte(key)))
}

func (c *Consistent) AddKey(key string) {
	c.Add(key, nil)
}

func AddKey(key string) {
	default_Consistent.Add(key, nil)
}

func (c *Consistent) Add(key string, data interface{}) {
	c.AddWithWeight(key, data, 1)
}

func Add(key string, data interface{}) {
	default_Consistent.Add(key, data)
}

func (c *Consistent) AddWithWeight(key string, data interface{}, weight int) {
	c.members_lock.Lock()
	defer c.members_lock.Unlock()

	if _, ok := c.members[key]; ok {
		return
	}

	for n := 0; n < c.virtualNode*weight; n++ {
		idx := hashKey(keyNum(key, n))
		c.circle[idx] = &entry{idx, key, data}
	}

	c.members[key] = weight
}

func AddWithWeight(key string, data interface{}, weight int) {
	default_Consistent.AddWithWeight(key, data, weight)
}

func (c *Consistent) Remove(key string) {
	c.members_lock.Lock()
	defer c.members_lock.Unlock()

	weight, ok := c.members[key]
	if !ok {
		return
	}

	for n := 0; n < c.virtualNode*weight; n++ {
		delete(c.circle, hashKey(keyNum(key, n)))
	}
	delete(c.members, key)
}

func (c *Consistent) Remove(key string) {
	default_Consistent.Remove(key)
}

func (c *Consistent) GetKey(name string) (string, error) {
	c.objs_lock.RLock()
	defer c.objs_lock.RUnlock()

	if c.objs == nil {
		return "", ErrEmptyCircle
	}
	idx := hashKey(name)

	return c.objs[idx].Key, nil
}

func GetKey(name string) (string, error) {
	return default_Consistent.GetKey(name)
}

func (c *Consistent) Get(name string) (string, interface{}, error) {
	c.objs_lock.RLock()
	defer c.objs_lock.RUnlock()

	if c.objs == nil {
		return "", nil, ErrEmptyCircle
	}
	idx := hashKey(name)

	return c.objs[idx].Key, c.objs[idx].Data, nil
}

func Get(name string) (string, interface{}, error) {
	default_Consistent.Get(name)
}

func (c *Consistent) Update() {
	c.members_lock.Lock()
	defer c.members_lock.Unlock()
	mb_len := len(c.circle)
	var objs []*entry = nil

	if mb_len > 0 {
		list := make([]*entry, 0, mb_len)
		for _, v := range c.circle {
			list = append(list, v)
		}

		sort.Sort(entrys(list))

		objs = make([]*entry, maxObjSize, maxObjSize)
		begin := 0
		end := 0

		for _, ey := range list {
			end = int(ey.Idx)
			for i := begin; i <= end && i < maxObjSize; i++ {
				objs[i] = ey
			}
			begin = end + 1
		}

		if begin < maxObjSize {
			for i := begin; i < maxObjSize; i++ {
				objs[i] = list[0]
			}
		}
	}

	c.objs_lock.Lock()
	defer c.objs_lock.Unlock()
	c.objs = objs
}

func Update() {
	default_Consistent.Update()
}
