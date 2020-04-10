package main

import (
  "fmt"
	"github.com/kataras/iris/v12"

	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"

	cnf "github.com/rzrbld/zabbix-exporter-3000/config"
	hdl "github.com/rzrbld/zabbix-exporter-3000/handlers"
)

func main() {

	fmt.Println("\033[31m\r\n\r\n\r\n███████╗███████╗██████╗  ██████╗  ██████╗  ██████╗ \r\n╚══███╔╝██╔════╝╚════██╗██╔═████╗██╔═████╗██╔═████╗ \r\n  ███╔╝ █████╗   █████╔╝██║██╔██║██║██╔██║██║██╔██║ \r\n ███╔╝  ██╔══╝   ╚═══██╗████╔╝██║████╔╝██║████╔╝██║ \r\n███████╗███████╗██████╔╝╚██████╔╝╚██████╔╝╚██████╔╝ \r\n╚══════╝╚══════╝╚═════╝  ╚═════╝  ╚═════╝  ╚═════╝  \r\n\033[m")
	fmt.Println("\033[33mZabbix Exporter for Prometheus")
	fmt.Println("version  : 0.1 ")
	fmt.Println("Author   : rzrbld")
	fmt.Println("License  : MIT")
	fmt.Println("Git-repo : https://github.com/rzrbld/zabbix-exporter-3000 \033[m \r\n")

	app := iris.New()

	app.Logger().SetLevel("debug")

	app.Use(recover.New())
	app.Use(logger.New())

  // prometheus metrics

	m := prometheusMiddleware.New("ze3000", 0.3, 1.2, 5.0)
	hdl.RecordMetrics()
	app.Use(m.ServeHTTP)
	app.Get("/metrics", iris.FromStd(promhttp.Handler()))


	app.Get("/liveness", func(ctx iris.Context) {
		ctx.WriteString("ok")
	})

	app.Get("/readiness", func(ctx iris.Context) {
		ctx.WriteString("ok")
	})

	app.Run(iris.Addr(cnf.MainHostPort), iris.WithoutServerError(iris.ErrServerClosed))
}
