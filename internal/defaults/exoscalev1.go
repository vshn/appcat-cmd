package defaults

import (
	exoscalev1 "github.com/vshn/component-appcat/apis/exoscale/v1"
)

var exoscaleApiVersion string = "exoscale.appcat.vshn.io/v1"

var EXOSCALE_TYPES = map[string]func() interface{}{
	"exoscalepostgresql": getExoscalePostgreSQLDefault,
	"exoscaleredis":      getExoscaleRedisDefault,
	"exoscalekafka":      getExoscaleKafkaDefault,
	"exoscalemysql":      getExoscaleMySQLdefault,
	"exoscaleopensearch": getExoscaleOpenSearchDefault,
}

func getExoscalePostgreSQLDefault() interface{} {
	var postgreSQLdefault exoscalev1.ExoscalePostgreSQL
	postgreSQLdefault.APIVersion = exoscaleApiVersion
	postgreSQLdefault.Kind = "ExoscalePostgreSQL"
	postgreSQLdefault.Spec.Parameters.Service.MajorVersion = "15"
	postgreSQLdefault.Spec.Parameters.Service.Zone = "ch-dk-2"
	postgreSQLdefault.Spec.Parameters.Size.Plan = "hobbyist-2 "
	postgreSQLdefault.Spec.Parameters.Backup.TimeOfDay = "12:00:00"
	postgreSQLdefault.Spec.Parameters.Maintenance.DayOfWeek = "Sunday"
	postgreSQLdefault.Spec.Parameters.Maintenance.TimeOfDay = "24:00:00"

	postgreSQLdefault.SetGenerateName("my-postgres-example")
	postgreSQLdefault.SetNamespace("my-namespace")
	return &postgreSQLdefault
}

func getExoscaleRedisDefault() interface{} {
	var redisDefault exoscalev1.ExoscaleRedis
	redisDefault.APIVersion = exoscaleApiVersion
	redisDefault.Kind = "ExoscaleRedis"
	redisDefault.Spec.Parameters.Maintenance.DayOfWeek = "Sunday"
	redisDefault.Spec.Parameters.Maintenance.TimeOfDay = "24:00:00"
	redisDefault.Spec.Parameters.Service.Zone = "ch-dk-2"
	redisDefault.SetGenerateName("redis-app1-prod")
	redisDefault.SetNamespace("prod-app")
	return &redisDefault
}

func getExoscaleKafkaDefault() interface{} {
	var kafkaDefault exoscalev1.ExoscaleKafka
	kafkaDefault.APIVersion = exoscaleApiVersion
	kafkaDefault.Kind = "ExoscaleKafka"
	kafkaDefault.Spec.Parameters.Service.Version = "3.4.0"
	kafkaDefault.Spec.Parameters.Service.Zone = "ch-dk-2"
	kafkaDefault.Spec.Parameters.Size.Plan = "startup-2 "
	kafkaDefault.Spec.Parameters.Maintenance.DayOfWeek = "Sunday"
	kafkaDefault.Spec.Parameters.Maintenance.TimeOfDay = "24:00:00"

	kafkaDefault.SetGenerateName("my-kafka-example ")
	kafkaDefault.SetNamespace("my-namespace")
	return &kafkaDefault
}

func getExoscaleMySQLdefault() interface{} {
	var mySQLdefault exoscalev1.ExoscaleMySQL
	mySQLdefault.APIVersion = exoscaleApiVersion
	mySQLdefault.Kind = "ExoscaleMySQL"
	mySQLdefault.Spec.Parameters.Service.MajorVersion = "8"
	mySQLdefault.Spec.Parameters.Service.Zone = "ch-dk-2"
	mySQLdefault.Spec.Parameters.Size.Plan = "hobbyist-2 "
	mySQLdefault.Spec.Parameters.Backup.TimeOfDay = "12:00:00"
	mySQLdefault.Spec.Parameters.Maintenance.DayOfWeek = "Sunday"
	mySQLdefault.Spec.Parameters.Maintenance.TimeOfDay = "24:00:00"

	mySQLdefault.SetGenerateName("my-mysql-example")
	mySQLdefault.SetNamespace("my-namespace")
	return &mySQLdefault
}

func getExoscaleOpenSearchDefault() interface{} {
	var openSearchDefault exoscalev1.ExoscaleOpenSearch
	openSearchDefault.APIVersion = exoscaleApiVersion
	openSearchDefault.Kind = "ExoscaleOpenSearch"
	openSearchDefault.Spec.Parameters.Service.MajorVersion = "2"
	openSearchDefault.Spec.Parameters.Service.Zone = "ch-dk-2"
	openSearchDefault.Spec.Parameters.Size.Plan = "hobbyist-2 "
	openSearchDefault.Spec.Parameters.Backup.TimeOfDay = "12:00:00"
	openSearchDefault.Spec.Parameters.Maintenance.DayOfWeek = "Sunday"
	openSearchDefault.Spec.Parameters.Maintenance.TimeOfDay = "24:00:00"

	openSearchDefault.SetGenerateName("my-openSearch-example")
	openSearchDefault.SetNamespace("my-namespace")
	return &openSearchDefault
}
