package forms

import (
	"fmt"
	"net/url"
	"strings"
	"regexp"
	"golang.org/x/exp/slices"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

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


func (f *Form) Required(fields ...string){
	for _,field:= range fields {
		value:= f.Get(field)
		if strings.TrimSpace(value)==""{
			f.Errors.Add(field, "This field is required")
		}
	}
}

func (f *Form) MaxLength(field string, d int){
	value:= f.Get(field)
	if len([]rune(value))>d{
		f.Errors.Add(field, fmt.Sprintf("This field is too long. Maximum value is %d", d))
	}
}
func (f *Form) MinLength(field string, d int){
	value:= f.Get(field)
	if len([]rune(value))<d{
		f.Errors.Add(field, fmt.Sprintf("This field is too short. Minimum value is %d", d))
	}
}


func (f *Form) PermittedValues (field string, options ...string){
	value:= f.Get(field)
	if !slices.Contains(options, value){
		f.Errors.Add(field, "This field is invalid")
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp){
	value:= f.Get(field)
	if value==""{
		return
	}
	if !pattern.MatchString(value){
		f.Errors.Add(field, "This field is invalid")
	}
}


func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}