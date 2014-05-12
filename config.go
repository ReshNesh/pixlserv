package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v1"
)

const (
	defaultThrottlingRate             = 60 // Requests per min
	defaultAllowCustomTransformations = true
	defaultAllowCustomScale           = true
	defaultLocalPath                  = "local-images"
)

type Config struct {
	throttlingRate                               int
	allowCustomTransformations, allowCustomScale bool
	localPath                                    string
	transformations                              map[string]Params
}

func configInit(configFilePath string) (Config, error) {
	config := Config{defaultThrottlingRate, defaultAllowCustomTransformations, defaultAllowCustomScale, defaultLocalPath, make(map[string]Params)}

	if configFilePath == "" {
		return config, nil
	}

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return config, err
	}

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(data), &m)

	throttlingRate, ok := m["throttling-rate"].(int)
	if ok && throttlingRate >= 0 {
		config.throttlingRate = throttlingRate
	}

	allowCustomTransformations, ok := m["allow-custom-transformations"].(bool)
	if ok {
		config.allowCustomTransformations = allowCustomTransformations
	}

	allowCustomScale, ok := m["allow-custom-scale"].(bool)
	if ok {
		config.allowCustomScale = allowCustomScale
	}

	localPath, ok := m["local-path"].(string)
	if ok {
		config.localPath = localPath
	}

	transformations, ok := m["transformations"].([]interface{})
	if ok {
		for _, transformationInterface := range transformations {
			transformation, ok := transformationInterface.(map[interface{}]interface{})
			if ok {
				parametersStr, ok := transformation["parameters"].(string)
				if ok {
					params, err := parseParameters(parametersStr)
					if err != nil {
						return config, fmt.Errorf("invalid transformation parameters: %s (%s)", parametersStr, err)
					}
					name, ok := transformation["name"].(string)
					if ok {
						config.transformations[name] = params
					}
				}
			}
		}
	}

	return config, nil
}
