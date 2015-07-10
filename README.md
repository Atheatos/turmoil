Turmoil
=======
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
   
Run Turmoil with the configuration file:
```
./turmoil -config=params.ini
```
* * *
### Docker
For the container to run successfully, Turmoil must be compiled as a statically-linked binary
```
$ CGO_ENABLED=0 go build -a -installsuffix cgo
$ ldd turmoil
	not a dynamic executable
```
Build the container image using ```docker build``` or retrieve with ```docker pull atheatos/turmoil```  

Use ```-v``` to mount the local time file. Turmoil will use ```params.ini``` from the root directory if it is not found in ```/mnt/mesos/sandbox```
```
$ docker run --rm -it \
> -v /etc/localtime:/etc/localtime:ro \
> atheatos/turmoil:dev
```  
  
Run the container on Marathon:
```
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
	"cpus": 0.1,
	"id": "turmoil",
	"instances": 1,
	"mem": 16,
	"uris": [],
	"disk": 4
}
```
```curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" http://<marathon_url>:8080/v2/apps -d@turmoil.json```

### Dependencies
+ [iniflags](https://github.com/vharitonsky/iniflags)
+ [glog](https://github.com/golang/glog)
+ [go-marathon](http://github.com/gambol99/go-marathon)
