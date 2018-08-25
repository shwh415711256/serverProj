package match

import (
	"database/sql"
	"DBServer/db"
	"DBServer/db/mysqldb"
	"fmt"
)

type MatchDefaultProvider struct {
	conn *sql.DB
}

func (p *MatchDefaultProvider) Init (){
	p.conn = db.GetDBManager().GetConn()
}

func (p *MatchDefaultProvider) LoadMatchConfigData() (*[]map[string][]byte, error) {
	tblName := "match_config"
	queryStr := fmt.Sprintf("select * from %s", tblName)
	return mysqldb.QueryRows(p.conn, queryStr)
}
