package libstar

var (
	Version string
	Date    string
	Commit  string
)

func init() {
	Info("version is %s", Version)
	Info("built on %s", Date)
	Info("commit at %s", Commit)
}
