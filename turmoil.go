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
	"fmt"
	"flag"
	marathon "github.com/gambol99/go-marathon"
	"github.com/golang/glog"
	"github.com/vharitonsky/iniflags"
	"os"
	"strings"
)

var (
	// Client config
	blacklistString = flag.String("blacklist", "turmoil", "Application names to remove from target lists (separate by comma)")
	marathonURL = flag.String("hostURL", "http://127.0.0.1:8080", "the url for the marathon endpoint")
	runStart = flag.Int("start", 1000, "The start time for Turmoil")
	runStop = flag.Int("stop", 1600, "The stop time for Turmoil")
	// Kill one task
	taskFrequency = flag.Float64("taskFrequency", 0.1, "Number of hours between attempts to kill a single random task")
	taskProbability = flag.Float64("taskProbability", 0.5, "Probability that a single task kill attempt succeeds")
	// Kill one application
	appFrequency = flag.Float64("appFrequency", 0.5, "Number of hours between attempts to kill a single random application")
	appProbability = flag.Float64("appProbability", 0.2, "Probability that a single task kill attempt succeeds")
	// Kill a fraction of tasks
	killFraction = flag.Float64("fraction", 0.25, " of tasks to be killed (e.g. 0.25 kills 25 percent of all tasks)")
	fractionFrequency = flag.Float64("fractionFrequency", 0.5, "Number of hours between attempts to kill a fraction of tasks at random")
	fractionProbability = flag.Float64("fractionProbability", 0.2, "Probability that a single task kill attempt succeeds")
)

func main() {
	// Parse config and initialize the client (Marathon interface)
	iniflags.Parse()
	blacklist := strings.Split(*blacklistString, ",")
	config := marathon.NewDefaultConfig()
	config.URL = *marathonURL
	config.LogOutput = os.Stdout
	client, err := marathon.NewClient(config)
	Assert(err)

	// Log
	glog.Info(fmt.Sprintf("Single task: %.2f probability every %.2f hour(s)", *taskProbability, *taskFrequency))
	glog.Info(fmt.Sprintf("Single application: %.2f probability every %.2f hour(s)", *appProbability, *appFrequency))
	glog.Info(fmt.Sprintf("All tasks * %.2f: %.2f probability every %.2f hour(s)", *killFraction, *fractionProbability, *fractionFrequency))

	// Timing
	taskQuit := make(chan int)
	appQuit := make(chan int)
	fractionQuit := make(chan int)
	for {
		startTimers(client, blacklist, taskQuit, appQuit, fractionQuit)
	}
	<-taskQuit
	<-appQuit
	<-fractionQuit
}

func startTimers(client marathon.Marathon, blacklist []string, task, app, frac chan int) {
	if (*taskFrequency*3600.0 >= 1) {
		go taskTimer(client, blacklist, task)
	}
	if (*appFrequency*3600.0 >= 1) {
		go appTimer(client, blacklist, app)
	}
	if (*fractionFrequency*3600.0 >= 1) {
		go fractionTimer(client, blacklist, frac)
	}
}

func stopTimers(task, app, frac chan int) {
	task <- 1
	app  <- 1
	frac <- 1
}
