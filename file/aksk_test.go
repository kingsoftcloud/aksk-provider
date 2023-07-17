package file

import (
	"testing"
)

const (
	cipherKey        = "8cca0smDmR478v8F"
	rawAK            = "AKTRQxqRY0SdCw31S46rrcMA"
	rawSK            = "ODPedeQvrIo2BF6QkzkZ1HZdhkjH648cOF0fVXGt"
	rawSecurityToken = "V2qaSoc7yqHMt4mkuSkh2FAhVIiQkqMQ6+7IDSeTHCPx95MXJRtMRz1ArWrIXgalkEIywuFUlit9RlWbLiuonoECvCj7Hh5QdasUFIcCVaGviotq++9mUdVX3n6wUFi7hfvb/Trq0R0Tkq5R7ysqgS6irGUrwZi11vUGmWiK4ISVVmTT3kzR7nS9P6kav0uLzboK1YnJwShRgzrknr3jAer7P80RjwPAv0lDutD1p7D1Dp+WdJ1Hy13wC16pN9+xXBKwgyID6e5unL9JLkL+ixvABEdVt/g3B2eJQph9BLCZzfRTmsuai0nGL+EcRwUqGvjIV7U5NbvGkJ/w4J9FiBl8bHK4azGT4MgyydoCYvPh6gzTi3S4BzaYPaI4kdNLBx3i1/S9KukpOaLVw53NyTB27IgxbxNgbN6e4t4V6CR6/AqxWzy6GzlN5BV1Y6chwMhgUUaWL/KPjgBoB80/cMAgHPtHvmhQkeI+UDSRire/SxjgCFf4Dd9yW8E6JJB9b10F9nUNK9EKtSL7KqWuGH1T+qG7gYVrBs4BB7jQdG2v2CSEW0DyVtlUZGM2tksvrffQ+5vx1+ycKKqoIrmN93TWF/qOHOed3sZxRCuSoqAxFNXNZvKHOn/xyzPR9XBGwXy3Lw0f+2CdlSgXYrvnHcuA=="
)

func TestFileProvider(t *testing.T) {
	provider := NewFileAKSKProvider("./aksk", cipherKey)
	aksk, err := provider.GetAKSK()
	if err != nil {
		t.Errorf("get aksk failed: %v", err)
	}

	if aksk.AK != rawAK {
		t.Errorf("get ak failed")
	}
	if aksk.SK != rawSK {
		t.Errorf("get sk failed")
	}
	if aksk.SecurityToken != rawSecurityToken {
		t.Errorf("get securityToken failed")
	}
}
