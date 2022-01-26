package dealbinlog

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func DealBinLog(mybinlogpath string, masterip string, masterbinlog string, position string) string {
	//get binlog path
	SplitBinlogpath := strings.Split(mybinlogpath, "/")
	List_num := len(SplitBinlogpath) - 1
	NewBinLogPathList := SplitBinlogpath[0:List_num]
	NewBinLogPath := strings.Join(NewBinLogPathList, "/")
	//fmt.Printf("\nThe Binlog Path :%v\n", NewBinLogPath)
	log.Printf("The Binlog Path :%v\n", NewBinLogPath)
	//save binlog
	Getorchleadercmd := "hostname -i "
	orchleader := execMyCmd(Getorchleadercmd)
	//fmt.Printf("Binlog %v will be scp from %v to orch leader %v\n", masterbinlog, masterip, orchleader)
	log.Printf("Binlog %v will be scp from %v to orch leader %v\n", masterbinlog, masterip, orchleader)

	savebinlogpath := "/tmp" + masterip + "_" + masterbinlog + "_" + position
	scp_cmd := "scp mysql@" + masterip + ":" + NewBinLogPath + "/" + masterbinlog + "  " + savebinlogpath
	//fmt.Printf("Save Cmd: %v\n", scp_cmd)
	log.Printf("Save Cmd: %v\n", scp_cmd)
	execMyCmd(scp_cmd)

	//fmt.Printf("Save binlog to %V ok\n", savebinlogpath)
	log.Printf("Save binlog to %V ok\n", savebinlogpath)

	new_parse_binlog := fmt.Sprintf("/tmp/Parsebinlog_%v_%v", masterbinlog, position)
	MysqlBinlogParse := fmt.Sprintf("mysqlbinlog --start-position=%v --stop-nerver %v > %v", position, savebinlogpath, new_parse_binlog)

	//fmt.Printf("MysqlbinlogGet : %v\n", MysqlBinlogParse)
	log.Printf("MysqlbinlogGet : %v\n", MysqlBinlogParse)
	execMyCmd(MysqlBinlogParse)
	//fmt.Printf("Parse binlog ok %v\n", new_parse_binlog)
	log.Printf("Parse binlog ok %v\n", new_parse_binlog)
	return new_parse_binlog
}

func execMyCmd(command string) string {
	out, err := exec.Command("bash", "-c", command).Output()
	checkErr(err)
	return string(out)
}

func checkErr(errinfo error) {
	if errinfo != nil {
		fmt.Printf(errinfo.Error())
	}
}
