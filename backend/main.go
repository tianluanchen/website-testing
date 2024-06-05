package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"website-testing/config"
	"website-testing/internal/server"
	"website-testing/internal/server/services"
	"website-testing/internal/tc"
	"website-testing/pkg"
	"website-testing/wt"
)

var logger = pkg.NewLogger()

func main() {
	var (
		addr      string
		open      bool
		immediate bool
		v         bool
	)
	flag.StringVar(&addr, "addr", "127.0.0.1:8080", "Listen address")
	flag.BoolVar(&open, "open", true, "Open browser")
	flag.BoolVar(&immediate, "immediate", true, "Immediately start testing")
	flag.BoolVar(&v, "v", false, "Print version information")
	flag.Parse()
	if v {
		fmt.Println("website-testing:", config.Version)
		return
	}
	srv := http.Server{
		Addr:    addr,
		Handler: server.New(),
	}
	logger.Infoln("Listening and serving HTTP on", addr)
	for _, v := range config.Origins {
		logger.Infoln("Ayouth的个人网站：", (*url.URL)(v).String())
	}
	go func() {
		if immediate {
			tc.Test(&config.Conf{
				TimeoutSeconds: 10,
				UserAgent:      wt.UAFirefoxWin64,
			}, services.CallbackOption)
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()
	if open {
		go openBrowser(addr)
	}
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Warnln("Shutdown Server ...")
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Fatal(err)
	}
	logger.Infoln("Server exiting")
}

func openBrowser(addr string) {
	s := strings.Split(addr, ":")
	if len(s) >= 2 {
		port := s[len(s)-1]
		u := fmt.Sprintf("http://%s:%s", strings.Join(s[0:len(s)-1], ":"), port)
		logger.Debugln("Opening browser:", u)
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "linux":
			cmd = exec.Command("xdg-open", u)
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", u)
		case "darwin":
			cmd = exec.Command("open", u)
		}
		if cmd == nil {
			logger.Warnln("Can't auto open browser:", u)
		} else {
			if err := cmd.Run(); err == nil {
				logger.Debugln("Open browser:", u, "success")
			} else {
				logger.Warnln("Open browser:", u, "failed")
			}
		}
	} else {
		logger.Warnln("Can't parse address:", addr)
	}

}
