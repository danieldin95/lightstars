package libstar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"reflect"
	"runtime"
	"time"
)

func GenToken(n int) string {
	letters := []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	buffer := make([]byte, n)

	size := len(letters)
	rand.Seed(time.Now().UnixNano())
	for i := range buffer {
		buffer[i] = letters[rand.Int63()%int64(size)]
	}

	return string(buffer)
}

func GenEthAddr(n int) []byte {
	if n == 0 {
		n = 6
	}

	data := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range data {
		data[i] = byte(rand.Uint32() & 0xFF)
	}

	data[0] &= 0xfe

	return data
}

func Marshal(v interface{}, pretty bool) (string, error) {
	str, err := json.Marshal(v)
	if err != nil {
		Error("Marshal error: %s", err)
		return "", err
	}

	if !pretty {
		return string(str), nil
	}

	var out bytes.Buffer

	if err := json.Indent(&out, str, "", "  "); err != nil {
		return string(str), nil
	}

	return out.String(), nil
}

func MarshalSave(v interface{}, file string, pretty bool) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		Error("MarshalSave: %s", err)
		return err
	}

	str, err := Marshal(v, true)
	if err != nil {
		Error("MarshalSave error: %s", err)
		return err
	}

	if _, err := f.Write([]byte(str)); err != nil {
		Error("MarshalSave: %s", err)
		return err
	}

	return nil
}

func UnmarshalLoad(v interface{}, file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return NewErr("UnmarshalLoad: file:<%s> does not exist", file)
	}

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return NewErr("UnmarshalLoad: file:<%s> %s", file, err)
	}

	if err := json.Unmarshal([]byte(contents), v); err != nil {
		return NewErr("UnmarshalLoad: %s", err)
	}

	return nil
}

func FunName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func Netmask2Len(s string) int {
	mask := net.IPMask(net.ParseIP(s).To4())
	prefixSize, _ := mask.Size()

	return prefixSize
}

func PrettySecs(t uint64) string {
	mins := t / 60
	if mins < 60 {
		return fmt.Sprintf("%dm%ds", mins, t%60)
	}
	hours := mins / 60
	if hours < 24 {
		return fmt.Sprintf("%dh%dm", hours, mins%60)
	}
	days := hours / 24
	return fmt.Sprintf("%dd%dh", days, hours%24)
}

func PrettyBytes(b uint64) string {
	split := func(_v uint64, _m uint64) (i uint64, d int) {
		v := float64(_v%_m) / float64(_m)
		return _v / _m, int(v * 100) //move two decimal to integer
	}

	if b < 1024 {
		return fmt.Sprintf("%dB", b)
	}
	k, d := split(b, 1024)
	if k < 1024 {
		return fmt.Sprintf("%d.%02dK", k, d)
	}
	m, d := split(k, 1024)
	if m < 1024 {
		return fmt.Sprintf("%d.%02dM", m, d)
	}
	g, d := split(m, 1024)
	return fmt.Sprintf("%d.%02dG", g, d)
}

func PrettyKBytes(k uint64) string {
	split := func(_v uint64, _m uint64) (i uint64, d int) {
		v := float64(_v%_m) / float64(_m)
		return _v / _m, int(v * 100) //move two decimal to integer
	}

	if k < 1024 {
		return fmt.Sprintf("%dK", k)
	}
	m, d := split(k, 1024)
	if m < 1024 {
		return fmt.Sprintf("%d.%02dM", m, d)
	}
	g, d := split(m, 1024)
	return fmt.Sprintf("%d.%02dG", g, d)
}
