package configtailor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/santiago-labs/configtailor/internal"
)

type ConfigTailor struct {
	Mappings map[string][]string
	Dirs     []string
	RootPath string
}

// Compile takes all the Dirs in the RootPath and generates a config.json for
// every directory along with a config.json per mapping.
func (t *ConfigTailor) Compile() error {
	genPathRoot := path.Join(path.Dir(t.RootPath), path.Base(t.RootPath)+"_generated")
	if err := os.RemoveAll(genPathRoot); err != nil {
		return err
	}

	for _, dir := range t.Dirs {
		mergedConfig, err := t.loadAndMergeConfigs(t.RootPath, dir)
		if err != nil {
			return err
		}

		for k, config := range mergedConfig {
			out, err := json.MarshalIndent(*config.Config, "", "  ")
			if err != nil {
				return err
			}

			genPath := path.Join(genPathRoot, k)
			if err := os.MkdirAll(genPath, 0755); err != nil {
				return err
			}

			if err := ioutil.WriteFile(path.Join(genPath, "config.json"), out, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *ConfigTailor) loadAndMergeConfigs(rootPath, dir string) (map[string]*ConfigWriter, error) {
	mergedConfig := Config{}

	cellToConfig := make(map[string]*ConfigWriter)

	cellToConfig[dir] = &ConfigWriter{
		Config:       &mergedConfig,
		Substitution: make(map[string]string),
	}

	err := filepath.Walk(filepath.Join(rootPath), func(currPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		trimmedPathPrefix := trimPathPrefix(rootPath, currPath)
		if !strings.Contains(trimmedPathPrefix, "/") &&
			info.Name() != dir &&
			info.IsDir() {
			return nil
		}

		if strings.Contains(trimmedPathPrefix, "/") && strings.Split(trimmedPathPrefix, "/")[0] != dir {
			return nil
		}

		// If we find a mapping then we need to branch to create a new Config for each mapping.
		if result, ok := t.Mappings[info.Name()]; ok {
			// We then need to make a copy of the current parent merged config.
			parentKeys := configWriterParent(cellToConfig, trimmedPathPrefix)
			for _, parentKey := range parentKeys {
				for _, r := range result {
					copy := Config(internal.CopyableMap(*cellToConfig[parentKey].Config).DeepCopy())

					relativePath := trimPathPrefix(rootPath, currPath)

					// Replace from backwards to ensure we are matching at the right spot
					// This should also be substituted
					substitutedRelative := strings.Replace(relativePath, info.Name(), r, 1)

					copySubs := internal.MergeMaps(cellToConfig[parentKey].Substitution, map[string]string{
						info.Name(): r,
					})
					copySubs[info.Name()] = r

					cw := &ConfigWriter{
						Config:       &copy,
						Substitution: copySubs,
					}

					cellToConfig[cw.substitutePath(substitutedRelative)] = cw
				}

				delete(cellToConfig, parentKey)
			}
		}

		if !info.IsDir() && filepath.Ext(currPath) == ".json" {
			for _, config := range cellToConfig {
				newMergedConfig := config.Config
				fileBytes, err := ioutil.ReadFile(currPath)
				if err != nil {
					return fmt.Errorf("read file error path: %s, error: %w", currPath, err)
				}

				currentConfig := Config{}
				if err := json.Unmarshal(fileBytes, &currentConfig); err != nil {
					return fmt.Errorf("unmarshal path: (%s) error: %w", currPath, err)
				}

				mergeConfigs(newMergedConfig, currentConfig)
				config.Config.substitute(config.Substitution)
			}
		}
		return nil
	})

	return cellToConfig, err
}
