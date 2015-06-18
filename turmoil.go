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
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	// Client config
   blacklistString = flag.String("blacklist", "turmoil", "Application names to remove from target lists (separate by comma)")
   marathonURL = flag.String("hostURL", "http://127.0.0.1:8080", "the url for the marathon endpoint")
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
	glog.Info(fmt.Sprintf("Single task: %.2f probability every %.2f", *taskProbability, *taskFrequency))
	glog.Info(fmt.Sprintf("Single application: %.2f probability every %.2f", *appProbability, *appFrequency))
	glog.Info(fmt.Sprintf("All tasks * %.2f: %.2f probability every %.2f", *killFraction, *fractionProbability, *fractionFrequency))

	// Timing
	taskTicker := time.NewTicker(time.Duration(*taskFrequency*3600.0) * time.Second)
  appTicker := time.NewTicker(time.Duration(*appFrequency*3600.0) * time.Second)
  fractionTicker := time.NewTicker(time.Duration(*fractionFrequency*3600.0) * time.Second)

	for {
		select {
			case <- taskTicker.C:
				rand.Seed(time.Now().UnixNano())
				glog.Info("Attempting to kill a random task")
				if (rand.Float64()<=*taskProbability) {
					KillRandomTask(client, blacklist)
					glog.Info("Killed a random task")
				} else {
					glog.Info("Did not kill a random task")
				}
			case <- appTicker.C:
				rand.Seed(time.Now().UnixNano())
				glog.Info("Attempting to kill a random application")
				if (rand.Float64()<=*appProbability) {
					KillRandomApp(client, blacklist)
					glog.Info("Killed a random application")
				} else {
          glog.Info("Did not kill a random application")
        }
			case <- fractionTicker.C:
				rand.Seed(time.Now().UnixNano())
				glog.Info("Attempting to kill a fraction of running tasks")
				if (rand.Float64()<=*fractionProbability) {
					KillTaskFraction(client, blacklist, *killFraction)
					glog.Info("Killed a fraction of running tasks")
				} else {
          glog.Info("Did not kill a fraction of running tasks")
        }
		}
	}

}
