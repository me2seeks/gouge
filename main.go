package main

import (
	"flag"
	"log"
	"os"
	"strconv"

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
	flags.Parse(args)
	s, err := goserver.NewServer(&config)
	if err != nil {
		log.Fatal(err)
	}
	generatePidFile()
	ctx := cos.InterruptContext()
	if err := s.StartContext(ctx, *host, *port); err != nil {
		log.Fatal(err)
	}
	if err := s.Wait(); err != nil {
		log.Fatal(err)
	}
}

func client([]string) {
	//TODO
}

func generatePidFile() {
	pid := []byte(strconv.Itoa(os.Getpid()))
	os.WriteFile("gouge.pid", pid, 0644)
}
