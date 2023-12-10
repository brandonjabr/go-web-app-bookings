package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		val := f.Get(field)
		if strings.TrimSpace(val) == "" {
			f.Errors.Add(field, "This field is required.")
		}
	}
}

func (f *Form) Has(field string) bool {
	val := f.Get(field)
	if val == "" {
		f.Errors.Add(field, "This field is required.")
		return false
	}

	return true
}

func (f *Form) MinLength(field string, length int) bool {
	val := f.Get(field)
	if len(val) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}

	return true
}

func (f *Form) IsValidEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
