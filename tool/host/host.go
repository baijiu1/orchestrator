package host

import (
    "fmt"
    "database/sql"
    "net"

    "time"

    _ "github.com/go-sql-driver/mysql"
)

type HostName struct {
    HostName string
}

var (
    h HostName
    hostlist []HostName
)

func CheckMachineDNS(ClusterName string, MetaDsn string) error {
    sqlStr := fmt.Sprintf("select hostname from dbinfo where cluster_name = '%v' and hostname not like '%%sz%%' and physical_ip not like '%%192.168%%'", ClusterName)

	// dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", MetaUser, MetaPasswd, network, r.Physical_ip, port, OrchDbName)
	DB, err := sql.Open("mysql", MetaDsn)
	if err != nil {
		return
	}
	defer DB.Close()
	rows, _ := DB.Query(sqlStr)
	for rows.Next() {
		err3 := rows.Scan(&h.HostName)
		if err3 != nil {
			return err3
		}
        hostlist = append(hostlist, h)
	}
	errFlag := 0
    for _, dns := range hostlist {
        _, err := net.LookupHost(dns.HostName)
        if err != nil {
            errFlag++
            return fmt.Errorf("cannot ping")
        }
    }
    if errFlag == 0 {
        fmt.println("success")
    } else {
        return fmt.Errorf("cannot ping")
    }
    return nil
}