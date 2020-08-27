#go-scheduler
It is a Golang package for running delayed jobs at any specified epoch.It uses redis for adding and executing jobs and has a very simple interface.

###Getting Started
* Start redis at localhost:6379 or change to appropriate host
* Push value to "SCHEDULE" using ZADD manually.