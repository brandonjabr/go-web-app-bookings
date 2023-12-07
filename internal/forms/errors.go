package forms

type errors map[string][]string

func (e errors) Add(field, msg string) {
	e[field] = append(e[field], msg)
}

func (e errors) Get(field string) string {
	errorStr := e[field]
	if len(errorStr) == 0 {
		return ""
	}

	return errorStr[0]
}
