package configtailor

import (
	"fmt"
	"strings"
)

type Config map[string]interface{}

// ConfigWriter is a wrapper around a config that keeps track of what
// Substitutions need to be performed.
type ConfigWriter struct {
	Config *Config

	Substitution map[string]string
}

func (c *Config) substitute(substitutions map[string]string) {
	for key := range *c {
		for subKey, subValue := range substitutions {
			if existingValue, exists := (*c)[key]; exists {
				switch t := existingValue.(type) {
				case map[string]interface{}:
					if existingMap, ok := existingValue.(map[string]interface{}); ok {
						newBase := Config(existingMap)
						newBasePointer := &newBase
						newBasePointer.substitute(substitutions)
						(*c)[key] = newBase
						continue
					}

				case string:
					(*c)[key] = strings.ReplaceAll(t, fmt.Sprintf("$%s", subKey), subValue)
				}
			}
		}
	}
}

func (c ConfigWriter) substitutePath(path string) string {
	for k, v := range c.Substitution {
		path = strings.ReplaceAll(path, k, v)
	}

	return path
}

// configWriterParent finds the ConfigWriter based on the parent path.
// For example,
// If you have a config file with config/region/cell/config.json and a mapping for:
// region: us-west-2, eu-west-1
// cell: us000, us001, eu001
// Then when handing a path of config/us-west-2/us000 we will use the region specific mapping.
func configWriterParent(writers map[string]*ConfigWriter, matchPath string) []string {
	var res []string
	for k, v := range writers {
		modifiedPath := v.substitutePath(matchPath)
		if strings.HasPrefix(modifiedPath, k) && modifiedPath != k {
			res = append(res, v.substitutePath(k))
		}
	}

	return res
}

// mergeConfigs takes a base Config and modifies it to take in the overlay
// Config.
func mergeConfigs(base *Config, overlay Config) {
	for key, value := range overlay {
		if existingValue, exists := (*base)[key]; exists {
			switch t := value.(type) {
			// If the key exists in the base and both the base and the overlay have
			// map values, then merge the maps.
			case map[string]interface{}:
				if existingMap, ok := existingValue.(map[string]interface{}); ok {
					newBase := Config(existingMap)
					mergeConfigs(&newBase, t)
					(*base)[key] = newBase
					continue
				}
			}
		}

		(*base)[key] = value
	}
}
