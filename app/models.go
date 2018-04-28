package app

import (
	"time"
)

type Appointment struct {
	ID                uint       `gorm:"primary_key" bson:"-"`
	TenantID          int        `bson:"tenantID"`
	DealerID          int        `bson:"dealerID"`
	AppointmentTime   *time.Time `bson:"appointmentTime"`
	AppointmentType   string     `bson:"appointmentType"`
	AppointmentSource string     `bson:"appointmentSource"`
}

type Report struct {
	ID                            string `gorm:"primary_key"`
	AppointmentDate               *time.Time
	AppointmentTypePDICount       uint
	AppointmentTypeScheduledCount uint
	AppointmentTypeWalkinCount    uint
}

type HReport struct {
	ID        string      `bson:"_id"`
	Date      time.Time   `bson:"date"`
	Stat      stat        `bson:"stat"`
	DailyStat []dailyStat `bson:"dailyStat"`
}

type dailyStat struct {
	DayNumber int       `bson:"dayNumber"`
	Date      time.Time `bson:"date"`
	Stat      stat      `bson:"stat"`
}

type stat struct {
	AppointmentTypePDICount       uint `bson:"appointmentTypePDICount"`
	AppointmentTypeScheduledCount uint `bson:"appointmentTypeScheduledCount"`
	AppointmentTypeWalkinCount    uint `bson:"appointmentTypeWalkinCount"`
}
