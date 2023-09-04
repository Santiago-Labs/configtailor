package configtailor

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func ParentDirs(rootPath string) ([]string, error) {
	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}
	return dirs, nil
}

func Mappings(mappingsFlag string) (map[string][]string, error) {
	if mappingsFlag == "" {
		return nil, nil
	}

	mappings := make(map[string][]string)
	mapSplit := strings.Split(mappingsFlag, ":")
	if len(mapSplit)%2 != 0 {
		return nil, fmt.Errorf("mappings must be in the format (key:value1,value2:key2:value1,value2) e.g. (cell:us1,us2,us3,us4:env:prod,dev,test)")
	}

	for i := 0; i < len(mapSplit); i += 2 {
		mappings[mapSplit[i]] = strings.Split(mapSplit[i+1], ",")
	}

	return mappings, nil
}
