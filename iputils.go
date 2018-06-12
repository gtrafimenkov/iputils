// SPDX-License-Identifier: MIT-0

package iputils

import (
	"bytes"
	"fmt"
	"net"
)

const (
	// IPv4Size contains size of IPv4 in bytes
	IPv4Size = 4

	// IPv6Size contains size of IPv6 in bytes
	IPv6Size = 16
)

var (
	// V4InV6Prefix describes the prefix of a IPv4 address in IPv6 structure
	V4InV6Prefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}

	// MaxIPv4 contains maximal possible IPv4
	MaxIPv4 = []byte{0xff, 0xff, 0xff, 0xff}

	// MaxIPv4In6 contains maximal possible IPv4 in IPv6 structure
	MaxIPv4In6 = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	// MaxIPv6 contains maximal possible IPv6
	MaxIPv6 = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

// CopyIP copies ip address
func CopyIP(ip net.IP) net.IP {
	size := len(ip)
	result := make([]byte, size, size)
	copy(result, ip)
	return result
}

// Next increments ip to the next sequental value if that's possible.
// If not possible, false is returned.
func Next(ip net.IP) bool {
	size := len(ip)

	if size == IPv4Size && bytes.Equal(ip, MaxIPv4) {
		return false
	}

	if size == IPv6Size && (bytes.Equal(ip, MaxIPv4In6) || bytes.Equal(ip, MaxIPv6)) {
		return false
	}

	for i := size - 1; i >= 0; i-- {
		ip[i]++
		// if no overflow, we are done.
		if ip[i] > 0 {
			break
		}
	}
	return true
}

// GetNetworkIPRange returns the first and the last address of the network
func GetNetworkIPRange(n *net.IPNet) (first, last net.IP) {
	size := len(n.IP)
	first = make([]byte, size, size)
	last = make([]byte, size, size)
	for i := range n.IP {
		first[i] = n.IP[i] & n.Mask[i]
		last[i] = n.IP[i] | ^n.Mask[i]
	}
	return
}

// CompareIPs compares two ip addresses and returns 0 if they are equal,
// -1 if the first one preceeds the second, +1 if the first one is bigger that
// the second.
//
// If ip addresses has different sizes, an error is returned.
func CompareIPs(a, b net.IP) (int, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("IP addresses %v and %v have different sizes", a, b)
	}
	return bytes.Compare(a, b), nil
}

// IPRangeIterator allows you to iterate over a range of IP addresses
type IPRangeIterator interface {

	// Next returns the next ip address in the range and true if the next ip exists.
	// If it doesn't exits, the last ip and false are returned.
	Next() (ip net.IP, ok bool)
}

// GetIPRangeIterator returns an interator over IP range.  The last ip will be included
// into the sequence produced.
func GetIPRangeIterator(first, last net.IP) IPRangeIterator {
	return &ipRangeIterator{first, last, CopyIP(first)}
}

type ipRangeIterator struct {
	first net.IP
	last  net.IP
	next  net.IP
}

func (iter *ipRangeIterator) Next() (ip net.IP, ok bool) {
	result := CopyIP(iter.next)
	check, err := CompareIPs(iter.next, iter.last)
	if err == nil && check <= 0 {
		Next(iter.next)
		return result, true
	}
	return result, false
}

func (iter *ipRangeIterator) String() string {
	if res, _ := CompareIPs(iter.last, iter.next); res < 0 {
		return fmt.Sprintf("IPRangeIterator(%v -> %v, next: none)", iter.first, iter.last)
	}
	return fmt.Sprintf("IPRangeIterator(%v -> %v, next: %v)", iter.first, iter.last, iter.next)
}
