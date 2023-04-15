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

func usage(cmd string, apps applications.AppMap) string {
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

	return out
}

func main() {
	apps := applications.MakeAppMap()
	code := Main(apps, os.Args, os.Stdin, os.Stdout)
	os.Exit(code)
}

func Main(apps applications.AppMap, args []string, in io.Reader, out io.Writer) int {
	if len(args) < 2 {
		fmt.Println(usage(args[0], apps))
		return 1
	}
	plainArgs := args[1:]

	parsedType := util.NormalizeName(plainArgs[0])
	app, ok := apps[parsedType]
	if !ok {
		logrus.Errorf("service Type %s is not supported", parsedType)
		usage(args[0], apps)
		os.Exit(1)
	}

	service := app.GetDefault()
	plainArgs, err := util.CleanInputArguments(plainArgs)
	if err != nil {
		logrus.Error(err)
		return 1
	}

	parameters := util.ParseArgs(plainArgs)

	_, err = util.DecorateType(service, parameters)
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
