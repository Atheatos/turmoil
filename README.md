#Turmoil
* * *
Turmoil is a [Chaos Monkey](https://github.com/Netflix/SimianArmy/wiki/Chaos-Monkey)-like tool for testing the recovery ability of applications running on [Marathon](https://mesosphere.github.io/marathon/).

Turmoil can currently perform four functions:  
  1. Kill a single task  
  2. Kill a single application  
  3. Kill a given fraction of currently running tasks  
  4. Kill all tasks on a single host  
  
Targets are selected pseudo-randomly and are killed via calls to Marathon's REST API.
* * *
### Configuration
Turmoil can be configured by modifying the fields listed in the ```params.ini``` file. This may be used to set the frequency and probability of kill attempts.

For example: every two hours, there is a 75% chance that half of all running tasks will be deleted
```
fraction = 0.5
fractionFrequency = 2
fractionProbability = 0.75
```   
   
The start and stop times can also be set:
```
start = "10:00"
stop = "16:00"
```  
##### Logging  
Log outputs can be enabled using ```-logtostderr```
```
./turmoil -logtostderr
```  
* * *
### Docker  
  
##### Build  
For the container to run successfully, Turmoil must be compiled as a statically-linked binary.
```
$ CGO_ENABLED=0 go build -a -installsuffix cgo
$ ldd turmoil
	not a dynamic executable
```
Build the container image using ```docker build``` or retrieve with ```docker pull atheatos/turmoil```  
  
##### Run  
Use ```-v``` to mount the local time file
```
docker run --rm -it \
    -v /etc/localtime:/etc/localtime:ro \
    atheatos/turmoil:dev
```  
  
##### Configure  
Use the environment flag, ```-e, --env``` to set the variables within the container
+ ```TURMOIL_PARAM``` specifies the name of the parameter file; Turmoil will use ```params.ini``` if this is not set
+ ```MESOS_SANDBOX``` specifies the directory that Turmoil will first check for the parameter file (default: ```/mnt/mesos/sandbox```). If the file is not found, Turmoil uses ```/params.ini```
  
A parameter file can be placed in the mesos sandbox by providing Marathon a URI from which to pull the file.  
  
For example, the file can be served by a local fileserver:  
```go
package main

import "net/http"

func main() {
  http.Handle("/", http.FileServer(http.Dir(".")))
  http.ListenAndServe(":8108",nil)
}
```  
```
$ ls
custom_filename.ini  server.go
$ go run server.go

```  
 
##### Marathon  
Now, run the container on Marathon:
```json
{
	"container": {
		"docker": {
			"image": "atheatos/turmoil:dev"
		},
		"type": "DOCKER",
		"volumes": [
			{
				"containerPath": "/etc/localtime",
				"hostPath": "/etc/localtime",
				"mode": "RO"
			}
		]
	},
	"env": {
		"MESOS_SANDBOX": "$MESOS_SANDBOX",
		"TURMOIL_PARAM": "custom_filename.ini"
	},
	"cpus": 0.1,
	"id": "turmoil",
	"instances": 1,
	"mem": 16,
	"uris": [ "http://127.0.0.1:8108/custom_filename.ini" ],
	"disk": 4
}

```  
```
curl -X POST \
  -H "Accept: application/json" \
  -H "Content-Type: application/json" \
  http://127.0.0.1:8080/v2/apps \
  -d@turmoil.json
```  
Turmoil will log to standard error which can be viewed in the task's sandbox in Mesos.

### Dependencies
+ [iniflags](https://github.com/vharitonsky/iniflags)
+ [glog](https://github.com/golang/glog)
+ [go-marathon](http://github.com/gambol99/go-marathon)
