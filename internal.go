// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package confloader

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// flatten takes an interface and extract all of its values and put them in a map.
func flatten(obj interface{}, prefix ...string) Config {
	fields := make(Config)

	var pre string
	if len(prefix) > 0 {
		pre = pre + prefix[0]
	}

	switch obj.(type) {
	case map[interface{}]interface{}:
		for key, value := range obj.(map[interface{}]interface{}) {
			res := flatten(value, pre+key.(string)+".")
			for k, v := range res {
				fields[strings.TrimRight(k, ".")] = v
			}
		}
	case map[string]interface{}:
		for key, value := range obj.(map[string]interface{}) {
			res := flatten(value, pre+key+".")
			for k, v := range res {
				fields[strings.TrimRight(k, ".")] = v
			}
		}
	case []interface{}:
		switch obj.([]interface{})[0].(type) {
		case string:
			arr := make([]string, len(obj.([]interface{})))
			for i, k := range obj.([]interface{}) {
				arr[i] = getEnvValue(k.(string))
			}
			fields[pre] = arr
		case int:
			arr := make([]float64, len(obj.([]interface{})))
			for i, k := range obj.([]interface{}) {
				arr[i] = float64(k.(int))
			}
			fields[pre] = arr
		case float64:
			arr := make([]float64, len(obj.([]interface{})))
			for i, k := range obj.([]interface{}) {
				arr[i] = k.(float64)
			}
			fields[pre] = arr
		case bool:
			arr := make([]bool, len(obj.([]interface{})))
			for i, k := range obj.([]interface{}) {
				arr[i] = k.(bool)
			}
			fields[pre] = arr
		}
		for index, value := range obj.([]interface{}) {
			res := flatten(value, pre+strconv.Itoa(index)+".")
			for k, v := range res {
				fields[strings.TrimRight(k, ".")] = v
			}
		}
	case int:
		fields[strings.TrimRight(pre, ".")] = float64(obj.(int))
	case float64:
		fields[strings.TrimRight(pre, ".")] = obj.(float64)
	case string:
		v := getEnvValue(obj.(string))
		fields[strings.TrimRight(pre, ".")] = v
	case bool:
		fields[strings.TrimRight(pre, ".")] = obj.(bool)
	}

	return fields
}

// unmarshal unmarshals data depending on the format.
func unmarshal(format string, data []byte, v interface{}) error {
	if format == ".json" {
		return json.Unmarshal(data, v)
	} else if format == ".yml" || format == ".yaml" {
		return yaml.Unmarshal(data, v)
	}
	return errors.New("Unrecognized file format  " + format)
}

// get environment variable value if v is in the form ${xxx} or $xxx
func getEnvValue(v string) string {
	if strings.HasPrefix(v, "$") {
		v = strings.Replace(v, "$", "", -1)
		v = strings.Replace(v, "{", "", -1)
		v = strings.Replace(v, "}", "", -1)
		v = os.Getenv(v)
	}
	return v
}
