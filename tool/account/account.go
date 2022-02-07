package account

type ReplicaInfo struct {
	Physical_ip string
	Port        string
	ClusterName string
}

type OrchDiscoverPriv struct {
	SelectPriv     string
	ProcessPriv    string
	SuperPriv      string
	ReplSlavePriv  string
	ReplClientPriv string
}

type OrchDiscoverDbPriv struct {
	SelectPriv string
}

type OrchDiscoverDbPseudoPriv struct {
	DropPriv string
}

type SupportGtidMode struct {
	SupportsOracleGtid string
}

var (
	r ReplicaInfo
	p OrchDiscoverPriv
	d OrchDiscoverDbPriv
	e OrchDiscoverDbPseudoPriv
	s SupportGtidMode
	network = "tcp"
	OrchdbName = "orchestrator"

	OrchUser = "orch_discover"
	OrchDbName = "mysql"
	OrchEncodePasswd = "xx"
	OrchPasswd, _ = base64.StdEncoding.DecodeString(OrchEncodePasswd)

	MetaEncodePasswd = "xxxx"
	MetaPasswd, _ = base64.StdEncoding.DecodeString(MetaEncodePasswd)
	MetaUser = "dbmgr"
)

func DetectOrchAccount(ClusterName string, MetaDsn string) (*ReplicaInfo, error) {
	sqlStr := fmt.Sprintf("select physical_ip, port, cluster_name from dbinfo where cluster_name = '%v' and db_role = 'PRIMARY'", ClusterName)
	DB, err := sql.Open("mysql", MetaDsn)
	if err != nil {
		return
	}
	defer DB.Close()
	rows := DB.Query(sqlStr)
	if err := rows.Scan(&r.Physical_ip, &r.Port, &r.ClusterName); err != nil {
		return &r, err
	}
	return &r, nil
}

func OrchDiscoverAccountCheck() error {
	port, _ := strconv.Atoi(r.Port)

	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", OrchUser, OrchPasswd, network, r.Physical_ip, port, OrchDbName)
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	defer DB.Close()
	if err := DB.Ping(); err != nil {
		return err
	}
	return nil
}

func CheckOrchDiscoverPrivilege() error {
	sqlStr := fmt.Sprintf("select Select_priv, Process_prvi,Super_priv,Repl_slave_priv,Repl_client_priv from user where user = '%v' ", OrchUser)
	port, _ := strconv.Atoi(r.Port)
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", MetaUser, MetaPasswd, network, r.Physical_ip, port, OrchDbName)
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	defer DB.Close()
	err := DB.QueryRow(sqlStr).Scan(&p.SelectPriv, &p.ProcessPriv, &p.SuperPriv, &p.ReplSlavePriv, &p.ReplClientPriv)
	if err != nil {
		return err
	}
	if p.SelectPriv == "Y" && p.ProcessPriv == "Y" && p.SuperPriv == "Y" && p.ReplSlavePriv == "Y" && p.ReplClientPriv == "Y" {
		fmt.Println("orch_discover privileges success")
	} else {
		return fmt.Errorf("orch_discover privileges failed")
	}
}

func CheckOrchDiscoverDbMetaPrivilege() {
	sqlStr := fmt.Sprintf("select Select_priv from user where user = '%v' and Db = 'meta' ", OrchUser)
	port, _ := strconv.Atoi(r.Port)
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", MetaUser, MetaPasswd, network, r.Physical_ip, port, OrchDbName)
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	defer DB.Close()
	err := DB.QueryRow(sqlStr).Scan(&d.SelectPriv)
	if err != nil {
		return err
	}
	if d.SelectPriv == "Y" {
		fmt.Println("orch_discover privileges success")
	} else {
		return fmt.Errorf("orch_discover privileges failed")
	}
}

func CheckOrchDiscoverDbPseudoGtidPrivilege( SupportGtid string) error {
	if SupportGtid == "0" {
		sqlStr := fmt.Sprintf("select Drop_priv from db where User = '%v' and Db = '_pseudo_gtid_' ", OrchUser)
		port, _ := strconv.Atoi(r.Port)
		dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", MetaUser, MetaPasswd, network, r.Physical_ip, port, OrchDbName)
		DB, err := sql.Open("mysql", dsn)
		if err != nil {
			return
		}
		defer DB.Close()
		err := DB.QueryRow(sqlStr).Scan(&e.DropPriv)
		if err != nil {
			return err
		}
		if e.DropPriv == "Y" {
			fmt.Println("orch_discover privileges success")
		} else {
			return fmt.Errorf("orch_discover privileges failed")
		}
	} else if SupportGtid == "1" {
		fmt.Println("need't PseudoGtid")
	} else {
		return fmt.Errorf("orch_discover privileges failed")
	}
	return nil
}

func GetClusterInUSEGtidMode(ClusterName string, HostIP string, BackendPort string) (string, error) {
	sqlStr := fmt.Sprintf("select supports_oracle_gtid from database_instance where suggested_cluster_alias = '%v' ", ClusterName)
	port, _ := strconv.Atoi(r.BackendPort)
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", MetaUser, MetaPasswd, network, HostIP, port, OrchDbName)
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
		err3 := rows.Scan(&s.SupportsOracleGtid)
		if err3 != nil {
			return "0", err3
		}
	}
	if s.SupportsOracleGtid == "" {
		return "0", fmt.Errorf(" failed")
	} else if s.SupportsOracleGtid == "1" {
		return "1", nil
	}
	return "0", nil
}