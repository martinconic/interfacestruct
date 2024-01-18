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
			value, err := getAssertedTypedValue(data[i], fieldType)
			if err != nil {
				return dt.Data, err
			}
			fieldValue.Set(value)
		} else {
			fieldValue.Set(reflect.ValueOf(data[i]))
		}
	}

	return dt.Data, nil
}

func getAssertedTypedValue(data interface{}, fieldType reflect.StructField) (reflect.Value, error) {
	if fieldType.Type.Kind() == reflect.Int {
		d, _ := strconv.Atoi(data.(string))
		return reflect.ValueOf(d), nil
	} else if fieldType.Type.Kind() == reflect.Uint64 {
		d, _ := strconv.ParseUint(data.(string), 10, 64)
		return reflect.ValueOf(d), nil
	} else if fieldType.Type.Kind() == reflect.Float32 {
		d, _ := strconv.ParseFloat(data.(string), 32)
		return reflect.ValueOf(d), nil
	} else if fieldType.Type.Kind() == reflect.Float64 {
		d, _ := strconv.ParseFloat(data.(string), 64)
		return reflect.ValueOf(d), nil
	} else if fieldType.Type == reflect.TypeOf(time.Time{}) {
		layout := "2006-01-02T15:04:05.999999999Z07:00"
		d, _ := time.Parse(layout, data.(string))
		return reflect.ValueOf(d), nil
	} else if fieldType.Type.Kind() == reflect.Bool {
		d, _ := strconv.ParseBool(data.(string))
		return reflect.ValueOf(d), nil
	}

	return reflect.ValueOf(nil), fmt.Errorf("type mismatch for field %s", fieldType.Name)
}
