## Setup

Make sure both `git` and `hg` is installed.

There are two servers: `checkin-api` and `checkin-websocket`.

### Check In Web API

Install dependencies and start the server:

    cd $GOPATH/src/tickit/checkin-api
    go get
    go install && checkin-api


### Web API

Install dependencies and start the server:

    cd $GOPATH/src/tickit/checkin-websocket
    go get
    go install && checkin-websocket


### Deployment

    vagrant up
    ssh vagrant@127.0.0.1 -p 2222 -i ~/.vagrant.d/insecure_private_key

Install Go > 1.1 (on Ubuntu: http://stackoverflow.com/a/17566846/559596)

    mkdir ~/go
    export GOPATH=~/go
    mkdir -p $GOPATH/src/tickit

    cd $GOPATH/src/tickit
    git clone git@github.com:alexdunae/tickit-api.git .

    cd tickit-api
    go get
    go install
