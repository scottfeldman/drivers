// This example uses TLS to send an HTTPS request to retrieve a webpage
//
// You shall see "strict-transport-security" header in the response,
// this confirms communication is indeed over HTTPS
//
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security
//
// Root certificate for httpbin.org built with:
//   openssl s_client -showcerts -connect httpbin.org:443 </dev/null 2>/dev/null |
//       openssl x509 -outform PEM > httpbin.crt

//go:build ninafw || wioterminal

package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"fmt"
	"io"
	"log"
	"machine"
	"net"
	"runtime"
	"strings"
	"time"

	"tinygo.org/x/drivers/netlink"
	"tinygo.org/x/drivers/netlink/probe"
)

var (
	ssid string
	pass string
	// HTTPS server address to hit with a GET / request
	address string = "httpbin.org:443"
	// NTP server to get current time
	ntpHost string = "0.pool.ntp.org:123"
	// TLS config to hold CA certificate(s)
	tlsConfig *tls.Config
	conn      net.Conn
)

// Wait for user to open serial console
func waitSerial() {
	for !machine.Serial.DTR() {
		time.Sleep(100 * time.Millisecond)
	}
}

func check(err error) {
	if err != nil {
		println("Hit an error:", err.Error())
		panic("BYE")
	}
}

func readResponse() {
	r := bufio.NewReader(conn)
	resp, err := io.ReadAll(r)
	check(err)
	println(string(resp))
}

func closeConnection() {
	conn.Close()
}

func dialConnection() {
	var err error

	println("\r\n---------------\r\nDialing TLS connection")
	conn, err = tls.Dial("tcp", address, tlsConfig)
	println("\r\n---------------\r\nDialing TLS done", conn, err)
	for ; err != nil; conn, err = tls.Dial("tcp", address, tlsConfig) {
		println("Connection failed:", err.Error())
		time.Sleep(5 * time.Second)
	}
	println("Connected!\r")
}

func makeRequest() {
	print("Sending HTTPS request...")
	w := bufio.NewWriter(conn)
	fmt.Fprintln(w, "GET /get HTTP/1.1")
	fmt.Fprintln(w, "Host:", strings.Split(address, ":")[0])
	fmt.Fprintln(w, "User-Agent: TinyGo")
	fmt.Fprintln(w, "Connection: close")
	fmt.Fprintln(w)
	check(w.Flush())
	println("Sent!\r\n\r")
}

const NTP_PACKET_SIZE = 48

var response = make([]byte, NTP_PACKET_SIZE)

func getCurrentTime(conn net.Conn) (time.Time, error) {
	if err := sendNTPpacket(conn); err != nil {
		return time.Time{}, err
	}

	n, err := conn.Read(response)
	if err != nil && err != io.EOF {
		return time.Time{}, err
	}
	if n != NTP_PACKET_SIZE {
		return time.Time{}, fmt.Errorf("expected NTP packet size of %d: %d", NTP_PACKET_SIZE, n)
	}

	return parseNTPpacket(response), nil
}

func sendNTPpacket(conn net.Conn) error {
	var request = [48]byte{
		0xe3,
	}

	_, err := conn.Write(request[:])
	return err
}

func parseNTPpacket(r []byte) time.Time {
	// the timestamp starts at byte 40 of the received packet and is four bytes,
	// this is NTP time (seconds since Jan 1 1900):
	t := uint32(r[40])<<24 | uint32(r[41])<<16 | uint32(r[42])<<8 | uint32(r[43])
	const seventyYears = 2208988800
	return time.Unix(int64(t-seventyYears), 0)
}

func setTime() error {
	println("Requesting NTP time...")

	conn, err := net.Dial("udp", ntpHost)
	if err != nil {
		return err
	}

	now, err := getCurrentTime(conn)
	if err != nil {
		return fmt.Errorf("Error getting current time: %v", err)
	}

	conn.Close()

	fmt.Printf("NTP time: %v\r\n", now)
	runtime.AdjustTimeOffset(-1 * int64(time.Since(now)))

	return nil
}

//go:embed httpbin.crt
var pem []byte

func setCACerts() error {
	// Create a CA certificate pool and add the embedded CA certificate
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(pem) {
		return fmt.Errorf("Failed to append CA certificate")
	}

	// Configure TLS with the custom CA pool
	tlsConfig = &tls.Config{
		RootCAs: caCertPool,
	}

	return nil
}

func main() {
	waitSerial()

	link, _ := probe.Probe()

	err := link.NetConnect(&netlink.ConnectParams{
		Ssid:       ssid,
		Passphrase: pass,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := setTime(); err != nil {
		log.Fatal(err)
	}

	if err := setCACerts(); err != nil {
		log.Fatal(err)
	}

	for i := 0; ; i++ {
		dialConnection()
		makeRequest()
		readResponse()
		closeConnection()
		println("--------", i, "--------\r\n")
		time.Sleep(10 * time.Second)
	}
}
