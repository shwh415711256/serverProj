package game

import (
	"database/sql"
	"DBServer/db"
	"fmt"
	"DBServer/db/mysqldb"
)

type GameDefaultProvider struct{
	conn *sql.DB
}

func (p *GameDefaultProvider)Init () {
	p.conn = db.GetDBManager().GetConn()
}

func (p *GameDefaultProvider) LoadOneMatchHis(gameid string, openid string) (*[]map[string][]byte, error){
	tblName := "game_his_" + gameid
	queryStr := fmt.Sprintf("select *from %s where `openid` <> '%s' limit 10", tblName, openid)
	return mysqldb.QueryRows(p.conn, queryStr)
}
func (p *GameDefaultProvider) InsertOneMatchHis(gameid string, openid string, data string) error{
	tblName := "game_his_" + gameid
	insertStr := fmt.Sprintf("insert into %s (`openid`,`hisdata`) value(?,?)", tblName)
	_, err := mysqldb.Exec(p.conn, insertStr, openid, data)
	return err
}
