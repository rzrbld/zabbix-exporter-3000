package handlers

import (
	iris "github.com/kataras/iris/v12"
)

var Probes = func(ctx iris.Context) {
	var res = iris.Map{"Success": "OK"}
	ctx.JSON(res)
}
