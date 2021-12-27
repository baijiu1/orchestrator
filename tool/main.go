package main

/*
########################################################
##@File      :   orch_ha
# Author     :   baijie
#
# Version History:
#
#            v1.0. 2021/05/13 10:23:18 baijie Start writing, the first verion.
#######################################################

*/

import (
    "fmt"
    "precheck/access"
    "precheck/account"
    "precheck/cut"
    "precheck/host"
    "precheck/journal"
    "precheck/lag"
    "precheck/myflag"
    "precheck/orch"
    "precheck/tab"
    "precheck/topology"
    "precheck/meta"
    "precheck/acq"
)

func main() {
    Init, err := myflag.Newflag()
    Meta, _ := myflag.NewMetaInfo()
    if err != nil {
        fmt.Printf("初始化失败")
    }
    switch Init.OpType {
        case "precheck":
            // 检查给定集群是否在dbinfo元数据表中
            err1 := GetMetaClusterNameInfo(Init.ClusterName, Meta.MetaDsn)
            if err1 != nil {
                return
            }
            // 检查给定集群属于哪个orchestrator集群
            OrchHost, err2 := orch.GetClusterInwhichOrch(Init.ClusterName,Meta.MetaDsn)
            if err2 != nil {
                return
            }
            //检查元数据表的的vip和ha_type是否正确
            if !meta.CheckDbtoolVipAndHaType(Meta.MetaDsn, Init.ClusterName) {
                return
            }
            // 检查及联
            _, err3 := topology.DetectReplica(Init.ClusterName,OrchHost[0].HostIP, OrchHost[0].BackendPort)
            if err3 != nil {
                return
            }
            // 检查账号 连通性 权限等
            infolist, err4 := account.DetectOrchAccount(Init.ClusterName,Meta.MetaDsn)
            if err4 != nil {
                return
            }
            err5 := access.OrchDiscoverAccountCheck()
            if err5 != nil {
                return
            }
            

        case "switchover":
            OrchHist, _ := orch.GetClusterInwhichOrch(Init.ClusterName, Meta.MetaDsn)
            isCoMas, _ := topology.DetectReplica(Init.ClusterName, OrchHost[0].HostIP, OrchHost[0].BackendPort)
            if Init.ForceFailOver {
                cut.FailOver(Init.ClusterName, Meta.MetaDsn, OrchHost[0].OrchHostName, OrchHost[0].OrchWebPort, isCoMas)
            } else {
                cut.SwitchOver(Init.ClusterName, Meta.MetaDsn, OrchHost[0].OrchHostName, OrchHost[0].OrchWebPor, isCoMas)
            }
        default:
            acq.ByDefaultAction()
    }
}