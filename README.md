# xudp

xudp is a small UDP packet generator for sending fixed-size payloads to IPv4 or
IPv6 destinations at a configurable packet rate.

## Installation

Install the latest version directly from this repository:

```sh
go install github.com/yz4230/xudp@latest
```

Or build from a local checkout:

```sh
go build .
```

## Usage

Send 1,000 UDP packets with a 64-byte payload to `::1:4321` at 1,000 packets per
second:

```sh
xudp
```

Send 10,000 UDP packets to an IPv4 address at 5,000 packets per second:

```sh
xudp --dst 192.0.2.10 --port 4321 --count 10000 --rate 5000 --len 128
```

Disable rate limiting by setting `--rate` to `0`:

```sh
xudp --dst 127.0.0.1 --port 4321 --count 1000 --rate 0
```

## Flags

```text
  -c, --count int   Number of packets to send (default 1000)
  -d, --dst ip      Destination IP address (default ::1)
  -h, --help        help for xudp
  -l, --len int     Length of the UDP payload (default 64)
  -p, --port int    Destination port (default 4321)
  -r, --rate int    Packets per second (default 1000)
  -v, --verbose     Verbose output
```

## Development

Run the CLI locally:

```sh
go run . --help
```

Check the package:

```sh
go test ./...
```
