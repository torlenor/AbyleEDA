# AbyleEDA
This is AbyleEDA, the Abyle Event Driven Automation programming suite.
It consists of client/server architecture to send events, which can be sensor values, 
images, switch states, over the network to a server which handles the events by 
either showing them or activating other events. The event graph will be customizable
by a config file (probably xml based). We will provide examples (some of them exist already)
on how to implement a server or a client based on AbyleEDA.

Working already is sending event messages from a client to a server via UDP, which are encrypted, 
hashed JSON packages. As an example we provide a fakedata client which generates fake data 
and sends it to a server, a simple client which reads temperature data from a sensors file
(e.g., /sys/class/hwmon/hwmon0/temp2_input) and the server which receives the data,
prints the events to stdout and can also spawn a webserver to show the values received
for sensor events.

Make sure your $GOPATH is set.

Checkout AbyleEDA with

    go get github.com/torlenor/AbyleEDA

Checkout all dependencies

    go get -v ./...

Build all binaries (udpserver, udpclient at the moment) with

    go install ./.../AbyleEDA/...

The binaries should land in $GOPATH/bin
