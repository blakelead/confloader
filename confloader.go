// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package confloader is a simple configuration file loader that
// I use for my personal Go projects.
// It accepts JSON and YAML formats.
package confloader

import (
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"time"
)

// Config is a map of parameters in a configuration file.
// The key is the absolute path of the parameter.
// The value is the raw value of the parameter.
type Config map[string]interface{}

// Load loads a configuration file and returns a Config object, or an error
// if file could not be read or unmarshalled, or if the file doesn't exist.
func Load(filename string) (Config, error) {
	blob, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	var raw interface{}
	err = unmarshal(path.Ext(filename), blob, &raw)
	if err != nil {
		return Config{}, err
	}

	cnf := flatten(raw)

	// for k, v := range cnf {
	// 	fmt.Println(k, v, reflect.TypeOf(v))
	// }

	return cnf, nil
}

// Get gets value of parameter p. p should be the absolute path to the parameter.
// Example: { "param1": { "param2": 3.14 } }; to access param2, p should be
// "param1.param2".
func (c *Config) Get(p string) interface{} {
	return (*c)[p]
}

// GetString gets string value of parameter p. Uses Get().
// If parameter is a number, the number is converted to a string.
// If parameter is a boolean, the string will be "true" or "false".
func (c *Config) GetString(p string) (s string) {
	switch v := c.Get(p).(type) {
	case string:
		s = v
	case float64:
		s = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		s = strconv.FormatBool(v)
	case []string:
		s = strings.Join(v, ",")
	case []float64:
		arr := make([]string, len(v))
		for i, k := range v {
			arr[i] = strconv.FormatFloat(k, 'f', -1, 64)
		}
		s = strings.Join(arr, ",")
	case []bool:
		arr := make([]string, len(v))
		for i, k := range v {
			arr[i] = strconv.FormatBool(k)
		}
		s = strings.Join(arr, ",")
	}
	return s
}

// GetFloat gets float value of parameter p. Uses Get().
// If parameter is a boolean, the number will be 1.0 if true, 0.0 if false.
func (c *Config) GetFloat(p string) (f float64) {
	switch v := c.Get(p).(type) {
	case float64:
		f = v
	case bool:
		if v {
			f = 1.0
		}
	case []float64:
		if len(v) > 0 {
			f = v[0]
		}
	case []bool:
		if len(v) > 0 && v[0] {
			f = 1.0
		}
	}
	return f
}

// GetInt gets int value of parameter p.
func (c *Config) GetInt(p string) int {
	return int(c.GetFloat(p))
}

// GetDuration gets duration value of parameter p. p can have
// suffixes like s, ms, h, etc. In fact the same as standard time.ParseDuration().
func (c *Config) GetDuration(p string) (d time.Duration) {
	d, _ = time.ParseDuration(c.GetString(p))
	return d
}

// GetBool gets number value of parameter p. Uses Get().
// If parameter is a number, the boolean will be true if parameter is not 0,
// false otherwise.
func (c *Config) GetBool(p string) (b bool) {
	switch v := c.Get(p).(type) {
	case bool:
		b = v
	case float64:
		if v != 0 {
			b = true
		}
	case []bool:
		if len(v) > 0 {
			b = v[0]
		}
	case []float64:
		if len(v) > 0 && v[0] != 0 {
			b = true
		}
	}
	return b
}

// GetStringArray gets a string slice from parameter p.
func (c *Config) GetStringArray(p string) (a []string) {
	switch v := c.Get(p).(type) {
	case []string:
		a = v
	case []float64:
		arr := make([]string, len(v))
		for i, k := range v {
			arr[i] = strconv.FormatFloat(k, 'f', -1, 64)
		}
		a = arr
	case []bool:
		arr := make([]string, len(v))
		for i, k := range v {
			arr[i] = strconv.FormatBool(k)
		}
		a = arr
	case string:
		a = []string{v}
	case float64:
		a = []string{strconv.FormatFloat(v, 'f', -1, 64)}
	case bool:
		a = []string{strconv.FormatBool(v)}
	}
	return a
}

// GetFloatArray gets a float64 slice from parameter p.
func (c *Config) GetFloatArray(p string) (a []float64) {
	switch v := c.Get(p).(type) {
	case []float64:
		a = v
	case []bool:
		arr := make([]float64, len(v))
		for i, k := range v {
			if k {
				arr[i] = 1.0
			} else {
				arr[i] = 0.0
			}
		}
		a = arr
	case float64:
		a = []float64{v}
	case bool:
		if v {
			a = []float64{1.0}
		} else {
			a = []float64{0.0}
		}
	}
	return a
}

// GetIntArray gets a int slice from parameter p.
func (c *Config) GetIntArray(p string) []int {
	arr := c.GetFloatArray(p)
	a := make([]int, len(arr))
	for i, k := range arr {
		a[i] = int(k)
	}
	return a
}

// GetDurationArray gets a duration slice from parameter p.
func (c *Config) GetDurationArray(p string) []time.Duration {
	arr := c.GetStringArray(p)
	a := make([]time.Duration, len(arr))
	for i, k := range arr {
		a[i], _ = time.ParseDuration(k)
	}
	return a
}

// GetBoolArray gets a bool slice from parameter p.
func (c *Config) GetBoolArray(p string) (a []bool) {
	switch v := c.Get(p).(type) {
	case []bool:
		a = v
	case []float64:
		arr := make([]bool, len(v))
		for i, k := range v {
			if k != 0 {
				arr[i] = true
			} else {
				arr[i] = false
			}
		}
		a = arr
	case bool:
		a = []bool{v}
	case float64:
		if v != 0 {
			a = []bool{true}
		} else {
			a = []bool{false}
		}
	}
	return a
}
