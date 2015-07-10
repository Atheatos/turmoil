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
	"flag"
	"fmt"
	marathon "github.com/gambol99/go-marathon"
	"github.com/golang/glog"
	"github.com/vharitonsky/iniflags"
	"os"
	"strings"
)

var (
	// Client config
	blacklistString = flag.String("blacklist", "turmoil", "Application names to remove from target lists (separate by comma)")
	marathonURL     = flag.String("hostURL", "http://127.0.0.1:8080", "the url for the marathon endpoint")
	runStart        = flag.String("start", "10:00", "The start time for Turmoil")
	runStop         = flag.String("stop", "16:00", "The stop time for Turmoil")
	// Kill one task
	taskFrequency   = flag.Float64("taskFrequency", 0.1, "Number of hours between attempts to kill a single random task")
	taskProbability = flag.Float64("taskProbability", 0.5, "Probability that a single task kill attempt succeeds")
	// Kill one application
	appFrequency   = flag.Float64("appFrequency", 0.5, "Number of hours between attempts to kill a single random application")
	appProbability = flag.Float64("appProbability", 0.2, "Probability that a single task kill attempt succeeds")
	// Kill a fraction of tasks
	killFraction        = flag.Float64("fraction", 0.25, " of tasks to be killed (e.g. 0.25 kills 25 percent of all tasks)")
	fractionFrequency   = flag.Float64("fractionFrequency", 0.5, "Number of hours between attempts to kill a fraction of tasks at random")
	fractionProbability = flag.Float64("fractionProbability", 0.2, "Probability that a single task kill attempt succeeds")
	hostFrequency       = flag.Float64("hostFrequency", 2, "Number of hours between attempts to kill all tasks on a random host")
	hostProbability     = flag.Float64("hostProbability", 0.25, "Probability that a host tasks kill attempt succeeds")

	blacklist []string
	client    marathon.Marathon
)

func main() {
	// Parse config and initialize the client (Marathon interface)
	var err error
	if _, err = os.Stat("/mnt/mesos/sandbox/params.ini"); err == nil {
		iniflags.SetConfigFile("/mnt/mesos/sandbox/params.ini")
	} else {
		iniflags.SetConfigFile("/params.ini")
	}
	iniflags.Parse()
	glog.Warningln(err)
	if err != nil {
		glog.Info("Using default configuration")
	} else {
		glog.Info("Using custom configuration")
	}
	blacklist = strings.Split(*blacklistString, ",")
	config := marathon.NewDefaultConfig()
	config.URL = *marathonURL
	client, _ = marathon.NewClient(config)

	// Log
	glog.Info("Kill settings:")
	glog.Info(fmt.Sprintf(" | Single task: %.2f probability every %.2f hour(s)", *taskProbability, *taskFrequency))
	glog.Info(fmt.Sprintf(" | Single application: %.2f probability every %.2f hour(s)", *appProbability, *appFrequency))
	glog.Info(fmt.Sprintf(" | Task Fraction (%.2f): %.2f probability every %.2f hour(s)", *killFraction, *fractionProbability, *fractionFrequency))
	glog.Info(fmt.Sprintf(" | Single host: %.2f probability every %.2f hour(s)", *hostProbability, *hostFrequency))

	// Timing
	quitChans := []chan int{make(chan int), make(chan int), make(chan int), make(chan int)}
	start := ParseTime(*runStart)
	stop := ParseTime(*runStop)
	RunScheduler(start, stop, quitChans)
}

// Assert an error, if any
func Assert(err error) {
	if err != nil {
		glog.Fatalf("Failed, error: %s", err)
	}
}
