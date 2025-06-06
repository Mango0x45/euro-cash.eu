package main

/* TODO: Customize logger format when running in a debug state */
/* TODO: Set production logger to the syslog */

import (
	"flag"
	"log"
	"os"
	"syscall"
	"time"

	"git.thomasvoss.com/euro-cash.eu/src"
	"git.thomasvoss.com/euro-cash.eu/src/email"
)

func main() {
	port := flag.Int("port", 8080, "port number")
	debugp := flag.Bool("debug", false, "run in debug mode")
	flag.BoolVar(&email.Config.Disabled, "no-email", false,
		"disables email support")
	flag.StringVar(&email.Config.Host, "smtp-host", "smtp.migadu.com",
		"SMTP server hostname")
	flag.IntVar(&email.Config.Port, "smtp-port", 465,
		"SMTP server port number")
	flag.StringVar(&email.Config.ToAddr, "email-to", "bugs@euro-cash.eu",
		"address to send error messages to")
	flag.StringVar(&email.Config.FromAddr, "email-from", "noreply@euro-cash.eu",
		"address to send error messages from")
	flag.StringVar(&email.Config.Password, "email-password", "",
		"password to authenticate the email client")
	flag.Parse()

	if *debugp {
		go watch()
	}
	src.Run(*port)
}

func watch() {
	path, err := os.Executable()
	if err != nil {
		die(err)
	}

	ostat, err := os.Stat(path)
	if err != nil {
		die(err)
	}

	for {
		nstat, err := os.Stat(path)
		if err != nil {
			die(err)
		}

		if nstat.ModTime() != ostat.ModTime() {
			if err := syscall.Exec(path, os.Args, os.Environ()); err != nil {
				die(err)
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func die(err error) {
	log.Fatal(err)
	os.Exit(1)
}
