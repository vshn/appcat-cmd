package util

import (
	"fmt"
	"strings"
)

// Checks and parses input parameters of type "--foo=bar" and "--foo 1 1 1 1"
// returns fixed argument list
func CleanInputArguments(arguments []string) ([]string, error) {
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
func ParseArgs(args []string) [][]string {
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
