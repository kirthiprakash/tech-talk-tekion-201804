package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kirthiprakash/tech-talk-tekion-201804/app"
)

var _db *gorm.DB

func init() {
	username := app.MysqlUsername
	password := app.MysqlPassword
	host := app.MysqlHost
	databaseName := app.MysqlDBName
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, host, databaseName)
	var openErr error
	_db, openErr = gorm.Open("mysql", dsn)
	if openErr != nil {
		err := fmt.Errorf("error initializing mysql DB connection %s; err: %v", host, openErr)
		panic(err.Error())
	}
	//defer _db.Close()

}

func GetConnection() (*gorm.DB, error) {
	if _db == nil {
		err := fmt.Errorf("GetConnection: mysql connection is not initialised")
		return nil, err
	}
	return _db, nil
}

func CloseConnection() {
	_db.Close()
}
