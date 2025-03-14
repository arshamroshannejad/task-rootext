package helpers

import (
	"net/url"
	"strconv"
)

type Map map[string]string

type Validator struct {
	Errors Map
}

func NewValidator() *Validator {
	return &Validator{
		Errors: make(Map),
	}
}

func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) Add(key, value string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = value
	}
}

func (v *Validator) Check(ok bool, key, value string) {
	if !ok {
		v.Add(key, value)
	}
}

func (v *Validator) In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

func (v *Validator) ReadQsString(qs url.Values, key, defaultValue string) string {
	k := qs.Get(key)
	if k == "" {
		return defaultValue
	}
	return k
}

func (v *Validator) ReadQsInt(qs url.Values, key string, defaultValue int) int {
	k := qs.Get(key)
	if k == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(k)
	if err != nil {
		v.Add(key, "must be an integer value")
		return defaultValue
	}
	return i
}
