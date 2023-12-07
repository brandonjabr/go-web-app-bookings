package models

import "github.com/brandonjabr/go-web-app-bookings/internal/forms"

type TemplateData struct {
	StringData map[string]string
	IntData    map[string]int
	FloatData  map[string]float32
	OtherData  map[string]interface{}
	CSRFToken  string
	Flash      string
	Warning    string
	Error      string
	Form       *forms.Form
}
