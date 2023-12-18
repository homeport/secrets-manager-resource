// Copyright © 2023 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package smr

import (
	"encoding/json"
	"io"
	"time"
)

type Version struct {
	Timestamp time.Time `json:"timestamp"`
}

type Config struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InConfig struct {
	Config
	Params map[string]string `json:"params"`
}

type Source struct {
	EndpointURL   string `json:"endpointURL"`
	ApiKey        string `json:"apikey"`
	SecretName    string `json:"secretName"`
	SecretGroupID string `json:"secretGroupID"`
}

type CheckResult []Version

type KeyVal struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Output struct {
	Version  Version  `json:"version"`
	Metadata []KeyVal `json:"metadata,omitempty"`
}

func LoadConfig(in io.Reader) (*Config, error) {
	data, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadInConfig(in io.Reader) (*InConfig, error) {
	data, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}

	var inConfig InConfig
	if err := json.Unmarshal(data, &inConfig); err != nil {
		return nil, err
	}

	return &inConfig, nil
}
