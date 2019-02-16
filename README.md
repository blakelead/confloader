# Configuration Loader

[![Build Status](https://travis-ci.org/blakelead/confloader.svg?branch=master)](https://travis-ci.org/blakelead/confloader)
[![Coverage Status](https://coveralls.io/repos/github/blakelead/confloader/badge.svg?branch=master)](https://coveralls.io/github/blakelead/confloader?branch=master)
[![Software License](https://img.shields.io/badge/license-MIT-green.svg)](/LICENSE.txt)

## Presentation

**Confloader** is a minimal configuration file loader without any fancy features. It accepts JSON and YAML formats.

## How to use it

The following example uses all of the methods available in **confloader**:

Here are the configuration file that will be used for our example. You can write it in YAML or JSON, whichever you prefer:

```json
{
    "paramString": "foo",
    "paramInt": 42,
    "paramFloat": 42.1,
    "paramBool": true,
    "paramDuration": "20s",
    "paramObj": {
        "paramStringArray": ["foo", "bar", "baz"],
        "paramIntArray": [42, 43, 44],
        "paramFloatArray": [42.2, 43.4, 44.6],
        "paramBoolArray": [true, false, true],
        "paramDurationArray": ["42ns", "5m", "10h10m"]
    },
    "paramEnv": "${ENV_FOO}"
}
```

```yaml
paramString: foo
paramInt: 42
paramFloat: 42.1
paramBool: true
paramDuration: 20s
paramObj:
  paramStringArray: ["foo", "bar", "baz"]
  paramIntArray: [42, 43, 44]
  paramFloatArray: [42.2, 43.4, 44.6]
  paramBoolArray: # you can also write it this way
   - true
   - false
   - true
  paramDurationArray: [42ns, 5m, 10h10m]
paramEnv: ${ENV_FOO}
```

Then, you can use them in your code like so:

```go
package main

import (
    "fmt"
    cl "github.com/blakelead/confloader"
)

func main() {
    config, err := cl.Load("conf.json") // or cl.Load("conf.yml")
    if err != nil {
        panic(err)
    }

    // Examples for all methods provided by the library

    ps := config.GetString("paramString")
    fmt.Println(ps) // foo

    pi := config.GetInt("paramInt")
    fmt.Println(pi) // 42

    pd := config.GetDuration("paramDuration")
    fmt.Println(pd) // 20s

    pf := config.GetFloat("paramFloat")
    fmt.Println(pf) // 42.1

    pb := config.GetBool("paramBool")
    fmt.Println(pb) // true

    psa := config.GetStringArray("paramObj.paramStringArray")
    fmt.Println(psa) // [foo bar baz]

    pia := config.GetIntArray("paramObj.paramIntArray")
    fmt.Println(pia) // [42 43 44]

    pda := config.GetDurationArray("paramObj.paramDurationArray")
    fmt.Println(pda) // [42ns 5m 10h10m0s]

    pfa := config.GetFloatArray("paramObj.paramFloatArray")
    fmt.Println(pfa) // [42.2 43.4 44.6]

    pba := config.GetBoolArray("paramObj.paramBoolArray")
    fmt.Println(pba) // [true false true]

    // It is also possible to access elements in an array with the following syntax

    pa1 := config.GetInt("paramObj.paramIntArray.1")
    fmt.Println(pa1) // 43

    // Some basic conversions are done automatically

    // bool -> int: true=1 false=0
    pbi := config.GetInt("paramObj.paramBoolArray.1")
    fmt.Println(pbi) // 0

    // int -> bool 0=false else=true
    pib := config.GetBool("paramObj.paramIntArray.1")
    fmt.Println(pib) // true

    // float64 -> string
    pfs := config.GetString("paramObj.paramFloatArray.1")
    fmt.Println(pfs) // 43.4

    // []string -> string
    pss := config.GetString("paramObj.paramStringArray")
    fmt.Println(pss) // foo,bar,baz

    // Finally, environment variables can be used inside configuration file,
    // but only for string values

    // export ENV_FOO=fooz
    pe := config.GetString("paramEnv")
    fmt.Println(pe) // fooz
}
```

## Author Information

Adel Abdelhak
