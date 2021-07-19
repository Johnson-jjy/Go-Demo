package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	//Given the above hash function,this will give replicas with "hashes":
	//2, 4, 6, 12, 14, 16,22, 24, 26
	hash.Add("6", "4", "2")

	testCase := map[string]string{
		"2": "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCase {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yield %s", k, v)
		}
	}

	//Adds 8, 18, 28
	hash.Add("8")

	//27 should now map to 8.
	testCase["27"] = "8"

	for k, v := range testCase {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yield %s", k, v)
		}
	}
}
