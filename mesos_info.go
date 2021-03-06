/*
The MIT License (MIT)

Copyright (c) 2015 Atheatos

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func ListSlaves() []string {
	// Prepare request
	slavelist := make([]string, 0)
	leader := MesosLeader()
	// Make request
	resp, err := http.Get(fmt.Sprintf("http://%s/slaves", leader))
	Assert(err)
	defer resp.Body.Close()
	// Process response
	var slaves Slaves
	Assert(json.NewDecoder(resp.Body).Decode(&slaves))
	for _, slave := range slaves.Slaves {
		slavelist = append(slavelist, slave.Hostname)
	}
	return slavelist
}

func MesosLeader() string {
	leader, err := client.Leader()
	Assert(err)
	stateURL := fmt.Sprintf("http://%s:5050/master/state.json", strings.Split(leader, ":")[0])
	resp, err := http.Get(stateURL)
	Assert(err)
	defer resp.Body.Close()
	var state MesosState
	Assert(json.NewDecoder(resp.Body).Decode(&state))
	return strings.Split(state.Leader, "master@")[1]
}

type Slaves struct {
	Slaves []Slave `json:"slaves"`
}

type Slave struct {
	Active         bool       `json:"active",omitempty`
	Attributes     *Attribute `json:"attributes",omitempty`
	Hostname       string     `json:"hostname",omitempty`
	ID             string     `json:"id",omitempty`
	PID            string     `json:"pid"omitempty`
	RegisteredTime float64    `json:"registered_time",omitemtpy`
	Resources      *Resource  `json:"resources",omitempty`
}

type Resource struct {
	CPUs  int    `json:"cpus",omitempty`
	Disk  int    `json:"disk"`
	Mem   int    `json:"mem"`
	Ports string `json:"ports",omitemtpy`
}

type Attribute struct{}

type MesosState struct {
	Leader string `json:"leader",omitempty`
}
