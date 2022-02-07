package journal

import (
    "fmt"
    "time"
    "database/sql"

    _ "github.con/go-sql-driver/mysql"
)

func WritePrecheckTable(MetaDsn string, ClusterName string) error {
    sqlStr := fmt.Sprintf("insert into orch_precheck(cluster_name) values(?) ")
    
    DB, err := sql.Open("mysql", MetaDsn)
    if err != nil {
        return
    }
    defer DB.Close()
    stmt, _ := DB.Prepare(sqlStr)
    _, err1 := stmt.Exec(ClusterName)
    if err1 != nil {
        return fmt.Errorf("insert failed")
    }
    return nil
}

func WritePrecheckStepTable(MetaDsn string, ClusterName string, step string, messages string, status string) error {
    sqlStr := fmt.Sprintf("insert into orch_precheck_step(cluster_name, step, messages, status) values(?, ?, ?, ?) ")
    
    DB, err := sql.Open("mysql", MetaDsn)
    if err != nil {
        return
    }
    defer DB.Close()
    stmt, _ := DB.Prepare(sqlStr)
    _, err1 := stmt.Exec(ClusterName, step, messages, status)
    if err1 != nil {
        return fmt.Errorf("insert failed")
    }
    return nil
}