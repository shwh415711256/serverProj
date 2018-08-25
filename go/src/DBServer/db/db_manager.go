package db

import (
	"database/sql"
	"github.com/cihub/seelog"
	"DBServer/conf"
	"DBServer/db/mysqldb"
	"DBServer/db/redisdb"
)

type DbManager struct {
	Addr string
	UserName string
	PassWd string
	DbName string
	MaxOpenConns int
	MaxIdleConns int

	conn *sql.DB
}

var (
	dbManager *DbManager
)

func Init(){
	dbManager = &DbManager{}
	dbManager.Connect(conf.ServerData.DBAddr, conf.ServerData.DBUserName, conf.ServerData.DBPassword,
		conf.ServerData.DBName, conf.ServerData.MaxOpenConns, conf.ServerData.MaxIdleConns)
	redisdb.Init()
}

func GetDBManager() *DbManager{
	return dbManager
}

func (m *DbManager) GetConn() *sql.DB {
	return m.conn
}

func (m *DbManager) Connect(addr string, userName string, passWd string, dbName string, maxOpenConns int, maxIdleConns int) *sql.DB{
	if m.conn != nil {
		return m.conn
	}
	m.Addr = addr
	m.UserName = userName
	m.PassWd = passWd
	m.DbName = dbName
	m.MaxOpenConns = maxOpenConns
	m.MaxIdleConns = maxIdleConns

	conn := mysqldb.Dial(addr, userName, passWd, dbName, "utf8")
	if conn == nil {
		seelog.Errorf("connet to mysql err")
		return nil
	}
	seelog.Debugf("mysql connet success")
	conn.SetMaxOpenConns(maxOpenConns)
	conn.SetMaxIdleConns(maxIdleConns)

	m.conn = conn
	return m.conn
}