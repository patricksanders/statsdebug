[![CircleCI](https://circleci.com/gh/patricksanders/statsdebug/tree/master.svg?style=svg)](https://circleci.com/gh/patricksanders/statsdebug/tree/master)

# statsdebug

A small service to listen for statsd over UDP, parse the stats, and serve information about the metrics over HTTP.

## Metric Summaries

The service keeps track of the following for each unique metric name that is reported:
* `count` - the number of times the metric was sent
* `tags` - maps unique tag names to a list of the values for each tag. 
* `tag_sets` - provides the number of times that each unique **set** of tags was reported together. This may be useful for 
estimating cardinality: each set of tags corresponds to a unique timeseries

## Running

```bash
# Run the container
docker run --rm -p 8080:8080 -p 8125:8125/udp patricksanders/statsdebug

# Send lots of stats
for i in {1..1000}; do printf "bar.foo.baz:5|c#foo:bar,baz:bang" | socat -t 0 - UDP:localhost:8125; done

# See what we got!
curl localhost:8080/metric/bar.foo.baz | jq .
# {
#   "count": 1000,
#   "tags": {
#     "baz": [
#       "bang"
#     ],
#     "foo": [
#       "bar"
#     ]
#   },
#   "tag_sets": {
#     "baz:bang,foo:bar": 1000
#   }
# }

# Try another metric name
for i in {1..500}; do printf "hello.world:5|c" | socat -t 0 - UDP:localhost:8125; done

# Get counts for all metrics
curl localhost:8080/all | jq .
# {
#   "bar.foo.baz": 1000,
#   "hello.world": 500
# }

# Get details about all metrics
curl localhost:8080/all/details | jq .
# {
#   "bar.foo.baz": {
#     "count": 1000,
#     "tags": {
#       "baz": [
#         "bang"
#       ],
#       "foo": [
#         "bar"
#       ]
#     },
#     "tag_sets": {
#       "baz:bang,foo:bar": 1000
#     }
#   },
#   "hello.world": {
#     "count": 500,
#     "tags": {},
#     "tag_sets": {
#       "": 500
#     }
#   }
# }

# Reset the count
curl localhost:8080/reset
curl localhost:8080/all | jq .
# {}
```
