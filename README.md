# appcat-cli

`appcat-cli` generates yaml files for custom kubernetes resources from command line arguments.

**Warning:** This tool is currently under heavy development. Stuff may change & break between releases!


## Goal & Purpose
`appcat-cli` is developed as a converter tool for [k8ify](https://github.com/vshn/k8ify).
It's purpose is to generate yaml files for custom kubernetes resources to be processed by `k8ify` into kubernetes manifests.

### Non-Goals

Out of scope(for now) of this project are:
- Checking scope and availability of custom resources on target cluster

## Mode of Operation

### Command Line Arguments

`appcatcli` supports the following command line arguments:

- Argument #1: `Name of custom resource type`. Mandatory.
- Argument #2...N: `Parameters of custom resource`. Optional.

## Testing

To test the tool with standart parameters, in the root of the repository run:
```shell
go run . "VSHNPostgreSQL" "--spec.parameters.service.majorVersion" "1" "--spec.parameters.size.CPU" "6m" "--spec.parameters.size.disk" "8Gi" "--spec.parameters.size.memory" "35Mi" "--Spec.Parameters.Size.Requests.CPU" "3m" "--spec.parameters.size.requests.memory" "1Mi" "--spec.parameters.backup.schedule" "3 2 * * *" "--spec.parameters.backup.retention" "1"
```

## License

This project is licensed under the [BSD 3-Clause License](LICENSE)
