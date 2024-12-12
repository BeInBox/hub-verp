package verp

import (
	"errors"
	"testing"
)

func TestCreateOneAndCheckDecode(t *testing.T) {
	publicIp = GetPublicIp()
	Token = "32fd3f3490b91945d6f76d5c9da910e6ab7757c8427e83c5232274747c0e7486"
	verp := NewMinVerp(1)
	verp = NewMinVerp(1)
	uniqueID := verp.Encode()
	newVerp := DecodeMiniVerp(uniqueID)
	newVerp.Encode()
	if verp.Chain != newVerp.Chain {
		t.Fatalf("Init failled  = %v", errors.New("verp Missmatch"))
	}
}
func TestDecode(t *testing.T) {
	newVerp := DecodeMiniVerp("dadc8c30baa9f59bb4adedfc0dXXXX")
	newVerp.Encode()
	if newVerp.Inc != 2 || newVerp.User != 1 {
		t.Fatalf("Init failled  = %v", errors.New("verp Missmatch"))
	}

}
