package maps

import "testing"

func TestSafeMap_Mapi(t *testing.T) {
	runMapiTests[SafeMap[string, int]](t)
}
