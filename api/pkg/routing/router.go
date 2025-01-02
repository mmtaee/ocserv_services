package routing

import (
	"api/internal/routes"
	"api/pkg/config"
	"api/pkg/database"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	LabstackLog "github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	defer database.Close()
}

func serve() {
	appConf := config.GetApp()
	server := fmt.Sprintf("%s:%s", appConf.Host, appConf.Port)

	engine = echo.New()

	engine.Pre(middleware.RemoveTrailingSlash())
	engine.Use(middleware.Logger())
	engine.Use(middleware.Recover())
	engine.Use(TimeoutMiddleware(10 * time.Second))
	engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: appConf.AllowOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodHead,
			http.MethodOptions,
		},
	}))

	routes.Register(engine)

	if appConf.Debug {
		engine.Debug = true
		engine.Logger.SetLevel(LabstackLog.DEBUG)
		verboseLog(server)
		engine.GET("/swagger/*", echoSwagger.WrapHandler)
	} else {
		engine.Logger.SetLevel(LabstackLog.WARN)
		engine.HideBanner = true
	}

	//engine.Use(middleware.GzipWithConfig(middleware.GzipConfig{
	//	Skipper: func(c echo.Context) bool {
	//		if strings.Contains(c.Request().URL.Path, "swagger") {
	//			return true
	//		}
	//		return false
	//	},
	//}))
	if err := engine.Start(server); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func verboseLog(service string) {
	paths := engine.Routes()
	sort.SliceStable(paths, func(i, j int) bool {
		return paths[i].Path < paths[j].Path
	})
	maxNameLen := len("ROUTE NAME")
	maxPathLen := len("PATH")
	maxMethodLen := len("METHOD")
	for _, path := range paths {
		if len(path.Name) > maxNameLen {
			maxNameLen = len(path.Name)
		}
		if len(path.Path) > maxPathLen {
			maxPathLen = len(path.Path)
		}
		if len(path.Method) > maxMethodLen {
			maxMethodLen = len(path.Method)
		}
	}

	headerFormat := fmt.Sprintf("\n%%-%ds %%-%ds %%-%ds\n", maxNameLen+5, maxMethodLen, maxPathLen)
	log.Printf(headerFormat, "ROUTE NAME", "METHOD", "PATH")
	log.Println(strings.Repeat("-", maxNameLen+maxPathLen+maxMethodLen+3))

	rowFormat := fmt.Sprintf("%%-%ds %%-%ds %%-%ds\n", maxNameLen+5, maxMethodLen, maxPathLen)
	for _, path := range paths {
		if !strings.HasSuffix(path.Name, ".init.func1") {
			log.Printf(
				rowFormat,
				strings.TrimSuffix(path.Name, "-fm"),
				path.Method,
				fmt.Sprintf("http://%s%s/", service, path.Path),
			)
		}
	}
}
