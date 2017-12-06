package server

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"
	"github.com/labstack/echo"
)
const mrModelTag = "mrModel"
const pciPlanTag = "pciPlan"
const pciEvaluate = "pciEvaluate"

const mrTaskStatus = "TaskStatus"

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
	s.docs.Definitions = spec.Definitions{"ErrMsg": errMsgDefinition(),mrTaskStatus:mrTaskStatusDefinition()}
	//  /mr/decode  post
	//  /mr/matrix/   post+get -->enodeb_id,ci,neib_enodeb,neib_ci,mr_total,noise_ratio,begin_datetime,end_datetime
	// /mr/matrix/merge post 多个mr_matrix_范围_开始时间_截止时间.xlsx合并

	s.docs.Paths =&spec.Paths{Paths:map[string]spec.PathItem{
		"/mr/decode": {PathItemProps: spec.PathItemProps{
			Post:newOperationFull(
				mrModelTag,
				fmt.Sprintf("解码MR文件"),
				fmt.Sprintf("一个或多个MR文件解码后输出"),
				[]spec.Parameter{
					{
						SimpleSchema: spec.SimpleSchema{
							Type: "file",
						},
						ParamProps: spec.ParamProps{
							In:          "formData",
							Name:        "mr",
							Required:    false,
							Description: "MR文件",
						}},
				},
				fmt.Sprintf("返回解码后的文件"),
				&spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"string"}},
					SwaggerSchemaProps: spec.SwaggerSchemaProps{Example: ""},
				},[]string{"multipart/form-data"},[]string{"application/json","application/octet-stream"}),
		}},
		"/mr/matrix/": {PathItemProps: spec.PathItemProps{Post: newOperationFull(
			mrModelTag,
			fmt.Sprintf("添加统计MR模型任务"),
			fmt.Sprintf("传入要统计的MR文件,时间维度(60分,24小时,是否夸天合并)"),
			[]spec.Parameter{
				{
					SimpleSchema: spec.SimpleSchema{
						Type: "file",
					},
					ParamProps: spec.ParamProps{
						In:          "formData",
						Name:        "mr",
						Required:    false,
						Description: "MR文件",
						}},
				{
					SimpleSchema: spec.SimpleSchema{
						Type: "file",
					},
					ParamProps: spec.ParamProps{
						In:          "formData",
						Name:        "cell",
						Required:    false,
						Description: "cell工参文件",
					}},
				{
					SimpleSchema: spec.SimpleSchema{
						Type: "file",
					},
					ParamProps: spec.ParamProps{
						In:          "formData",
						Name:        "neib_cell",
						Required:    false,
						Description: "neib_cell工参文件",
					}},
			},
			fmt.Sprintf("返回创建的任务id"),
			&spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"}},
				SwaggerSchemaProps: spec.SwaggerSchemaProps{Example: ""},
			},[]string{"multipart/form-data"},[]string{"application/json","application/octet-stream"}),
			Get:newOperation(mrModelTag,
				fmt.Sprintf("查询统计MR模型任务"),
				fmt.Sprintf("根据任务id查询任务进度)"),
				[]spec.Parameter{
					{
						SimpleSchema: spec.SimpleSchema{
							Type: "string",
						},
						ParamProps: spec.ParamProps{
							In:          "query",
							Name:        "tid",
							Required:    true,
							Description: "MR任务标识",
						}},
				},
				fmt.Sprintf("返回任务进度"),
				&spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"object"},
						Ref: getModelSwaggerRef(mrTaskStatus),
						},
				}),
		}},
		"/mr/matrix/merge":{PathItemProps: spec.PathItemProps{Post: newOperationFull(
			mrModelTag,
			fmt.Sprintf("合并MR模型"),
			fmt.Sprintf("按天合并,按小时合并等"),
			[]spec.Parameter{
				{
					SimpleSchema: spec.SimpleSchema{
						Type: "file",
					},
					ParamProps: spec.ParamProps{
						In:          "formData",
						Name:        "mfiles",
						Required:    false,
						Description: "MR Matrix文件压缩包",
					}},
			},
			fmt.Sprintf("返回合并的结果"),
			&spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"}},
				SwaggerSchemaProps: spec.SwaggerSchemaProps{Example: ""},
			},[]string{"multipart/form-data"},[]string{"application/json","application/octet-stream"})}},
		}}

	s.GET("/docs/", func(c echo.Context) error{
		return c.HTML(http.StatusOK, DOCS_HTML)
	}).Name = "Docs UI"
	s.GET("/docs.json",  func(c echo.Context) error{
		s.docs.Schemes = []string{c.Scheme()}
		return c.JSON(http.StatusOK, s.docs)
	}).Name = "Docs Infomation"
	return s;
}

func errMsgDefinition() (schema spec.Schema) {
	//ErrorMessage
	schema.Type = spec.StringOrArray{"object"}
	schema.Title = "错误消息"
	schema.Description = "意外的错误时的消息"
	schema.SchemaProps = spec.SchemaProps{
		Required: []string{"errType"},
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
func mrTaskStatusDefinition() (schema spec.Schema) {
	//mrTaskStatus
	schema.Type = spec.StringOrArray{"object"}
	schema.Title = "错误消息"
	schema.Description = "意外的错误时的消息"
	schema.SchemaProps = spec.SchemaProps{
		Required: []string{"tid"},
		Properties: map[string]spec.Schema{
			"tid": spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:        spec.StringOrArray{"string"},
					Description: "任务标识",
				},
			},
			"progress": spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:        spec.StringOrArray{"string"},
					Description: "任务进度百分比",
				},
			},
			"completed": spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:        spec.StringOrArray{"string"},
					Description: "已完成数量",
				},
			},
			"undone": spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:        spec.StringOrArray{"string"},
					Description: "剩余数量",
				},
			},

		},
	}
	return
}

func newOperation(tagName,summary, opDescribetion string, params []spec.Parameter,responseDescription string, respSchema *spec.Schema) (op *spec.Operation) {
	op = newOperationFull(tagName,summary,opDescribetion,params,responseDescription,respSchema,[]string{"application/json","application/octet-stream"},[]string{"application/json","application/octet-stream"})
	return
}
func newOperationFull(tagName,summary, opDescribetion string, params []spec.Parameter,responseDescription string, respSchema *spec.Schema,consumes []string,produces []string) (op *spec.Operation) {
	op = &spec.Operation{
		spec.VendorExtensible{}, spec.OperationProps{
			Summary:summary,
			Description: opDescribetion,
			Consumes:consumes,
			Produces:produces,
			Tags:        []string{tagName},
			Parameters:  params,
			Responses: &spec.Responses{
				spec.VendorExtensible{},
				spec.ResponsesProps{
					&spec.Response{
						ResponseProps:spec.ResponseProps{
							Description:"错误消息",
							Schema: &spec.Schema{
								SchemaProps:spec.SchemaProps{
									Ref:getModelSwaggerRef("ErrMsg"),
								},
							},
						},
					},
					map[int]spec.Response{
						200: {
							ResponseProps: spec.ResponseProps{
								Description: responseDescription,
								Schema: respSchema,
							},
						},
						401:{
							ResponseProps: spec.ResponseProps{
								Description: "未认证",
							},
						},
						403:{
							ResponseProps: spec.ResponseProps{
								Description: "未授权",
							},
						},
					},
				},
			},
		},
	}
	return
}
func getModelSwaggerRef(t string) (ref spec.Ref) {
	ref = spec.Ref{}
	ref.Ref, _ = jsonreference.New(fmt.Sprintf("#/definitions/%s", t))
	return
}