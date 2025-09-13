//go:generate goversioninfo
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/askasoft/pango/gog"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/cmd/schema"
	"github.com/askasoft/pangox/xwa"
)

func usage() {
	help := `
Usage: %s <command> [options]
  <command>:
    version             print the version information.
    help | usage        print the usage information.
    generate [dbtype] [output] generate database schema DDL.
      [dbtype]          specify the database type.
      [output]          specify the output DDL file.
    migrate [schema]... migrate database schemas.
      [schema]...       specify schemas to migrate.
  <options>:
    -h | -help          print the help message.
    -v | -version       print the version message.
    -dir                set the working directory.
    -debug              print the debug log.
`
	fmt.Printf(help, filepath.Base(os.Args[0]))
}

func chdir(workdir string) {
	if workdir != "" {
		if err := os.Chdir(workdir); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to change directory: %v\n", err)
			os.Exit(app.ExitErrCMD)
		}
	}
}

func initConfigs() {
	if err := xwa.InitConfigs(); err != nil {
		log.Fatal(app.ExitErrCFG, err)
	}
}

func main() {
	var (
		debug   bool
		version bool
		workdir string
	)

	flag.BoolVar(&version, "v", false, "print version message.")
	flag.BoolVar(&version, "version", false, "print version message.")
	flag.BoolVar(&debug, "debug", false, "print debug log.")
	flag.StringVar(&workdir, "dir", "", "set the working directory.")

	flag.CommandLine.Usage = usage
	flag.Parse()

	chdir(workdir)

	if version {
		fmt.Println(app.Versions())
		os.Exit(0)
	}

	cw := log.NewConsoleWriter()
	cw.SetFormat("%t [%p] - %m%n%T")

	log.SetWriter(cw)
	log.SetLevel(gog.If(debug, log.LevelDebug, log.LevelInfo))

	arg := flag.Arg(0)
	switch arg {
	case "generate":
		dbtype, outfile := "", ""
		for _, a := range flag.Args()[1:] {
			if str.EndsWith(a, ".sql") {
				outfile = a
			} else {
				dbtype = a
			}
		}

		initConfigs()
		if err := schema.GenerateDDL(dbtype, outfile); err != nil {
			log.Fatal(app.ExitErrCMD, err)
		}
	case "migrate":
		initConfigs()
		if err := schema.MigrateSchemas(flag.Args()[1:]...); err != nil {
			log.Fatal(app.ExitErrCMD, err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid command %q\n\n", arg)
		usage()
		os.Exit(app.ExitErrCMD)
	}
}
