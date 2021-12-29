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
    dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/information_schema", myUser, myPasswd, netWork, conf.OldMaster, portInt)
    DB, err := sql.Open("mysql", dsn)
    if err != nil {
        return
    }
    if err4 := DB.Ping(); err4 != nil {
        logger.Printf("connot connect old master, maybe shutdown or node is down")
    } else {
        defer DB.Close()
        rows, _ := DB.Query(MessSelectSQL)
        for rows.Next() {
            err2 := rows.Scan(&d.Id)
            if err2 != nil {
                return
            }
            // begin kill
            killsql := fmt.Sprintf("kill %v", d.Id)
            _, err3 := DB.Exec(killsql)
            if err3 != nil {
                logger.Printf("kill failed. error: %v \n", err3)
            }
        }
    }
}