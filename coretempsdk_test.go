// Copyright 2022 Jeremy Edwards
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build windows
// +build windows

package coretempsdk

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed testdata/example.json
	exampleJSON []byte
	//go:embed testdata/example.yaml
	exampleYAML []byte
)

func ExampleGetCoreTempInfo() {
	info, err := GetCoreTempInfo()
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
	fmt.Printf("GetCoreTempInfo: %+v", info)
}

func TestMarshalJSON(t *testing.T) {
	data := &CoreTempInfo{}
	if err := json.Unmarshal(exampleJSON, data); err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff([]int{0, 1, 2, 3}, data.Load); diff != "" {
		t.Errorf("json.Unmarshal() mismatch (-want +got):\n%s", diff)
	}
	m, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(exampleJSON, m); diff != "" {
		t.Logf("json.Marshal\n------------\n%s", string(m))
		t.Errorf("json.Marshal() mismatch (-want +got):\n%s", diff)
	}
}

func TestMarshalYAML(t *testing.T) {
	data := &CoreTempInfo{}
	if err := yaml.Unmarshal(exampleYAML, data); err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff([]int{0, 1, 2, 3}, data.Load); diff != "" {
		t.Errorf("yaml.Unmarshal() mismatch (-want +got):\n%s", diff)
	}
	m, err := yaml.Marshal(data)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(exampleYAML, m); diff != "" {
		t.Logf("yaml.Marshal\n------------\n%s", string(m))
		t.Errorf("yaml.Marshal() mismatch (-want +got):\n%s", diff)
	}
}
