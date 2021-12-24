package cut

import (
    "fmt"
    "os/exec"
    "precheck/tab"
    "strings"

    _ "github.com/go-sql-driver/mysql"
)

const (
    CliCmd       = "orchestrator-client"
    switchtype   = "graceful-master-takeover-auto"
    cswitchtype  = "cascade-master-takeover-auto"
    dswitchtype  = "force-master-takeover"
    alias        = "alias"
)

type HostName struct {
    HostName string
}

var (
    h      HostName
    hostlist  []*HostName
)

func SwitchOver(ClusterName string, MetaDsn string, OrchHost string, OrchWebPort string, isCoMas bool) {
    if isCoMas {
        fmt.Println("Begin: cascade-master-takeover-auto switchover...")
        // 双主切换
        cmd := fmt.Sprintf("source /etc/profile;source ~/.bash_profile;which %v", CliCmd)
        out, err := exec.Command("bash", "-c", cmd).Output()
        if err != nil {
            fmt.Printf("which failed, err: =====> %v", err)
        }
        c := strings.Replace(string(out), "\n", "", -1)
        switchcmd := fmt.Sprintf("%v -c %v -%v %v", c, cswitchtype, alias, ClusterName)
        SwitchoverCmd := fmt.Sprintf("exporter ORCHESTRTOR_API=\"%v:%v/api\";%v", OrchHost, OrchWebPort, switchcmd)
        fmt.Printf("SwitchOverCmd: %v \n", SwitchoverCmd)
        switchOut, err1 := exec.Command("bash", "-c", SwitchoverCmd).Output()
        if err1 != nil {
            fmt.Printf("switchover failed, err: =====> %v", err)
        }
        fmt.Println(string(switchOut))
        // ack
        ackCmd := fmt.Sprintf("orchestrator-client -c ack-cluster-recoveries -a %v -reason='%v'", ClusterName, ClusterName)
        fmt.Printf("ack command: %v \n", ackCmd)
        ackOut, err1 := exec.Command("bash", "-c", ackCmd).Output()
        fmt.Printf("ackOut: %v \n", string(ackOut))

        //打印拓补结构
        printTopology := fmt.Sprintf("orchestrator-client -c topology-tabulated -a %v ", ClusterName)
        fmt.Printf("print topology command: %v \n", printTopology)
        printOut, err1 := exec.Command("bash", "-c", printTopology).Output()
        fmt.Printf("print topology: %v \n", string(printOut))
    } else {
        // 优雅切换
        
    }
}