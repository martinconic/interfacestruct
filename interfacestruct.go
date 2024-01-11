package interfacestruct

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Data generic struct used for converting to an []interface{} into structures like Member, Child, OrderNumber
type DataGeneric[T any] struct {
	Data T
}

type RequestData struct {
	Values [][]interface{}
}

func (dt *DataGeneric[T]) ConvertToInterfaceRequest() RequestData {
	var requestData RequestData
	v := reflect.ValueOf(dt.Data)

	datai := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		datai[i] = v.Field(i).Interface()
	}

	requestData.Values = append(requestData.Values, datai)

	return requestData
}

func (dt *DataGeneric[T]) ConvertToStruct(data []interface{}) (T, error) {
	structType := reflect.TypeOf(dt.Data)
	structValue := reflect.ValueOf(&dt.Data).Elem()

	if len(data) != structType.NumField() {
		return dt.Data, fmt.Errorf("invalid input length")
	}

	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		fieldValue := structValue.Field(i)

		if i >= len(data) {
			return dt.Data, fmt.Errorf("not enough data to populate all fields")
		}

		// Convert the interface value to the type of the struct field
		if reflect.TypeOf(data[i]) != fieldType.Type {
			if fieldType.Type.Kind() == reflect.Int {
				d, _ := strconv.Atoi(data[i].(string))
				fieldValue.Set(reflect.ValueOf(d))
			} else if fieldType.Type == reflect.TypeOf(time.Time{}) {
				layout := "2006-01-02T15:04:05.999999999Z07:00"
				d, _ := time.Parse(layout, data[i].(string))
				fieldValue.Set(reflect.ValueOf(d))
			} else if fieldType.Type.Kind() == reflect.Bool {
				d, _ := strconv.ParseBool(data[i].(string))
				fieldValue.Set(reflect.ValueOf(d))
			} else {
				return dt.Data, fmt.Errorf("type mismatch for field %s", fieldType.Name)
			}
		} else {
			fieldValue.Set(reflect.ValueOf(data[i]))
		}
	}

	return dt.Data, nil
}
