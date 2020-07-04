package http

import (
	"github.com/danieldin95/lightstar/http/api"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/schema"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type FileInfo struct {
	FileName string
	Size     string
	ModTime  string
	IsDir    bool
}

type ListFileInfo struct {
	schema.List
	Items []FileInfo `json:"items"`
}

type condResult int

const (
	condNone condResult = iota
	condTrue
	condFalse
)

type fileHandler struct {
	root http.FileSystem
}

func FileServer(root http.FileSystem) http.Handler {
	return &fileHandler{root}
}

func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	serveFile(w, r, f.root, path.Clean(upath), true)
}

func serveFile(w http.ResponseWriter, r *http.Request, fs http.FileSystem, name string, redirect bool) {
	// TODO 判断文件类型
	libstar.Info("name is %s", name)
	f, err := fs.Open(name)
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}

	if redirect {
		// redirect to canonical path: / at end of directory url
		// r.URL.Path always begins with /
		url := r.URL.Path
		if d.IsDir() {
			if url[len(url)-1] != '/' {
				localRedirect(w, r, path.Base(url)+"/")
				return
			}
		} else {
			if url[len(url)-1] == '/' {
				localRedirect(w, r, "../"+path.Base(url))
				return
			}
		}
	}

	// redirect if the directory name doesn't end in a slash
	if d.IsDir() {
		url := r.URL.Path
		if url[len(url)-1] != '/' {
			localRedirect(w, r, path.Base(url)+"/")
			return
		}
	}

	// Still a directory? (we didn't find an index.html file)
	if d.IsDir() {
		if checkIfModifiedSince(r, d.ModTime()) == condFalse {
			writeNotModified(w)
			return
		}
		w.Header().Set("Last-Modified", d.ModTime().UTC().Format(http.TimeFormat))
		if r.Method == "POST" {
			dirList(w, r, f)
		} else {
			data := struct {
				Error string
			}{}

			file := api.GetFile("ui/files.html")
			if err := api.ParseFiles(w, file, data); err != nil {
				libstar.Error("files load %s", err)
			}
		}
		return
	}

	// serveContent will check modification time
	//sizeFunc := func() (int64, error) { return d.Size(), nil }
	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

func dirList(w http.ResponseWriter, r *http.Request, f http.File) {
	dirs, err := f.Readdir(-1)

	if err != nil {
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}
	list := ListFileInfo{
		Items: make([]FileInfo, 0, 32),
	}

	for _, d := range dirs {
		info := &FileInfo{
			d.Name(),
			formatFileSize(d.Size()),
			GetTime(d.ModTime()),
			d.IsDir(),
		}
		libstar.Info("%s", d.ModTime().Year())
		list.Items = append(list.Items, *info)
	}
	sort.Slice(list.Items, func(i, j int) bool {
		return list.Items[i].FileName < list.Items[j].FileName
	})
	list.Metadata.Size = len(list.Items)
	list.Metadata.Total = len(list.Items)
	api.ResponseJson(w, list)
}

func localRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}
	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)
}

func toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}
	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}

func writeNotModified(w http.ResponseWriter) {
	// RFC 7232 section 4.1:
	// a sender SHOULD NOT generate representation metadata other than the
	// above listed fields unless said metadata exists for the purpose of
	// guiding cache updates (e.g., Last-Modified might be useful if the
	// response does not have an ETag field).
	h := w.Header()
	delete(h, "Content-Type")
	delete(h, "Content-Length")
	if h.Get("Etag") != "" {
		delete(h, "Last-Modified")
	}
	w.WriteHeader(http.StatusNotModified)
}

func checkIfModifiedSince(r *http.Request, modtime time.Time) condResult {
	if r.Method != "GET" && r.Method != "HEAD" {
		return condNone
	}
	ims := r.Header.Get("If-Modified-Since")
	if ims == "" || isZeroTime(modtime) {
		return condNone
	}
	t, err := http.ParseTime(ims)
	if err != nil {
		return condNone
	}
	// The Date-Modified header truncates sub-second precision, so
	// use mtime < t+1s instead of mtime <= t to check for unmodified.
	if modtime.Before(t.Add(1 * time.Second)) {
		return condFalse
	}
	return condTrue
}

var unixEpochTime = time.Unix(0, 0)

func isZeroTime(t time.Time) bool {
	return t.IsZero() || t.Equal(unixEpochTime)
}

func formatFileSize(size int64) string {
	s := int(size)
	KB := 1024
	MB := 1024 * KB
	GB := 1024 * MB

	if s < KB {
		return strconv.Itoa(s)
	} else if KB < s && s < MB {
		return strconv.Itoa(s/KB) + "KB"
	} else if MB < s && s < GB {
		return strconv.Itoa(s/MB) + "MB"
	} else {
		return strconv.Itoa(s/GB) + "GB"
	}
}

func GetTime(t time.Time) string {
	tmp1 := "2006-01-02 15:04:05"
	return t.Format(tmp1)
}
