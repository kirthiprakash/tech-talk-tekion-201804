package main

import (
	"fmt"
	"github.com/kirthiprakash/tech-talk-tekion-201804/app"
	"github.com/kirthiprakash/tech-talk-tekion-201804/mongo"
	"github.com/kirthiprakash/tech-talk-tekion-201804/mysql"
	"math/rand"
	"strconv"
	"time"
)

func main() {

	mysqlDB, connErr := mysql.GetConnection()
	if connErr != nil {
		panic(connErr)
	}

	mysqlDB.AutoMigrate(&app.Appointment{})
	mysqlDB.AutoMigrate(&app.Report{})

	mongoDB := mongo.GetConnection()

	fmt.Print("How may rows should I insert?: ")
	var input string
	fmt.Scanln(&input)
	noOfRows, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	for i := 1; i < noOfRows; i++ {
		apptDate := time.Date(app.Random(2016, 2018), time.Month(app.Random(1, 12)), app.Random(1, 28),
			app.Random(1, 24), app.Random(1, 60), app.Random(1, 60), 0, time.UTC)
		appointmentSource := app.SourceList[app.Random(0, 3)]
		appointmentType := app.AppointmentTypeList[app.Random(0, 3)]

		app.ProduceAppointmentForMongo(mongoDB, apptDate, appointmentSource, appointmentType)
		app.ProduceAppointmentForMysql(mysqlDB, apptDate, appointmentSource, appointmentType)
	}

	end := time.Now()
	fmt.Println("Total execution time: ", end.Sub(start).Minutes(), "mins")

	mongo.CloseConnection()
	mysql.CloseConnection()
}
