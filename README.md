# sphero-rallye-server

Server for sphero-rallye game written in golang.

## Usage

To start the server just run:

        sudo go run main.go

Or to build an executable and run that:

        go build main.go
        sudo ./main

To compile executable for ARM to run server on a Raspberry Pi:

        GOARM=7 GOARCH=arm GOOS=linux go build main.go

The server uses [Gobot library](https://gobot.io/documentation/platforms/sprkplus/) to control sphero robots of players. In order to connect to all spheros correctly the server needs to know the bluetooth id's of all spheros. One need to configure them in config file (see below) under `rallye.players[0-n].bluetoothId`.

The server needs a running mqtt broker to communicate with. The address of the broker can be configured in
config files `./conf/conf.json` or `/etc/rallye-rallye-server/conf.json`.

In config file one can define the length of one round of the game in seconds.

The server reads and writes to/from several mqtt topics:

### MQTT Writes

#### End time of the next round

In the end of each round the server will send the end time of each round to the topic `<global topic prefix>/roundEnd`

Were
* `<global topic prefix>` can be configured in config file under path `mqtt.topicPrefix`

#### Next possible moves for each player

In the end of each round the server will produce a number of random ActionTypes for each of the participatory players.
As well number of actions per round as the participatory players can be configured in config file.

The number of possible ActionTypes will be published to different topics for each player: `<global topic prefix>/<player name>/possibleActionTypes`

Where
* `<global topic prefix>` can be configured in config file under path `mqtt.topicPrefix`
* `<player name>` can be configured under path `rallye.players[0-n].name`

#### Error informations for players

While processing next actions, send by players (see next section) errors may occur. Informations about these errors are
sent to the mqtt topic `<global topic prefix>/<player name>/errors`.

Where
* `<global topic prefix>` can be configured in config file under path `mqtt.topicPrefix`
* `<player name>` can be configured under path `rallye.players[0-n].name`

### MQTT Reads

Player clients may send their next actions to the topic `<global topic prefix>/<player name>/plannedActions`

Where
* `<global topic prefix>` can be configured in config file under path `mqtt.topicPrefix`
* `<player name>` can be configured under path `rallye.players[0-n].name`