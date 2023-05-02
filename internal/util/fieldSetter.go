package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// Iterates over all parameters and tries to set them
func DecorateType(serviceType interface{}, inputs []Input) (interface{}, error) {
	for _, input := range inputs {
		_, err := setParameter(serviceType, input)
		if err != nil {
			return nil, err
		}
	}

	return serviceType, nil
}

// Iterates through service Type to find and set the required parameters
// logs error and exits if parameter does not match any in Type specified field names
func setParameter(serviceType interface{}, input Input) (interface{}, error) {
	reflectedServiceType := reflect.ValueOf(serviceType).Elem()
	var parameterName string
	var err error
	for _, parameterName = range input.ParameterHierarchy {
		if !reflectedServiceType.FieldByName(parameterName).IsValid() {

			parameterName, err = getStringCase(getAllFieldNames(reflectedServiceType), parameterName)
			if err != nil {
				err = fmt.Errorf("%w %s contains field with name %s : %t",
					err,
					reflectedServiceType.Type().Name(),
					parameterName,
					reflectedServiceType.FieldByName(parameterName).IsValid(),
				)
				return nil, err
			}
		}

		reflectedServiceType = reflectedServiceType.FieldByName(parameterName)
	}

	err = SetFields(reflectedServiceType, input)
	if err != nil {
		err = fmt.Errorf(
			"%w cannot assign value %s to field %s with field Type %s",
			err,
			input.Value,
			strings.Join(input.ParameterHierarchy, HIERARCHY_DELIMITER),
			reflectedServiceType.Kind(),
		)
		err = fmt.Errorf("%w %s contains field with name %s : %t",
			err,
			reflectedServiceType.Type().Name(),
			parameterName,
			reflectedServiceType.IsValid(),
		)
		return nil, err
	}
	info := fmt.Sprintf("setting field: %s value: %s", strings.Join(input.ParameterHierarchy, HIERARCHY_DELIMITER), input.Value)
	logrus.Info(info)

	return serviceType, nil
}

// Sets value of reflected field with type checking
func SetFields(field reflect.Value, input Input) error {
	if input.Unset {
		field.Set(reflect.Zero(field.Type()))
	} else if input.IsJson {
		field.Set(reflect.Zero(field.Type()))
		var jsonInput map[string]interface{}
		//err := json.Unmarshal([]byte(input.Value), field.Addr().Interface())
		err := json.Unmarshal([]byte(input.Value), &jsonInput)
		if err != nil {
			return fmt.Errorf("Json value could not be Unmarshalled: %s", err)
		}
		for key, value := range jsonInput {
			err := SetFields(field.FieldByName(key), Input{Value: fmt.Sprintf("%v", value)})
			if err != nil {
				return err
			}
		}
	} else if field.Kind() == reflect.String {
		field.SetString(input.Value)
	} else if field.Kind() >= reflect.Int && field.Kind() <= reflect.Int64 {
		intValue, err := strconv.ParseInt(input.Value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	} else if field.Kind() >= reflect.Uint && field.Kind() <= reflect.Uint64 {
		intValue, err := strconv.ParseUint(input.Value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(intValue)
	} else if field.Kind() == reflect.Float32 || field.Kind() == reflect.Float64 {
		floatValue, err := strconv.ParseFloat(input.Value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatValue)
	} else if field.Kind() == reflect.Bool {
		boolValue, err := strconv.ParseBool(input.Value)
		if err != nil {
			return err
		}
		field.SetBool(boolValue)
	} else {
		return fmt.Errorf("setFields failed with field Type: %T and value: %s", reflect.TypeOf(field), input.Value)
	}
	return nil
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
