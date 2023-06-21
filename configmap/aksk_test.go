package configmap

import (
	"testing"
)

const (
	rawAK            = "AKTRgGLIy2oS0SMnC2dHo8cU"
	rawSK            = "OHeJQvMoubjsILI0aYqc5N1EBGI8bX2LJk0sao4J"
	rawSecurityToken = "V2hSWzZKucRnHdpO3GBEoin2l3rGt4BIQ+8RgxlhBGfHvNK3KfNCoBKowGsQ0IGRXlyBqKw6XNHP4XxqNOHydS4hMacJJfJsLRoqbjv3j0PZazqin7jwRW/kYxAc6UXkGkfxyS0sWXO8LRcA6GCtT29i1q6LTWtSfiuEVK+BsS0RxDfHt/Wd22qXUNr6eVey0vIhY4lp3fQIY0osFzJZW59trZno4s+amefB2FgOP0T88OJe7Nn9/xe5BtPL2ILT717dY0hzXvR4bpJaOWD2U2c0qvO1+YSkJleJZP7YmIcUBerR0NZ6rDn41KLw1sB9k4jLcYTtLfI/CCLkeQOPsa58GJBXo51fBJz0d6ZlmVnU+GoHjIpUYVMcYfktgA38REF4MNDbFEs9NxILfflisZmouLAtdV2j3ma9rFRlKC7XZyQeLNewAcv78OK3wRmI/cchiuJJml7BnS+x8y4VrQ8bCGuqlLp500eWhh1lyi+3ZmIjxoAjw/jqngxBglsP1zpLZbDg6fqMziVnyiY9QIFIqf1jYUkcwO0CYOyIwirP7mXsl0w+ISMdA24MtT/AYDbbpIcoLV2kl4SyZwiwo0zfOXLmz/fepMsa8ZJay/gc8F/F2rh2Y/M8RzepgF2b6px/uuzuXESXER7Hdqj5Ecow=="
)

func TestCMProvider(t *testing.T) {
	provider := NewCMAKSKProvider("./aksk")
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
