package app

import (
	"time"
)

type Appointment struct {
	ID    uint	`gorm:"primary_key"`
	TenantID	int
	DealerID	int
	AppointmentTime  *time.Time
	AppointmentType string
	AppointmentSource string
}
