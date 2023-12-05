package models

type TemplateData struct {
	StringData map[string]string
	IntData    map[string]int
	FloatData  map[string]float32
	OtherData  map[string]interface{}
	CSRFToken  string
	Flash      string
	Warning    string
	Error      string
}
