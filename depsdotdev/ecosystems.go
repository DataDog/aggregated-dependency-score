package depsdotdev

import (
	"fmt"

	api "deps.dev/api/v3"
)

func depsdotdevEcosystem(x string) (api.System, error) {
	switch x {
	case "pypi":
		return api.System_PYPI, nil
	case "npm":
		return api.System_NPM, nil
	case "maven":
		return api.System_MAVEN, nil
	case "cargo":
		return api.System_CARGO, nil
	case "go":
		return api.System_GO, nil
	default:
		return api.System_SYSTEM_UNSPECIFIED, fmt.Errorf("unknown ecosystem: %q", x)
	}
}

func depsdotdevEcosystemString(x api.System) (string, error) {
	switch x {
	case api.System_PYPI:
		return "pypi", nil
	case api.System_NPM:
		return "npm", nil
	case api.System_MAVEN:
		return "maven", nil
	case api.System_CARGO:
		return "cargo", nil
	case api.System_GO:
		return "go", nil
	default:
		return "", fmt.Errorf("unknown ecosystem: %v", x)
	}
}
