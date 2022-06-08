package utils

import (
	"errors"
	"log"
	"reflect"
	"strconv"
)

func GetStructFieldNames(targetStruct interface{}) []string {
	t := reflect.TypeOf(targetStruct)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Fatal(ErrWrapOrWithMessage(true, errors.New("input is not a struct")))
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

func GetStructValuesAsString(targetStruct interface{}) []string {
	v := reflect.ValueOf(targetStruct)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		log.Fatal(ErrWrapOrWithMessage(true, errors.New("input is not a struct")))
	}

	fieldNum := v.NumField()
	result := make([]string, 0, fieldNum)
	var value reflect.Value
	for i := 0; i < fieldNum; i++ {
		value = v.Field(i)
		// fmt.Println(value.Kind().String())
		switch value.Kind().String() {
		case "uint8":
			result = append(result, strconv.FormatUint(value.Uint(), 10))
		case "int", "int64":
			result = append(result, strconv.FormatInt(value.Int(), 10))
		case "float64":
			result = append(result, strconv.FormatFloat(value.Float(), 'f', -1, 64))
		case "bool":
			result = append(result, strconv.FormatBool(value.Bool()))
		default:
			log.Fatal(ErrWrapOrWithMessage(true, errors.New("unrecognized fields: "+value.Kind().String())))
		}
	}
	return result
}
