package storage

import "strings"

const Location = "/lightstar/"
const DataStore = "datastore@"

type StorePath struct{}

var PATH = StorePath{}

func (p StorePath) Fmt(path string) string {
	newPath := path

	// datastore
	if strings.HasPrefix(path, Location) {
		newPath = strings.Split(path, Location)[1]
		if strings.HasPrefix(newPath, "datastore") {
			splits := strings.SplitN(newPath, "/", 3)
			if len(splits) == 2 {
				newPath = splits[0] + "@" + splits[1]
			} else if len(splits) > 2 {
				newPath = splits[0] + "@" + splits[1] + ":/" + splits[2]
			}
		}
		return newPath
	}
	// dev
	if strings.HasPrefix(path, "/dev/") {
		newPath = strings.Split(path, "/dev")[1]
		return "device:" + newPath
	}

	return newPath
}

func (p StorePath) Unix(path string) string {
	newPath := path
	if strings.HasPrefix(path, "datastore@") {
		splits := strings.SplitN(newPath, "/", 2)
		if len(splits) >= 2 {
			stores := strings.SplitN(splits[0], "@", 2)
			suffix := splits[1]
			storeN := strings.TrimRight(stores[1], ":")
			newPath = Location + stores[0] + "/" + storeN + "/" + suffix
		} else {
			stores := strings.SplitN(splits[0], "@", 2)
			newPath = Location + stores[0] + "/" + strings.TrimRight(stores[1], ":")
		}
	}
	if strings.HasPrefix(path, "device:/") {
		splits := strings.SplitN(newPath, "device:", 2)
		return "/dev" + splits[1]
	}

	return newPath
}

func (p StorePath) Root() string {
	return Location + "datastore" + "/"
}

func (p StorePath) RootXML() string {
	return Location + "datastore" + "/" + "libvirt/"
}
