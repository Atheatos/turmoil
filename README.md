Turmoil
=======
Turmoil is a [Chaos Monkey](https://github.com/Netflix/SimianArmy/wiki/Chaos-Monkey)-like tool for testing the recovery ability of applications running on [Marathon](https://mesosphere.github.io/marathon/).

Turmoil can currently perform three functions:  
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
For the container to build and run successfully, Turmoil must be compiled with statically-linked binaries
```
$ CGO_ENABLED=0 go build -a -installsuffix cgo
$ ldd turmoil
	not a dynamic executable
```
Build the container image using ```docker build``` or retrieve with ```docker pull atheatos/turmoil```  

Use ```-v``` to mount the parameter and local time files at runtime
```
$ docker run --rm -it \
> -v /etc/localtime:/etc/localtime:ro \
> -v /path/to/params.ini:/params.ini:ro \
> atheatos/turmoil:dev
```  
  
### Dependencies
+ [iniflags](https://github.com/vharitonsky/iniflags)
+ [glog](https://github.com/golang/glog)
+ [go-marathon](http://github.com/gambol99/go-marathon)
