package routing

import (
	"api/internal/route"
	"api/pkg/config"
	"api/pkg/postgres"
	"api/pkg/rabbitmq"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	LabstackLog "github.com/labstack/gommon/log"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
)

var engine *echo.Echo

func Serve() {
	go serve()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := engine.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
	defer postgres.Close()
	defer rabbitmq.Close()
}

func serve() {
	appConf := config.GetApp()
	engine = echo.New()
	engine.Pre(middleware.AddTrailingSlash())
	engine.Use(middleware.Logger())
	engine.Use(middleware.Recover())
	engine.Use(TimeoutMiddleware(10 * time.Second))

	engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: appConf.AllowOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodHead, http.MethodOptions},
	}))

	route.Register(engine)

	if appConf.Debug {
		engine.Debug = true
		engine.Logger.SetLevel(LabstackLog.DEBUG)
		verboseLog(engine)
	} else {
		engine.Logger.SetLevel(LabstackLog.WARN)
		engine.HideBanner = true
	}
	server := fmt.Sprintf("%s:%s", appConf.Host, appConf.Port)
	if err := engine.Start(server); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func verboseLog(e *echo.Echo) {
	paths := e.Routes()
	sort.SliceStable(paths, func(i, j int) bool {
		return paths[i].Path < paths[j].Path
	})
	maxNameLen := len("ROUTE NAME")
	maxPathLen := len("PATH")
	maxMethodLen := len("METHOD")
	for _, route := range paths {
		if len(route.Name) > maxNameLen {
			maxNameLen = len(route.Name)
		}
		if len(route.Path) > maxPathLen {
			maxPathLen = len(route.Path)
		}
		if len(route.Method) > maxMethodLen {
			maxMethodLen = len(route.Method)
		}
	}

	headerFormat := fmt.Sprintf("\n%%-%ds %%-%ds %%-%ds\n", maxNameLen+5, maxMethodLen, maxPathLen)
	log.Printf(headerFormat, "ROUTE NAME", "METHOD", "PATH")
	log.Println(strings.Repeat("-", maxNameLen+maxPathLen+maxMethodLen+3))

	rowFormat := fmt.Sprintf("%%-%ds %%-%ds %%-%ds\n", maxNameLen+5, maxMethodLen, maxPathLen)
	for _, route := range paths {
		if !strings.HasSuffix(route.Name, ".init.func1") {
			log.Printf(rowFormat, strings.TrimSuffix(route.Name, "-fm"), route.Method, route.Path)
		}
	}
}
