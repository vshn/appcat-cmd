package defaults

import (
	"fmt"
	"strings"
)

// Selects the appropriate service type and initializes it with default values
func FindServiceType(parsedType string) (interface{}, error) {
	parsedType = strings.ToLower(parsedType)
	if strings.Contains(parsedType, "vshn") {
		getDefault, ok := VSHN_TYPES[parsedType]
		if !ok {
			return nil, fmt.Errorf("vshn does not support service Type %s", parsedType)
		}
		return getDefault(), nil
	} else if strings.Contains(parsedType, "exoscale") {
		getDefault, ok := EXOSCALE_TYPES[parsedType]
		if !ok {
			return nil, fmt.Errorf("exoscale does not support service Type %s", parsedType)
		}
		return getDefault(), nil
	} else {
		return nil, fmt.Errorf("service Type %s could not be found", parsedType)
	}
}
