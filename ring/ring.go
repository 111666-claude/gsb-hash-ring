// Package ring 维护一致性哈希环。
package ring

import (
	"hash/fnv"
	"sort"
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
	set    map[string]struct{}
}

// New 建一个空环。
func New() *Ring { return &Ring{set: make(map[string]struct{})} }

// HashOf 算一个字符串的哈希值（环上条目与 key 都用它）。
func HashOf(text string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(text))
	return h.Sum32()
}

func hashOf(text string) uint32 { return HashOf(text) }

// Add 把一个节点加到环上。
func (r *Ring) Add(node string) {
	if _, ok := r.set[node]; ok {
		return
	}
	r.set[node] = struct{}{}
	for i := 0; i < VNodeCount; i++ {
		h := hashOf(node + "#" + strconv.Itoa(i))
		at := sort.Search(len(r.hashes), func(j int) bool { return r.hashes[j] >= h })
		r.hashes = append(r.hashes, 0)
		r.nodes = append(r.nodes, "")
		copy(r.hashes[at+1:], r.hashes[at:])
		copy(r.nodes[at+1:], r.nodes[at:])
		r.hashes[at] = h
		r.nodes[at] = node
	}
}

// Remove 把一个节点从环上摘掉。
func (r *Ring) Remove(node string) {
	if _, ok := r.set[node]; !ok {
		return
	}
	delete(r.set, node)
	keptHashes := r.hashes[:0]
	keptNodes := r.nodes[:0]
	for i, item := range r.nodes {
		if item != node {
			keptHashes = append(keptHashes, r.hashes[i])
			keptNodes = append(keptNodes, item)
		}
	}
	r.hashes = keptHashes
	r.nodes = keptNodes
}

// Size 返回环上的条目数。
func (r *Ring) Size() int { return len(r.nodes) }

// Get 找 key 对应的节点。
func (r *Ring) Get(key string, c *Counter) string {
	if len(r.hashes) == 0 {
		return ""
	}
	want := hashOf(key)
	at := sort.Search(len(r.hashes), func(i int) bool { return r.hashes[i] >= want })
	if at == len(r.hashes) {
		at = 0
	}
	c.Scanned++
	return r.nodes[at]
}
