# AbyleEDA

## Description

This is AbyleEDA, the Abyle Event Driven Automation programming suite.

It is based on a client/server architecture to send events, which can be floating point numbers, sensor values, switch states, etc., over the network. The server handles the events by either showing them or triggering other events.

The event graph will be customizable and examples on how to use it will be provided. At the moment the event system can be used to call a custom callback function based on client and event ID.

The data is transmitted as encrypted (AES256) JSON packages over UDP. In the feature it is planed to support also TCP and maybe other encryptions.

## How to download/install

Make sure your $GOPATH is set.

Checkout AbyleEDA with

    go get github.com/torlenor/AbyleEDA

Checkout all dependencies

    go get -v ./...

Build all examples with

    go install ./.../AbyleEDA/...

The binaries will land in $GOPATH/bin

## Examples

### simplefakedataserver/simplefakedataclient

These examples demonstrate the simplest way to setup a server and a client, send floating point events from client to server and how to setup the custom event callback of the AEDAevents system.

### temperaturesensorserver/temperaturesensorclient

The client sends temperature sensor data read from a file, e.g., /sys/class/hwmon/hwmon0/temp2_input, to a server which provides a simple webserver to show the data (webserver and templates are work in progress).

### pingserver/pingclient

On client start hosts to ping and the ping interval can be specified. Those hosts are pinged using github.com/SewanDevs/go-ping and sent to a server which will, at some point in the future, provide a web interface to the results (work in progress).