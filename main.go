package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	goclient "github.com/me2seeks/gouge/internal/client"
	goserver "github.com/me2seeks/gouge/internal/server"
	"github.com/me2seeks/gouge/share/cos"
)

func main() {
	flag.Parse()
	args := flag.Args()
	subcmd := "server"
	if len(args) > 0 {
		subcmd = args[0]
		args = args[1:]
	}
	switch subcmd {
	case "server":
		server(args)
	case "client":
		client(args)
	default:
		os.Exit(0)
	}
}

func server(args []string) {
	flags := flag.NewFlagSet("server", flag.ContinueOnError)
	config := goserver.Config{}

	flags.StringVar(&config.KeyFile, "keyfile", "", "")
	flags.StringVar(&config.Proxy, "proxy", "", "")
	flags.StringVar(&config.Proxy, "backend", "", "")

	host := flags.String("h", "", "")
	port := flags.String("p", "", "")
	verbose := flags.Bool("v", false, "")

	flags.Parse(args)
	s, err := goserver.NewServer(&config)
	if err != nil {
		log.Fatal(err)
	}
	s.Debug = *verbose
	generatePidFile()
	ctx := cos.InterruptContext()
	if err := s.StartContext(ctx, *host, *port); err != nil {
		log.Fatal(err)
	}
	if err := s.Wait(); err != nil {
		log.Fatal(err)
	}
}

func client(args []string) {
	flags := flag.NewFlagSet("client", flag.ContinueOnError)
	config := goclient.Config{Headers: http.Header{}}
	flags.StringVar(&config.Fingerprint, "fingerprint", "", "")
	flags.DurationVar(&config.KeepAlive, "keepalive", 25*time.Second, "")
	flags.IntVar(&config.MaxRetryCount, "max-retry-count", -1, "")
	flags.DurationVar(&config.MaxRetryInterval, "max-retry-interval", 0, "")
	flags.StringVar(&config.Proxy, "proxy", "", "")
	flags.Var(&headerFlags{config.Headers}, "header", "")
	hostname := flags.String("hostname", "", "")
	verbose := flags.Bool("v", false, "")
	flags.Parse(args)
	if *hostname != "" {
		config.Headers.Set("Host", *hostname)
	}
	c, err := goclient.NewClient(&config)
	if err != nil {
		log.Fatal(err)
	}
	c.Debug = *verbose

	config.Server = args[0]
	config.Remotes = args[1:]
	generatePidFile()
	ctx := cos.InterruptContext()
	if err := c.StartContext(ctx); err != nil {
		log.Fatal(err)
	}
	if err := c.Wait(); err != nil {
		log.Fatal(err)
	}

}

type headerFlags struct {
	http.Header
}

func (f *headerFlags) String() string {
	out := ""
	for k, v := range f.Header {
		out += fmt.Sprintf("%s: %s\n", k, strings.Join(v, ","))
	}
	return out
}
func (f *headerFlags) Set(arg string) error {
	index := strings.Index(arg, ":")
	if index < 0 {
		return fmt.Errorf("invalid header %s", arg)
	}
	if f.Header == nil {
		f.Header = http.Header{}
	}
	key := arg[:index]
	value := arg[index+1:]
	f.Header.Set(key, strings.TrimSpace(value))
	return nil
}

func generatePidFile() {
	pid := []byte(strconv.Itoa(os.Getpid()))
	os.WriteFile("gouge.pid", pid, 0644)
}
