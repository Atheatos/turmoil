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
	marathon "github.com/gambol99/go-marathon"
	"math/rand"
	"strings"
	"time"
)

/*  Kills a fraction of existing tasks
 * 	client: 	Marathon client interface
 * 	blacklist: 	do not kill tasks from these applications
 * 	fraction: 	fraction of the total tasks to be killed
 */
func KillTaskFraction(client marathon.Marathon, blacklist []string, fraction float64) {
	// Get tasks from marathon
	tasks, err := client.AllTasks()
	Assert(err)
	// Prevent suicide
	tasklist := ExtractTaskIDs(tasks.Tasks)
	tasklist = EnforceBlacklist(tasklist, blacklist)
	// Random permutate the array and kill the first `numTargets` tasks
	rand.Seed(time.Now().UnixNano())
	numTasks := float64(len(tasklist))
	numTargets := int(numTasks*fraction)
	indices := rand.Perm(int(numTasks))[0:numTargets]
	targets := make([]string, numTargets)
	for i, randi := range(indices) {
		targets[i] = tasklist[randi]
	}
	// Execute
	Assert(client.KillTasks(targets, true))
}

/*  Kills one random application
 * 	client: 	Marathon client interface
 * 	blacklist: 	do not kill these applications
 */
func KillRandomApp(client marathon.Marathon, blacklist []string) {
	// Get applications from marathon
	applications, err := client.ListApplications()
	Assert(err)
	// Prevent suicide
	applist := EnforceBlacklist(applications, blacklist)
	// Kill random
	rand.Seed(time.Now().UnixNano())
	app := applist[rand.Intn(len(applist))]
	// Execute
	_, delerr := client.KillApplicationTasks(app, "", false)
	Assert(delerr)
}

/*  Kill one random task
 *  	client: 	Marathon client interface
 *  	blacklist:  do not kill tasks from these applications
 */
func KillRandomTask(client marathon.Marathon, blacklist []string) {
	// Get all of the running tasks from marathon
	tasks, err := client.AllTasks()
	Assert(err)
	// Remove turmoil from the list of targets
	tasklist := ExtractTaskIDs(tasks.Tasks)
	tasklist = EnforceBlacklist(tasklist, blacklist)
	// Choose a random task
	rand.Seed(time.Now().UnixNano())
	task := tasklist[rand.Intn(len(tasklist))]
	app := task[0:strings.LastIndex(task, ".")] // extract the app name from the task ID
	// Tell Marathon to delete chosen task
	_, delerr := client.KillTask(app, task, false)
	Assert(delerr)
}

/*  Remove blacklisted applications or tasks from blacklisted applications from a list of potential targets 
 *  	targets: 	array of application or task ids for potential targets
 * 	blacklist: 	array of applications to not kill
 */
func EnforceBlacklist(targets, blacklist []string) []string {
	for i, target := range(targets) {
		glog.Info(target)
		for _, blacklisted := range(blacklist) {
			if (strings.HasPrefix(target, "/")) && (strings.TrimPrefix(target, "/")==blacklisted) {
				targets = append(targets[:i], targets[i+1:]...)
			} else if (!strings.HasPrefix(target, "/")) && (target[0:strings.LastIndex(target, ".")]==blacklisted) {
				targets = append(targets[:i], targets[i+1:]...)
			}
		}
	}
	return targets
}

/*  Extract a string array of task IDs from an array of Marathon Task structs
 *  	tasks: 	array of Marathon Task structs
 */
func ExtractTaskIDs (tasks []marathon.Task) []string {
	stringTasks := make([]string, len(tasks))
	for i, task := range(tasks) {
		stringTasks[i] = task.ID
	}
	return stringTasks
}

// Assert an error, if any
func Assert(err error) {
	if err != nil {
		glog.Fatalf("Failed, error: %s", err)
	}
}
