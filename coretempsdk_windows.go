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
	"fmt"
	"os"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

// https://www.alcpu.com/CoreTemp/developers.html

type coreTempSharedDataEx struct {
	// Original structure (CoreTempSharedData)
	uiLoad      [256]uint32
	uiTjMax     [128]uint32
	uiCoreCnt   uint32
	uiCPUCnt    uint32
	fTemp       [256]float32
	fVID        float32
	fCPUSpeed   float32
	fFSBSpeed   float32
	fMultiplier float32
	sCPUName    [100]byte
	// If ucFahrenheit is true, the temperature is reported in Fahrenheit.
	ucFahrenheit byte
	// If ucDeltaToTjMax is true, the temperature reported represents the distance from TjMax.
	ucDeltaToTjMax byte

	// uiStructVersion = 2

	// If ucTdpSupported is true, processor TDP information in the uiTdp array is valid.
	ucTdpSupported byte
	// If ucPowerSupported is true, processor power consumption information in the fPower array is valid.
	ucPowerSupported byte
	uiStructVersion  uint32
	uiTdp            [128]uint32
	fPower           [128]float32
	fMultipliers     [256]float32
}

var (
	globalFnGetCoreTempInfo *windows.LazyProc
	globalLock              sync.Mutex
)

const (
	dllNameGetCoreTempInfoDLL      = "GetCoreTempInfo.dll"
	dllFuncfnGetCoreTempInfoAlt    = "fnGetCoreTempInfoAlt"
	dllSourceURIGetCoreTempInfoDLL = "https://www.alcpu.com/CoreTemp/developers.html"
)

func GetCoreTempInfo() (*CoreTempInfo, error) {
	rawInfo, err := getCoreTempInfoAlt()
	if err != nil {
		return nil, err
	}

	coreCount := int(rawInfo.uiCoreCnt)

	return &CoreTempInfo{
		Load:         intList(rawInfo.uiLoad[:], coreCount),
		TJMax:        intList(rawInfo.uiTjMax[:], int(rawInfo.uiCPUCnt)),
		CoreCount:    coreCount,
		Temperature:  float32List(rawInfo.fTemp[:], coreCount),
		VID:          rawInfo.fVID,
		CPUSpeed:     rawInfo.fCPUSpeed,
		FSBSpeed:     rawInfo.fFSBSpeed,
		Multiplier:   rawInfo.fMultiplier,
		CPUName:      cleanString(string(rawInfo.sCPUName[:])),
		Fahrenheit:   byteToBool(rawInfo.ucFahrenheit),
		DeltaToTJMax: byteToBool(rawInfo.ucDeltaToTjMax),
	}, nil
}

func byteToBool(b byte) bool {
	return b != 0
}

func intList[T uint32 | int32](input []T, size int) []int {
	result := make([]int, size)
	for i := 0; i < int(size); i++ {
		result[i] = int(input[i])
	}
	return result
}

func float32List[T float32 | float64](input []T, size int) []float32 {
	result := make([]float32, size)
	for i := 0; i < int(size); i++ {
		result[i] = float32(input[i])
	}
	return result
}

func cleanString(input string) string {
	return strings.TrimSpace(strings.Trim(input, string("\x00")))
}

type DLLLoaderError struct {
	Err error
}

func (d DLLLoaderError) Error() string {
	dir, err := os.Getwd()
	if err != nil {
		dir = "."
	}
	return fmt.Sprintf("Make sure that '%s' is in directory '%s'. And the version is at least 1.2.0.0. You can download the DLL from '%s'. Error= %s", dllNameGetCoreTempInfoDLL, dir, dllSourceURIGetCoreTempInfoDLL, d.Err.Error())
}

func getFnGetCoreTempInfo() (*windows.LazyProc, error) {
	globalLock.Lock()
	defer globalLock.Unlock()

	if globalFnGetCoreTempInfo == nil {
		coretempDLL := windows.NewLazyDLL(dllNameGetCoreTempInfoDLL)
		globalFnGetCoreTempInfo = coretempDLL.NewProc(dllFuncfnGetCoreTempInfoAlt)
	}
	if err := globalFnGetCoreTempInfo.Find(); err != nil {
		globalFnGetCoreTempInfo = nil

		return nil, err
	}

	return globalFnGetCoreTempInfo, nil
}

func getCoreTempInfoAlt() (*coreTempSharedDataEx, error) {
	fnGetCoreTempInfo, err := getFnGetCoreTempInfo()
	if err != nil {
		return nil, err
	}

	rawInfo := &coreTempSharedDataEx{}
	r1, _, err := fnGetCoreTempInfo.Call(uintptr(unsafe.Pointer(rawInfo)))

	if r1 != 1 {
		return nil, err
	}
	return rawInfo, nil
}
