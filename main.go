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
	out := fmt.Sprintf(`usage: %s <serviceName> --kind <serviceKind> [options]

Generate AppCat YAML manifests

Known service kinds:
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
	code := Main(apps, os.Args, os.Stdout)
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
func Main(apps applications.AppMap, args []string, out io.Writer) int {
	if len(args) < 3 {
		printUsage(args[0], apps)
		return 1
	}

	var err error
	var parameters []util.Input

	// args[0] is the binary
	resourceName := args[1]
	parameters, err = util.ParseArgs(args[2:])

	if err != nil {
		logrus.Errorf("Failed parsing arguments: %s", err)
		return 1
	}

	serviceKind, err := util.FilterServiceKind(parameters)
	if err != nil {
		logrus.Errorf("%s", err)
		printUsage(args[0], apps)
		return 1
	}
	serviceKind = util.NormalizeName(serviceKind)
	app, ok := apps[serviceKind]
	if !ok {
		logrus.Errorf("service type '%s' is not supported", serviceKind)
		printUsage(args[0], apps)
		return 1
	}
	service := app.GetDefault()
	parameters = append(parameters, util.Input{ParameterHierarchy: []string{"ObjectMeta", "Name"}, Value: resourceName, Unset: false})
	parameters = append([]util.Input{{ParameterHierarchy: []string{"Spec", "WriteConnectionSecretToRef", "Name"}, Value: resourceName + "-creds", Unset: false}}, parameters...)

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
