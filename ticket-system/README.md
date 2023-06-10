# 問題解決
1. 遇到無法引入 `github.com/metaverse/truss/deftree/googlethirdparty`
```bash
../pb/user/user.pb.go:10:2: missing go.sum entry for module providing package github.com/metaverse/truss/deftree/googlethirdparty (imported by github.com/POABOB/go-microservice/ticket-system/pb/user); to add:
        go get github.com/POABOB/go-microservice/ticket-system/pb/user
```

```bash
go mod tidy
go get -t .
```
2. 遇到 grpc 版本太高
```bash
go: finding module for package github.com/POABOB/go-microservice/ticket-system/user-service/svc
go: finding module for package google.golang.org/grpc/naming
github.com/POABOB/go-microservice/ticket-system/user-service/client/grpc imports
        github.com/POABOB/go-microservice/ticket-system/user-service/svc: module github.com/POABOB/go-microservice@latest found (v0.0.0-20230603091548-e4c793bb5e6b), but does not contain package github.com/POABOB/go-microservice/ticket-system/user-service/svc
github.com/POABOB/go-microservice/ticket-system/pkg/config imports
        github.com/coreos/etcd/clientv3 tested by
        github.com/coreos/etcd/clientv3.test imports
        github.com/coreos/etcd/integration imports
        github.com/coreos/etcd/proxy/grpcproxy imports
        google.golang.org/grpc/naming: module google.golang.org/grpc@latest found (v1.55.0), but does not contain package google.golang.org/grpc/naming
```
```bash
go get google.golang.org/grpc@v1.26.0
```


# 指令

## 產生程式碼

```bash
cd ./pb/{service}
truss *.proto --svcout ../../
# 注入 tag 
protoc-go-inject-tag -input=./{service}.pb.go
```

## 產生 Doc