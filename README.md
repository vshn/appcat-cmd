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

`appcatcli` expects input in the form of:<br>
`appcatcli ServiceName --kind ServiceKind [Options]`<br>

#### Defaults:
Global:
| Key  | Value  |
| ------ | -------|
|`Spec.WriteConnectionSecretToRef.Name`| ServiceName+refSlug |

Exoscale:
|ServiceKind |Key  | Value  |
| ------ | -------|-------|
|ExoscalePostgreSQL|`Spec.Parameters.Service.MajorVersion`| 14|
|ExoscalePostgreSQL|`Spec.Parameters.Service.Zone`| ch-dk-2|
|ExoscalePostgreSQL|`Spec.Parameters.Size.Plan`| hobbyist-2|
|ExoscalePostgreSQL|`Spec.Parameters.Backup.TimeOfDay`| 12:00:00|
|ExoscalePostgreSQL|`Spec.Parameters.Maintenance.DayOfWeek`| sunday|
|ExoscalePostgreSQL|`Spec.Parameters.Maintenance.TimeOfDaye`| 00:00:00|
|ExoscaleRedis|`Spec.Parameters.Maintenance.DayOfWeek`|sunday|
|ExoscaleRedis|`Spec.Parameters.Maintenance.TimeOfDay`|00:00:00|
|ExoscaleRedis|`Spec.Parameters.Service.Zone`|ch-dk-2|
|ExoscaleKafka|`Spec.Parameters.Service.Version`|3.4.0|
|ExoscaleKafka|`Spec.Parameters.Service.Zone`|ch-dk-2|
|ExoscaleKafka|`Spec.Parameters.Size.Plan`|startup-2|
|ExoscaleKafka|`Spec.Parameters.Maintenance.DayOfWeek`|sunday|
|ExoscaleKafka|`Spec.Parameters.Maintenance.TimeOfDay`|00:00:00|
|ExoscaleMySQL|`Spec.Parameters.Service.MajorVersion`|8|
|ExoscaleMySQL|`Spec.Parameters.Service.Zone`|ch-dk-2|
|ExoscaleMySQL|`Spec.Parameters.Size.Plan`|hobbyist-2|
|ExoscaleMySQL|`Spec.Parameters.Backup.TimeOfDay`|12:00:00|
|ExoscaleMySQL|`Spec.Parameters.Maintenance.DayOfWeek`|sunday|
|ExoscaleMySQL|`Spec.Parameters.Maintenance.TimeOfDay`|00:00:00|
|ExoscaleOpenSearch|`Spec.Parameters.Service.MajorVersion`|2|
|ExoscaleOpenSearch|`Spec.Parameters.Service.Zone`|ch-dk-2|
|ExoscaleOpenSearch|`Spec.Parameters.Size.Plan`|hobbyist-2|
|ExoscaleOpenSearch|`Spec.Parameters.Backup.TimeOfDay`|12:00:00|
|ExoscaleOpenSearch|`Spec.Parameters.Maintenance.DayOfWeek`|sunday|
|ExoscaleOpenSearch|`Spec.Parameters.Maintenance.TimeOfDay`|00:00:00|

VSHN:
|ServiceKind |Key  | Value  |
| ------ | -------|-------|
|VSHNPostgreSQL|`Spec.Parameters.Service.MajorVersion`|14|
|VSHNPostgreSQL|`Spec.Parameters.Size.CPU`|600m|
|VSHNPostgreSQL|`Spec.Parameters.Size.Disk`|80Gi|
|VSHNPostgreSQL|`Spec.Parameters.Size.Memory`|3500Mi|
|VSHNPostgreSQL|`Spec.Parameters.Size.Requests.CPU`|300m|
|VSHNPostgreSQL|`Spec.Parameters.Size.Requests.Memory`|1000Mi|
|VSHNPostgreSQL|`Spec.Parameters.Backup.Schedule`||30 23 * * *|
|VSHNPostgreSQL|`Spec.Parameters.Backup.Retention`|12|
|VSHNPostgreSQL|`Spec.Parameters.Scheduling.NodeSelector`|{"appuio.io/node-class": "plus"}|
|VSHNRedis|`Spec.Parameters.TLS.TLSAuthClients`|true|
|VSHNRedis|`Spec.Parameters.TLS.TLSEnabled`|true|
|VSHNRedis|`Spec.Parameters.Service.Version`|7.0|
|VSHNRedis|`Spec.Parameters.Service.RedisSettings`||activedefrag yes|
|VSHNRedis|`Spec.Parameters.Size.Disk`|80Gi|
|VSHNRedis|`Spec.Parameters.Size.CPULimits`|1000m|
|VSHNRedis|`Spec.Parameters.Size.CPURequests`|500m|
|VSHNRedis|`Spec.Parameters.Size.MemoryRequests`|500Mi|
|VSHNRedis|`Spec.Parameters.Size.MemoryLimits`|1Gi|

## Testing

For testing we use unit test as well as golden tests both can be run via:
```shell
go test .
```

## License

This project is licensed under the [BSD 3-Clause License](LICENSE)
