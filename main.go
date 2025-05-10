package main

import (
	"github.com/00unnmd/pills_parser/handlers"
	"github.com/00unnmd/pills_parser/internals/utils"
)

func parseNow() {
	data := handlers.ReqAllPills()
	utils.GenerateXLSX(data)
}

func main() {
	parseNow()
}
