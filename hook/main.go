package main

import (
	"vfailover/myflag"
	"vfailover/logic"
	"vfailover/journal"
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