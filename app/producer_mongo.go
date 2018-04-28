package app

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func ProduceAppointmentForMongo(db *mgo.Session, apptDate time.Time, appointmentSource, appointmentType string) {
	appt := Appointment{
		AppointmentSource: appointmentSource,
		AppointmentTime:   &apptDate,
		AppointmentType:   appointmentType,
		DealerID:          1,
		TenantID:          2,
	}
	db.DB(MongoDBName).C(MongoAppointmentsCollection).Insert(&appt)
	precomputeMongo(db, appt)
}

func precomputeMongo(db *mgo.Session, appt Appointment) {
	appointmentDate := time.Date(appt.AppointmentTime.Year(), appt.AppointmentTime.Month(), appt.AppointmentTime.Day(),
		0, 0, 0, 0, time.UTC)
	appointmentMonthStr := appt.AppointmentTime.Format("200601")

	report := new(HReport)

	err := db.DB(MongoDBName).C(MongoReportCollection).Find(bson.M{"_id": appointmentMonthStr}).One(report)
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	}
	if err == mgo.ErrNotFound {
		report = GenerateTemplate(*appt.AppointmentTime)
		err := db.DB(MongoDBName).C(MongoReportCollection).Insert(report)
		if err != nil {
			panic(err)
		}
	}
	selector := bson.M{
		"_id":            appointmentMonthStr,
		"dailyStat.date": appointmentDate,
	}

	var updateFieldStr string
	switch appt.AppointmentType {
	case "PDI":
		updateFieldStr = "appointmentTypePDICount"
	case "SCHEDULED":
		updateFieldStr = "appointmentTypeScheduledCount"
	case "WALKIN":
		updateFieldStr = "appointmentTypeWalkinCount"
	}

	updateParams := bson.M{
		"$inc": bson.M{
			fmt.Sprintf("stat.%s", updateFieldStr):             1,
			fmt.Sprintf("dailyStat.$.stat.%s", updateFieldStr): 1,
		},
	}
	db.DB(MongoDBName).C(MongoReportCollection).Update(selector, updateParams)

}

func GenerateTemplate(appointmentTime time.Time) *HReport {

	appointmentDate := time.Date(appointmentTime.Year(), appointmentTime.Month(), appointmentTime.Day(),
		0, 0, 0, 0, time.UTC)
	appointmentMonth := time.Date(appointmentTime.Year(), appointmentTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	appointmentMonthStr := appointmentDate.Format("200601")

	report := HReport{
		ID:   appointmentMonthStr,
		Date: appointmentMonth,
		Stat: stat{
			AppointmentTypePDICount:       0,
			AppointmentTypeWalkinCount:    0,
			AppointmentTypeScheduledCount: 0,
		},
	}

	var dailyStatList []dailyStat
	for d := appointmentMonth; d.Before(appointmentMonth.AddDate(0, 1, 0)); d = d.AddDate(0, 0, 1) {
		dailyStatList = append(dailyStatList, dailyStat{
			Date:      d,
			DayNumber: d.Day(),
			Stat: stat{AppointmentTypePDICount: 0,
				AppointmentTypeWalkinCount:    0,
				AppointmentTypeScheduledCount: 0,
			},
		})
	}

	report.DailyStat = dailyStatList

	return &report

}
