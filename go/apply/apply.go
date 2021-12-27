package apply

import (
    "fmt"
    "os/exec"
    "encoding/base64"
    "github.com/openark/orchestrator/go/inst"
)

var (
    myEncodePasswd = "xxxx"
    myPasswd, _ = base64.StdEndoding.DecodeString(myEncodePasswd
)

fun ApplyBinLog(parseBinlog string, instanceKey *inst.InstanceKey) {
    execCmd := fmt.Sprintf("mysql -uroot -p%s -h%s -P%s -e '%s'", myPasswd, instanceKey.Hostname, instanceKey.Port, parseBinlog)
    execRes := exec.Command(execCmd).Output()
    fmt.Println(execRes)
}