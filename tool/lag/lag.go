package lag

import (
    "fmt"
    "time"
    "database/sql"
    "strconv"
    "encoding/base64"

    _ "github.con/go-sql-driver/mysql"
)

var (
    myUser = "dbmgxx"
    myPort = 3306
    myHost = "192003"
    netWork = "tcp"
    dbName = "meta"
    myEncodePasswd = "xxx"
    myPasswd, _ = base64.StdEncoding.DecodeString(myEncodePasswd)

    OrchdbName = "orchestrator"

    HostDBName = "mysql"
    ReplicationLagSeconds = 60

    s SlaveIP
    SlaveIPList []SlaveIP
    SlaveStatusQuerys = [1]string{"SHOW SLAVE STATUS"}
    l LagColumn
    LagColumnList []int
)

type SlaveIP struct {
    Physical_ip string
    Port string
}

type LagColumn struct {
    ClusterAlias string
    SecondsBehindMaster string
}

func CheckSlaveLag(ClusterName string, HostIP string, BackendPort string) error {
    sqlStr := fmt.Sprintf("select suggested_cluster_alias, ifnull(slave_lag_seconds, '-1') from database_instance where suggested_cluster_alias = '%v' ", ClusterName)
	port, _ := strconv.Atoi(r.BackendPort)
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", myUser, myPasswd, network, HostIP, port, OrchDbName)
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	defer DB.Close()
	rows, err := DB.Query(sqlStr)
	if err != nil {
		return err
	}
	for rows.Next() {
		err3 := rows.Scan(&l.ClusterAlias, &l.SecondsBehindMaster)
		if err3 != nil {
			return "0", err3
		}
        SQLLag, _ := strconv.Atoi(l.SecondsBehindMaster)
        LagColumnList = append(LagColumnList, SQLLag)
        if SQLLag >= 60 {
            return fmt.Errorf(" failed")
        } else {
            fmt.Println("lag < 60")
        }
	}
	lagFlag := 0
    for _, lag := range LagColumnList {
        if lag == -1 {
            lagFlag++
        }
    }
    if lagFlag >= 2 {
        return fmt.Errorf(" failed")
    } else {
        fmt.Println("no lag")
    }
    return nil
}

