package main

/* TODO: Customize logger format when running in a debug state */
/* TODO: Set production logger to the syslog */

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"git.thomasvoss.com/euro-cash.eu/pkg/atexit"
	. "git.thomasvoss.com/euro-cash.eu/pkg/try"
	"git.thomasvoss.com/euro-cash.eu/pkg/watch"

	"git.thomasvoss.com/euro-cash.eu/src"
	"git.thomasvoss.com/euro-cash.eu/src/dbx"
	"git.thomasvoss.com/euro-cash.eu/src/email"
	"git.thomasvoss.com/euro-cash.eu/src/i18n"
)

func main() {
	Try(os.Chdir(filepath.Dir(os.Args[0])))

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
	flag.StringVar(&dbx.DBName, "db-name", "eurocash.db",
		"database name or ‘:memory:’ for an in-memory database")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		atexit.Exec()
		os.Exit(0)
	}()

	if *debugp {
		path := Try2(os.Executable())
		go watch.File(path, func() {
			atexit.Exec()
			Try(syscall.Exec(path, os.Args, os.Environ()))
		})
	}

	i18n.Init()
	dbx.Init(Try2(os.OpenRoot("src/dbx/sql")).FS())
	app.BuildTemplates(Try2(os.OpenRoot("src/templates")).FS(), *debugp)
	app.Run(*port)
}
