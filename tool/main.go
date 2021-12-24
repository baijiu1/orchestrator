package main

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
            err1 := GetMetaClusterNameInfo(Init.ClusterName, Meta.MetaDsn)
            if err1 != nil {
                return
            }
            OrchHist, err2 := 
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