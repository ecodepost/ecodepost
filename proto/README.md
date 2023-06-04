# egoctl
## 1 Requirements
- Go version >= 1.16.

## 2 Installation
```
git clone https://hub.fastgit.org/gotomicro/egoctl.git
go install

# 【必要】格式化proto
buf format -w

# lint proto
buf lint .
```