package ring

import "testing"

func TestAddGivesVNodes(t *testing.T) {
	r := New()
	r.Add("n1")
	if r.Size() != VNodeCount {
		t.Fatalf("一个节点应有 %d 个虚拟节点，得到 %d", VNodeCount, r.Size())
	}
}

func TestRemoveReducesSize(t *testing.T) {
	r := New()
	r.Add("n1")
	r.Add("n2")
	before := r.Size()
	r.Remove("n1")
	if r.Size() >= before {
		t.Fatalf("摘掉后条目应减少，得到 %d（原 %d）", r.Size(), before)
	}
}

func TestEmptyRingReturnsEmpty(t *testing.T) {
	r := New()
	if got := r.Get("k", &Counter{}); got != "" {
		t.Fatalf("空环应返回空串，得到 %q", got)
	}
}

func TestThreeNodesGiveThreeTimesVNodes(t *testing.T) {
	r := New()
	r.Add("a")
	r.Add("b")
	r.Add("c")
	if r.Size() != 3*VNodeCount {
		t.Fatalf("三个节点应有 %d 个虚拟节点，得到 %d", 3*VNodeCount, r.Size())
	}
}

func TestSizeCountsEntries(t *testing.T) {
	r := New()
	r.Add("a")
	r.Add("b")
	if r.Size() != 2*VNodeCount {
		t.Fatalf("两个节点应有 %d 个虚拟节点，得到 %d", 2*VNodeCount, r.Size())
	}
}
