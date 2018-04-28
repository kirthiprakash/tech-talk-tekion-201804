package app

import "math/rand"

func Random(min int, max int) int {
	return rand.Intn(max-min) + min
}

var SourceList = []string{"WEB", "MOBILE", "PHONE"}
var AppointmentTypeList = []string{"SCHEDULED", "WALKIN", "PDI"}
