package main

import (
	"github.com/openark/orchestrator/hook/journal"
	"github.com/openark/orchestrator/hook/logic"
	"github.com/openark/orchestrator/hook/myflag"
)

func main() {
	InitCfg, _ := myflag.NewFlagArgs()
	logger := journal.InitLog(InitCfg.ClusterName)
	logic.UmountVip(InitCfg, logger)
	logic.KillMySQLConnection(InitCfg, logger)
	logic.MountVip(InitCfg, logger)
	CheckMatch, _ := logic.PingVIP(InitCfg.CmdVipStat, InitCfg.MaxWaitPing, logger)
	if CheckMatch {
		logger.Println("vip mount success")
	} else {
		logger.Println("There may be some problems with the vip status")
	}
}
