# go-cfglayer 架构说明

## 核心类型

- `Merger`：层栈生命周期与合并入口（**非** diffpack `Packer`、**非** task10 `Engine`）
- `Layer`：单层键值 + ID / Priority / Source
- `LayerStore`：栈序存储（internal/layerstore）

## 与 go-diffpack 差异

| 维度 | diffpack | cfglayer |
|------|----------|----------|
| 主数据 | 二进制 op 流 + Bundle | 字符串键值层栈 |
| 核心操作 | BuildDelta / ApplyDelta | PushLayer / MergeStack / Resolve |
| 溯源 | Hunk 块图 | ExplainKey 赋值链 |
| 守护进程 | diffpackd 8240 | cfglayerd 8241 |

## 包结构

```text
merger.go · layer_admin.go · resolve.go · explain.go · export.go · import.go
internal/merge · layerstore · path · parse · provenance · validate · codec · audit
cmd/cfglayerd · web/
```
