package cache

import (
	"fmt"
	"hash/fnv"
	"sort"
	"sync"
)

type LoadBalancer interface {
	SelectNode(key string) CacheNode
	GetReplicationNodes(key string) []CacheNode
	UpdateNodes(nodes []CacheNode)
}

type consistentHashBalancer struct {
	nodes        []CacheNode
	hashRing     []uint32
	nodeMap      map[uint32]CacheNode
	mu           sync.RWMutex
	virtualNodes int
}

func newLoadBalancer() LoadBalancer {
	return &consistentHashBalancer{
		nodeMap:      make(map[uint32]CacheNode),
		virtualNodes: 100, // Number of virtual nodes per physical node
	}
}

func (b *consistentHashBalancer) SelectNode(key string) CacheNode {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.nodes) == 0 {
		return CacheNode{}
	}

	hash := b.hashKey(key)
	idx := b.findNode(hash)
	return b.nodes[idx]
}

func (b *consistentHashBalancer) GetReplicationNodes(key string) []CacheNode {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.nodes) == 0 {
		return nil
	}

	hash := b.hashKey(key)
	idx := b.findNode(hash)

	// Get next healthy nodes for replication
	result := make([]CacheNode, 0)
	seen := make(map[string]bool)

	for i := 0; i < len(b.nodes) && len(result) < 3; i++ {
		nodeIdx := (idx + i) % len(b.nodes)
		node := b.nodes[nodeIdx]

		if !seen[node.ID] && node.Health {
			result = append(result, node)
			seen[node.ID] = true
		}
	}

	return result
}

func (b *consistentHashBalancer) UpdateNodes(nodes []CacheNode) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.nodes = nodes
	b.rehash()
}

func (b *consistentHashBalancer) rehash() {
	b.hashRing = make([]uint32, 0, len(b.nodes)*b.virtualNodes)
	b.nodeMap = make(map[uint32]CacheNode)

	for _, node := range b.nodes {
		for i := 0; i < b.virtualNodes; i++ {
			hash := b.hashKey(fmt.Sprintf("%s-%d", node.ID, i))
			b.hashRing = append(b.hashRing, hash)
			b.nodeMap[hash] = node
		}
	}

	sort.Slice(b.hashRing, func(i, j int) bool {
		return b.hashRing[i] < b.hashRing[j]
	})
}

func (b *consistentHashBalancer) hashKey(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func (b *consistentHashBalancer) findNode(hash uint32) int {
	idx := sort.Search(len(b.hashRing), func(i int) bool {
		return b.hashRing[i] >= hash
	})

	if idx == len(b.hashRing) {
		idx = 0
	}

	return idx
}
