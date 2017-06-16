# Iftopbeat

Welcome to Iftopbeat.

Ensure that this folder is at the following location:
`${GOPATH}/github.com/yuemin-li`

## Getting Started with Iftopbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Install iftop
On Mac
```
brew install iftop
```
iftop requires root privileges so you will need to run `sudo iftop`.

On Ubuntu 16.04
```
sudo apt install iftop
```
This will install iftop, version 1.0pre4 on your system.
You will need iftop version 1.0pre4 to have the `-t` option, the iftop 1.0pre2 packaged with Ubuntu 14.04 won't have this option.

There is a `.Vagrantfile` included, which I use for my local development env.

### More on iftop
Read more on its man page: https://linux.die.net/man/8/iftop

iftop cmd in use

-t          Use  text interface without ncurses and print the output to STD-OUT

-s num      print one single text output afer num seconds, then quit

-L num      number of lines to print

-n          don't do hostname lookups


`sudo iftop -t -s 10 -L 10 -i eth0 -n`

Wait for 10 sec and print 10 lines iftop result on interface eth0 to STD-OUT


By default iftop requires root privileges(which make sense, since you don't want any user application can sniff on your traffic).
So if you don't wanna start iftopbeat with sudo, add iftop to userGroup. Otherwise you will have to start iftopbeat with sudo, and provide your root credential. 


### Init Project
To get running with Iftopbeat and also install the dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Iftopbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/yuemin-li/iftopbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Iftopbeat run the command below. This will generate a binary
in the same directory with the name iftopbeat.

```
make
```

After dev work, use `go clean` to clean up your local binary.

### Configure iftopbeat
In your iftopbeat.yml file, config the following options.
```
iftopbeat:
  # Defines how often an event is sent to the output (in seconds)
  period: 10
  # Defines how often the iftop gathers results
  interval: 10
  # Defines which interface iftop is monitoring
  listenOn: "eth0"
  # Defines the limitation of lines iftop outputs each time
  numLines: 10
```
If there is anything wrong parsing config from this file, a default config will be used for your beater process.

####TODO: cmd line args support.####


### Run
To run elasticsearch docker locally for dev, run:
```
docker run --name elasticsearch -p 9200:9200 -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" -e "http.host=0.0.0.0" -e "transport.host=127.0.0.1" docker.elastic.co/elasticsearch/elasticsearch:5.4.0
```

To run kibana docker locally for dev, run:
```
docker run --name kibana --link elasticsearch:elasticsearch -p 5601:5601 -d docker.elastic.co/kibana/kibana:5.4.0
```

By default, use `elastic/changeme` for login credentials.

To run Iftopbeat with debugging output enabled, run:

```
sudo ./iftopbeat -c iftopbeat.yml -e -d "*"
```

Get iftopbead index.
```
curl -u elastic:changeme http://localhost:9200/_cat/indices?v
```

The index is in iftopbeat-DATE format. 

Get records of iftopbeat.
```
curl -u elastic:changeme http://localhost:9200/iftopbeat-<DATE>/_search?pretty=true&q=*:*
```

Get Kibana GUI.

Visit your kinana at http://HOSTNAME:5601.

Create a index pattern like `iftopbeat-*`, that matches multiple iftopbeat indices with different date.

### Test

To test Iftopbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/iftopbeat.template.json and etc/iftopbeat.asciidoc

```
make update
```


### Cleanup

To clean  Iftopbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Iftopbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/yuemin-li
cd ${GOPATH}/github.com/yuemin-li
git clone https://github.com/yuemin-li/iftopbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.
