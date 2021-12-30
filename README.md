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

useage:

`orchestrator-client -c cascade-master-takeover-auto -i slave_inst:port -d dest_inst:port`


