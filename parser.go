package waypoint

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strconv"
	"strings"
)

func ParseYaml(reader io.Reader) (interface{}, error) {
	var m map[string]interface{}

	decoder := yaml.NewDecoder(reader)

	err := decoder.Decode(&m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func ReadFile(path string) (interface{}, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, os.ErrNotExist
	}

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	return ParseYaml(file)
}

func PathLookup(path string, m interface{}) (string, error) {

	keys := strings.Split(path, "/")

	if len(keys) > 0 && keys[len(keys)-1] == "" {
		keys = keys[:len(keys)-1]
	}

	if len(keys) > 0 && keys[0] == "" {
		keys = keys[1:]
	}

	if len(keys) == 0 {
		return "", os.ErrNotExist
	}

	return lookup(keys, m)
}

func lookup(key []string, lookupTree interface{}) (string, error) {
	if _, err := strconv.Atoi(key[0]); err == nil {
		arr, ok := lookupTree.([]interface{})

		if !ok {
			return "", os.ErrNotExist
		}

		return indexLookup(key, arr)
	} else {
		m, ok := lookupTree.(map[string]interface{})

		if !ok {
			return "", os.ErrNotExist
		}

		return namedLookup(key, m)
	}
}

func namedLookup(keys []string, lookupTree map[string]interface{}) (string, error) {
	key := keys[0]

	subMap, ok := lookupTree[key]

	if !ok {
		return "", os.ErrNotExist
	}

	if len(keys) == 1 {
		return subMap.(string), nil
	} else {
		return lookup(keys[1:], subMap)
	}
}

func indexLookup(keys []string, lookupTree []interface{}) (string, error) {
	index, err := strconv.Atoi(keys[0])

	if err != nil {
		return "", err
	}

	if index >= len(lookupTree) {
		return "", os.ErrNotExist
	}

	subMap := lookupTree[index]

	if len(keys) == 1 {
		return subMap.(string), nil
	} else {
		return lookup(keys[1:], subMap)
	}
}
