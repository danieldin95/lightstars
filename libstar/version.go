package libstar

import "fmt"

var (
	Version string
	Date    string
	Commit  string
)

func init() {
	fmt.Printf("%s %s %s\n", Version, Date, Commit)
}