package libstar

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"
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

func FunName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func Netmask2Len(s string) int {
	mask := net.IPMask(net.ParseIP(s).To4())
	prefixSize, _ := mask.Size()

	return prefixSize
}

func PrettySecs(t uint64) string {
	min := t / 60
	if min < 60 {
		return fmt.Sprintf("%dm%ds", min, t%60)
	}
	hours := min / 60
	if hours < 24 {
		return fmt.Sprintf("%dh%dm", hours, min%60)
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

func ToBytes(v, u string) uint64 {
	value, _ := strconv.Atoi(v)
	switch u {
	case "Bytes":
		return uint64(value)
	case "KiB":
		return uint64(value) * 1024
	case "MiB":
		return uint64(value) * 1024 * 1024
	case "GiB":
		return uint64(value) * 1024 * 1024 * 1024
	case "TiB":
		return uint64(value) * 1024 * 1024 * 1024 * 1024
	default:
		return 0
	}
}

func ToKiB(v, u string) uint64 {
	value, _ := strconv.Atoi(v)
	switch u {
	case "Bytes":
		return uint64(value / 1024)
	case "KiB":
		return uint64(value)
	case "MiB":
		return uint64(value) * 1024
	case "GiB":
		return uint64(value) * 1024 * 1024
	case "TiB":
		return uint64(value) * 1024 * 1024 * 1024
	default:
		return 0
	}
}

type jsonUtils struct {
}

var JSON = jsonUtils{}

func (j jsonUtils) Marshal(v interface{}, pretty bool) (string, error) {
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

func (j jsonUtils) MarshalSave(v interface{}, file string, pretty bool) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		Error("MarshalSave: %s", err)
		return err
	}
	str, err := j.Marshal(v, true)
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

func (j jsonUtils) UnmarshalLoad(v interface{}, file string) error {
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

type xmlUtils struct {
}

var XML = xmlUtils{}

func (x xmlUtils) Marshal(v interface{}, pretty bool) (string, error) {
	if !pretty {
		str, err := xml.Marshal(v)
		if err != nil {
			Error("Marshal error: %s", err)
			return "", err
		}
		return string(str), nil
	} else {
		str, err := xml.MarshalIndent(v, "", "  ")
		if err != nil {
			Error("Marshal error: %s", err)
			return "", err
		}
		return string(str), nil
	}
}

func (x xmlUtils) MarshalSave(v interface{}, file string, pretty bool) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		Error("MarshalSave: %s", err)
		return err
	}
	str, err := x.Marshal(v, true)
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

func (x xmlUtils) UnmarshalLoad(v interface{}, file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return NewErr("UnmarshalLoad: file:<%s> does not exist", file)
	}
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return NewErr("UnmarshalLoad: file:<%s> %s", file, err)
	}
	if err := xml.Unmarshal([]byte(contents), v); err != nil {
		return NewErr("UnmarshalLoad: %s", err)
	}
	return nil
}

type Dir struct {
}

var DIR = Dir{}

func (d Dir) ListFiles(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}

func (d Dir) ListDirs(dirPth string) (dirs []string, err error) {
	dirs = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if !fi.IsDir() {
			continue
		}
		dirs = append(dirs, dirPth+PthSep+fi.Name())
	}
	return dirs, nil
}

func H2D8(h string) uint8 {
	v, err := strconv.ParseUint(h, 0, 8)
	if err != nil {
		Warn("H2D16: %s", err)
		return 0
	}
	return uint8(v)
}

func H2D16(h string) uint16 {
	v, err := strconv.ParseUint(h, 0, 16)
	if err != nil {
		Warn("H2D16: %s", err)
		return 0
	}
	return uint16(v)
}

func PrettyPCI(dom, bus, slot, fun uint16) string {
	return fmt.Sprintf("%04x:%02x:%02x.%x", dom, bus, slot, fun)
}

func PrettyDrive(ctl, bus, tgt, unit uint16) string {
	return fmt.Sprintf("%04x:%02x:%02x.%x", ctl, bus, tgt, unit)
}

func IsDigit(s string) bool {
	for _, v := range s {
		if unicode.IsDigit(v) {
			continue
		}
		return false
	}
	return true
}

func BasicAuth(username, password string) string {
	auth := username
	if password != "" {
		auth += ":" + password
	}
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func Wait() {
	x := make(chan os.Signal)
	signal.Notify(x, os.Interrupt, syscall.SIGTERM)
	signal.Notify(x, os.Interrupt, syscall.SIGKILL)
	signal.Notify(x, os.Interrupt, syscall.SIGQUIT) //CTL+/
	signal.Notify(x, os.Interrupt, syscall.SIGINT)  //CTL+C

	<-x
	fmt.Println("done")
}

func GetJSON(r io.ReadCloser, v interface{}) error {
	defer r.Close()
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	Debug("GetJSON %s", body)
	if err := json.Unmarshal([]byte(body), v); err != nil {
		return err
	}
	return nil
}
