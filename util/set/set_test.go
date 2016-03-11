package set

import (
	"testing"
	"fmt"
)

func TestSet(t *testing.T) {
	set := New()
	set.Add("string")
	if !set.Has("string") {
		t.Errorf("wrong")
	}
	set.Add("string")
	set.Add("foo")
	set.Add("bar")
	if set.Len() != 3 {
		t.Errorf("len is %v", set.Len())
	}

	ipList := set.List()
	for _,item := range ipList {
		fmt.Printf("item = %v", item)
	}
}