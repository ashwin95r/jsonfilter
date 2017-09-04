### To run:

```
$ go get -u github.com/ashwin95r/jsonfilter
$ go install .
$ PORT=<port number> jsonfilter
```

### POST JSON:

The requests should be made to the root endpoint.
```
$ curl http://jsonfilter123.herokuapp.com/ -XPOST -d @<path to input file>
```

### Improvements (TODO):

Whole JSON document is unmarshaled at once. If it's too huge we could use `json.Decoder` to break it in parts if Memory is a bottleneck.
