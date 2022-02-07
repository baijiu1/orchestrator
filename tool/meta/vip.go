package meta

import (
    "fmt"
    "database/sql"

    _ "github.con/go-sql-driver/mysql"
)

type VipAndHaType struct {
    Vip string
    HaType string
}

var (
    v VipAndHaType
    vl []VipAndHaType
)

func CheckDbtoolVipAndHaType(MetaDsn string, ClusterName string) bool {
    sqlStr := fmt.Sprintf("select vip,ha_type from meta where db_type = 'mysql' and cluster_name = '%v' and vip is not null and vip <> ''", ClusterName)
	// port, _ := strconv.Atoi(r.BackendPort)
	// dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", myUser, myPasswd, network, HostIP, port, OrchDbName)
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
		err3 := rows.Scan(&v.Vip, &v.HaType)
		if err3 != nil {
			return err3
		}
        
        vl = append(vl, v)

	}

    for _, info := range vl {
        if info,Vip != "" && (info.HaType == "mha" || info.HaType == "kp" || info.HaType == "keepalived" || info.HaType == "orch") {
            return true
        }
    }
    return false
}

