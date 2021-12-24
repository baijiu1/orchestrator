package myflag

import (
    "encoding/base64"
    "flag"
    "fmt"
    "os"
    "strings"
)

type Args struct {
    OpType  string
    ClusterName string
    ForceFailOver string
}

type MetaDBInfo struct {
    MetaUser string
    MetaPasswd []byte
    MetaDBName string
    MetaHost string
    MetaPort int
    MetaNetWork string
    MetaEncodePasswd string
    MetaDsn string
}

type OrcheMapping struct {
    Id string
    DbClusterName string
}

var (
    a Args
    m MetaDBInfo
    o OrcheMapping
)

const (
    ver = "-*-==================================== orch_ha version 1.0 ====================================-*-"
    exam = `Example:
    
./orch_ha --OpType=[precheck|switchover] --ClusterName=[cluster name] [--force]

`
)

func NewFlag() (*Args, error) {
    flag.StringVar(&a.OpType, "OpType", "none", "operation type(precheck, switchover). default none")
    flag.StringVar(&a.ClusterName, "ClusterName", "none", "input cluster name. default none")
    forceFailover := flag.Bool("force", false, "force failover.")
    flag.Parse()
    if a.OpType == "none" || a.ClusterName == "none" {
        fmt.Printf("\n")
        fmt.Println(ver)
        fmt.Printf("\n")
        fmt.PrintDefaults()
        fmt.Printf("\n")
        fmt.Println(exam)
    }
    return &Args{OpType: a.OpType, ClusterName: a.ClusterName, ForceFailOver: *forceFailover}, nil
}