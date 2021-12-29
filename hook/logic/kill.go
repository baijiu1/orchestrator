package logic

import (
    "database/sql"
    "encoding/base64"
    "fmt"
    "log"
    "strconv"
    "vfailover/myflag"

    _ "github.com/go-sql/driver/mysql"
)

type ConnectID struct {
    Id string
}

var (
    myUser = "xxx"
    netWork = "tcp"
    myEncodePasswd = "xxx"
    myPasswd, _ = base64.StdEncoding,DecodeString(myEncodePasswd)
    d ConnectID
)

func KillMySQLConnection(conf *myflag.OrchCfg, logger *log.Logger) {
    // kill connection
    logger.Printf("Being kill old master connection on %v", conf.OldMaster)
    MessSelectSQL := fmt.Sprintf("select id from information_schema.processlist where user not in ('root', 'xxx', 'system user', 'xxx')")
    portInt, _ := strconv.Atoi(conf.OldMasterPort)
    dsn := fmt.Sprintf("")
}