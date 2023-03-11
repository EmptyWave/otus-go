package hw09structvalidator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	reDigit = regexp.MustCompile(`([0-9-]+)`)
	reIn    = regexp.MustCompile(`(in:|,)`)
	reRe    = regexp.MustCompile(`[^regexp:].*`)
)

type Field struct {
	fieldName string
	fieldVal  reflect.Value
	tagVal    reflect.StructTag
}

func (field Field) validateKeyLen(key string) error {
	slice := reflect.ValueOf([]string{})

	switch field.fieldVal.Kind() { //nolint:exhaustive
	case reflect.String:
		slice = reflect.Append(slice, field.fieldVal)
	case reflect.Slice:
		slice = reflect.AppendSlice(slice, field.fieldVal.Slice(0, field.fieldVal.Len()))
	default:
		return ErrValidateWrong
	}

	valKey, err := strconv.Atoi(reDigit.FindString(key))
	if err != nil {
		return err
	}

	for i := 0; i < slice.Len(); i++ {
		v := slice.Index(i)

		switch {
		case valKey < 0:
			return ErrValidateWrong
		case len(v.String()) != valKey:
			return ErrValidate
		}
	}
	return nil
}

func (field Field) validateKeyIn(key string) error {
	slice := reflect.ValueOf([]string{})

	switch field.fieldVal.Kind() { //nolint:exhaustive
	case reflect.String:
		slice = reflect.Append(slice, reflect.ValueOf(field.fieldVal.String()))
	case reflect.Int:
		slice = reflect.ValueOf([]int{})
		slice = reflect.Append(slice, field.fieldVal)
	case reflect.Slice:
		slice = reflect.AppendSlice(slice, field.fieldVal.Slice(0, field.fieldVal.Len()))
	default:
		return ErrValidateWrong
	}

	arrayKey := reIn.Split(key, -1)

	if len(arrayKey) < 2 {
		return ErrValidateWrong
	}

	arrayKey = arrayKey[1:]

	var valResult string

	for i := 0; i < slice.Len(); i++ {
		val := slice.Index(i)

		switch val.Kind() { //nolint:exhaustive
		case reflect.String:
			valResult = val.String()
		case reflect.Int:
			valResult = strconv.Itoa(int(val.Int()))
		default:
			return ErrValidateWrong
		}

		keyOK := true
		for _, k := range arrayKey {
			if k == valResult {
				keyOK = false
				break
			}
		}

		if keyOK {
			return ErrValidate
		}
	}
	return nil
}

func (field Field) validateKeyRegexp(key string) error {
	if field.fieldVal.Kind() == reflect.String {
		re := reRe.FindString(key)
		result, err := regexp.MatchString(re, field.fieldVal.String())
		if err != nil {
			return err
		}

		if !result {
			return ErrValidate
		}
	}
	return nil
}

func (field Field) validateKeyMinMax(key string, less bool) error {
	if field.fieldVal.Kind() != reflect.Int {
		return ErrValidateWrong
	}

	valKey, err := strconv.Atoi(reDigit.FindString(key))
	if err != nil {
		return err
	}

	switch {
	case valKey < 0:
		return ErrValidateWrong
	case less && field.fieldVal.Int() < int64(valKey):
		return ErrValidate
	case !less && field.fieldVal.Int() > int64(valKey):
		return ErrValidate
	}
	return nil
}

func (field Field) parsing(errValid *ValidationErrors) {
	var err error

	if keys, ok := field.tagVal.Lookup("validate"); ok && keys != "" {
		for _, key := range strings.Split(keys, "|") {
			switch {
			case strings.Contains(key, "len:"):
				err = field.validateKeyLen(key)
			case strings.Contains(key, "regexp:"):
				err = field.validateKeyRegexp(key)
			case strings.Contains(key, "min:"):
				err = field.validateKeyMinMax(key, true)
			case strings.Contains(key, "max:"):
				err = field.validateKeyMinMax(key, false)
			case strings.Contains(key, "in:"):
				err = field.validateKeyIn(key)
			default:
				err = nil
			}

			if err != nil {
				*errValid = append(
					*errValid,
					ValidationError{
						Field: field.fieldName,
						Err:   err,
					},
				)
			}
		}
	}
}

func Validate(v interface{}) error {
	errValid := ValidationErrors{}
	valStruct := reflect.ValueOf(v)

	if valStruct.Kind() != reflect.Struct {
		return ErrNoStruct
	}

	typeStruct := valStruct.Type()
	for i := 0; i < typeStruct.NumField(); i++ {
		Field{
			fieldName: typeStruct.Field(i).Name,
			fieldVal:  valStruct.Field(i),
			tagVal:    typeStruct.Field(i).Tag,
		}.parsing(&errValid)
	}

	if len(errValid) != 0 {
		return errValid
	}

	return nil
}
