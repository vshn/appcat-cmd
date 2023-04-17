package util

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// Iterates over all parameters and tries to set them
func DecorateType(serviceType interface{}, fieldNames [][]string) (interface{}, error) {
	for _, v := range fieldNames {
		// TODO: this function "re-parses" arguments again
		_, err := setParameter(serviceType, v[:len(v)-1], v[len(v)-1])
		if err != nil {
			return nil, err
		}
	}

	return serviceType, nil
}

// Iterates through service Type to find and set the required parameters
// logs error and exits if parameter does not match any in Type specified field names
func setParameter(serviceType interface{}, parameterHierarchy []string, value string) (interface{}, error) {
	reflectedServiceType := reflect.ValueOf(serviceType).Elem()
	for _, parameterName := range parameterHierarchy {
		if !reflectedServiceType.FieldByName(parameterName).IsValid() {
			var err error
			parameterName, err = getStringCase(getAllFieldNames(reflectedServiceType), parameterName)
			if err != nil {
				err = fmt.Errorf("%w\n%s contains field with name %s : %t",
					err,
					reflectedServiceType.Type().Name(),
					parameterName,
					reflectedServiceType.FieldByName(parameterName).IsValid(),
				)
				return nil, err
			}
		}

		reflectedServiceType = reflectedServiceType.FieldByName(parameterName)

		if reflectedServiceType.Kind() != reflect.Struct {
			err := SetFields(reflectedServiceType, value)
			if err != nil {
				err := fmt.Errorf(
					"%w\ncannot assign value %s to field %s with field Type %s",
					err,
					value,
					strings.Join(parameterHierarchy, "."),
					reflectedServiceType.FieldByName(parameterName).Kind(),
				)
				err = fmt.Errorf("%w\n%s contains field with name %s : %t",
					err,
					reflectedServiceType.Type().Name(),
					parameterName,
					reflectedServiceType.FieldByName(parameterName).IsValid(),
				)
				return nil, err
			}
			info := fmt.Sprintf("setting field: %s value: %s", strings.Join(parameterHierarchy, "."), value)
			logrus.Info(info)
		}

	}
	return serviceType, nil
}

// Returns all field names of a struct
func getAllFieldNames(field reflect.Value) []string {
	var fieldNames []string
	for i := 0; i < field.NumField(); i++ {
		if field.Type().Field(i).Anonymous {
			// Recursively finds all field names of anonymous/embedded structs
			fieldNames = append(fieldNames, getAllFieldNames(field.Field(i))...)
		} else {
			fieldNames = append(fieldNames, field.Type().Field(i).Name)
		}
	}
	return fieldNames
}

// Sets value of reflected field with type checking
func SetFields(field reflect.Value, value string) error {
	// TODO: Handle error cases on type conversion
	if field.Kind() == reflect.String {
		field.SetString(value)
	} else if field.Kind() >= reflect.Int && field.Kind() <= reflect.Int64 {
		intValue, _ := strconv.ParseInt(value, 10, 64)
		field.SetInt(intValue)
	} else if field.Kind() >= reflect.Uint && field.Kind() <= reflect.Uint64 {
		intValue, _ := strconv.ParseUint(value, 10, 64)
		field.SetUint(intValue)
	} else if field.Kind() == reflect.Float32 || field.Kind() == reflect.Float64 {
		floatValue, _ := strconv.ParseFloat(value, 64)
		field.SetFloat(floatValue)
	} else if field.Kind() == reflect.Bool {
		boolValue, _ := strconv.ParseBool(value)
		field.SetBool(boolValue)
	} else {
		return fmt.Errorf("setFields failed with field Type: %T and value: %s", field.Type(), value)
	}
	return nil
}

// Checks for all field names of a struct if any match the parameter name under Unicode case-folding
// returns the correct field name if successfull, nil if unsuccessful
func getStringCase(fieldNames []string, parameterName string) (string, error) {
	for _, fieldName := range fieldNames {
		if strings.EqualFold(fieldName, parameterName) {
			return fieldName, nil
		}
	}
	return parameterName, fmt.Errorf("could not find field with name %s", parameterName)
}
