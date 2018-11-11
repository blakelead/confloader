// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package confloader

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	generateTestFiles(t)
	defer deleteTestFiles(t)

	os.Setenv("ENV_STRING", "foo")
	os.Setenv("ENV_INT", "42")
	os.Setenv("ENV_FLOAT", "42.3")
	os.Setenv("ENV_BOOL", "true")
	os.Setenv("ENV_ARR_0", "fooz")

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name:    "Load Non-Existent File",
			args:    args{filename: "non-existent-conf.json"},
			want:    Config{},
			wantErr: true,
		}, {
			name:    "Load Unhandled Type File",
			args:    args{filename: "conf.unhandled"},
			want:    Config{},
			wantErr: true,
		}, {
			name:    "Load Invalid JSON File",
			args:    args{filename: "invalid-conf.json"},
			want:    Config{},
			wantErr: true,
		}, {
			name:    "Load Invalid YAML File",
			args:    args{filename: "invalid-conf.yaml"},
			want:    Config{},
			wantErr: true,
		}, {
			name:    "Load Simple JSON File",
			args:    args{filename: "simple-conf.json"},
			want:    Config{"paramString": "foo", "paramInt": 42.0, "paramFloat": 42.1, "paramBool": true},
			wantErr: false,
		}, {
			name:    "Load Simple YAML File",
			args:    args{filename: "simple-conf.yaml"},
			want:    Config{"paramString": "foo", "paramInt": 42.0, "paramFloat": 42.1, "paramBool": true},
			wantErr: false,
		}, {
			name: "Load Complex JSON File",
			args: args{filename: "complex-conf.json"},
			want: Config{
				"paramString": "foo", "paramInt": 42.0, "paramFloat": 42.1, "paramBool": true,
				"paramObj.paramIntArray": []float64{0.0, 1.0, 2.0}, "paramObj.paramIntArray.0": 0.0, "paramObj.paramIntArray.1": 1.0, "paramObj.paramIntArray.2": 2.0,
				"paramObj.paramFloatArray": []float64{0.1, 1.1, 2.1}, "paramObj.paramFloatArray.0": 0.1, "paramObj.paramFloatArray.1": 1.1, "paramObj.paramFloatArray.2": 2.1,
				"paramObj.paramStringArray": []string{"foo", "bar", "baz"}, "paramObj.paramStringArray.0": "foo", "paramObj.paramStringArray.1": "bar", "paramObj.paramStringArray.2": "baz",
				"paramObj.paramBoolArray": []bool{true, false, true}, "paramObj.paramBoolArray.0": true, "paramObj.paramBoolArray.1": false, "paramObj.paramBoolArray.2": true,
			},
			wantErr: false,
		}, {
			name: "Load Complex YAML File",
			args: args{filename: "complex-conf.yaml"},
			want: Config{
				"paramString": "foo", "paramInt": 42.0, "paramFloat": 42.1, "paramBool": true,
				"paramObj.paramIntArray": []float64{0, 1, 2}, "paramObj.paramIntArray.0": 0.0, "paramObj.paramIntArray.1": 1.0, "paramObj.paramIntArray.2": 2.0,
				"paramObj.paramFloatArray": []float64{0.1, 1.1, 2.1}, "paramObj.paramFloatArray.0": 0.1, "paramObj.paramFloatArray.1": 1.1, "paramObj.paramFloatArray.2": 2.1,
				"paramObj.paramStringArray": []string{"foo", "bar", "baz"}, "paramObj.paramStringArray.0": "foo", "paramObj.paramStringArray.1": "bar", "paramObj.paramStringArray.2": "baz",
				"paramObj.paramBoolArray": []bool{true, false, true}, "paramObj.paramBoolArray.0": true, "paramObj.paramBoolArray.1": false, "paramObj.paramBoolArray.2": true,
			},
			wantErr: false,
		}, {
			name: "Load JSON File With Environment Variables",
			args: args{filename: "conf-withenv.json"},
			want: Config{
				"paramString": "foo", "paramInt": "42", "paramFloat": "42.3", "paramBool": "true",
				"paramStringArray": []string{"fooz", "bar", "baz"}, "paramStringArray.0": "fooz", "paramStringArray.1": "bar", "paramStringArray.2": "baz",
			},
			wantErr: false,
		},
		{
			name: "Load YAML File With Environment Variables",
			args: args{filename: "conf-withenv.yaml"},
			want: Config{
				"paramString": "foo", "paramInt": "42", "paramFloat": "42.3", "paramBool": "true",
				"paramStringArray": []string{"fooz", "bar", "baz"}, "paramStringArray.0": "fooz", "paramStringArray.1": "bar", "paramStringArray.2": "baz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetString(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name  string
		c     *Config
		args  args
		wantS string
	}{
		{
			name:  "Get String",
			args:  args{p: "paramString"},
			c:     &Config{"paramString": "foo"},
			wantS: "foo",
		}, {
			name:  "Get String From Int",
			args:  args{p: "paramInt"},
			c:     &Config{"paramInt": 42.0},
			wantS: "42",
		}, {
			name:  "Get String From Float",
			args:  args{p: "paramFloat"},
			c:     &Config{"paramFloat": 42.1},
			wantS: "42.1",
		}, {
			name:  "Get String From Bool",
			args:  args{p: "paramBool"},
			c:     &Config{"paramBool": true},
			wantS: "true",
		}, {
			name:  "Get String From String Array",
			args:  args{p: "paramStringArray"},
			c:     &Config{"paramStringArray": []string{"foo", "bar", "baz"}},
			wantS: "foo,bar,baz",
		}, {
			name:  "Get String From Float Array",
			args:  args{p: "paramFloatArray"},
			c:     &Config{"paramFloatArray": []float64{0.1, 1.1, 2.1}},
			wantS: "0.1,1.1,2.1",
		}, {
			name:  "Get String From Bool Array",
			args:  args{p: "paramBoolArray"},
			c:     &Config{"paramBoolArray": []bool{true, false, true}},
			wantS: "true,false,true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := tt.c.GetString(tt.args.p); gotS != tt.wantS {
				t.Errorf("Config.GetString() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestConfig_GetInt(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name  string
		c     *Config
		args  args
		wantI int
	}{
		{
			name:  "Get Int",
			args:  args{p: "paramInt"},
			c:     &Config{"paramInt": 42.0},
			wantI: 42,
		}, {
			name:  "Get Float",
			args:  args{p: "paramFloat"},
			c:     &Config{"paramFloat": 42.1},
			wantI: 42,
		}, {
			name:  "Get Bool",
			args:  args{p: "paramBool"},
			c:     &Config{"paramBool": true},
			wantI: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := tt.c.GetInt(tt.args.p); gotI != tt.wantI {
				t.Errorf("Config.GetInt() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestConfig_GetFloat(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name  string
		c     *Config
		args  args
		wantF float64
	}{
		{
			name:  "Get Int",
			args:  args{p: "paramInt"},
			c:     &Config{"paramInt": 42.0},
			wantF: 42.0,
		}, {
			name:  "Get Float",
			args:  args{p: "paramFloat"},
			c:     &Config{"paramFloat": 42.1},
			wantF: 42.1,
		}, {
			name:  "Get Bool",
			args:  args{p: "paramBool"},
			c:     &Config{"paramBool": true},
			wantF: 1,
		}, {
			name:  "Get Float From Float Array",
			args:  args{p: "paramFloatArray"},
			c:     &Config{"paramFloatArray": []float64{0.1, 1.1, 2.1}},
			wantF: 0.1,
		}, {
			name:  "Get Float From Bool Array",
			args:  args{p: "paramBoolArray"},
			c:     &Config{"paramBoolArray": []bool{true, false, true}},
			wantF: 1.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotF := tt.c.GetFloat(tt.args.p); gotF != tt.wantF {
				t.Errorf("Config.GetFloat() = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}

func TestConfig_GetBool(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name  string
		c     *Config
		args  args
		wantB bool
	}{
		{
			name:  "Get Boolean",
			args:  args{p: "paramBool"},
			c:     &Config{"paramBool": true},
			wantB: true,
		}, {
			name:  "Get Boolean From Float",
			args:  args{p: "paramFloat"},
			c:     &Config{"paramFloat": 42.1},
			wantB: true,
		}, {
			name:  "Get Bool From Bool Array",
			args:  args{p: "paramBoolArray"},
			c:     &Config{"paramBoolArray": []bool{true, false, true}},
			wantB: true,
		}, {
			name:  "Get Bool From Float Array",
			args:  args{p: "paramFloatArray"},
			c:     &Config{"paramFloatArray": []float64{0.1, 1.1, 2.1}},
			wantB: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotB := tt.c.GetBool(tt.args.p); gotB != tt.wantB {
				t.Errorf("Config.GetBool() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestConfig_GetStringArray(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name  string
		c     *Config
		args  args
		wantA []string
	}{
		{
			name:  "Get String Array",
			args:  args{p: "paramStringArray"},
			c:     &Config{"paramStringArray": []string{"foo", "bar", "baz"}},
			wantA: []string{"foo", "bar", "baz"},
		}, {
			name:  "Get String Array From Float Array",
			args:  args{p: "paramFloatArray"},
			c:     &Config{"paramFloatArray": []float64{0.1, 1.1, 2.1}},
			wantA: []string{"0.1", "1.1", "2.1"},
		}, {
			name:  "Get String Array From Bool Array",
			args:  args{p: "paramBoolArray"},
			c:     &Config{"paramBoolArray": []bool{true, false, true}},
			wantA: []string{"true", "false", "true"},
		}, {
			name:  "Get String Array From String",
			args:  args{p: "paramString"},
			c:     &Config{"paramString": "foo"},
			wantA: []string{"foo"},
		}, {
			name:  "Get String Array From Float",
			args:  args{p: "paramFloat"},
			c:     &Config{"paramFloat": 0.1},
			wantA: []string{"0.1"},
		}, {
			name:  "Get String Array From Bool",
			args:  args{p: "paramBool"},
			c:     &Config{"paramBool": true},
			wantA: []string{"true"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotA := tt.c.GetStringArray(tt.args.p); !reflect.DeepEqual(gotA, tt.wantA) {
				t.Errorf("Config.GetStringArray() = %v, want %v", gotA, tt.wantA)
			}
		})
	}
}

func TestConfig_GetFloatArray(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name  string
		c     *Config
		args  args
		wantA []float64
	}{
		{
			name:  "Get Float Array",
			args:  args{p: "paramFloatArray"},
			c:     &Config{"paramFloatArray": []float64{0.1, 1.1, 2.1}},
			wantA: []float64{0.1, 1.1, 2.1},
		}, {
			name:  "Get Float Array From Bool Array",
			args:  args{p: "paramBoolArray"},
			c:     &Config{"paramBoolArray": []bool{true, false, true}},
			wantA: []float64{1.0, 0.0, 1.0},
		}, {
			name:  "Get Float Array From Float",
			args:  args{p: "paramFloat"},
			c:     &Config{"paramFloat": 42.1},
			wantA: []float64{42.1},
		}, {
			name:  "Get Float Array From Bool (True)",
			args:  args{p: "paramBool"},
			c:     &Config{"paramBool": true},
			wantA: []float64{1.0},
		}, {
			name:  "Get Float Array From Bool (False)",
			args:  args{p: "paramBool"},
			c:     &Config{"paramBool": false},
			wantA: []float64{0.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotA := tt.c.GetFloatArray(tt.args.p); !reflect.DeepEqual(gotA, tt.wantA) {
				t.Errorf("Config.GetFloatArray() = %v, want %v", gotA, tt.wantA)
			}
		})
	}
}

func TestConfig_GetIntArray(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name string
		c    *Config
		args args
		want []int
	}{
		{
			name: "Get Int Array",
			args: args{p: "paramFloatArray"},
			c:    &Config{"paramFloatArray": []float64{0, 1, 2}},
			want: []int{0, 1, 2},
		}, {
			name: "Get Int Array From Float Array",
			args: args{p: "paramFloatArray"},
			c:    &Config{"paramFloatArray": []float64{0.2, 1.4, 2.6}}, // floor casting
			want: []int{0, 1, 2},
		}, {
			name: "Get Int Array From Bool Array",
			args: args{p: "paramBoolArray"},
			c:    &Config{"paramBoolArray": []bool{true, false, true}},
			want: []int{1, 0, 1},
		}, {
			name: "Get Int Array From Float",
			args: args{p: "paramFloat"},
			c:    &Config{"paramFloat": 42.1},
			want: []int{42},
		}, {
			name: "Get Int Array From Bool (True)",
			args: args{p: "paramBool"},
			c:    &Config{"paramBool": true},
			want: []int{1},
		}, {
			name: "Get Int Array From Bool (False)",
			args: args{p: "paramBool"},
			c:    &Config{"paramBool": false},
			want: []int{0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetIntArray(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.GetIntArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetBoolArray(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name  string
		c     *Config
		args  args
		wantA []bool
	}{
		{
			name:  "Get Bool Array",
			args:  args{p: "paramBoolArray"},
			c:     &Config{"paramBoolArray": []bool{true, false, true}},
			wantA: []bool{true, false, true},
		}, {
			name:  "Get Bool Array From Float Array",
			args:  args{p: "paramFloatArray"},
			c:     &Config{"paramFloatArray": []float64{0.0, 1.1, 2.1}},
			wantA: []bool{false, true, true},
		}, {
			name:  "Get Bool Array From Bool",
			args:  args{p: "paramBool"},
			c:     &Config{"paramBool": true},
			wantA: []bool{true},
		}, {
			name:  "Get Bool Array From Float (Non Zero)",
			args:  args{p: "paramFloat"},
			c:     &Config{"paramFloat": 42.1},
			wantA: []bool{true},
		}, {
			name:  "Get Bool Array From Float (Zero)",
			args:  args{p: "paramFloat"},
			c:     &Config{"paramFloat": 0.0},
			wantA: []bool{false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotA := tt.c.GetBoolArray(tt.args.p); !reflect.DeepEqual(gotA, tt.wantA) {
				t.Errorf("Config.GetBoolArray() = %v, want %v", gotA, tt.wantA)
			}
		})
	}
}

func generateTestFiles(t *testing.T) {
	// simple-conf.json
	simpleConfJSON := []byte(`{
	"paramString": "foo",
	"paramInt": 42,
	"paramFloat": 42.1,
	"paramBool": true
}`)
	err := ioutil.WriteFile("simple-conf.json", simpleConfJSON, 0644)
	if err != nil {
		t.Error("Could not generate test file simple-conf.json")
	}

	// simple-conf.yaml
	simpleConfYAML := []byte(`paramString: foo
paramInt: 42
paramFloat: 42.1
paramBool: true`)
	err = ioutil.WriteFile("simple-conf.yaml", simpleConfYAML, 0644)
	if err != nil {
		t.Error("Could not generate test file simple-conf.yaml")
	}

	//complex-conf.json
	complexConfJSON := []byte(`{
    "paramString": "foo",
    "paramInt": 42,
    "paramFloat": 42.1,
    "paramBool": true,
    "paramObj": {
        "paramIntArray": [0, 1, 2],
        "paramFloatArray": [0.1, 1.1, 2.1],
        "paramStringArray": ["foo", "bar", "baz"],
        "paramBoolArray": [true, false, true]
    }
}`)
	err = ioutil.WriteFile("complex-conf.json", complexConfJSON, 0644)
	if err != nil {
		t.Error("Could not generate test file complex-conf.json")
	}

	// complex-conf.yaml
	complexConfYAML := []byte(`{
    "paramString": "foo",
    "paramInt": 42,
    "paramFloat": 42.1,
    "paramBool": true,
    "paramObj": {
        "paramIntArray": [0, 1, 2],
        "paramFloatArray": [0.1, 1.1, 2.1],
        "paramStringArray": ["foo", "bar", "baz"],
        "paramBoolArray": [true, false, true]
    }
}`)
	err = ioutil.WriteFile("complex-conf.yaml", complexConfYAML, 0644)
	if err != nil {
		t.Error("Could not generate test file complex-conf.yaml")
	}

	// invalid-conf.json
	invalidConfJSON := []byte(`{{
    "paramString": "foo",
    "paramInt": 42,
    "paramFloat": 42.1,
    "paramBool": true
}`)
	err = ioutil.WriteFile("invalid-conf.json", invalidConfJSON, 0644)
	if err != nil {
		t.Error("Could not generate test file invalid-conf.json")
	}

	// invalid-conf.yaml
	invalidConfYAML := []byte(`paramString: foo
  paramInt: 42
paramFloat: 42.1
paramBool: true`)
	err = ioutil.WriteFile("invalid-conf.yaml", invalidConfYAML, 0644)
	if err != nil {
		t.Error("Could not generate test file invalid-conf.yaml")
	}

	// conf.unhandled
	unhandledConf := []byte(`{
    "paramString": "foo",
    "paramInt": 42,
    "paramFloat": 42.1,
    "paramBool": true
}`)
	err = ioutil.WriteFile("conf.unhandled", unhandledConf, 0644)
	if err != nil {
		t.Error("Could not generate test file conf.unhandled")
	}

	// conf-withenv.json
	confWithEnvJSON := []byte(`{
    "paramString": "${ENV_STRING}",
    "paramInt": "${ENV_INT}",
    "paramFloat": "${ENV_FLOAT}",
    "paramBool": "${ENV_BOOL}",
    "paramStringArray": ["${ENV_ARR_0}", "bar", "baz"]
}`)
	err = ioutil.WriteFile("conf-withenv.json", confWithEnvJSON, 0644)
	if err != nil {
		t.Error("Could not generate test file conf-withenv.json")
	}

	// conf-withenv.yaml
	confWithEnvYAML := []byte(`paramString: ${ENV_STRING}
paramInt: ${ENV_INT}
paramFloat: ${ENV_FLOAT}
paramBool: ${ENV_BOOL}
paramStringArray: ["${ENV_ARR_0}", "bar", "baz"]`)
	err = ioutil.WriteFile("conf-withenv.yaml", confWithEnvYAML, 0644)
	if err != nil {
		t.Error("Could not generate test file conf-withenv.yaml")
	}
}

func deleteTestFiles(t *testing.T) {
	files := []string{
		"simple-conf.json",
		"simple-conf.yaml",
		"complex-conf.json",
		"complex-conf.yaml",
		"invalid-conf.json",
		"invalid-conf.yaml",
		"conf-withenv.json",
		"conf-withenv.yaml",
		"conf.unhandled",
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			t.Error("Could not delete test file " + file)
		}
	}
}
