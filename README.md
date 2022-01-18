# orchestrator

#### License

`orchestrator` is free and open sourced under the [Apache 2.0 license](LICENSE).

#### 二次开发

该版本基于《https://github.com/openark/orchestrator》 3.2.4版本开发

新增功能：
```shell
1、 生态建设： 一键precheck与switchover工具
2、 生态建设： VIP漂移工具
3、 改进： 从CMDB元数据中心取值，避免在集群中建表，避免使用sql_log_bin语句造成主备gtid不一致
4、 修复： 优雅切换时，无法选择到正确的datacenter实例，导致切换失败的情况
5、 新增： 日志补齐，在延迟情况下补齐缺失数据
6、 新增： orchestrator-client DataCenter显示
7、 新增： 双主模式切换(该模块时根据我们实际环境单独开发，不影响其他切换函数或切换方式)
```

一、 **一键工具使用方式：** 

构建：

```shell
cd ~/orchestrator/tool/
go build main.go

./main --OpType=[precheck|switchover] --ClusterName=[cluster name] [--force]

说明：
--OpType: 选择是进行预先检查还是进行切换，切换可选--force进行强制切换
--ClusterNmae： 选择要进行检查或切换的集群
```

二、 **vip漂移脚本使用：** 
```shell
cd ~/orchestrator/hook/

go build main.go

useage:

./main --ClusterName={failureCluasterAlias} --DeadStatus={failureType} --OldMaster={failedHost} --OldMasterPort={failedPort} --NewMaster={successorHost} --NewMasterPort={successorPort}

说明：
该脚本是通过外部的orchestrator传入的环境变量来做到vip切换控制的。
{failureCluasterAlias}： 失败的集群名称（别名）。
{failureType}： 探测到失败的类型，一般为deadmaster，还有其他类型，可以看源码logic/topology_recovery.go下的getCheckAndRecoverFunction()函数里面的switch case语句，里面记录了所有的失败类型。return 后面true和false来控制是否进入恢复过程。至于如何判断是否恢复的，可以看我csdn博客或者公众号：收获不止Oracle
{failedHost}： 老主
{failedPort}： 老主的端口
{successorHost}： 新主
{successorPort}： 新主的端口

当然还有很多环境变量，具体可以查看官网文档。
```

三： **以CMDB为中心取值：** 
新增三个配置项，示例配置文件：conf/orchestrator-simple.conf.json：
```shell
MetaDBHost： 该参数意思是到哪个元数据库去取值
MetaDBPort： 对应的端口
MetaDBName： 对应的数据库名称
```

比如，我要去元数据库取对等切换节点，让这两个对等切换节点处于同一个数据中心（datacenter）或取集群别名，或者取vip，那可以这样配置：
```sql
DetectDataCenterQuery: "select dc_vaild from dbinfo where hostname = 'dc_vaild_host_flag' and port = 'dc_vaild_port_flag'"
DetectClusterDomainQuery: "select vip as cluster_domain from dbinfo where hostname = 'dc_vaild_host_flag' and port = 'dc_vaild_port_flag'"
DetectClusterAliasQuery: "select cluster_name as cluster_alias from dbinfo where hostname = 'dc_vaild_host_flag' and port = 'dc_vaild_port_flag'"
```

这里是通过hostname和port字段去确定一个实例的，同理physical_ip和port去确定一个实例也一样。

其中， **_dc_vaild_host_flag_** 和 **_dc_vaild_port_flag_** 是占位符，程序里去替换的东西，不要动。

这里你只需要替换dc_vaild 、dbinfo 、hostname、 port字段就好，改成你们表里自己定义的字段或表明。

字段说明：
```shell
1. dc_vaild： 标识某两个实例处于同一个数据中心的字段，比如A->B->C三个实例，我需要标识A和B处于同一数据中心，那么在表里只需要将这两个实例的dc_vaild设置为Y就好，或者设置为相同的值
2. dbinfo： 表名。也就是取数据的元数据表的表明。
3. hostname： 这个主要是实例的主机名或者IP地址。具体看HostnameResolveMethod设置。
4. port： 顾名思义就是端口字段。
```

把元数据放到配置文件中来做的一点是因为：如果要在提供服务的集群上做的话，就需要在除对等实例外的实例上用sql_log_bin去更新dc_vaild这个字段，是有侵入的。

当然也可以把这份元数据维护在orchestrator自身的那个数据库当中，我这边是有一个总的cmdb，所以放在这里。


四、 **修复优雅切换**

由于我们需要在同一个数据中心的实例相互切换，所以设置了： PreventCrossDataCenterMasterFailover = True，所以在进行切换时，发现选择不到正确的数据中心，导致切换失败

修复如下：

```go
func chooseCandidateReplica(...){
    ...
    priorityMajorVersion, _ := getPriorityMajorVersionForCandidate(replicas)
	priorityBinlogFormat, _ := getPriorityBinlogFormatForCandidate(replicas)
    // 新增getPriorityDataCenterForCandidate函数，返回datacenter字段，具体代码请查看inst/instance_topology.go
    priorityDataCenter, _ := getPriorityDataCenterForCandidate(replicas)

    // 新增对参数：PreventCrossDataCenterMasterFailover 判断
    if config.Config.PreventCrossDataCenterMasterFailover {
        for _, replica := range replicas {
            replica := replica
            if isGenerallyValidAsCandidateReplica(replica) &&
                // 新增数据中心判断 没有改判断，所以切换失败
                IsDataCenterCandidateReplica(priorityDataCenter, replica) &&
                !IsBannedFromBeingCandidateReplica(replica) &&
                !IsSmallerMajorVersion(priorityMajorVersion, replica.MajorVersionString()) &&
                !IsSmallerBinlogFormat(priorityBinlogFormat, replica.Binlog_format) {
                // this is the one
                candidateReplica = replica
                break
            }
	    }
    } else {
        // 还是原来的切换方案
        for _, replica := range replicas {
            replica := replica
            if isGenerallyValidAsCandidateReplica(replica) &&
                !IsBannedFromBeingCandidateReplica(replica) &&
                !IsSmallerMajorVersion(priorityMajorVersion, replica.MajorVersionString()) &&
                !IsSmallerBinlogFormat(priorityBinlogFormat, replica.Binlog_format) {
                // this is the one
                candidateReplica = replica
                break
            }
	    }
    }
    ...
}
```

以上代码修复了在优雅切换中，配置了数据中心参数，但是实例选择不正确的问题，同时也不会影响不配置该参数的情况。该bug已提交至github官方。


五、  **日志补齐** 

日志补齐系统的原理和MHA日志补齐原理相似，都是用到了  **show slave status\G**  里的  **execute_master_position** 这个位点来做的。

只不过orchestrator默认是超过配置的 **ReasonableReplicationLagSeconds** 秒后，切换会直接退出，所以新增了一个配置： **SlaveBinLogEnableMaxLagSeconds** 。

日志补齐系统具体介入时机：

```go
if (secondbehindmaster > ReasonableReplicationLagSeconds && secondbehindmaster < SlaveBinLogEnableMaxLagSeconds) {
    ...
    return instance, true, nil
} else if (secondbehindmaster < ReasonableReplicationLagSeconds) {
    // sql线程等待代码
} else {
    return instance, false, nil
}
```

以上代码就是日志补齐系统的介入时机。也就是延迟大于ReasonableReplicationLagSeconds配置并且小于SlaveBinLogEnableMaxLagSeconds配置。

具体怎么做？
```shell
首先拿到execute_master_position，通过它找到binlog日志文件，拉取到orchestrator主节点上，通过mysqlbinlog解析为sql文件，然后在新主库上应用。
```


六、 **orchestrator-client改进** 

因为我们需要用到datacenter来配置同机房实例来互相切换，所以需要在orchestrator-client中显示datacenter的配置情况，以便切换前做最后检查

改进代码如下：
```go
func (this *Instance) descriptionTokens() {
    ...
    tokens = append(tokens, this,DataCenter)
    ...
}
```

使用方式如下：
```shell
orchestrator-client -c topology-tabulated -a <cluster_name>

可以看到datacenter显示
```



七、 **双主模式切换** 
```shell
useage:

`orchestrator-client -c cascade-master-takeover-auto -i slave_inst:port -d dest_inst:port`

```
