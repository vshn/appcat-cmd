package util

import (
	"fmt"
	"strings"
)

type Input struct {
	parameterHierarchy []string
	value              string
	unset              bool
}

// Takes the input arguments and outputs them as separate Input structs
// returns a list of Input structs and an error if the arguments are invalid
func ParseArgs(args []string) ([]Input, error) {
	cleanArgs := FormatInputArguments(args)

	err := CheckForMissingValues(cleanArgs)
	if err != nil {
		return nil, fmt.Errorf("Input is missing a value: %s", err)
	}
	inputList := mapArgsToInput(cleanArgs)
	return inputList, nil
}

// Takes the input arguments and outputs them as separate Input structs
func mapArgsToInput(args []string) []Input {
	var inputList []Input
	input := Input{}
	for _, arg := range args {
		if isParameter(arg) {
			input.parameterHierarchy = strings.Split(strings.TrimPrefix(arg, "--"), ".")
		} else if isParamToUnset(arg) {
			param := strings.TrimPrefix(arg, "--")
			param = strings.TrimSuffix(param, "-")
			input.parameterHierarchy = strings.Split(param, ".")
			input.unset = true
			inputList = append(inputList, input)
			input = Input{}
		} else {
			input.value = arg
			inputList = append(inputList, input)
			input = Input{}
		}

	}
	return inputList
}

// Parses raw cli input parameters and returns a fixed argument list
// if a = is used in any form of key=value pair, the = needs to be at least the suffix of the key
// otherwise it's impossible to differentiate between a key=value pair and a key with a value(starting with =)
func FormatInputArguments(arguments []string) []string {
	var fixedArguments []string
	value := ""
	for _, argument := range arguments {
		if isParameterValuePair(argument) {
			if value != "" {
				fixedArguments = append(fixedArguments, value)
				value = ""
			}
			param, cutValue, _ := strings.Cut(argument, "=")
			fixedArguments = append(fixedArguments, param)
			value = cutValue
		} else if isParameter(argument) {
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
	lastIndex := len(arguments) - 1
	lastArgument := arguments[lastIndex]
	if isParameter(lastArgument) && !isParamToUnset(lastArgument) {
		return fmt.Errorf("parameter '%s' is missing a value", lastArgument)
	} else if isValue(arguments[0]) {
		return fmt.Errorf("parameter '%s' is missing a value", arguments[0])
	}
	var prevArgument string
	for index, argument := range arguments[1:] {
		if isValue(argument) {
			prevArgument = arguments[index]
			if !isParameter(prevArgument) {
				return fmt.Errorf("value '%s' is missing a key", prevArgument)
			}
		} else if isParamToUnset(argument) {
			prevArgument = arguments[index]
			if isParameter(prevArgument) {
				return fmt.Errorf("parameter '%s' is missing a value", prevArgument)
			}
		}
	}

	return nil
}
func isValue(arg string) bool {
	if !isParamToUnset(arg) && !isParameter(arg) && !isParameterValuePair(arg) {
		return true
	}
	return false
}

func isParameterValuePair(arg string) bool {
	if strings.HasPrefix(arg, "--") && strings.Contains(arg, "=") && !isParamToUnset(arg) {
		return true
	}
	return false
}

func isParameter(arg string) bool {
	if strings.HasPrefix(arg, "--") && !strings.Contains(arg, "=") && !isParamToUnset(arg) {
		return true
	}
	return false
}

func isParamToUnset(arg string) bool {
	if strings.HasPrefix(arg, "--") && strings.HasSuffix(arg, "-") {
		return true
	}
	return false
}
