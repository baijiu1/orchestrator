package apply

import (
    "fmt"
    "os/exec"
    "encoding/base64"
    "github.com/openark/orchestrator/go/inst"
)

fun ApplyBinLog(parseBinlog string, instanceKey *inst.InstanceKey) {
    execCmd := fmt.Sprintf("mysql -u%s -p%s -h%s -P%s -e 'source %s'", config.Config.MySQLTopologyUser, config.Config.MySQLTopologyPassword, instanceKey.Hostname, instanceKey.Port, parseBinlog)
    execRes := exec.Command(execCmd).Output()
    fmt.Println(execRes)
}