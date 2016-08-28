package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"strings"

	"errors"

	"github.com/codegangsta/cli"
	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/setting"
	"github.com/containerops/vessel/web"
	"gopkg.in/macaron.v1"
)

// GetCmdWeb get a client command
func GetCmdWeb() cli.Command {
	return cli.Command{
		Name:        "web",
		Usage:       "start vessel web service",
		Description: "vessel is a CI module.",
		Action:      runWeb,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "conf",
				Value: "./conf/global.yaml",
				Usage: "configuraion file path; defaults to './conf/global.yaml'",
			},
		},
	}
}

func runWeb(c *cli.Context) {
	if err := parseConf(c); err != nil {
		log.Println(err)
		return
	}
	if err := models.InitEtcd(); err != nil {
		log.Println(err)
		return
	}
	if err := models.InitK8S(); err != nil {
		log.Println(err)
		return
	}
	if err := models.InitDatabase(); err != nil {
		log.Println(err)
		return
	}

	m := macaron.New()

	//Set Macaron Web Middleware And Routers
	web.SetVesselMacaron(m)
	switch setting.RunTime.HTTP.ListenMode {
	case "http":
		listenaddr := fmt.Sprintf("%s:%s", setting.RunTime.HTTP.Host, setting.RunTime.HTTP.Port)
		if err := http.ListenAndServe(listenaddr, m); err != nil {
			log.Printf("Start http service error: %v", err.Error())
		}
		break
	case "https":
		listenaddr := fmt.Sprintf("%s:%s", setting.RunTime.HTTP.Host, setting.RunTime.HTTP.Port)
		server := &http.Server{Addr: listenaddr, TLSConfig: &tls.Config{MinVersion: tls.VersionTLS10}, Handler: m}
		if err := server.ListenAndServeTLS(setting.RunTime.HTTP.HTTPSCertFile, setting.RunTime.HTTP.HTTPSKeyFile); err != nil {
			log.Printf("Start https service error: %v", err.Error())
		}
		break
	default:
		break
	}
}

func parseConf(c *cli.Context) (err error) {
	globalConf := c.String("conf")
	if !strings.HasSuffix(globalConf, "global.yaml") {
		err = errors.New("Conf file must be named 'global.yaml'")
	} else {
		err = setting.InitGlobalConf(globalConf)
	}
	return err
}
