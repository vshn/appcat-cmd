package applications

import (
	"reflect"
	"sort"

	"github.com/vshn/appcat-cli/internal/defaults"
	"github.com/vshn/appcat-cli/internal/util"
	exoscalev1 "github.com/vshn/component-appcat/apis/exoscale/v1"
	vshnv1 "github.com/vshn/component-appcat/apis/vshn/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ExoscaleApiVersion = "exoscale.appcat.vshn.io/v1"
	VshnApiVersion     = "vshn.appcat.vshn.io/v1"
)

var (
	// List of supported apps
	//
	// This is mostly used to generate `AppMap` (via `MakeAppMap`).
	//
	// In order to avoid allocations of (potentially) large types, pass in a nil pointer (`(*type)(nil)`) to `NewApp`.
	Apps = []App{
		// Exoscale
		NewApp(ExoscaleApiVersion, "ExoscalePostgreSQL", (*exoscalev1.ExoscalePostgreSQLSpec)(nil)),
		NewApp(ExoscaleApiVersion, "ExoscaleRedis", (*exoscalev1.ExoscaleRedisSpec)(nil)),
		NewApp(ExoscaleApiVersion, "ExoscaleKafka", (*exoscalev1.ExoscaleKafkaSpec)(nil)),
		NewApp(ExoscaleApiVersion, "ExoscaleMySQL", (*exoscalev1.ExoscaleMySQLSpec)(nil)),
		NewApp(ExoscaleApiVersion, "ExoscaleOpenSearch", (*exoscalev1.ExoscaleOpenSearchSpec)(nil)),

		// VSHN
		NewApp(VshnApiVersion, "VSHNPostgreSQL", (*vshnv1.VSHNPostgreSQLSpec)(nil)),
		NewApp(VshnApiVersion, "VSHNRedis", (*vshnv1.VSHNRedisSpec)(nil)),
	}
)

// App describes a supported AppCat application.
//
// The data in here is used to generate command-line help, and some scaffolding
// around instantiating new applications,
type App struct {
	// Metadata "template" used to instantiate new objects as well as for lookup
	metav1.TypeMeta
	spec reflect.Type
}

func NewApp(apiversion, kind string, spec interface{}) App {
	return App{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiversion,
			Kind:       kind,
		},
		spec: reflect.TypeOf(spec).Elem(),
	}
}

// Call the appropriate "GetDefault" function on the `defaults` package
//
// "GetDefault" functions are expected to be
// * named like `Get<Kind>Default`
// * belong to the `defaults.Defaults` struct
func (app *App) GetDefault() interface{} {
	v := (&defaults.Defaults{}).GetDefaultFor(app.Kind)

	v.Elem().FieldByName("TypeMeta").Set(reflect.ValueOf(app.TypeMeta))

	return v.Interface()
}

// AppMap is a helper type
//
// It is a map of normalized application names to `App` instances.
type AppMap (map[string]App)

func MakeAppMap() AppMap {
	apps := make(AppMap, len(Apps))

	for _, app := range Apps {
		apps[util.NormalizeName(app.Kind)] = app
	}

	return apps
}

// Names returns a sorted list of all known application names
func (m AppMap) Names() []string {
	names := make([]string, len(m))

	i := 0
	for name := range m {
		names[i] = name
		i++
	}

	sort.Strings(names)
	return names
}
