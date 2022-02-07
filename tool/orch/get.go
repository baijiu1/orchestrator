package orch

import (
    "fmt"
    "time"
    "database/sql"
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

    o OrchHost
)

type OrchHost struct {
    HostIP string
    BackendPort string
    OrchHostName string
    OrchWebPort string
}

func GetClusterInwhichOrch(ClusterName string, MetaDsn string) ([]OrchHost ,error) {
    sqlStr := fmt.Sprintf("select oi.orch_physical_ip, oi.orch_backend_db_port, oi.orch_hostname, oi.orch_web_port from meta oi left join orchmapping om on oi.orch_cluster_name - om.orch_cluster_name where om.db_cluster_name = '%v' ", ClusterName)
	var OrchHostList []OrchHost
	DB, err := sql.Open("mysql", MetaDsn)
	if err != nil {
		return
	}
	defer DB.Close()
	rows, err := DB.Query(sqlStr)
	if err != nil {
		return err
	}
	for rows.Next() {
		err3 := rows.Scan(&o.HostIP, &o.BackendPort, &o.OrchWebPort)
		if err3 != nil {
			return OrchHostList, err3
		}
        
        OrchHostList = append(OrchHostList, o)
	}
	if len(OrchHostList) > 3 {
        return OrchHostList, fmt.Errorf("none info")
    } else if len(OrchHostList) == 0 {
        return OrchHostList, fmt.Errorf("none info")
    }
    return OrchHostList, nil
}

