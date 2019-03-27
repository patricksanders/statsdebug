[![Build Status](https://travis-ci.org/patricksanders/statsdebug.svg?branch=master)](https://travis-ci.org/patricksanders/statsdebug)

# statsdebug

A small service to listen for statsd over UDP, parse the stats, and serve information about the metrics over HTTP.

## Running

```bash
# Run the container
docker run --rm -p 8080:8080 -p 8125:8125/udp patricksanders/statsdebug

# Send lots of stats
for i in {1..1000}; do printf "bar.foo.baz:5|c#foo:bar,baz:bang" | socat -t 0 - UDP:localhost:8125; done

# See what we got!
curl localhost:8080/metric/bar.foo.baz
# {"count":1000}


# Try another metric name
for i in {1..500}; do printf "hello.world:5|c" | socat -t 0 - UDP:localhost:8125; done

# Get counts of all metrics
curl localhost:8080/all
# {"bar.foo.baz":1000,"hello.world":500}
```
