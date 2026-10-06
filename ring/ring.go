// Package ring 维护一致性哈希环。
package ring

import (
	"hash/fnv"
	"strconv"
)

// VNodeCount 是每个节点的虚拟节点数。
const VNodeCount = 8

// Counter 记录本轮扫描过的环上条目数。
type Counter struct {
	Scanned int
}

// Ring 是一致性哈希环。
type Ring struct {
	hashes []uint32
	nodes  []string
}

// New 建一个空环。
func New() *Ring { return &Ring{} }

// HashOf 算一个字符串的哈希值（环上条目与 key 都用它）。
func HashOf(text string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(text))
	return h.Sum32()
}

func hashOf(text string) uint32 { return HashOf(text) }

// Add 把一个节点加到环上。
func (r *Ring) Add(node string) {
	for i := 0; i < VNodeCount; i++ {
		r.hashes = append(r.hashes, hashOf(node+"#"+strconv.Itoa(i)))
		r.nodes = append(r.nodes, node)
	}
}

// Remove 把一个节点从环上摘掉。
func (r *Ring) Remove(node string) {
	for i, item := range r.nodes {
		if item == node {
			r.hashes = append(r.hashes[:i], r.hashes[i+1:]...)
			r.nodes = append(r.nodes[:i], r.nodes[i+1:]...)
			return
		}
	}
}

// Size 返回环上的条目数。
func (r *Ring) Size() int { return len(r.nodes) }

type pair struct {
	h    uint32
	node string
}

// Get 找 key 对应的节点。
func (r *Ring) Get(key string, c *Counter) string {
	want := hashOf(key)
	pairs := make([]pair, len(r.nodes))
	for i := range r.nodes {
		c.Scanned++
		pairs[i] = pair{r.hashes[i], r.nodes[i]}
	}
	_ = pairs
	for i, item := range r.hashes {
		if item >= want {
			return r.nodes[i]
		}
	}
	return ""
}
