package consistenthash

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sort"
)

// 一个映射空间
type Consistent struct {
	numOfVirtualNode int // 虚拟节点的个数 -> 用户设置 -> 可否系统根据情况配置?
	hashSortedNodes []uint32 // 虚拟节点的有序排列
	circle map[uint32]string // 哈希值空间组成一个虚拟圆环 -- 虚拟节点映射<->真实节点
	nodes map[string]bool // 保存的真实节点
}

func New() *Consistent {
	return &Consistent{
		numOfVirtualNode: 20,
		circle: make(map[uint32]string),
		nodes: make(map[string]bool),
	}
}

// Get the nearby node -- key映射--虚拟节点--哈希节点
func (c *Consistent) Get(key string) (string,error) {
	if len(c.nodes) == 0 {
		return "", errors.New("no host added")
	}
	nearbyIndex := c.searchNearbyIndex(key)
	nearHost := c.circle[c.hashSortedNodes[nearbyIndex]]
	return nearHost, nil
}

// Add the node -- 向Consistent中添加节点
func (c *Consistent) Add(node string) error {
	if _, ok := c.nodes[node]; ok {
		return errors.New("host already existed")
	}
	c.nodes[node] = true
	// add Virtual node -> 每个节点都会映射出一定数量的虚拟节点
	for i := 0; i < c.numOfVirtualNode; i++ {
		virtualKey := getVirtualKey(i, node)
		c.circle[virtualKey] = node
		c.hashSortedNodes = append(c.hashSortedNodes, virtualKey)
	}

	// 将虚拟节点添加进数组并按virtualKey排序 -> 本质上虚拟节点是随机散落在换上的
	sort.Slice(c.hashSortedNodes, func(i, j int) bool {
		return c.hashSortedNodes[i] < c.hashSortedNodes[j]
	})

	return nil
}

// Remove the node
func (c *Consistent) Remove(node string) error {
	if _, ok := c.nodes[node]; ok {
		return errors.New("host is not existed")
	}
	delete(c.nodes, node)

	for i := 0; i < c.numOfVirtualNode; i++ {
		virtualKey := getVirtualKey(i, node)
		delete(c.circle, virtualKey)
	}

	c.refreshHashSlice()

	return nil
}

// ListNodes lists the nodes already existed
func (c *Consistent) ListNodes() []string {
	var nodes []string
	for node := range c.nodes {
		nodes = append(nodes, node)
	}

	return nodes
}

// 求hashKey -> 使用hash算法计算每台机器的位置
func hashKey(host string) uint32 {
	scratch := []byte(host)
	return crc32.ChecksumIEEE(scratch) // 由官方库提供的crc32+host字段求哈希值
}

func getVirtualKey(index int, node string) uint32 {
	return hashKey(fmt.Sprintf("%s#%d", node, index))
}

// searchNearbyIndex -> 找到最近的(虚拟)节点
func (c *Consistent) searchNearbyIndex(key string) int {
	hashKey := hashKey(key)
	targetIndex := sort.Search(len(c.hashSortedNodes), func(i int) bool {
		return c.hashSortedNodes[i] >= hashKey
	})

	if targetIndex >= len(c.hashSortedNodes) {
		targetIndex = 0
	}

	return targetIndex
}

// 删除节点后, 对环上的虚拟节点进行重拍 -- 方便真实节点到虚拟节点的重新散列
func (c *Consistent) refreshHashSlice() {
	c.hashSortedNodes = nil
	for virtualKey := range c.circle {
		c.hashSortedNodes = append(c.hashSortedNodes, virtualKey)
	}
	sort.Slice(c.hashSortedNodes, func(i, j int) bool {
		return c.hashSortedNodes[i] < c.hashSortedNodes[j]
	})
}
