package server

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/gog"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/srv"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pangox-xdemo/app"
	"github.com/askasoft/pangox-xdemo/app/server/gormdb"
	"github.com/askasoft/pangox-xdemo/data"
	"github.com/askasoft/pangox-xdemo/tpls"
	"github.com/askasoft/pangox-xdemo/txts"
	"github.com/askasoft/pangox-xdemo/web"
	"github.com/askasoft/pangox/xwa/xcpts"
	"github.com/askasoft/pangox/xwa/xschs"
)

// -----------------------------------
// srv.Cmd implement

// Flag process optional command flag
func (s *service) Flag() {
	flag.StringVar(&s.outdir, "out", "", "specify the output directory.")
	flag.BoolVar(&s.debug, "debug", false, "print debug log.")
}

// Usage print command line usage
func (s *service) Usage() {
	fmt.Println("Usage: " + s.Name() + " <command> [options]")
	fmt.Println("  <command>:")
	srv.PrintDefaultCommand(os.Stdout)
	fmt.Println("    migrate <target> [schema]...")
	fmt.Println("      target=schema      migrate database schema.")
	fmt.Println("      target=settings    migrate database settings.")
	fmt.Println("      target=super       migrate database super users.")
	fmt.Println("        [schema]...      specify schemas to migrate.")
	fmt.Println("    schema <action> [schema]...")
	fmt.Println("      action=init        initialize the schema.")
	fmt.Println("      action=check       check schema tables.")
	fmt.Println("      action=update      apply schema update scripts.")
	fmt.Println("      action=vacuum      vacuum schema tables (postgresql only).")
	fmt.Println("        [schema]...      specify schemas to execute.")
	fmt.Println("      action=genddl      generate schema DDL script.")
	fmt.Println("    execsql <file> [schema]...")
	fmt.Println("      <file>             execute sql file.")
	fmt.Println("        [schema]...      specify schemas to execute sql.")
	fmt.Println("    exectask <task>      execute task [ " + str.Join(xschs.Schedules.Keys(), ", ") + " ].")
	fmt.Println("    export <target> ...")
	fmt.Println("      target=settings    export settings from database.")
	fmt.Println("        [schema]...      specify schemas to export.")
	fmt.Println("      target=assets      export embedded assets.")
	fmt.Println("        <password>       specify the develop password.")
	fmt.Println("    import <target> <source>")
	fmt.Println("      target=settings    import settings to database.")
	fmt.Println("        <source>         specify setting files ({schema}.csv) directory to import.")
	fmt.Println("    encrypt [key] <str>  encrypt string.")
	fmt.Println("    decrypt [key] <str>  decrypt string.")
	fmt.Println("  [options]:")
	srv.PrintDefaultOptions(os.Stdout)
	fmt.Println("    -out                 specify the output directory.")
	fmt.Println("    -debug               print the debug log.")
}

// Exec execute optional command except the internal command
// Basic: 'help' 'usage' 'version'
// Windows only: 'service' (install | remove | start | stop | debug)
func (s *service) Exec(cmd string) {
	cw := log.NewConsoleWriter(true)
	cw.SetFormat("%t{15:04:05} [%p] - %m%n%T")

	log.SetWriter(cw)
	log.SetLevel(gog.If(s.debug, log.LevelDebug, log.LevelInfo))

	switch cmd {
	case "migrate":
		s.doMigrate()
	case "schema":
		s.doSchemas()
	case "execsql":
		s.doExecSQL()
	case "exectask":
		s.doExecTask()
	case "encrypt":
		s.doEncrypt()
	case "decrypt":
		s.doDecrypt()
	case "export":
		s.doExport()
	case "import":
		s.doImport()
	default:
		flag.CommandLine.SetOutput(os.Stdout)
		fmt.Fprintf(os.Stderr, "Invalid command %q\n\n", cmd)
		s.Usage()
		os.Exit(app.ExitErrARG)
	}
}

func (s *service) doMigrate() {
	sub := flag.Arg(1)
	if sub == "" {
		fmt.Fprintln(os.Stderr, "Missing migrate <target>.")
		os.Exit(app.ExitErrARG)
	}

	args := flag.Args()[2:]

	switch sub {
	case "schema":
		initConfigs()
		if err := gormdb.MigrateSchemas(args...); err != nil {
			log.Fatal(app.ExitErrDB, err)
		}
	case "settings":
		initConfigs()
		initDatabase()
		if err := dbMigrateSettings(args...); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(app.ExitErrDB)
		}
	case "super":
		initConfigs()
		initDatabase()
		if err := dbMigrateSupers(args...); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(app.ExitErrDB)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid migrate <target>: %q", sub)
		os.Exit(app.ExitErrARG)
	}

	log.Info("DONE.")
}

func (s *service) doSchemas() {
	sub := flag.Arg(1)
	if sub == "" {
		fmt.Fprintln(os.Stderr, "Missing schema <command>.")
		os.Exit(app.ExitErrARG)
	}

	args := flag.Args()[2:]

	switch sub {
	case "init":
		initConfigs()
		initDatabase()
		if err := dbSchemaInit(args...); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(app.ExitErrDB)
		}
	case "check":
		initConfigs()
		initDatabase()
		if err := dbSchemaCheck(args...); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(app.ExitErrDB)
		}
	case "update":
		initConfigs()
		initDatabase()
		if err := dbSchemaUpdate(args...); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(app.ExitErrDB)
		}
	case "vacuum":
		initConfigs()
		initDatabase()
		if app.DBType() != "postgres" {
			fmt.Printf("Database %s does not support vacuum", app.DBType())
			os.Exit(app.ExitErrDB)
		}
		if err := dbSchemaVacuum(args...); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(app.ExitErrDB)
		}
	case "genddl":
		initConfigs()

		if err := gormdb.GenerateDDL(s.outdir); err != nil {
			log.Fatal(app.ExitErrDB, err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid schema <command>: %q", sub)
		os.Exit(app.ExitErrARG)
	}

	log.Info("DONE.")
}

func (s *service) doExecSQL() {
	file := flag.Arg(1)
	if file == "" {
		fmt.Fprintln(os.Stderr, "Missing SQL file.")
		os.Exit(app.ExitErrARG)
	}

	args := flag.Args()[2:]

	initConfigs()
	initDatabase()

	if err := dbExecSQL(file, args...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(app.ExitErrDB)
	}

	log.Info("DONE.")
}

func (s *service) doExecTask() {
	tn := flag.Arg(1)
	tf, ok := xschs.Schedules.Get(tn)
	if !ok {
		fmt.Fprintf(os.Stderr, "Invalid Task %q\n", tn)
		os.Exit(app.ExitErrARG)
	}

	initConfigs()
	initCaches()
	initDatabase()

	tf()

	log.Info("DONE.")
}

func (s *service) doExport() {
	sub := flag.Arg(1)
	if sub == "" {
		fmt.Fprintln(os.Stderr, "Missing export <target>.")
		os.Exit(app.ExitErrARG)
	}

	args := flag.Args()[2:]

	switch sub {
	case "settings":
		initConfigs()
		initMessages()
		initDatabase()
		if err := dbExportSettings(s.outdir, args...); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(app.ExitErrDB)
		}
	case "assets":
		pwd := asg.First(args)
		if pwd != "askadevelop" {
			fmt.Fprintln(os.Stderr, "invalid develop password!")
			os.Exit(app.ExitErrARG)
		}
		if err := exportAssets(s.outdir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid export <target>: %q", sub)
		os.Exit(app.ExitErrARG)
	}

	log.Info("DONE.")
}

func exportAssets(dir string) error {
	if err := os.CopyFS(filepath.Join(dir, "data"), data.FS); err != nil {
		return err
	}
	if err := os.CopyFS(filepath.Join(dir, "txts"), txts.FS); err != nil {
		return err
	}
	if err := os.CopyFS(filepath.Join(dir, "tpls"), tpls.FS); err != nil {
		return err
	}

	if err := os.CopyFS(filepath.Join(dir, "web"), web.FS); err != nil {
		return err
	}
	for key, fs := range web.Statics {
		if err := os.CopyFS(filepath.Join(dir, "web", "static", key), fs); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) doImport() {
	sub := flag.Arg(1)
	if sub == "" {
		fmt.Fprintln(os.Stderr, "Missing import <source>.")
		os.Exit(app.ExitErrARG)
	}

	switch sub {
	case "settings":
		initConfigs()
		initMessages()
		initDatabase()
		if err := dbImportSettings(flag.Arg(2)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(app.ExitErrDB)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid import <target>: %q", sub)
		os.Exit(app.ExitErrARG)
	}

	log.Info("DONE.")
}

func (s *service) doEncrypt() {
	k, v := cryptFlags()
	if es, err := xcpts.Encrypt(k, v); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Println(es)
	}
}

func (s *service) doDecrypt() {
	k, v := cryptFlags()
	if ds, err := xcpts.Decrypt(k, v); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Println(ds)
	}
}

func cryptFlags() (k, v string) {
	k, v = flag.Arg(1), flag.Arg(2)
	if v == "" {
		initConfigs()
		k, v = app.Secret(), k
	}
	return
}
