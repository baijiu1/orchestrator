package tab

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
	"tool/account"
)

type MetaCluster struct {
	Anchor        string
	ClusterName   string
	ClusterDomain string
	DCVaild       string
}
type ClusterALLInfo struct {
	Physical_ip string
	Port        string
}

type DCVaild struct {
	DCVaild string
}

var (
	myUser         = "dbmgr"
	network        = "tcp"
	dbName         = "meta"
	myEncodePasswd = "xxxxx"
	myPasswd, _    = base64.StdEncoding.DecodeString(myEncodePasswd)

	m       MetaCluster
	c       ClusterALLInfo
	ALLInfo []ClusterALLInfo
	d       DCVaild
	dcinfo  []DCVaild
)

func checkMetaClusterTable(infolist *account.ReplicaInfo, ClusterName string) error {
	MessSelectSQL := fmt.Sprintf("select anchor,cluster_name,cluster_domain from cluster ")
	port, _ := strconv.Atoi(infolist.Port)
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", myUser, myPasswd, network, infolist.Physical_ip, port, dbName)
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("\x1b[%dm open: dbtool conn filed,err: %v \x1b[0m\n", 31, err)
	}
	defer DB.Close()
	rows, err1 := DB.Query(MessSelectSQL)
	if err1 != nil {
		fmt.Printf("\x1b[%dm check: %v ==> meta.cluster is not exists. err:%v \x1b[0m\n", 31, infolist.Physical_ip, err)
		return err1
	}
	for rows.Next() {
		err := rows.Scan(&m.Anchor, &m.ClusterName, &m.ClusterDomain)
		if err != nil {
			fmt.Printf("Scan Failed , err :%v", err)
		}
	}
	DB.SetConnMaxLifetime(100 * time.Second)
	DB.SetMaxOpenConns(2)
	DB.SetMaxIdleConns(2)

	if m.ClusterName == ClusterName {
		fmt.Printf("\x1b[%dm Check: %v  --> meta.cluster.ClusterName == %v \x1b[0m\n", 32, infolist.Physical_ip, ClusterName)
		fmt.Println("-------------------------------------------------------------------------")
	} else {
		fmt.Printf("\x1b[%dm Check: %v  --> meta.cluster.ClusterName != %v \x1b[0m\n", 32, infolist.Physical_ip, ClusterName)
		fmt.Printf("\x1b[%dm Process Now Exit!!!! \x1b[0m\n", 31)
		return fmt.Errorf("Check:%v --> meta.cluster.ClusterName failed.", ClusterName)
	}

	return nil
}
