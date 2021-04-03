package provider

import (
	"hash/crc32"
	"net"
)

// Credits
// https://github.com/hashicorp/terraform-provider-dns/blob/main/internal/hashcode/hashcode.go
// https://github.com/hashicorp/terraform-provider-dns/blob/main/internal/provider/hash_ip_string.go

// String hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func hashcodeString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

// Strings hashes a list of strings to a unique hashcode.
/*
func hashcodeStrings(strings []string) string {
	var buf bytes.Buffer

	for _, s := range strings {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", hashcodeString(buf.String()))
}
*/

func hashIPString(v interface{}) int {
	addr := v.(string)
	ip := net.ParseIP(addr)
	if ip == nil {
		return hashcodeString(addr)
	}
	return hashcodeString(ip.String())
}
