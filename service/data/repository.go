package data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
	"fmt"
)

type DataRepository interface {
	Begin() (*sql.Tx, error)
	query(sql string)([]map[string]interface{},error)
	execWithTx(sql string,tx *sql.Tx)(int64,int64, error)
}

type DefatultDataRepository struct {
	DB *sql.DB
}

func (repo *DefatultDataRepository)Begin()(*sql.Tx, error){
	return repo.DB.Begin()
}

func (repo *DefatultDataRepository)execWithTx(sql string,tx *sql.Tx)(int64,int64, error){
	log.Println(sql)
	res,err:=tx.Exec(sql)
	if err!=nil {
		log.Println(err)
		return 0,0,err
	}

	rowCount,err:=res.RowsAffected()
	if err!=nil {
		log.Println(err)
		return 0,0,err 
	}

	//获取最后插入数据的ID	
	id,err:=res.LastInsertId()
	if err!=nil {
		log.Println(err)
		return 0,0,err 
	}
		
	return id,rowCount,nil
}

func (repo *DefatultDataRepository)rowsToMap(rows *sql.Rows)([]map[string]interface{},error){
	cols,_:=rows.Columns()
	columns:=make([]interface{},len(cols))
	colPointers:=make([]interface{},len(cols))
	for i,_:=range columns {
		colPointers[i] = &columns[i]
	}

	var list []map[string]interface{}
	for rows.Next() {
		err:= rows.Scan(colPointers...)
		if err != nil {
			log.Println(err)
			return nil,err
		}
		row:=make(map[string]interface{})
		for i,colName :=range cols {
			val:=colPointers[i].(*interface{})
			switch (*val).(type) {
			case []byte:
				row[colName]=string((*val).([]byte))
			default:
				row[colName]=*val
			} 
		}
		list=append(list,row)
	}
	return list,nil
}

func (repo *DefatultDataRepository)query(sql string)([]map[string]interface{},error){
	log.Println(sql)
	rows, err := repo.DB.Query(sql)
	if err != nil {
		log.Println(err)
		log.Println(sql)
		return nil,err
	}
	defer rows.Close()
	//结果转换为map
	return repo.rowsToMap(rows)
}

func (repo *DefatultDataRepository)Connect(
	server,user,password,dbName string,
	connMaxLifetime,maxOpenConns,maxIdleConns int,tls string){
	// Capture connection properties.
    /*cfg := mysql.Config{
        User:   user,
        Passwd: password,
        Net:    "tcp",
        Addr:   server,
        DBName: dbName,
		AllowNativePasswords:true,
    }*/
    // Get a database handle.
	dsn:=fmt.Sprintf("%s:%s@tcp(%s)/%s?allowNativePasswords=true&tls=%s",user,password,server,dbName,tls)
	log.Println("connect to mysql server "+dsn)
    var err error
    repo.DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Println(err)
    }

    pingErr := repo.DB.Ping()
    if pingErr != nil {
        log.Println(pingErr)
    }
		repo.DB.SetConnMaxLifetime(time.Minute * time.Duration(connMaxLifetime))
		repo.DB.SetMaxOpenConns(maxOpenConns)
		repo.DB.SetMaxIdleConns(maxIdleConns)
    log.Println("connect to mysql server "+server)
}

