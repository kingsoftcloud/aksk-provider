package env

import (
	"os"
	"testing"
)

const (
	cipherKey        = "8cca0smDmR478v8F"
	rawAK            = "AKTRgGLIy2oS0SMnC2dHo8cU"
	rawSK            = "ODPedeQvrIo2BF6QkzkZ1HZdhkjH648cOF0fVXGt"
	rawSecurityToken = "V2qaSoc7yqHMt4mkuSkh2FAhVIiQkqMQ6+7IDSeTHCPx95MXJRtMRz1ArWrIXgalkEIywuFUlit9RlWbLiuonoECvCj7Hh5QdasUFIcCVaGviotq++9mUdVX3n6wUFi7hfvb/Trq0R0Tkq5R7ysqgS6irGUrwZi11vUGmWiK4ISVVmTT3kzR7nS9P6kav0uLzboK1YnJwShRgzrknr3jAer7P80RjwPAv0lDutD1p7D1Dp+WdJ1Hy13wC16pN9+xXBKwgyID6e5unL9JLkL+ixvABEdVt/g3B2eJQph9BLCZzfRTmsuai0nGL+EcRwUqGvjIV7U5NbvGkJ/w4J9FiBl8bHK4azGT4MgyydoCYvPh6gzTi3S4BzaYPaI4kdNLBx3i1/S9KukpOaLVw53NyTB27IgxbxNgbN6e4t4V6CR6/AqxWzy6GzlN5BV1Y6chwMhgUUaWL/KPjgBoB80/cMAgHPtHvmhQkeI+UDSRire/SxjgCFf4Dd9yW8E6JJB9b10F9nUNK9EKtSL7KqWuGH1T+qG7gYVrBs4BB7jQdG2v2CSEW0DyVtlUZGM2tksvrffQ+5vx1+ycKKqoIrmN93TWF/qOHOed3sZxRCuSoqAxFNXNZvKHOn/xyzPR9XBGwXy3Lw0f+2CdlSgXYrvnHcuA=="
	aesSK            = "hFL/syRB0pIvDMk1wWDS/xHS9paLPEaHWrH8JFL6U19DN/vnsHIlI+kPTrUNGI+1"
	aesSecurityToken = "OI70w98RsPE/hwbcUTNUQ+koDRglhbDV0OJvTDgTPWtRsfpp2tNAk9EsRTlZO0YKk0FxF19VpkZv/GJKHzHzy3wFmO9TvksXymU5hke0caDJtygsAHzYaDVIw5AJzbnOKnyenGNbJDtWGPsbyQqRKKJSstjpHKfDRmKq+28Xgo1fzLdfZVbkirL1qPDfpmO9W005omxWRVu1ZKmqjB0hx3duX9E4u6zVFko9jOKXjj+O4VpbEMhstSTCu7Z2lqe6LBv1KKB/Gshng1eExl3dPqg3CqDHMgd1Z+roMe1nzhdeLY7ZWXZKJS6hEu6bAsXDr74dBc8DIZUqY5yd84xhbG1Eqrk5WWPu6VOm9zYcqV89TvBZRkw1KCvWXyjgpLbLeB66nf+b4ea2FTYb1ARa8wr8TpGqzaE6GGE6B8q8LzWNsAK4sR/mwEbvOWgMdMeOQF9VnDkJsD1sY1C0NZ5FnrB9n79e77rksqjaovjaDc9TYgtsZJgrXWi296u4pvQhrgc9ZJvaclIDvL9aBQn30sr1Q9iVhFU7viXFElpxoqh0ySpCRjCkBSmRv1IYZ6Jnbg/AFz+fJ7P0LYYS42JFE9BYJU2hYVwtLJZ/GuZloLfRpnzvwswaeeYc5BOV64lfzqbO31ILQxfQDTqY01AD6/bL/JgkYhtq/oALOs3gSVBeM4OJd338JDMIrPaAYn7pj4a1sAG3B7HXaVDvRNwU/jflviWIFqNpW1nSnr/qGpbAqQrHXv1Rq59GC9FgszT7w7v1vrUgPH+LBzM44uhevU1c0aCm3X30wanteOddTWUOfP8CdEZL1I+68h1X2eL9J3zbiwEEiG9kzROvy3kHiGzUbYHHK9z8TE+0JaQQJ8c7lUk9Ud5SqS4ONeJVMOs6"
)

func TestEnvProviderRaw(t *testing.T) {
	os.Setenv("AK", rawAK)
	os.Setenv("SK", rawSK)
	os.Setenv("SECURITY_TOKEN", rawSecurityToken)
	provider := NewEnvAKSKProvider(false, "")
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

func TestEnvProviderEncrypt(t *testing.T) {
	os.Setenv("AK", rawAK)
	os.Setenv("SK", aesSK)
	os.Setenv("SECURITY_TOKEN", aesSecurityToken)
	provider := NewEnvAKSKProvider(true, cipherKey)
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
