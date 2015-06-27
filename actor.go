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
	"github.com/golang/glog"
	"math/rand"
	"strings"
	"time"
)

/*  Kills a fraction of existing tasks
 *  	fraction: 	fraction of the total tasks to be killed
 */
func KillTaskFraction(fraction float64) []string {
	// Get tasks from marathon and enforce blacklist
	tasklist, err := client.ListTasks()
	Assert(err)
	tasks := EnforceBlacklist(tasklist)
	// Randomly permute the array and kill the first `numTargets` tasks
	rand.Seed(time.Now().UnixNano())
	numTasks := float64(len(tasks))
	numTargets := int(numTasks * fraction)
	indices := rand.Perm(int(numTasks))[0:numTargets]
	targets := make([]string, numTargets)
	for i, randi := range indices {
		targets[i] = tasks[randi]
	}
	// Execute
	Assert(client.KillTasks(targets, false))
	return targets
}

/*  Kills one random application */
func KillRandomApp() string {
	// Get applications from marathon
	applications, err := client.ListApplications()
	Assert(err)
	// Prevent suicide
	applist := EnforceBlacklist(applications)
	// Kill random
	rand.Seed(time.Now().UnixNano())
	app := applist[rand.Intn(len(applist))]
	// Execute
	_, err = client.KillApplicationTasks(app, "", false)
	Assert(err)
	return app
}

/*  Kill one random task */
func KillRandomTask() string {
	// Get all of the running tasks from marathon
	tasks, err := client.ListTasks()
	Assert(err)
	// Remove turmoil from the list of targets
	tasks = EnforceBlacklist(tasks)
	// Choose a random task
	rand.Seed(time.Now().UnixNano())
	glog.Info(tasks)
	task := tasks[rand.Intn(len(tasks))]
	// Tell Marathon to delete chosen task
	_, err = client.KillTask(task, false)
	Assert(err)
	return task
}

/*  Remove blacklisted applications or tasks from blacklisted applications from a list of potential targets
 *  	targets: 	array of application or task ids for potential targets
 */
func EnforceBlacklist(targets []string) []string {
	for i, target := range targets {
		for _, blacklisted := range blacklist {
			if (strings.HasPrefix(target, "/")) && (strings.TrimPrefix(target, "/") == blacklisted) {
				targets = append(targets[:i], targets[i+1:]...)
			} else if (!strings.HasPrefix(target, "/")) && (target[0:strings.LastIndex(target, ".")] == blacklisted) {
				targets = append(targets[:i], targets[i+1:]...)
			}
		}
	}
	return targets
}

// Assert an error, if any
func Assert(err error) {
	if err != nil {
		glog.Fatalf("Failed, error: %s", err)
	}
}
