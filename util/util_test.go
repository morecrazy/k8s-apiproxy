package util

import (
	"testing"
	"encoding/json"
)

func TestBytesToDNSRequestJSON(t *testing.T)  {
	ipRaw1 := "192.168.1.10\n192.168.1.11\n192.168.1.12"
	ipRaw2 := "192.168.1.10\n192.168.1.10"


	cases := []struct{
		in_arg1, in_arg2, want string
	}{
		{ipRaw1, "test", "{\"url\":\"test\",\"rrs\":[{\"host\":\"192.168.1.10\"},{\"host\":\"192.168.1.11\"},{\"host\":\"192.168.1.12\"}]}"},
		{ipRaw2, "test", "{\"url\":\"test\",\"rrs\":[{\"host\":\"192.168.1.10\"}]}"},
	}

	for _, c := range cases {
		ips := []byte(c.in_arg1)
		reqJson := BytesToDNSRequestJSON(ips, c.in_arg2)
		b, _ := json.Marshal(reqJson)
		got := string(b)
		if c.want != got {
			t.Errorf("Revert(%q) == %q, want %q", c.in_arg1, got, c.want)
		}
	}
}
