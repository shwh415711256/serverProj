package account

import (
	"database/sql"
	"DBServer/db"
	"fmt"
	"DBServer/db/mysqldb"
	"Common/model"
	"github.com/pkg/errors"
	"github.com/cihub/seelog"
)

type AccountDefaultProvider struct {
	conn *sql.DB
}

func (p *AccountDefaultProvider) Init(){
	p.conn = db.GetDBManager().GetConn()
}

func (p *AccountDefaultProvider) LoadAccountInfo(gameid string, openid string) (*[]map[string][]byte, error){
	tblName := fmt.Sprintf("user_info_tbl_%s", gameid)
	queryStr := fmt.Sprintf("select * from %s where openid=?", tblName)
	return mysqldb.QueryRows(p.conn, queryStr, openid)
}

func (p *AccountDefaultProvider) InsertAccountInfo(gameid string, openid string) error {
	tblName := fmt.Sprintf("user_info_tbl_%s", gameid)
	insertStr := fmt.Sprintf("insert into %s(`openid`)values(?)",tblName)
	_, err := mysqldb.Exec(p.conn, insertStr, openid)
	return err
}

func (p *AccountDefaultProvider) LoadWxInfo(gameid string, openid string) (*[]map[string][]byte, error){
	tblName := fmt.Sprintf("wechat_user_tbl_%s", gameid)
	queryStr := fmt.Sprintf("select * from %s where openid=?", tblName)
	return mysqldb.QueryRows(p.conn, queryStr, openid)
}

func (p *AccountDefaultProvider) UpdateWxInfo(gameid string, data *model.WxUserInfo) error{
	tblName := fmt.Sprintf("wechat_user_tbl_%s", gameid)
	insertStr := fmt.Sprintf("insert into %s(`openid`, `nick_name`, `avatar_url`,`gender`,`city`,`province`,`country`, `language`) " +
		"values(?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `nick_name`=?,avatar_url=?,gender=?, city= ?, province=?, country=?, language=?", tblName)
	_, err := mysqldb.Exec(p.conn, insertStr, data.OpenId, data.NickName, data.AvatarUrl, data.Gendar, data.City, data.Province, data.Country, data.Lang,
		data.NickName, data.AvatarUrl, data.Gendar, data.City, data.Province, data.Country, data.Lang)
	return err
}

func (p *AccountDefaultProvider) UpdateAccountValue(gameid string, openid string, args ...interface{}) error{
	l := len(args)
	if l % 2 != 0 {
		seelog.Errorf("%d, %v",l, args)
		return errors.New("args length error")
	}
	tblName := fmt.Sprintf("user_info_tbl_%s", gameid)
	updateStr := fmt.Sprintf("update %s set ", tblName)
	for i := 0; i < l / 2; i ++ {
		updateStr += args[i].(string) + "= ?"
		if i < l /2 -1 {
			updateStr += ","
		}
	}
	updateStr += " where `openid`=?"
	args = append(args, openid)
	_, err := mysqldb.Exec(p.conn, updateStr, args[l/2:]...)
	return err
}

func (p *AccountDefaultProvider) UpdateQDInfo(gameid string, chanid string, num string) error{
	tblName := "qd_login_info"
	insertStr := fmt.Sprintf("insert into %s(`chanid`,`gameid`,`num`)values(?,?,?) ON DUPLICATE KEY UPDATE `num` = ?", tblName)
	_, err := mysqldb.Exec(p.conn, insertStr, chanid, gameid, num, num)
	return err
}