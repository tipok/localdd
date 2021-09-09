package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	log "github.com/go-pkgz/lgr"
)

const DefaultConfigFolder = "/localdd/"

type Proxy struct {
	Url    *url.URL
	Proxy  *httputil.ReverseProxy
	Domain string
}

func createDirIfNotExists(dir string) error {
	s, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("could not create folder: %w", err)
		}
		s, err = os.Stat(dir)
	}

	if !s.IsDir() {
		return fmt.Errorf("%s exists but is not a directory", dir)
	}
	return nil
}

func createProxy(name, content string) (*Proxy, error) {
	fileContent := strings.TrimSpace(content)
	if fileContent == "" {
		return nil, fmt.Errorf("empty file")
	}

	p, err := strconv.Atoi(fileContent)
	if err == nil {
		fileContent = fmt.Sprintf("http://127.0.0.1:%d", p)
	}

	if !strings.Contains(fileContent, "//") {
		fileContent = fmt.Sprintf("http://%s", fileContent)
	}

	u, err := url.Parse(fileContent)
	if err != nil {
		return nil, fmt.Errorf("could not parse url: %w", err)
	}

	return &Proxy{Proxy: httputil.NewSingleHostReverseProxy(u), Url: u, Domain: name}, nil
}

func watchConfigFolder(ctx context.Context, confDir string, proxies *[]*Proxy) {
	ticker := time.NewTicker(500 * time.Millisecond)

	updateProxies := func() {
		files, err := ioutil.ReadDir(confDir)
		if err != nil {
			log.Fatalf("[ERROR] could not list files in folder %s", confDir)
		}
		var ps []*Proxy
		for _, file := range files {
			if !file.IsDir() {
				filePath := path.Join(confDir, file.Name())
				dc, err := ioutil.ReadFile(filePath)
				if err != nil {
					log.Printf("[ERROR] could not read config file %s: %v", filePath, err)
					continue
				}

				proxy, err := createProxy(file.Name(), string(dc))
				if err != nil {
					log.Fatalf("[ERROR] could not load proxy definition: %v", err)
					continue
				}
				ps = append(ps, proxy)
			}
		}
		if len(ps) > 0 {
			*proxies = ps
		}
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				log.Printf("[DEBUG] update proxies: %s", t)
				updateProxies()
			}
		}
	}()
}

func requestHandler(ps *[]*Proxy) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rhu, err := url.Parse("//" + r.Host)
		if err != nil {
			log.Fatalf("[ERROR] invalid host in request: %v", err)
			http.Error(w, "Invalid Host", http.StatusBadRequest)
		}
		rhn := rhu.Hostname()
		for _, p := range *ps {
			if p.Domain == rhn {
				r.Host = p.Url.Host
				p.Proxy.ServeHTTP(w, r)
				return
			}
		}
		log.Printf("[ERROR] no target found for domain %s", rhn)
		http.Error(w, "No Target Found", http.StatusBadGateway)
	}
}

func main() {
	logOpts := []log.Option{log.Msec, log.LevelBraces, log.CallerFile, log.CallerFunc}

	var listen string
	var confDir string
	var debug bool
	flag.StringVar(&listen, "listen", "127.0.0.1:20559", "Address to listen on")
	flag.StringVar(&confDir, "confDir", "", "Config directory path")
	flag.BoolVar(&debug, "debug", false, "Enable debug logs")
	flag.Parse()

	if debug {
		logOpts = append(logOpts, log.Debug)
	}

	log.Setup(logOpts...)

	sig := make(chan os.Signal, 1)
	signal.Notify(
		sig,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer signal.Stop(sig)

	if confDir == "" {
		userConfDir, _ := os.UserConfigDir()
		if userConfDir == "" {
			userConfDir, _ = os.UserHomeDir()
		}
		confDir = path.Join(userConfDir, DefaultConfigFolder)
	}

	err := createDirIfNotExists(confDir)
	if err != nil {
		log.Fatalf("[ERROR] config folder %s does not exists and cannot be created: %v", confDir, err)
	}

	var proxies []*Proxy
	bctx, backgroundCancel := context.WithCancel(context.Background())
	watchConfigFolder(bctx, confDir, &proxies)

	srv := &http.Server{
		Addr:    listen,
		Handler: http.HandlerFunc(requestHandler(&proxies)),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("[PANIC] while listening to %s: %v", listen, err)
		}
	}()

	<-sig
	backgroundCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("[ERROR] during shutdown: %v", err)
	}

}
