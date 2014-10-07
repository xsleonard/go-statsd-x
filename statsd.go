package statsdx

import (
	"errors"
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
	// Namespace to prepend to all statsd calls
	Namespace string
}

// New returns a pointer to a new Client and an error.
// addr must have the format "hostname:port"
func New(addr string) (*Client, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}
	client := &Client{conn: conn}
	return client, nil
}

// send handles sampling and sends the message over UDP.
// It also adds global namespace prefixes.
func (c *Client) send(name string, value string) error {
	if name == "" {
		return errors.New("Name required")
	}

	if c == nil {
		return nil
	}

	if c.Namespace != "" {
		name = fmt.Sprintf("%s.%s", c.Namespace, name)
	}

	data := fmt.Sprintf("%s:%s", name, value)
	_, err := c.conn.Write([]byte(data))
	return err
}

// Gauges measure the value of a metric at a particular time
func (c *Client) Gauge(name string, value int64) error {
	if value < 0 {
		return errors.New("Gauge value must be >= 0")
	}
	stat := fmt.Sprintf("%d|g", value)
	return c.send(name, stat)
}

// GaugeDelta adjust a gauge's value incrementally
func (c *Client) GaugeDelta(name string, value int64) error {
	stat := fmt.Sprintf("%+d|g", value)
	return c.send(name, stat)
}

func (c *Client) GaugeF(name string, value float64) error {
	if value < 0 {
		return errors.New("GaugeF value must be >= 0")
	}
	stat := fmt.Sprintf("%f|g", value)
	return c.send(name, stat)
}

func (c *Client) GaugeDeltaF(name string, value float64) error {
	stat := fmt.Sprintf("%+f|g", value)
	return c.send(name, stat)
}

// Counters track how many times something happened per second
func (c *Client) Count(name string, value int64) error {
	stat := fmt.Sprintf("%d|c", value)
	return c.send(name, stat)
}

// SampledCount allows subsecond reporting
func (c *Client) SampledCount(name string, value int64, rate float64) error {
	stat := fmt.Sprintf("%d|c|@%f", value, rate)
	return c.send(name, stat)
}

// Timings measure the duration of an event
func (c *Client) Timing(name string, ms int64) error {
	stat := fmt.Sprintf("%d|ms", ms)
	return c.send(name, stat)
}

// Timings measure the duration of an event
func (c *Client) TimingF(name string, ms float64) error {
	stat := fmt.Sprintf("%f|ms", ms)
	return c.send(name, stat)
}

// Sets count the number of unique elements in a group
func (c *Client) Set(name string, value string) error {
	stat := fmt.Sprintf("%s|s", value)
	return c.send(name, stat)
}

// Close the client connection
func (c *Client) Close() error {
	if c == nil {
		return nil
	}
	return c.conn.Close()
}
