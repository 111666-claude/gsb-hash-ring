// Command hashring 是哈希环的场景入口。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"example.com/hash-ring/ring"
)

func emit(v any) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func main() {
	sample := flag.String("sample", "", "dup|remove|wrap|work")
	flag.Parse()

	c := &ring.Counter{}
	switch *sample {
	case "dup":
		r := ring.New()
		r.Add("n1")
		r.Add("n1")
		r.Add("n2")
		emit(map[string]int{"size": r.Size()})
	case "remove":
		r := ring.New()
		r.Add("n1")
		r.Add("n2")
		r.Remove("n1")
		emit(map[string]int{"size": r.Size()})
	case "wrap":
		r := ring.New()
		r.Add("n1")
		max := uint32(0)
		for i := 0; i < ring.VNodeCount; i++ {
			if h := ring.HashOf(fmt.Sprintf("n1#%d", i)); h > max {
				max = h
			}
		}
		node := ""
		for i := 0; i < 2000000; i++ {
			key := fmt.Sprintf("probe-%d", i)
			if ring.HashOf(key) > max {
				node = r.Get(key, c)
				break
			}
		}
		emit(map[string]any{"node": node, "expected": "n1"})
	case "work":
		r := ring.New()
		for i := 0; i < 1200; i++ {
			r.Add(fmt.Sprintf("node-%04d", i))
		}
		for i := 0; i < 2000; i++ {
			r.Get(fmt.Sprintf("key-%04d", i), c)
		}
		emit(map[string]int{"vnodes": r.Size(), "scanned": c.Scanned})
	default:
		fmt.Fprintln(os.Stderr, "未知场景："+*sample)
		os.Exit(2)
	}
}
