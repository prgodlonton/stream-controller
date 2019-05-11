# README

Tracks the number of streams that users are viewing concurrently. 

## Configuration

The service has been written to read its configuration from a Consul service KV store. The `startup.sh` and `shutdown.sh`
scripts found in the `dev` directory startup a dev Consul server, add the sample configuration in `config.json` to the 
KV store and shutdown the Consul container respectively. 

The two environment variables `CONSUL_ADDRESS` and `CONSUL_KEY` determine the address of the Consul server and the key 
under which the configuration for this service is stored. If these are not provided then the default values of 
`http://localhost:8500` and `services/stream-control` are used.

 