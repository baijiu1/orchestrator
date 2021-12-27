package account

import (
    "fmt"
    "encoding/base64"
    "database/sql"
    "strconv"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

type ReplicaInfo struct {
    Physical_ip string
    Port string
    ClusterName string
}

type OrchDiscoverPriv struct {
    SelectPriv string
    ProcessPriv string
    SuperPriv string
    ReplSlavePriv string
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