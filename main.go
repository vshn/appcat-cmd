package main

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/vshn/appcat-cli/internal/defaults"
	"github.com/vshn/appcat-cli/internal/util"
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
	service, err := defaults.FindServiceType(plainArgs[0])
	if err != nil {
		logrus.Error(err)
		return 1
	}
	plainArgs, err = util.CleanInputArguments(plainArgs)
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
