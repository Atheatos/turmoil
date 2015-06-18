Turmoil
=======
Turmoil is a tool for testing the recovery ability of applications running on [Marathon](https://mesosphere.github.io/marathon/).

Turmoil can currently perform three functions:  
  1. Kill a single task  
  2. Kill a single application  
  3. Kill a given fraction of currently running tasks  

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
_Note: a frequency of zero is not currently supported_  
  
  
Run Turmoil with the configuration file:
```
./turmoil -config=params.ini
```
* * *
### Dependencies
+ [iniflags](https://github.com/vharitonsky/iniflags) 

+ [glog](https://github.com/golang/glog)

+ [go-marathon](http://github.com/gambol99/go-marathon)
