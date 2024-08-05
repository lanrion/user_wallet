# Getting Started

# Run 

## wallet - 钱包服务
`go run cmd/wallet/main.go`
## wb - 回写服务
`go run cmd/wb/main.go`

# Dump
固定dump路径为："datadir/user_wallet_dump/"

# Kafka

本机安装kafka，确保 `127.0.0.1:9092` 服务正常。代码固定死了此brokers地址。

# 架构设计图

doc/dataflow.png

# 使用
启动默认为 0分片的用户
## 查看余额
```
curl http://localhost:9102/user/info/\?user_id\=10
```
## 充值
```
curl http://localhost:9102/user/deposit/\?user_id\=1120\&amount\=3500
```
## 提现
```
curl http://localhost:9102/user/withdraw/\?user_id\=1120\&amount\=1
```
## 转账

```shell
curl http://localhost:9102/user/transfer/\?user_id\=1120\&amount\=1\&to_uid\=20
```
# TODO
* trace
* metric
