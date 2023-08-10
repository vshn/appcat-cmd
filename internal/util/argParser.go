package util

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type Input struct {
	ParameterHierarchy []string
	Value              string
	Unset              bool
}

const (
	UNSET_PARAMETER_SUFFIX = "-"
	PARAMETER_PREFIX       = "--"
	HIERARCHY_DELIMITER    = "."
	PARAM_VALUE_INFIX      = "="
	K8S_SERVICE_KIND       = "kind"
)

// Takes the input arguments and outputs them as separate Input structs
// returns a list of Input structs and an error if an argument is invalid
func ParseArgs(args []string) ([]Input, error) {
	cleanArgs := FormatInputArguments(args)
	err := CheckForMissingValues(cleanArgs)
	if err != nil {
		return nil, fmt.Errorf("Input is missing a value: %v", err)
	}
	inputList := mapArgsToInput(cleanArgs)

	return inputList, nil
}

// Parses raw cli input parameters and returns a list of arguments
// if a "=" is used in any form of "key=value" pair, the "=" needs to be the suffix of the key or the separator between key and value
// if the "=" is the prefix of the value it is interpreted as part of the value
// otherwise it's impossible to differentiate between a "key=value" pair and a key with a value(starting with "=")
func FormatInputArguments(arguments []string) []string {
	var fixedArguments []string
	value := ""
	for _, argument := range arguments {
		if isParameterValuePair(argument) {
			if value != "" {
				fixedArguments = append(fixedArguments, value)
			}

			param, cutValue, _ := strings.Cut(argument, PARAM_VALUE_INFIX)
			fixedArguments = append(fixedArguments, param)
			value = cutValue
		} else if isParameter(argument) {
			if value != "" {
				fixedArguments = append(fixedArguments, value)
				value = ""
			}

			fixedArguments = append(fixedArguments, argument)
		} else if isParamToUnset(argument) {
			if value != "" {
				fixedArguments = append(fixedArguments, value)
				value = ""
			}
			fixedArguments = append(fixedArguments, argument)
		} else {
			value += argument
		}
	}
	lastElement := arguments[len(arguments)-1]
	if lastElement == value || isParameterValuePair(lastElement) {
		fixedArguments = append(fixedArguments, value)
	}
	return fixedArguments
}

// Takes a list of strings and checks if every parameter has a value
// returns an error if a parameter is missing a value
func CheckForMissingValues(arguments []string) error {
	lastArgument := arguments[len(arguments)-1]
	if isParameter(lastArgument) && !isParamToUnset(lastArgument) {
		return fmt.Errorf("parameter '%s' is missing a value", lastArgument)
	} else if isValue(arguments[0]) {
		return fmt.Errorf("parameter '%s' is missing a value", arguments[0])
	}
	var prevArgument string
	for i, argument := range arguments[1:] {
		if isValue(argument) {
			prevArgument = arguments[i]
			if !isParameter(prevArgument) {
				return fmt.Errorf("value '%s' is missing a key", prevArgument)
			}
		} else if isParamToUnset(argument) {
			prevArgument = arguments[i]
			if isParameter(prevArgument) {
				return fmt.Errorf("parameter '%s' is missing a value", prevArgument)
			}
		}
	}

	return nil
}
func FilterServiceKind(input []Input) (string, error) {
	serviceKindIndex := slices.IndexFunc(input, func(input Input) bool {
		return input.ParameterHierarchy[0] == K8S_SERVICE_KIND
	})
	if serviceKindIndex == -1 {
		return "", fmt.Errorf("ServiceKind is missing")
	} else {
		return input[serviceKindIndex].Value, nil
	}
}

// Takes the input arguments and outputs them as separate Input structs
func mapArgsToInput(args []string) []Input {
	var inputList []Input
	input := Input{}
	for _, arg := range args {
		if isParameter(arg) {
			input.ParameterHierarchy = strings.Split(strings.TrimPrefix(arg, PARAMETER_PREFIX), HIERARCHY_DELIMITER)
		} else if isParamToUnset(arg) {
			param := strings.TrimPrefix(arg, PARAMETER_PREFIX)
			param = strings.TrimSuffix(param, UNSET_PARAMETER_SUFFIX)
			input.ParameterHierarchy = strings.Split(param, HIERARCHY_DELIMITER)
			input.Unset = true
			inputList = append(inputList, input)
			input = Input{}
		} else {
			input.Value = arg
			inputList = append(inputList, input)
			input = Input{}
		}

	}
	return inputList
}

func isValue(arg string) bool {
	if !isParamToUnset(arg) && !isParameter(arg) && !isParameterValuePair(arg) {
		return true
	}
	return false
}

func isParameterValuePair(arg string) bool {
	if strings.HasPrefix(arg, PARAMETER_PREFIX) && strings.Contains(arg, PARAM_VALUE_INFIX) && !isParamToUnset(arg) {
		return true
	}
	return false
}

func isParameter(arg string) bool {
	if strings.HasPrefix(arg, PARAMETER_PREFIX) && !strings.Contains(arg, PARAM_VALUE_INFIX) && !isParamToUnset(arg) {
		return true
	}
	return false
}

func isParamToUnset(arg string) bool {
	if strings.HasPrefix(arg, PARAMETER_PREFIX) && strings.HasSuffix(arg, UNSET_PARAMETER_SUFFIX) {
		return true
	}
	return false
}
