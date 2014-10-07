## Overview

A statsd client for Go.

## Get the code

    $ go get github.com/xsleonard/go-statsd

## Usage

    // Create the client
    c, err := statsd.New("127.0.0.1:8125")
    if err != nil {
      log.Fatal(err)
    }
    // Prefix every metric with the app name
    c.Namespace = "app"
    // "app.requests:3|c" sent to the statsd server
    err = c.Count("requests", 3)

## Development

Run the tests with:

    $ go test

## Documentation

Please see: http://godoc.org/github.com/xsleonard/go-statsd

## License

go-statsd is released under the [MIT license](http://www.opensource.org/licenses/mit-license.php).
