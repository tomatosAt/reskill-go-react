package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin/v2"
	"github.com/tomatosAt/reskill-go-react/app"
	"github.com/tomatosAt/reskill-go-react/config"
	"github.com/tomatosAt/reskill-go-react/module"
)

var Version string

func main() {
	kp := kingpin.New(filepath.Base(os.Args[0]), fmt.Sprintf("%s", Version))
	kp.Version(Version)
	versionCmd := kp.Command("version", "Show application version.")
	startCmd := kp.Command("start", "Start application.")
	cfgFile := kp.Flag("config-file", "Set load config file (default: config.yml)").Default("config.yml").String()
	switch kingpin.MustParse(kp.Parse(os.Args[1:])) {
	case versionCmd.FullCommand():
		fmt.Println(Version)
	case startCmd.FullCommand():
		// Load configuration
		cfg := config.LoadConfig(*cfgFile, Version)
		a := app.New(cfg)
		a.InitFiberServer()
		l := a.NewLogger().WithField("package", "main")
		if err := module.Create(a.Context); err != nil {
			l.Errorln("[x] Start create module failed -:", err)
			os.Exit(1)
		}
		if err := a.StartHTTP(); err != nil {
			l.Errorln("[x] Start http server error -:", err)
			os.Exit(2)
		}
	}
}
