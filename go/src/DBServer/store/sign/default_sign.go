package sign

import (
	"database/sql"
	"DBServer/db"
	"fmt"
	"DBServer/db/mysqldb"
)

type SignDefaultProvider struct {
	conn *sql.DB
}

func (p *SignDefaultProvider) Init(){
	p.conn = db.GetDBManager().GetConn()
}

func (p *SignDefaultProvider) LoadSignRewardConfig() (*[]map[string][]byte, error){
	tblName := "sign_reward_config"
	queryStr := fmt.Sprintf("select * from %s", tblName)
	return mysqldb.QueryRows(p.conn, queryStr)
}

func (p *SignDefaultProvider) LoadSignInfo(gameId string, openId string, signType int) (*[]map[string][]byte, error){
	tblName := "sign_info_" + gameId
	queryStr := fmt.Sprintf("select * from %s where `openid`=? and sign_type=?", tblName)
	return mysqldb.QueryRows(p.conn, queryStr, openId, signType)
}