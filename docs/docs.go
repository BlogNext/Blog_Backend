// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://github.com/BlogNext",
        "contact": {
            "name": "Ly",
            "url": "https://github.com/FlashFeiFei?tab=repositories",
            "email": "51785816@qq.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/front/blog/detail": {
            "get": {
                "description": "博客详情",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "前台-博客"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "description": "博客id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "json格式",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/front/blog/get_list": {
            "get": {
                "description": "获取博客列表",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "前台-博客"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "description": "一页多少条",
                        "name": "per_page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "第几页",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "博客分类",
                        "name": "blog_type_id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "json格式",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/front/blog/get_list_by_sort": {
            "get": {
                "description": "按排序维度获取排序博客",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "前台-博客"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "description": "一页多少条，默认值5",
                        "name": "per_page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "排序维度，默认值browse_total",
                        "name": "sort_dimension",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "json格式",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/front/blog/get_stat": {
            "get": {
                "description": "blogInfo模块统计展示",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "前台-博客"
                ],
                "responses": {
                    "200": {
                        "description": "json格式",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/front/blog/search_blog": {
            "get": {
                "description": "搜素博客",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "前台-博客"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "description": "一页多少条",
                        "name": "per_page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "第几页",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "搜索等级，默认mysql搜索",
                        "name": "search_level",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "搜索关键字",
                        "name": "keyword",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "json格式",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/front/blog_type/get_list": {
            "get": {
                "description": "获取博客类型列表",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "前台-博客类型（知识库）"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "description": "每页多少条记录",
                        "name": "per_page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "第几页",
                        "name": "page",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "json格式",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/front/login/Login_by_yuque": {
            "post": {
                "description": "语雀账号登录",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "前台-登录"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "语雀login",
                        "name": "login",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "登录密码",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "json格式",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/front/person/blog_list": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "私人博客列表，这里的私人只的是登录的用户",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "前台-博客-登录用户"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "description": "一页多少条",
                        "name": "per_page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "第几页",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "博客分类",
                        "name": "blog_type_id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "json格式",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "x-access-token",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "blog.laughingzhu.com",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "晓琛博客",
	Description: "改swagger用于和前端联调用的，正在努力的测试中，目前的测试securityDefinitions.apikey ApiKeyAuth只能带一个请求头",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
