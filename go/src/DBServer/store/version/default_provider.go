package version
import (
	"database/sql"
	"fmt"
	"DBServer/db/mysqldb"
	"DBServer/db"
)

type VersionDefaultProvider struct {
	conn *sql.DB
}

func (p *VersionDefaultProvider) Init (){
	p.conn = db.GetDBManager().GetConn()
}

func (p *VersionDefaultProvider) LoadVersionConfigData() (*[]map[string][]byte, error) {
	tblName := "version_config"
	queryStr := fmt.Sprintf("select * from %s", tblName)
	return mysqldb.QueryRows(p.conn, queryStr)
}
