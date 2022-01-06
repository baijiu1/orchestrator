# orchestrator

#### License

`orchestrator` is free and open sourced under the [Apache 2.0 license](LICENSE).

#### 二次开发

该版本基于《https://github.com/openark/orchestrator》 3.2.4版本开发

新增日志补齐，双主切换等功能，兼容KP和MHA架构

生态建设： 一键precheck与switchover工具

一键工具使用方式：

构建：

```shell
cd ~/orchestrator/tool/
go build main.go

./main --OpType=[precheck|switchover] --ClusterName=[cluster name] [--force]

说明：
--OpType: 选择是进行预先检查还是进行切换，切换可选--force进行强制切换
--ClusterNmae： 选择要进行检查或切换的集群
```

vip漂移使用：
```shell
cd ~/orchestrator/hook/

go build main.go

useage:

./main --ClusterName={failureCluasterAlias} --DeadStatus={failureType} --OldMaster={failedHost} --OldMasterPort={failedPort} --NewMaster={successorHost} --NewMasterPort={successorPort}

说明：
该脚本是通过外部的orchestrator传入的环境变量来做到vip切换控制的。
{failureCluasterAlias}： 失败的集群名称（别名）。
{failureType}： 探测到失败的类型，一般为deadmaster，还有其他类型，可以看源码logic/topology_recovery.go下的getCheckAndRecoverFunction()函数里面的switch case语句，里面记录了所有的失败类型。
{failedHost}： 老主
{failedPort}： 老主的端口
{successorHost}： 新主
{successorPort}： 新主的端口

当然还有很多环境变量，具体可以查看官网文档。
```

useage:

`orchestrator-client -c cascade-master-takeover-auto -i slave_inst:port -d dest_inst:port`


