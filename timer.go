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
	marathon "github.com/gambol99/go-marathon"
	"github.com/golang/glog"
	"math/rand"
	"time"
)

/* Time/schedule task deletion
 *  client: Marthon interface
 *  blacklist: List of apps to protect
 *  quit: channel for telling fractionTimer to stop
 */
func taskTimer(client marathon.Marathon, blacklist []string, quit chan int) {
	glog.Info("taskTimer is running")
	stop := make(chan int)
	ticker := time.NewTicker(time.Duration(*taskFrequency*3600.0) * time.Second)
	for {
		select {
			case <- ticker.C:
				rand.Seed(time.Now().UnixNano())
				glog.Info("Attempting to kill a random task")
				if (rand.Float64() <= *taskProbability) {
					glog.Info("Killed task: %s", KillRandomTask(client, blacklist))
				} else {
					glog.Info("Did not kill a task")
				}
			case <- quit:
				ticker.Stop()
				close(stop)
				return
		}
	}
}

/* Time/schedule application deletion
 *  client: Marthon interface
 *  blacklist: List of apps to protect
 *  quit: channel for telling the timer to stop
 */
func appTimer(client marathon.Marathon, blacklist []string, quit chan int) {
	glog.Info("appTimer is running")
	stop := make(chan int)
	ticker := time.NewTicker(time.Duration(*appFrequency*3600.0) * time.Second)
	for {
		select {
			case <- ticker.C:
				rand.Seed(time.Now().UnixNano())
				glog.Info("Attempting to kill a random application")
				if (rand.Float64() <= *appProbability) {
					glog.Info("Kill application: %s", KillRandomApp(client, blacklist))
				} else {
					glog.Info("Did not kill an application")
				}
			case <- quit:
				ticker.Stop()
				close(stop)
				return
		}
	}
}

/* Time/schedule task fraction deletion
 *  client: Marthon interface
 *  blacklist: List of apps to protect
 *  quit: channel for telling the timer to stop
 */
func fractionTimer(client marathon.Marathon, blacklist []string, quit chan int) {
	glog.Info("fractionTimer is running")
	stop := make(chan int)
	ticker := time.NewTicker(time.Duration(*fractionFrequency*3600.0) * time.Second)
	for {
		select {
			case <- ticker.C:
				rand.Seed(time.Now().UnixNano())
				glog.Info("Attempting to kill a fraction of running tasks")
				if (rand.Float64() <= *fractionProbability) {
					victims := KillTaskFraction(client, blacklist, *killFraction)
					glog.Info("Killed %d tasks:\n %#v", len(victims), victims)
				} else {
					glog.Info("Did not kill any tasks")
				}
			case <- quit:
				ticker.Stop()
				close(stop)
				return
		}
	}
}
