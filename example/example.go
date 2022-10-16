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

// Package main is an example for using the Core Temp SDK.
package main

import (
	"log"

	"github.com/jeremyje/coretempsdk-go"
)

func main() {
	info, err := coretempsdk.GetCoreTempInfo()
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}
	log.Printf("CPU: %s", info.CPUName)
	log.Printf("Temperatures: %v", info.TemperatureCelcius)
	log.Printf("Full: %+v", info)
}
