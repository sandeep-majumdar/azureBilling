package observability

import "github.com/google/uuid"

var corrId string
var causationId string

// SetCorrId -
func SetCorrId(s string) {
	corrId = s
}

// GenCorrId -
func GenCorrId() {
	corrId = uuid.New().String()
}

// ClearCorrId -
func ClearCorrId() {
	corrId = ""
}

// GetCorrId -
func GetCorrId() string {
	return corrId
}

// SetCausationId -
func SetCausationId(s string) {
	causationId = s
}

// GenCausationId -
func GenCausationId() {
	causationId = uuid.New().String()
}

// ClearCausationId -
func ClearCausationId() {
	causationId = ""
}

// GetCausationId -
func GetCausationId() string {
	return causationId
}
