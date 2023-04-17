package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"

	"github.com/vshn/appcat-cli/internal/applications"
	"github.com/vshn/appcat-cli/internal/util"
)

func init() {
	logrus.SetOutput(os.Stderr)
}

func printUsage(cmd string, apps applications.AppMap) {
	out := fmt.Sprintf(`usage: %s <type> [options]

Generate AppCat YAML manifests

Known types:
`, cmd)

	names := apps.Names()
	longest := util.Longest(names)
	format := fmt.Sprintf("  %%-%ds (%%s.%%s)\n", longest+2)

	for _, name := range names {
		app := apps[name]
		out += fmt.Sprintf(format, name, app.Kind, app.APIVersion)
	}

	fmt.Fprintln(os.Stderr, out)
}

func main() {
	apps := applications.MakeAppMap()
	code := Main(apps, os.Args, os.Stdin, os.Stdout)
	os.Exit(code)
}

// Main function
//
// Separated from `main` for testing purposes.
//
// # Errors
//
// If during executions, error occur due to user input errors, an appropriate
// error message is logged, and a non-zero exit code is returned.
//
// # Panics
//
// If during execution, an unrecoverable error occurs (usually due to a bug),
// an error message is logged and the program will panic.
func Main(apps applications.AppMap, args []string, in io.Reader, out io.Writer) int {
	if len(args) < 2 {
		printUsage(args[0], apps)
		return 1
	}
	plainArgs := args[1:]

	fmt.Printf("%#v\n", args)

	parsedType := util.NormalizeName(plainArgs[0])
	app, ok := apps[parsedType]
	if !ok {
		logrus.Errorf("service type '%s' is not supported", parsedType)
		printUsage(args[0], apps)
		return 1
	}

	service := app.GetDefault()
	plainArgs, err := util.CleanInputArguments(plainArgs)
	if err != nil {
		logrus.Errorf("Invalid arguments: %s", err)
		return 1
	}

	parameters := util.ParseArgs(plainArgs)

	// TODO: Setting an invalid value just ignores it instead of erroring
	// example: `go run . VSHNPostgreSQL --spec.parameters.backup.retention asdf``
	_, err = util.DecorateType(service, parameters)
	if err != nil {
		logrus.Errorf("failed setting parameters: %s", err)
		return 1
	}

	err = writeYAML(service, out)
	if err != nil {
		logrus.Panicf("failed writing YAML: %s", err)
		return 1
	}

	return 0
}

func writeYAML(service interface{}, out io.Writer) error {
	outYaml, err := yaml.Marshal(service)
	if err != nil {
		return err
	}

	_, err = out.Write(outYaml)
	if err != nil {
		return err
	}

	return nil
}
