package virsh

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Run(uri string, args ...string) (string, error) {
	cmdArgs := make([]string, 0, len(args)+2)
	if uri != "" {
		cmdArgs = append(cmdArgs, "--connect", uri)
	}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("virsh", cmdArgs...)
	out, err := cmd.CombinedOutput()
	text := strings.TrimSpace(string(out))
	if err != nil {
		if text == "" {
			return "", fmt.Errorf("virsh %s: %w", strings.Join(cmdArgs, " "), err)
		}
		return "", fmt.Errorf("virsh %s: %v: %s", strings.Join(cmdArgs, " "), err, text)
	}
	return string(out), nil
}

func TempXML(prefix, data string) (string, func(), error) {
	f, err := os.CreateTemp("", prefix+"-*.xml")
	if err != nil {
		return "", nil, err
	}
	if _, err := f.WriteString(data); err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return "", nil, err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(f.Name())
		return "", nil, err
	}
	cleanup := func() {
		_ = os.Remove(f.Name())
	}
	return f.Name(), cleanup, nil
}

func ParseKV(text string) map[string]string {
	ret := make(map[string]string, 16)
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, ":")
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(strings.ToLower(line[:idx]))
		val := strings.TrimSpace(line[idx+1:])
		ret[key] = val
	}
	return ret
}

func ParseBytes(text string) uint64 {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return 0
	}
	v, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0
	}
	unit := "B"
	if len(parts) > 1 {
		unit = strings.ToUpper(strings.TrimSpace(parts[1]))
	}
	factor := float64(1)
	switch unit {
	case "B", "BYTE", "BYTES":
		factor = 1
	case "KB":
		factor = 1000
	case "MB":
		factor = 1000 * 1000
	case "GB":
		factor = 1000 * 1000 * 1000
	case "TB":
		factor = 1000 * 1000 * 1000 * 1000
	case "KIB":
		factor = 1024
	case "MIB":
		factor = 1024 * 1024
	case "GIB":
		factor = 1024 * 1024 * 1024
	case "TIB":
		factor = 1024 * 1024 * 1024 * 1024
	}
	return uint64(math.Round(v * factor))
}

func ParseSecondsNS(text string) uint64 {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}
	text = strings.TrimSuffix(text, "s")
	v, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
	if err != nil {
		return 0
	}
	return uint64(v * 1000 * 1000 * 1000)
}

func Lines(text string) []string {
	rows := strings.Split(text, "\n")
	ret := make([]string, 0, len(rows))
	for _, row := range rows {
		row = strings.TrimSpace(row)
		if row == "" {
			continue
		}
		ret = append(ret, row)
	}
	return ret
}
