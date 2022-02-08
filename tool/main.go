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
	"precheck/account"
	"precheck/acq"
	"precheck/cut"
	"precheck/host"
	"precheck/journal"
	"precheck/lag"
	"precheck/meta"
	"precheck/myflag"
	"precheck/orch"
	"precheck/tab"
	"precheck/topology"
	"tool/access"
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
		OrchHost, err2 := orch.GetClusterInwhichOrch(Init.ClusterName, Meta.MetaDsn)
		if err2 != nil {
			return
		}
		//检查元数据表的的vip和ha_type是否正确
		if !meta.CheckDbtoolVipAndHaType(Meta.MetaDsn, Init.ClusterName) {
			return
		}
		// 检查及联
		_, err3 := topology.DetectReplica(Init.ClusterName, OrchHost[0].HostIP, OrchHost[0].BackendPort)
		if err3 != nil {
			return
		}
		// 检查账号 连通性 权限等
		infolist, err4 := account.DetectOrchAccount(Init.ClusterName, Meta.MetaDsn)
		if err4 != nil {
			return
		}
		err5 := account.OrchDiscoverAccountCheck()
		if err5 != nil {
			return
		}
		err6 := account.CheckOrchDiscoverPrivilege()
		if err6 != nil {
			return
		}
		// 检查是否开启了了gtid
		SupportGtid, err7 := account.GetClusterInUSEGtidMode(Init.ClusterName, OrchHost[0].HostIP, OrchHost[0].BackendPort)
		if err7 != nil {
			return
		}
		err8 := account.CheckOrchDiscoverDbPseudoGtidPrivilege(SupportGtid)
		if err8 != nil {
			return
		}
		// 检查meta.cluster表中的集群名信息是否正确
		err9 := tab.CheckMetaClusterTable(infolist, Init.ClusterName)
		if err9 != nil {
			return
		}
		// 检查元数据表的数据中心标识字段是否正确
		_, _, err10 := tab.CHeckDbtoolVaild2(Init.ClusterName, Meta.MetaDsn)
		if err10 != nil {
			return
		}
		// 检查机器名是否可以ping通
		err11 := host.CheckMachineDBS(Init.ClusterName, Meta.MetaDsn)
		if err11 != nil {
			return
		}
		err12 := access.CheckOrchToClusterSSHIsalive(Init.ClusterName, OrchHost)
		if err12 != nil {
			return
		}
		// 检查延迟
		err13 := lag.CheckSlaveLag(Init.ClusterName, OrchHost[0].HostIP, OrchHost[0].BackendPort)
		if err13 != nil {
			return
		}
		//打印现在的结构
		err14 := topology.PrintTopologyTabulated(Init.ClusterName, OrchHost[0].HostIP, OrchHost[0].BackendPort)
		if err14 != nil {
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
