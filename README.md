# Lxdbeat

Welcome to Lxdbeat.

This is very much a work in progress and as of the time of this writing, performs the bare minimum to debug.
Most of the code was generated from the [beats developer guide](https://www.elastic.co/guide/en/beats/devguide/current/new-beat.html)

## Getting Started with Lxdbeat

To get lxdbeat running, run "make" and then run lxdbeat

### Configuring

The following parameters are required in lxdbeat.yml

    client_cert: |
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----

    client_key: |
    -----BEGIN RSA PRIVATE KEY-----
    ...
    -----END RSA PRIVATE KEY-----

    server_cert: |
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----

    hosts:
        - hostname/ip:port

The client cert will need added to each of the LXD hosts as container list and state calls are trusted APIs

## Run

To run Lxdbeat with debugging output enabled, run:

```
./lxdbeat -c lxdbeat.yml -e -d "*"
```
