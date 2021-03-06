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
	"github.com/golang/glog"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	H24 = time.Duration(24) * time.Hour
)

/* Starts and stops timers based on the given time constraints.
 *  start:  when to start running timers
 *  stop:   when to stop running timers
 *  quit:   quit channels for timers
 */
func RunScheduler(start, stop time.Duration, quit []chan int) {
	StabilizeTiming(start, stop, quit)
	if start > stop {
		// overnight
		for {
			StopTimers(quit)
			time.Sleep(start - CurrentTime())
			StartTimers(quit)
			time.Sleep(H24 - CurrentTime() + stop)
		}
	} else {
		// day
		for {
			StopTimers(quit)
			time.Sleep(H24 - CurrentTime() + start)
			StartTimers(quit)
			time.Sleep(stop - CurrentTime())
		}
	}
}

/* Aligns timing so that scheduling can be done in a single loop.
 * 	start: 	when to start running timers
 * 	stop: 	when to stop running timers
 * 	quit: 	quit channels for timers
 */
func StabilizeTiming(start, stop time.Duration, quit []chan int) {
	now := CurrentTime()
	if start > stop {
		glog.Info("Run overnight between ", *runStart, " and ", *runStop)
		switch {
		case (now > start):
			StartTimers(quit)
			time.Sleep(H24 + stop - CurrentTime())
		case (now < stop):
			StartTimers(quit)
			time.Sleep(stop - CurrentTime())
		default:
			glog.Info("Waiting until ", *runStart, " to start")
			time.Sleep(start - CurrentTime())
			StartTimers(quit)
			time.Sleep(H24 - CurrentTime() + stop)
		}
	} else {
		glog.Info("Run daily between ", *runStart, " and ", *runStop)
		// same day
		switch {
		case (now < start):
			glog.Info("Waiting until ", *runStart, " to start")
			time.Sleep(start - CurrentTime())
			StartTimers(quit)
			time.Sleep(stop - CurrentTime())
		case (now > stop):
			glog.Info("Waiting until ", *runStart, " to start")
			time.Sleep(H24 - CurrentTime() + start)
			StartTimers(quit)
			time.Sleep(stop - CurrentTime())
		default:
			StartTimers(quit)
			time.Sleep(stop - CurrentTime())
		}
	}
}

/* Start the kill timers in goroutines if frequency is above the precision threshold (one second)
 * 	task: 	quit channel for TaskTimer
 *  app:  quit channel for AppTimer
 *  frac:   quit channel for FractionTimer
 */
func StartTimers(quit []chan int) {
	glog.Info("Starting timers; Running until ", *runStop)
	if *taskFrequency*3600.0 >= 1 {
		go TaskTimer(quit[0])
		glog.Info(" | TaskTimer started")
	}
	if *appFrequency*3600.0 >= 1 {
		go AppTimer(quit[1])
		glog.Info(" | AppTimer started")
	}
	if *fractionFrequency*3600.0 >= 1 {
		go FractionTimer(quit[2])
		glog.Info(" | FractionTimer started")
	}
	if *hostFrequency*3600.0 >= 1 {
		go HostTimer(quit[3])
		glog.Info(" | HostTimer started")
	}
}

/* Send a stop signal
 * 	task: 	quit channel for TaskTimer
 *  app:  quit channel for AppTimer
 *  frac:   quit channel for FractionTimer
 */
func StopTimers(quit []chan int) {
	for _, q := range quit {
		q <- 1
	}
	glog.Info("All timers stopped; Restart at ", *runStart)
}

/* Time/schedule task deletion
 *  client: Marthon interface
 *  blacklist: List of apps to protect
 *  quit: channel for telling FractionTimer to stop
 */
func TaskTimer(quit chan int) {
	glog.Info("TaskTimer is running")
	stop := make(chan int)
	ticker := time.NewTicker(time.Duration(*taskFrequency*3600.0) * time.Second)
	for {
		select {
		case <-ticker.C:
			rand.Seed(time.Now().UnixNano())
			glog.Info("TaskTimer -- Attempting to kill a random task")
			if rand.Float64() <= *taskProbability {
				victim := KillRandomTask()
				glog.Info("TaskTimer -- Killed task: ", victim)
			} else {
				glog.Info("TaskTimer -- Did not kill a task")
			}
		case <-quit:
			ticker.Stop()
			close(stop)
			return
		}
	}
}

/* Time/schedule application deletion
 *  quit: channel for telling the timer to stop
 */
func AppTimer(quit chan int) {
	glog.Info("AppTimer is running")
	stop := make(chan int)
	ticker := time.NewTicker(time.Duration(*appFrequency*3600.0) * time.Second)
	for {
		select {
		case <-ticker.C:
			rand.Seed(time.Now().UnixNano())
			glog.Info("AppTimer -- Attempting to kill a random application")
			if rand.Float64() <= *appProbability {
				victim := KillRandomApp()
				glog.Info("AppTimer -- Killed application: ", victim)
			} else {
				glog.Info("AppTimer -- Did not kill an application")
			}
		case <-quit:
			ticker.Stop()
			close(stop)
			return
		}
	}
}

/* Time/schedule task fraction deletion
 *  quit: channel for telling the timer to stop
 */
func FractionTimer(quit chan int) {
	glog.Info("FractionTimer is running")
	stop := make(chan int)
	ticker := time.NewTicker(time.Duration(*fractionFrequency*3600.0) * time.Second)
	for {
		select {
		case <-ticker.C:
			rand.Seed(time.Now().UnixNano())
			glog.Info("FractionTimer -- Attempting to kill a fraction of running tasks")
			if rand.Float64() <= *fractionProbability {
				victims := KillTaskFraction(*killFraction)
				glog.Info(fmt.Sprintf("FractionTimer -- Killed %d tasks", len(victims)))
				glog.V(1).Info(fmt.Sprintf("%#v", victims))
			} else {
				glog.Info("FractionTimer -- Did not kill any tasks")
			}
		case <-quit:
			ticker.Stop()
			close(stop)
			return
		}
	}
}

/* Time/schedule host task deletion
 *  quit: channel for telling the timer to stop
 */
func HostTimer(quit chan int) {
	glog.Info("HostTimer is running")
	stop := make(chan int)
	ticker := time.NewTicker(time.Duration(*hostFrequency*3600.0) * time.Second)
	for {
		select {
		case <-ticker.C:
			rand.Seed(time.Now().UnixNano())
			glog.Info("HostTimer -- Attempting to kill all tasks on a random host")
			if rand.Float64() <= *hostProbability {
				victim := KillHostTasks()
				glog.Info(fmt.Sprintf("HostTimer -- Killed all tasks running on %s", victim))
			} else {
				glog.Info("HostTimer -- Did not kill any tasks")
			}
		case <-quit:
			ticker.Stop()
			close(stop)
			return
		}
	}
}

/* Converts the current time to the number of hours and minutes since 0000 of the same day */
func CurrentTime() time.Duration {
	h := time.Duration(time.Now().Hour()) * time.Hour
	m := time.Duration(time.Now().Minute()) * time.Minute
	s := time.Duration(time.Now().Second()) * time.Second
	return h + m + s
}

/* Converts time in "HH:MM" format to the number of hours and minutes since 0000 of the same day
 * 	t:	string containing a time in "HH:MM" format
 */
func ParseTime(t string) time.Duration {
	tSplit := strings.Split(t, ":")
	tHour, _ := strconv.Atoi(tSplit[0])
	tMin, _ := strconv.Atoi(tSplit[1])
	return (time.Duration(tHour) * time.Hour) + (time.Duration(tMin) * time.Minute)
}
