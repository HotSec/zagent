# zagent

一个精简的agent代码，方便二次开发各种分布式agent组件。

## 安装

环境信息

- Ubuntu 24.04 LTS (6.8.0-31-generic)
- go1.22.2

```bash
git clone https://github.com/HotSec/zagent.git
cd zagent
make
```

## 使用

```bash
./build/zagent -svc install
./build/zagent -svc start
```

## 亮点功能

- [ ] 开箱即用，具备agent所需要的基本功能，日志分割轮转、系统服务注册
- [ ] 多种通信方式可选（http/grpc/socket/redis/kafka/nsq等）
- [ ] 自我升级
- [ ] 资源限制
- [ ] 信息上报
- [ ] 任务执行

## Thank

> github.com/kardianos/service
> gopkg.in/natefinch/lumberjack.v2 

## Change Log

### v0.0.1 

- 初始化项目
- 添加日志记录
- 添加系统服务注册

