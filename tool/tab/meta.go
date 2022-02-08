package tab

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type InfoC struct {
	Physical_ip string
}

var (
	k InfoC
)

func GetMetaClusterNameInfo(ClusterName string, MetaDsn string) error {
	MessSelectSQL := fmt.Sprintf("select physical_ip from dbtool_dbinfo where cluster_name ='%v' ", ClusterName)
	DB, err := sql.Open("mysql", MetaDsn)
	if err != nil {
		fmt.Printf("\x1b[%dm open: dbtool conn filed,err: %v \x1b[0m\n", 31, err)
	}
	defer DB.Close()
	rows, err1 := DB.Query(MessSelectSQL)
	if err1 != nil {
		fmt.Printf("\x1b[%dm Query: dbtool_dbinfo not have %v \x1b[0m\n", 31, ClusterName)
	}
	if err3 := rows.Next(); err3 != true {
		for rows.Next() {
			err2 := rows.Scan(&k.Physical_ip)
			if err2 != nil {
				fmt.Printf("\x1b[%dm Check: dbtool_dbinfo --> %v none info exists \x1b[0m\n", 31, ClusterName)
				fmt.Printf("\x1b[%dm Process Now Exit!!!! \x1b[0m\n", 31)
				return err2
			} else {
				fmt.Printf("\x1b[%dm Check: dbtool_dbinfo --> %v exists \x1b[0m\n", 31, ClusterName)

			}
		}
	} else {
		fmt.Printf("\x1b[%dm Check: dbtool_dbinfo --> %v none info exists \x1b[0m\n", 31, ClusterName)
		fmt.Printf("\x1b[%dm Process Now Exit!!!! \x1b[0m\n", 31)
		return fmt.Errorf("Check: dbtool_dbinfo --> %v none info exists.", ClusterName)
	}
	DB.SetConnMaxLifetime(100 * time.Second)
	DB.SetMaxOpenConns(2)
	DB.SetMaxIdleConns(2)
	return nil
}
