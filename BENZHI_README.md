# go-cfglayer

基于 Go 实现的分层配置合并库组件，配套 cfglayerd 管理页，完成配置层压栈、优先级合并、键值溯源与冲突解释。

## Build / Test

```text
go build ./...
go test ./... -count=1
go run ./cmd/cfglayerd -addr :8241
```
