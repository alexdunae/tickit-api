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


