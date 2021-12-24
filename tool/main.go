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
    }
}