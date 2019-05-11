#!/usr/bin/env bash

# starts a detached consul instance
docker-compose up -d

# hack - delay allows Consul to initialise
sleep 2

# insert the dev configuration into consul kv store
curl --request PUT --data @config.json http://127.0.0.1:8500/v1/kv/services/stream-control
