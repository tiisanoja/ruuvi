# QuestDB
Supports InfluxDB (input) and PostgresDB (output) interfaces

Downloadable from https://hub.docker.com/r/questdb/questdb

## Staring database
$ docker run -p 9000:9000 -p 8812:8812 questdb/questdb

Persist database:
$ docker run -p 9000:9000 -p 8812:8812 -v local/dir:/var/lib/questdb questdb/questdb

Replace local/dir with the absolute path to the directory on your host machine where you want to persist the data.

This is the list of ports used by QuestDB:
* 9000 for the InfluxDB Line Protocol, REST API and the Web Console (accessible at localhost:9000⁠)
* 8812 for the Postgres wire protocol, Used to query the data

