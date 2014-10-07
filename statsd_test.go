package statsdx

import (
	"net"
	"testing"
)

func doTest(t *testing.T, conn *net.UDPConn, fn func() error, req string) {
	err := fn()
	if err != nil {
		t.Fatal(err)
	}

	bytes := make([]byte, 1024)
	n, _, err := conn.ReadFrom(bytes)
	if err != nil {
		t.Fatal(err)
	}

	recv := string(bytes[:n])
	if recv != req {
		t.Errorf("Expected: %s. Actual: %s", req, recv)
	}
}

func TestClient(t *testing.T) {
	addr := "localhost:1882"

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		t.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	c, err := New(addr)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := c.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	doTest(t, conn, func() error { return c.Gauge("test", 99) }, "test:99|g")
	doTest(t, conn, func() error { return c.GaugeDelta("test", 99) }, "test:+99|g")
	doTest(t, conn, func() error { return c.GaugeDelta("test", -99) }, "test:-99|g")
	doTest(t, conn, func() error { return c.GaugeF("test", 99.9) }, "test:99.900000|g")
	doTest(t, conn, func() error { return c.GaugeDeltaF("test", 99.9) }, "test:+99.900000|g")
	doTest(t, conn, func() error { return c.GaugeDeltaF("test", -99.9) }, "test:-99.900000|g")

	doTest(t, conn, func() error { return c.Count("test", 2) }, "test:2|c")
	doTest(t, conn, func() error { return c.Count("test", -2) }, "test:-2|c")
	doTest(t, conn, func() error { return c.SampledCount("test", 2, 0.8) }, "test:2|c|@0.800000")

	doTest(t, conn, func() error { return c.Timing("test", 5) }, "test:5|ms")
	doTest(t, conn, func() error { return c.TimingF("test", 5.1) }, "test:5.100000|ms")

	doTest(t, conn, func() error { return c.Set("test", "x") }, "test:x|s")

	// Namespace check
	c.Namespace = "ns"
	doTest(t, conn, func() error { return c.Set("test", "x") }, "ns.test:x|s")
}

func assertNotPanics(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatal(r)
		}
	}()
	f()
}

func TestNilSafe(t *testing.T) {
	var c *Client = nil

	checkNilSafe := func(fn func() error) {
		assertNotPanics(t, func() {
			err := fn()
			if err != nil {
				t.Fatal(err)
			}
		})
	}

	checkNilSafe(func() error { return c.Close() })
	checkNilSafe(func() error { return c.Gauge("t", 1) })
	checkNilSafe(func() error { return c.GaugeDelta("t", -1) })
	checkNilSafe(func() error { return c.GaugeF("t", 1.1) })
	checkNilSafe(func() error { return c.GaugeDeltaF("t", -1.1) })
	checkNilSafe(func() error { return c.Count("t", 2) })
	checkNilSafe(func() error { return c.SampledCount("t", 2, 0.4) })
	checkNilSafe(func() error { return c.Set("t", "x") })
	checkNilSafe(func() error { return c.send("t", "1|c") })
}

func TestErrors(t *testing.T) {
	addr := "localhost:1882"
	c, err := New(addr)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Gauge("test", -1)
	if err == nil || err.Error() != "Gauge value must be >= 0" {
		t.Fatalf("Expected error, got %v", err)
	}

	err = c.GaugeF("test", -1.2)
	if err == nil || err.Error() != "GaugeF value must be >= 0" {
		t.Fatalf("Expected error, got %v", err)
	}

	err = c.send("", "val")
	if err == nil || err.Error() != "Name required" {
		t.Fatalf("Expected error, got %v", err)
	}
}
