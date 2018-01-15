// Copyright (c) 2017 Tobias Kohlbau
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config contains the configuration of the just another youtube extractor.
type Config struct {
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
	Youtube struct {
		URL       string `json:"url"`
		Token     string `json:"token"`
		VideoPath string `json:"video_path"`
	} `json:"youtube"`
}

// FromFile returns a configuration parsed from the given file.
func FromFile(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
