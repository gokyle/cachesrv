package main

import (
	"flag"
	"fmt"
	"github.com/gokyle/filecache"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

const VERSION = "3.0.0"

var (
	cache        *filecache.FileCache
	path_regexp  = regexp.MustCompile("^/(.*)$")
	space_regexp = regexp.MustCompile(" ")
	srv_bin      string
)

func init() {
	cache = filecache.NewCache()
	filecache.DefaultExpireItem = 0
	cache.Start()
}

func Path(r *http.Request) (path string) {
	path = path_regexp.ReplaceAllString(r.URL.Path, "./$1")
	if len(path) == 0 {
		path = "."
	}
	return path
}

func Dispatch(w http.ResponseWriter, r *http.Request) {
	if r.URL.RawQuery == "cachestats" {
		go displayCacheStats()
	}
	fi, err := os.Stat(Path(r))
	if err != nil {
		http.ServeFile(w, r, Path(r))
	} else if fi.IsDir() {
		http.ServeFile(w, r, Path(r))
	} else if fi.Size() > cache.MaxSize {
		http.ServeFile(w, r, Path(r))
	} else {
	        cache.HttpWriteFile(w, r)
	}
	return
}

func displayCacheStats() {
	fmt.Printf("-----[ cache stats: %s ]-----\n",
		time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("files cached: ", cache.Size())
	fmt.Printf("cache size: %d bytes\n", cache.FileSize())
	cachedFiles := cache.StoredFiles()
	fmt.Println("[ cached files ]")
	for _, name := range cachedFiles {
		fmt.Printf("\t* %s\n", name)
	}
	fmt.Printf("\n\n")
}

func main() {
	var (
		srv_port int
		srv_ssl  bool
		srv_wd   = "."
	)

	srv_bin = filepath.Base(os.Args[0])
	empty := func(s string) bool {
		return len(s) == 0
	}

	fCert := flag.String("c", "", "TLS certificate file")
	fExpire := flag.Int("e", filecache.DefaultExpireItem,
		"maximum number of seconds between accesses a file can stay "+
			"in the cache")
	fGarbage := flag.Int("g", filecache.DefaultEvery,
		"scan the cache for expired items every <n> seconds")
	fKey := flag.String("k", "", "TLS key file")
	fItems := flag.Int("n", filecache.DefaultMaxItems, "max number of "+
		"files to store in the cache")
	fPort := flag.Int("p", 8080, "port to listen on")
	fChroot := flag.Bool("r", false, "chroot to the working directory")
	fSize := flag.Int64("s", filecache.DefaultMaxSize,
		"max file size to cache")
	fUser := flag.String("u", "", "user to run as")
	fVersion := flag.Bool("v", false, "print version information")
	flag.Parse()

	if *fVersion {
		version()
	}

	if flag.NArg() > 0 {
		srv_wd = flag.Arg(0)
	}

	if *fChroot {
		srv_wd = chroot(srv_wd)
	}

	if !empty(*fUser) {
		setuid(*fUser)
	}

	if !empty(*fCert) && !empty(*fKey) {
		srv_ssl = true
	}

	cache.MaxItems = *fItems
	cache.ExpireItem = *fExpire
	cache.Every = *fGarbage
	cache.MaxSize = *fSize
	srv_wd, err := filepath.Abs(srv_wd)
	checkFatal(err)
	err = os.Chdir(srv_wd)
	checkFatal(err)
	srv_port = *fPort
	srv_addr := fmt.Sprintf(":%d", srv_port)
	fmt.Printf("serving %s on %s\n", srv_wd, srv_addr)
	http.HandleFunc("/", Dispatch)
	cache.Start()
	if srv_ssl {
		log.Fatal(http.ListenAndServeTLS(srv_addr, *fCert, *fKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(srv_addr, nil))
	}
}

func checkFatal(err error) {
	if err != nil {
		fmt.Printf("[!] %s: %s\n", srv_bin, err.Error())
		os.Exit(1)
	}
}

func setuid(username string) {
	usr, err := user.Lookup(username)
	checkFatal(err)
	uid, err := strconv.Atoi(usr.Uid)
	checkFatal(err)
	err = syscall.Setreuid(uid, uid)
	checkFatal(err)
}

func chroot(path string) string {
	err := syscall.Chroot(path)
	checkFatal(err)
	return "/"
}

func version() {
	fmt.Printf("%s version %s\n", srv_bin, VERSION)
	os.Exit(0)
}
