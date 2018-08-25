package mysqldb

import "database/sql"
import (
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/cihub/seelog"
)

func Dial(addr string, userName string, passWord string, dbName string, charset string) *sql.DB{
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s",
		userName, passWord, addr, dbName, charset))
	if err != nil {
		seelog.Errorf("mysql open conn error[%s]", err)
		return nil
	}
	err = db.Ping()
	if err != nil {
		seelog.Errorf("mysql conn ping error[%s]", err)
		return nil
	}
	seelog.Debugf("connect to mysql success")
	return db
}

// return effectednum
func Exec(db *sql.DB, execStr string, args ...interface{}) (int64, error){
	stmtIn, err := db.Prepare(execStr)
	if err != nil {
		seelog.Errorf("Exec Prepare error[%v], execStr:%s", err, execStr)
		return 0, err
	}
	defer stmtIn.Close()
	ret, err := stmtIn.Exec(args...)
	if err != nil {
		seelog.Errorf("Exec error[%v], execStr:%s", err, execStr)
		return 0, err
	}
	return ret.RowsAffected()
}

func QueryRows(db *sql.DB, queryStr string, args ...interface{}) (*[]map[string][]byte, error){
	stmtOut, err := db.Prepare(queryStr)
	if err != nil{
		seelog.Errorf("QueryRows Prepare error[%v], queryStr:%s", err, queryStr)
		return nil, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(args...)
	if err != nil {
		seelog.Errorf("QueryRows Query error[%v]", err)
		return nil, err
	}
	defer rows.Close()
	colums, err := rows.Columns()
	if err != nil {
		seelog.Errorf("QueryRows GetColumns error[%v]", err)
		return nil, err
	}
	values := make([]sql.RawBytes, len(colums))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
 	}

	ret := make([]map[string][]byte, 0, 32)
	for rows.Next(){
		err = rows.Scan(scanArgs...)
		if err != nil {
			seelog.Errorf("QueryRows scan error[%v]", err)
			return nil, err
		}
		var value []byte
		vmap := make(map[string][]byte)
		for i, col := range values{
			if col == nil {
				value = []byte("NULL")
			}else{
				value = make([]byte, len(col))
				copy(value, col)
			}
			vmap[colums[i]] = value
		}
		ret = append(ret, vmap)
	}
	return &ret, nil
}
