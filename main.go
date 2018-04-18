package main

import (
	"fmt"
	"github.com/kirthiprakash/tech-talk-tekion-201804/mysql"
	"github.com/kirthiprakash/tech-talk-tekion-201804/app"
)

func main(){
	db, connErr := mysql.GetConnection()
	if connErr != nil{
		err := fmt.Errorf("error connecting to DB")
		panic(err)
	}
	db.AutoMigrate(&app.Appointment{})
}
