# Go iputils library

[![GoDoc](https://godoc.org/github.com/gtrafimenkov/iputils?status.svg)](https://godoc.org/github.com/gtrafimenkov/iputils)

A library of useful functions to work with IP addresses.

## Installation

```
go get -u github.com/gtrafimenkov/iputils
```

## Examples

Getting network ip range:

```
	_, network, _ := net.ParseCIDR("192.168.0.0/28")
	fmt.Println(GetNetworkIPRange(network))

	_, network, _ = net.ParseCIDR("beef::/64") // IPv6 network
	fmt.Println(GetNetworkIPRange(network))

	// Output:
	// 192.168.0.0 192.168.0.15
	// beef:: beef::ffff:ffff:ffff:ffff
```

Iterating over range of ip addresses:

```
	iter := GetIPRangeIterator(net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.5"))
	for ip, ok := iter.Next(); ok; ip, ok = iter.Next() {
		fmt.Println(ip)
	}

	// Output:
	// 192.168.0.1
	// 192.168.0.2
	// 192.168.0.3
	// 192.168.0.4
	// 192.168.0.5
```

## Documentation

[godoc.org/github.com/gtrafimenkov/iputils](https://godoc.org/github.com/gtrafimenkov/iputils)

## Tests

```
go test --cover
```

## License

MIT-0
