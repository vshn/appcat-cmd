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

`appcatcli` supports the following arguments:

- Argument #1: `Name of custom resource`. Mandatory.
- Argument #2: `--kind` with value `ServiceKind`. Mandatory
- Argument #2...N: `Parameters of custom resource`. Optional.

## Testing

To test the tool; in the root of the repository run:
```shell
go test .
```

## License

This project is licensed under the [BSD 3-Clause License](LICENSE)
