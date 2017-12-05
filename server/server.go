package server

import (
	"github.com/go-openapi/spec"
	"github.com/labstack/echo"
)

type PciServer struct {
	*echo.Echo // web service
	docs *spec.Swagger
}

func NewPciServer() (s *PciServer) {
	s = &PciServer{
		echo.New(),
		&spec.Swagger{},
	}
	// hide echo banner
	s.Echo.HideBanner = true
	return s;
}

var mrModelTag = "mrModel"
var pciPlanTag = "pciPlan"
var pciEvaluate = "pciEvaluate"
