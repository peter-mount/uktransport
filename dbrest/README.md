# dbrest

dbrest is a standalone utility for invoking PostgreSQL functions from either REST or messages from RabbitMQ.

It started off being part of the nre-feeds project when it was noticed it should be more generic for exposing other
opendata than just the NRE data feeds.

## Overview

The utility currently has 3 modes of operation:
1. Exposing functions over a rest api
1. Invoking functions with no parameters on a Cron schedule
1. Invoking functions when a message is received from a RabbitMQ queue

You can run the utility with any or all of these modes at the same time - although if you are receiving from RabbitMQ
as well as hosting REST I suggest you run them in two separate instances.

## Configuration

Configuration is done with a yaml file and contains the required db section defining the connection to the database and
the optional sections for the supports modes above.

### Database

The database section defines the connection to the database:

```yaml
db:
    url: "postgres://user:password@hostname/dbname?sslmode=disable&connect_timeout=3"
    maxOpen: 10
    maxIdle: 1
    maxLifetime: 3600
```

* maxOpen defines the maximum number of database connections to use, defaults to 1.
* maxIdle defines the maximum number of idle database connections, defaults to 1. If > maxOpen then it gets set to maxOpen.
* maxLifetime defines how long a connection will remain open in seconds. After that time the connection will be closed
and a new one opened. Default is 3600 (1 hour)

## Rest API

The rest API is simply an array of declared endpoints which are services by a PostgreSQL function.
An example of a pair of defined endpoints are:

```yaml
rest:
    -
      method: GET
      path: /api/crs/{crs}
      function: cif.getcrs
      parameters:
        - crs
      json: true
    -
      method: GET
      path: /api/crs
      function: cif.getallcrs
      json: true
```
In the above example we have two REST endpoints defined, one that takes a path parameter and the other that doesn't take any.
