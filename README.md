# hash-ring

一致性哈希环：把节点编成环，再按 key 找节点。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/hashring --sample=dup
go run ./cmd/hashring --sample=remove
go run ./cmd/hashring --sample=wrap
go run ./cmd/hashring --sample=work
```

## 口径

- **虚拟节点**：每个节点按 `节点名#序号` 打散成 8 个环上条目。
- **去重**：同一个节点名只编一份虚拟节点。
- **顺序**：环上的条目按哈希值升序排列。
- **回绕**：key 的哈希超过环上最大值时，落到环上第一个条目。
- **摘除**：摘掉节点要把它所有虚拟节点都从环上拿掉。

## 不变量

- 环上条目数等于节点数乘 8。
- 非空环对任何 key 都能返回一个节点。
- `scanned` 不随 key 数乘环上条目数放大：两千个 key 的 `scanned` 不超过 9000。

## 输出契约

每个场景打印一行 JSON；`--sample=work` 只打印 `{"vnodes": N, "scanned": N}`。
