// SPDX-License-Identifier: MIT-0

package iputils

import (
	"fmt"
	"net"
	"testing"
)

func TestNext(t *testing.T) {
	type testCase struct {
		input net.IP
		next  net.IP
		ok    bool
	}
	cases := []testCase{
		testCase{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.2"), true},
		testCase{[]byte{192, 168, 0, 1}, net.ParseIP("192.168.0.2"), true},
		testCase{net.IPv4(192, 168, 0, 1), net.ParseIP("192.168.0.2"), true},
		testCase{net.ParseIP("192.168.0.255"), net.ParseIP("192.168.1.0"), true},
		testCase{net.ParseIP("255.255.255.255"), net.ParseIP("255.255.255.255"), false},
		testCase{[]byte{255, 255, 255, 255}, net.ParseIP("255.255.255.255"), false},
		testCase{net.IPv4(255, 255, 255, 255), net.ParseIP("255.255.255.255"), false},
		testCase{net.ParseIP("::"), net.ParseIP("::1"), true},
		testCase{net.ParseIP("::1"), net.ParseIP("::2"), true},
		testCase{net.ParseIP("::ffff:ffff:ffff"), net.ParseIP("::ffff:ffff:ffff"), false},
		testCase{net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
			net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"), false},
	}
	for _, test := range cases {
		next := CopyIP(test.input)
		ok := Next(next)
		if test.ok != ok || !test.next.Equal(next) {
			t.Errorf("expecting (%v, %v), got (%v, %v)", test.next, test.ok, next, ok)
		}
	}
}

func TestGetIPRange(t *testing.T) {
	type testCase struct {
		network string
		first   net.IP
		last    net.IP
	}
	cases := []testCase{
		testCase{"192.168.0.0/24", net.ParseIP("192.168.0.0"), net.ParseIP("192.168.0.255")},
		testCase{"192.168.0.0/28", net.ParseIP("192.168.0.0"), net.ParseIP("192.168.0.15")},
		testCase{"192.168.0.0/32", net.ParseIP("192.168.0.0"), net.ParseIP("192.168.0.0")},
		testCase{"192.168.216.192/28", net.ParseIP("192.168.216.192"), net.ParseIP("192.168.216.207")},
	}
	for _, test := range cases {
		_, ipNet, err := net.ParseCIDR(test.network)
		if err != nil {
			t.Errorf("failed to parse network cidr %v: %v", test.network, err)
			continue
		}
		first, last := GetNetworkIPRange(ipNet)
		if !test.first.Equal(first) || !test.last.Equal(last) {
			t.Errorf("expecting (%v, %v), got (%v, %v)", test.first, test.last, first, last)
		}
	}
}

func TestCompareIPs(t *testing.T) {
	type testCase struct {
		a      net.IP
		b      net.IP
		result int
	}
	cases := []testCase{
		testCase{net.ParseIP("192.168.0.0"), net.ParseIP("192.168.0.1"), -1},
		testCase{net.ParseIP("192.168.0.0"), net.ParseIP("192.168.0.0"), 0},
		testCase{net.ParseIP("192.168.0.1"), net.ParseIP("192.168.0.0"), 1},
		testCase{net.ParseIP("192.168.0.0"), net.ParseIP("0.0.0.0"), 1},
		testCase{net.ParseIP("0.0.0.0"), net.ParseIP("255.255.255.255"), -1},
	}
	for _, test := range cases {
		result, err := CompareIPs(test.a, test.b)
		if err != nil {
			t.Errorf("unexpected error %v when comparing ip addresses %v and %v", err, test.a, test.b)
		}
		if test.result != result {
			t.Errorf("expecting %v, got %v when comparing ip addresses %v and %v", test.result, result, test.a, test.b)
		}
	}
}

func TestCompareIPsFaults(t *testing.T) {
	type faultCase struct {
		a net.IP
		b net.IP
	}
	faultCases := []faultCase{
		faultCase{[]byte{255, 255, 255, 255}, net.ParseIP("255.255.255.255")},
	}
	for _, test := range faultCases {
		_, err := CompareIPs(test.a, test.b)
		if err == nil {
			t.Errorf("didn't get an error when comparing %v and %v", test.a, test.b)
		}
	}
}

func TestIPRangeIterator(t *testing.T) {
	type testCase struct {
		first    net.IP
		last     net.IP
		sequence []net.IP
	}
	cases := []testCase{
		testCase{net.ParseIP("192.168.0.0"), net.ParseIP("192.168.0.1"),
			[]net.IP{net.ParseIP("192.168.0.0"), net.ParseIP("192.168.0.1")}},
		testCase{net.ParseIP("192.168.0.0"), net.ParseIP("192.168.0.0"), []net.IP{net.ParseIP("192.168.0.0")}},
		testCase{net.ParseIP("192.168.0.20"), net.ParseIP("192.168.0.10"), []net.IP{}},
	}
NEXT_CASE:
	for _, test := range cases {
		iter := GetIPRangeIterator(test.first, test.last)
		for i := 0; i < len(test.sequence); i++ {
			value, ok := iter.Next()
			if !ok {
				t.Errorf("iterator %v has not produced enough values; expecting sequence %v", iter, test.sequence)
				continue NEXT_CASE
			}

			if !test.sequence[i].Equal(value) {
				t.Errorf("iteration %v of %v produced %v, expecting %v", i+1, iter, value, test.sequence[i])
				continue NEXT_CASE
			}
		}
		_, ok := iter.Next()
		if ok {
			t.Errorf("iterator %v has produced more values than expected", iter)
			continue NEXT_CASE
		}
	}
}

func TestIPRangeIteratorStringConvertion(t *testing.T) {
	type testCase struct {
		first   net.IP
		last    net.IP
		results []string
	}
	cases := []testCase{
		testCase{
			net.ParseIP("192.168.0.0"),
			net.ParseIP("192.168.0.1"),
			[]string{
				fmt.Sprintf("IPRangeIterator(192.168.0.0 -> 192.168.0.1, next: 192.168.0.0)"),
				fmt.Sprintf("IPRangeIterator(192.168.0.0 -> 192.168.0.1, next: 192.168.0.1)"),
				fmt.Sprintf("IPRangeIterator(192.168.0.0 -> 192.168.0.1, next: none)"),
			},
		},
		testCase{
			net.ParseIP("::100"),
			net.ParseIP("::102"),
			[]string{
				fmt.Sprintf("IPRangeIterator(::100 -> ::102, next: ::100)"),
				fmt.Sprintf("IPRangeIterator(::100 -> ::102, next: ::101)"),
				fmt.Sprintf("IPRangeIterator(::100 -> ::102, next: ::102)"),
				fmt.Sprintf("IPRangeIterator(::100 -> ::102, next: none)"),
			},
		},
	}
NEXT_CASE:
	for _, test := range cases {
		iter := GetIPRangeIterator(test.first, test.last)
		for i := 0; i < len(test.results); i++ {
			if test.results[i] != fmt.Sprintf("%v", iter) {
				t.Errorf("after %v iterations expecting %v, got %v", i, test.results[i], iter)
				continue NEXT_CASE
			}
			iter.Next()
		}
	}
}

func ExampleNext() {
	ip := net.ParseIP("192.168.0.1")
	ok := Next(ip)
	fmt.Println(ok, ip)

	ip = net.ParseIP("::100") // IPv6 address
	ok = Next(ip)
	fmt.Println(ok, ip)

	ip = net.ParseIP("255.255.255.255")
	ok = Next(ip)
	fmt.Println(ok, ip)

	// Output:
	// true 192.168.0.2
	// true ::101
	// false 255.255.255.255
}

func ExampleGetNetworkIPRange() {
	_, network, _ := net.ParseCIDR("192.168.0.0/28")
	fmt.Println(GetNetworkIPRange(network))

	_, network, _ = net.ParseCIDR("beef::/64") // IPv6 network
	fmt.Println(GetNetworkIPRange(network))

	// Output:
	// 192.168.0.0 192.168.0.15
	// beef:: beef::ffff:ffff:ffff:ffff
}

func ExampleIPRangeIterator_Next() {
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
}
