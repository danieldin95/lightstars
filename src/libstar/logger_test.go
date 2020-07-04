package libstar

import (
	"fmt"
	"testing"
)

func TestLogger_List(t *testing.T) {
	Error("hhh")
	for m := range Log.List() {
		if m == nil {
			break
		}
		fmt.Println(m)
	}
	Warn("xxx")
	for m := range Log.List() {
		if m == nil {
			break
		}
		fmt.Println(m)
	}
	Info("nnnnn")
	for m := range Log.List() {
		if m == nil {
			break
		}
		fmt.Println(m)
	}
}
