package server

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/spec"
	"github.com/labstack/echo"
)

func (s *PciServer) InitDocs() *PciServer {
	s.docs.SwaggerProps = spec.SwaggerProps{}
	s.docs.Swagger = "2.0"
	//s.docs.Schemes = []string{"http"}
	s.docs.Tags = []spec.Tag{
		{TagProps: spec.TagProps{Name: mrModelTag, Description: "MR模型"}},
		{TagProps: spec.TagProps{Name: pciPlanTag, Description: "PCI规划"}},
		{TagProps: spec.TagProps{Name: pciEvaluate, Description: "PCI评估"}},
	}
	s.docs.Info = &spec.Info{spec.VendorExtensible{}, spec.InfoProps{
		Title:       fmt.Sprintf("PCI规划"),
		Version:     "1.0.0",
		Description: "依据MR话务,信噪比,工作日和非工作日等模型,采用遗传算法迭代规划,规划全网或局部PCI.",
	}}
	s.docs.Definitions = spec.Definitions{"ErrMsg": errorMessageDefinition(),}
	s.GET("/docs/", func(c echo.Context) error{
		return c.HTML(http.StatusOK, DOCS_HTML)
	}).Name = "Docs UI"
	s.GET("/docs.json",  func(c echo.Context) error{
		s.docs.Schemes = []string{c.Scheme()}
		return c.JSON(http.StatusOK, s.docs)
	}).Name = "Docs Infomation"
	return s;
}

func errorMessageDefinition() (schema spec.Schema) {
	//ErrorMessage
	schema.Type = spec.StringOrArray{"object"}
	schema.Title = "错误消息"
	schema.Description = "意外的错误时的消息"
	schema.SchemaProps = spec.SchemaProps{
		Required: []string{"error"},
		Properties: map[string]spec.Schema{
			"errType": spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:        spec.StringOrArray{"string"},
					Description: "消息标识",
				},
			},
			"errDescr": spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:        spec.StringOrArray{"string"},
					Description: "消息描述",
				},
			},
		},
	}
	return
}
