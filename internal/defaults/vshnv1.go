package defaults

import (
	vshnv1 "github.com/vshn/appcat/apis/vshn/v1"
)

func (d *Defaults) GetVSHNPostgreSQLDefault() *vshnv1.VSHNPostgreSQL {
	var postgreSQLdefault vshnv1.VSHNPostgreSQL

	postgreSQLdefault.Spec.Parameters.Service.MajorVersion = "14"
	postgreSQLdefault.Spec.Parameters.Size.CPU = "600m"
	postgreSQLdefault.Spec.Parameters.Size.Disk = "80Gi"
	postgreSQLdefault.Spec.Parameters.Size.Memory = "3500Mi"
	postgreSQLdefault.Spec.Parameters.Size.Requests.CPU = "300m"
	postgreSQLdefault.Spec.Parameters.Size.Requests.Memory = "1000Mi"
	postgreSQLdefault.Spec.Parameters.Backup.Schedule = "30 23 * * *"
	postgreSQLdefault.Spec.Parameters.Backup.Retention = 12
	postgreSQLdefault.Spec.Parameters.Scheduling.NodeSelector = map[string]string{"appuio.io/node-class": "plus"}

	return &postgreSQLdefault
}

func (d *Defaults) GetVSHNRedisDefault() *vshnv1.VSHNRedis {
	var redisDefault vshnv1.VSHNRedis

	redisDefault.Spec.Parameters.TLS.TLSAuthClients = true
	redisDefault.Spec.Parameters.TLS.TLSEnabled = true
	redisDefault.Spec.Parameters.Service.Version = "7.0"
	redisDefault.Spec.Parameters.Service.RedisSettings = "|activedefrag yes"
	redisDefault.Spec.Parameters.Size.Disk = "80Gi"
	redisDefault.Spec.Parameters.Size.CPULimits = "1000m"
	redisDefault.Spec.Parameters.Size.CPURequests = "500m"
	redisDefault.Spec.Parameters.Size.MemoryRequests = "500Mi"
	redisDefault.Spec.Parameters.Size.MemoryLimits = "1Gi"

	return &redisDefault
}
