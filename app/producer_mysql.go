package app

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

func ProduceAppointmentForMysql(db *gorm.DB, apptDate time.Time, appointmentSource, appointmentType string) {
	appt := Appointment{
		AppointmentSource: appointmentSource,
		AppointmentTime:   &apptDate,
		AppointmentType:   appointmentType,
		DealerID:          1,
		TenantID:          2,
	}
	db.Create(&appt)
	preCompute(db, appt)
}

func preCompute(db *gorm.DB, appt Appointment) {
	appointmentDate := time.Date(appt.AppointmentTime.Year(), appt.AppointmentTime.Month(), appt.AppointmentTime.Day(),
		0, 0, 0, 0, time.UTC)

	var pdiCount, walkinCount, scheduledCount uint
	var updateFieldStr string
	switch appt.AppointmentType {
	case "PDI":
		pdiCount = 1
		updateFieldStr = "appointment_type_pdi_count"
	case "SCHEDULED":
		scheduledCount = 1
		updateFieldStr = "appointment_type_scheduled_count"
	case "WALKIN":
		walkinCount = 1
		updateFieldStr = "appointment_type_walkin_count"
	}

	appointmentDateStr := appointmentDate.Format("20060102")
	var report Report
	if err := db.Where("id = ?", appointmentDateStr).First(&report).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			report := &Report{
				ID:                            appointmentDateStr,
				AppointmentDate:               &appointmentDate,
				AppointmentTypePDICount:       pdiCount,
				AppointmentTypeScheduledCount: scheduledCount,
				AppointmentTypeWalkinCount:    walkinCount,
			}
			db.Create(report)
		} else {
			panic(err)
		}

	} else {
		db.Model(Report{}).Where("id = ?", appointmentDateStr).UpdateColumn(updateFieldStr, gorm.Expr(fmt.Sprintf("%s + ?", updateFieldStr), 1))
	}
}
