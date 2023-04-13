package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vshn/appcat-cli/internal/defaults"
	"sigs.k8s.io/yaml"
)

func init() {
	logrus.SetOutput(os.Stderr)
}

func main() {
	code := Main(os.Args, os.Stdin, os.Stdout)
	os.Exit(code)
}

func Main(args []string, in io.Reader, out io.Writer) int {
	if len(args) < 2 {
		return 1
	}
	plainArgs := args[1:]
	service, err := findServiceType(plainArgs[0])
	if err != nil {
		logrus.Error(err)
		return 1
	}

	plainArgs, err = cleanInputArguments(plainArgs)
	if err != nil {
		logrus.Error(err)
		return 1
	}

	parameters := parseArgs(plainArgs)

	_, err = decorateType(service, parameters)
	if err != nil {
		logrus.Error(err)
		return 1
	}

	err = writeYAML(service, out)
	if err != nil {
		logrus.Error(err)
		return 1
	}

	return 0
}

func writeYAML(service interface{}, out io.Writer) error {
	outYaml, err := yaml.Marshal(service)
	out.Write(outYaml)
	if err != nil {
		return err
	}
	return nil
}

// Checks and parses input parameters of type "--foo=bar" and "--foo 1 1 1 1"
// returns fixed argument list
func cleanInputArguments(arguments []string) ([]string, error) {
	var fixedArguments []string
	copy(fixedArguments, arguments)
	for _, argument := range arguments {
		if strings.HasPrefix(argument, "--") && strings.Contains(argument, "=") {
			param, value, _ := strings.Cut(argument, "=")
			if value != "" && value != " " {
				fixedArguments = append(fixedArguments, param, value)
			} else {
				fixedArguments = append(fixedArguments, param)
			}
		} else {
			fixedArguments = append(fixedArguments, argument)
		}
	}
	//check for missing arguments
	for index, argument := range fixedArguments[:len(fixedArguments)-1] {
		if strings.HasPrefix(argument, "--") {
			if strings.HasPrefix(fixedArguments[index+1], "--") {
				return nil, fmt.Errorf("parameter %s is missing a value", argument)
			}
		}
	}
	if strings.HasPrefix(fixedArguments[len(fixedArguments)-1], "--") {
		return nil, fmt.Errorf("parameter %s is missing a value", fixedArguments[len(fixedArguments)-1])
	}
	return fixedArguments, nil
}

// Takes the input arguments and outputs them as separate values
// parameterName = "Parent.child1.child2. ... .childn"
// [parameterName-1 value-1, ...., parameterName-n value-n ] -> [[parent-1, child-1, ...child-n, value1],...]
func parseArgs(args []string) [][]string {
	var splitParameters [][]string
	argsLength := len(args)
	for index, parameterNames := range args[1:] {
		if index%2 != 0 || index+2 == argsLength {
			continue
		}
		parameterAndValue := strings.Split(strings.TrimPrefix(parameterNames, "--"), ".")
		parameterAndValue = append(parameterAndValue, args[index+2])
		splitParameters = append(splitParameters, parameterAndValue)
	}
	return splitParameters
}

// Iterates over all parameters and tries to set them
func decorateType(serviceType interface{}, fieldNames [][]string) (interface{}, error) {
	for _, v := range fieldNames {
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
			err := setFields(reflectedServiceType, value)
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
func setFields(field reflect.Value, value string) error {
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

// Selects the appropriate service type and initializes it with default values
func findServiceType(parsedType string) (interface{}, error) {
	parsedType = strings.ToLower(parsedType)
	if strings.Contains(parsedType, "vshn") {
		getDefault, ok := defaults.VSHN_TYPES[parsedType]
		if !ok {
			return nil, fmt.Errorf("vshn does not support service Type %s", parsedType)
		}
		return getDefault(), nil
	} else if strings.Contains(parsedType, "exoscale") {
		getDefault, ok := defaults.EXOSCALE_TYPES[parsedType]
		if !ok {
			return nil, fmt.Errorf("exoscale does not support service Type %s", parsedType)
		}
		return getDefault(), nil
	} else {
		return nil, fmt.Errorf("service Type %s could not be found", parsedType)
	}
}
