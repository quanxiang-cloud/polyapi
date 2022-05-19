package swagger

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestLoadSwagger(t *testing.T) {
	c := []byte(`
	{
    "definitions":{
        "handler.Response":{
            "properties":{
                "bar":{
                    "type":"string"
                }
            },
            "type":"object"
        }
    },
    "info":{
    	"version":"1.0",
        "contact":{

        }
    },
    "basePath":"/base/",
    "schemes": [
	    "http"
	],
    "paths":{
        "/api/":{
            "post":{
            	"x-open-request": true,
            	"operationId":"someName",
                "description":"handler description",
                "parameters":[
                    {
                        "description":"request",
                        "in":"body",
                        "name":"request",
                        "required":true,
                        "schema":{
                            "type":"string"
                        }
                    }
                ],
                "produces":[
                    "application/json"
                ],
                "responses":{
                    "200":{
                        "description":"OK",
                        "schema":{
                            "$ref":"#/definitions/handler.Response"
                        }
                    }
                },
                "summary":"handler summary",
                "tags":[
                    "handler"
                ]
            }
        }
    },
    "swagger":"2.0"
}
`)
	cfg := &APIServiceConfig{
		Schema:    "https",
		Host:      "api.xxx.com:8080",
		Service:   "/system/service",
		Namespace: "/system/namespace",
		AuthType:  "signature",
	}
	c = c
	apis, err := ParseSwagger([]byte(errSwag4), cfg)
	if err != nil {
		panic(err)
	}

	showObj(apis[0])

	fmt.Println(string(apis[0].Path))
	fmt.Println(apis[0].URL)
}

func showObj(obj interface{}) {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
	}
}

const okSwag2 = `
{
        "x-consts": [
        	{
				"name":"assertempty",
				"type":"string",
				"in":"body",
				"data":"foo"
			}
        ],
        "definitions": {
            "handler.Response": {
                "type": "object",
                "properties": {
                    "bar": {
                        "type": "string"
                    }
                }
            }
        },
        "host": "test_host",
        "swagger": "2.0",
        "info": {
            "title": "",
            "version": "",
            "description": "auto generate at 2021-11-17T06:59:22UTC",
            "contact": {
                "name": "",
                "url": "",
                "email": ""
            }
        },
        "schemes": [
            "http"
        ],
        "basePath": "",
        "paths": {
            "/": {
                "post": {
                    "x-consts": [],
                    "operationId": "minhj_float3_alpha011",
                    "parameters": [
                        {
                            "description": "request",
                            "name": "request",
                            "in": "body",
                            "required": true,
                            "schema": {
                                "type": "string"
                            }
                        }
                    ],
                    "responses": {
                        "200": {
                            "description": "OK",
                            "schema": {
                                "type": "object",
                                "properties": {
                                    "bar": {
                                        "type": "string"
                                    }
                                }
                            },
                            "headers": null
                        }
                    },
                    "consumes": [
                        "application/json"
                    ],
                    "produces": [
                        "application/json"
                    ],
                    "summary": "handlersummary",
                    "description": "handlerdescription"
                }
            }
        }
    }
`

const okSwag = `
{
    "swagger": "2.0",
    "info": {
        "version": "v1"
    },
    "basePath": "/",
    "schemes": [
        "http"
    ],
    "x-consts": [
        {
            "name": "const",
            "type": "string",
            "required": false,
            "data":"",
            "in":"body",
            "description": ""
        }
    ],
    "paths": {
        "/api/v1/test/:path": {
            "post": {
                "summary": "名字",
                "description": "描述",
                "operationId": "testapi",
                "consumes": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "name": "path",
                        "type": "string",
                        "required": true,
                        "description": "",
                        "in": "path"
                    },
                    {
                        "name": "formData",
                        "type": "string",
                        "required": true,
                        "description": "",
                        "in": "formData"
                    },
                    {
                        "name": "query",
                        "type": "string",
                        "required": false,
                        "description": "",
                        "in": "query"
                    },
                    {
                        "name": "header",
                        "type": "string",
                        "required": false,
                        "description": "",
                        "in": "header"
                    },
                    {
                        "name": "body",
                        "requiredX": ["x"],
                        "description": "",
                        "in": "body",
                        "schema": {
                        	"name":"x",
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "successful operation",
                        "schema": {
                            "type": "object",
                            "title": "api result",
                            "properties": {
                                "ret": {
                                    "type": "string",
                                    "requiredX": false,
                                    "description": ""
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}
`

const errSwag1 = `
{
    "swagger": "2.0",
    "info": {
        "version": "v1"
    },
    "basePath": "/",
    "schemes": [
        "http"
    ],
    "x-consts": [
        {
            "name": "x",
            "type": "string",
            "required": false,
            "in":"header",
            "data":"",
            "description": ""
        }
    ],
    "paths": {
        "/api": {
            "post": {
                "summary": "ss",
                "description": "",
                "operationId": "ttt",
                "consumes": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "name": "x",
                        "type": "string",
                        "required": false,
                        "description": "",
                        "in": "query"
                    },
                    {
                        "name": "x",
                        "type": "string",
                        "required": false,
                        "description": "",
                        "in": "header"
                    },
                    {
                        "name": "a",
                        "type": "string",
                        "required": false,
                        "description": "",
                        "schemaXXX": {
                            "type": "object",
                            "title": "api result",
                            "properties": {}
                        },
                        "in": "body"
                    }
                ],
                "response": {
                    "200": {
                        "description": "successful operation",
                        "schema": {
                            "type": "object",
                            "title": "api result",
                            "properties": {}
                        }
                    }
                }
            }
        }
    }
}
`

const errSwag2 = `
{
    "swagger": "2.0",
    "info": {
        "version": "v1"
    },
    "basePath": "/",
    "schemes": [
        "http"
    ],
    "x-consts": [
        {
            "name": "x-const",
            "description": "",
            "in": "body",
            "data": "aa"
        },
        {
            "name": "foo",
            "description": "",
            "in": "body",
            "data": "bar"
        }
    ],
    "paths": {
        "/api/v1/app/:appId/p2": {
            "post": {
                "summary": "test api",
                "description": "",
                "operationId": "getapp",
                "consumes": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "name": "appId",
                        "description": "",
                        "in": "path"
                    },
                    {
                        "name": "name",
                        "description": "",
                        "in": "query"
                    },
                    {
                        "name": "form-data",
                        "description": "",
                        "in": "body",
                        "schema": {
                            "type": "object",
                            "required": [
                                "name"
                            ],
                            "properties": {
                                "name": {
                                    "description": ""
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "schema": {
                            "type": "object",
                            "properties": {
                                "app_list": {
                                    "description": ""
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}
`

const errSwag3 = `
{
    "swagger": "2.0",
    "info": {
        "title": "\u7528\u6237\u7ba1\u7406",
        "version": "1.0.0",
        "description": "\u7528\u6237\u7ba1\u7406"
    },
    "host": "172.30.1.208:32639",
    "basePath": "/rest/s1/users",
    "schemes": [
        "http"
    ],
    "securityDefinitions": {
        "basicAuth": {
            "type": "basic",
            "description": "HTTP Basic Authentication"
        },
        "api_key": {
            "type": "apiKey",
            "name": "api_key",
            "in": "header",
            "description": "HTTP Header api_key"
        }
    },
    "consumes": [
        "application/json",
        "multipart/form-data"
    ],
    "produces": [
        "application/json"
    ],
    "tags": [],
    "paths": {
        "/": {
            "post": {
                "summary": "create User",
                "description": "\u521b\u5efa\u7528\u6237",
                "security": [
                    {
                        "basicAuth": []
                    },
                    {
                        "api_key": []
                    }
                ],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/userExtend.create#User.In"
                        }
                    }
                ],
                "responses": {
                    "401": {
                        "description": "Authentication required"
                    },
                    "403": {
                        "description": "Access Forbidden (no authz)"
                    },
                    "429": {
                        "description": "Too Many Requests (tarpit)"
                    },
                    "500": {
                        "description": "General Error"
                    }
                }
            },
            "get": {
                "summary": "query UserPage",
                "description": "\u67e5\u8be2\u7528\u6237\u5217\u8868",
                "security": [
                    {
                        "basicAuth": []
                    },
                    {
                        "api_key": []
                    }
                ],
                "parameters": [
                    {
                        "name": "orgId",
                        "in": "query",
                        "required": false,
                        "type": "string",
                        "format": "",
                        "description": "\u90e8\u95e8id"
                    }
                ],
                "responses": {
                    "401": {
                        "description": "Authentication required"
                    },
                    "403": {
                        "description": "Access Forbidden (no authz)"
                    },
                    "429": {
                        "description": "Too Many Requests (tarpit)"
                    },
                    "500": {
                        "description": "General Error"
                    }
                }
            }
        },
        "/:id": {
            "put": {
                "summary": "modify User",
                "description": "\u4fee\u6539\u7528\u6237",
                "security": [
                    {
                        "basicAuth": []
                    },
                    {
                        "api_key": []
                    }
                ],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": null
                    },
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/userExtend.modify#User.In"
                        }
                    }
                ],
                "responses": {
                    "401": {
                        "description": "Authentication required"
                    },
                    "403": {
                        "description": "Access Forbidden (no authz)"
                    },
                    "429": {
                        "description": "Too Many Requests (tarpit)"
                    },
                    "500": {
                        "description": "General Error"
                    }
                }
            },
            "delete": {
                "summary": "delete User",
                "description": "\u5220\u9664\u7528\u6237",
                "security": [
                    {
                        "basicAuth": []
                    },
                    {
                        "api_key": []
                    }
                ],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": null
                    }
                ],
                "responses": {
                    "401": {
                        "description": "Authentication required"
                    },
                    "403": {
                        "description": "Access Forbidden (no authz)"
                    },
                    "429": {
                        "description": "Too Many Requests (tarpit)"
                    },
                    "500": {
                        "description": "General Error"
                    }
                }
            },
            "get": {
                "summary": "find User",
                "description": "\u7528\u6237\u8be6\u60c5",
                "security": [
                    {
                        "basicAuth": []
                    },
                    {
                        "api_key": []
                    }
                ],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": null
                    }
                ],
                "responses": {
                    "401": {
                        "description": "Authentication required"
                    },
                    "403": {
                        "description": "Access Forbidden (no authz)"
                    },
                    "429": {
                        "description": "Too Many Requests (tarpit)"
                    },
                    "500": {
                        "description": "General Error"
                    }
                }
            }
        },
        "/:id/password": {
            "post": {
                "summary": "modify Password",
                "description": "\u4fee\u6539\u5bc6\u7801",
                "security": [
                    {
                        "basicAuth": []
                    },
                    {
                        "api_key": []
                    }
                ],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": null
                    },
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/userExtend.modify#Password.In"
                        }
                    }
                ],
                "responses": {
                    "401": {
                        "description": "Authentication required"
                    },
                    "403": {
                        "description": "Access Forbidden (no authz)"
                    },
                    "429": {
                        "description": "Too Many Requests (tarpit)"
                    },
                    "500": {
                        "description": "General Error"
                    }
                }
            }
        },
        "/:id/resetPassword": {
            "post": {
                "summary": "reset Password",
                "description": "\u91cd\u7f6e\u5bc6\u7801",
                "security": [
                    {
                        "basicAuth": []
                    },
                    {
                        "api_key": []
                    }
                ],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": null
                    },
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/userExtend.reset#Password.In"
                        }
                    }
                ],
                "responses": {
                    "401": {
                        "description": "Authentication required"
                    },
                    "403": {
                        "description": "Access Forbidden (no authz)"
                    },
                    "429": {
                        "description": "Too Many Requests (tarpit)"
                    },
                    "500": {
                        "description": "General Error"
                    }
                }
            }
        },
        "/:id/tenant": {
            "post": {
                "summary": "set UserTenant",
                "description": "\u8bbe\u7f6e\u7528\u6237\u79df\u6237",
                "security": [
                    {
                        "basicAuth": []
                    },
                    {
                        "api_key": []
                    }
                ],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "string",
                        "description": null
                    },
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/userExtend.set#UserTenant.In"
                        }
                    }
                ],
                "responses": {
                    "401": {
                        "description": "Authentication required"
                    },
                    "403": {
                        "description": "Access Forbidden (no authz)"
                    },
                    "429": {
                        "description": "Too Many Requests (tarpit)"
                    },
                    "500": {
                        "description": "General Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "userExtend.create#User.In": {
            "type": "object",
            "properties": {
                "username": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "realName": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "cellphone": {
                    "type": "string"
                },
                "telephone": {
                    "type": "string"
                },
                "orgId": {
                    "type": "string"
                },
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "photo": {
                    "type": "string"
                }
            },
            "required": [
                "username",
                "password"
            ]
        },
        "userExtend.modify#User.In": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "realName": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "cellphone": {
                    "type": "string"
                },
                "telephone": {
                    "type": "string"
                },
                "orgId": {
                    "type": "string"
                },
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            },
            "required": [
                "id"
            ]
        },
        "userExtend.modify#Password.In": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "oldPassword": {
                    "type": "string"
                },
                "newPassword": {
                    "type": "string"
                }
            },
            "required": [
                "id",
                "oldPassword",
                "newPassword"
            ]
        },
        "userExtend.reset#Password.In": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "newPassword": {
                    "type": "string"
                }
            },
            "required": [
                "id"
            ]
        },
        "userExtend.set#UserTenant.In": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "tenantId": {
                    "type": "string"
                }
            },
            "required": [
                "id",
                "tenantId"
            ]
        }
    }
}
`

const errSwag4 = `
{
    "schemes": [
        "http"
    ],
    "swagger": "2.0",
    "info": {
        "contact": {}
    },
    "paths": {
        "/": {
            "post": {
                "description": "handler 转化VirtualMachine结构至ecs",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "vsphere"
                ],
                "summary": "转化VirtualMachine结构至ecs",
                "parameters": [
                    {
                        "description": "请求参数",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/handler.EcsInstance"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.EcsInstance": {
            "type": "object",
            "properties": {
                "cpu": {
                    "type": "integer"
                },
                "hostname": {
                    "type": "string"
                },
                "instance_id": {
                    "type": "string"
                },
                "instance_name": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                },
                "memory": {
                    "type": "integer"
                },
                "os_type": {
                    "type": "string"
                },
                "specialized": {
                    "type": "object",
                    "additionalProperties": true
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "handler.Request": {
            "type": "object",
            "properties": {
                "virtualMachines": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handler.VirtualMachine"
                    }
                }
            }
        },
        "handler.Response": {
            "type": "object",
            "properties": {
                "data": {},
                "error": {
                    "type": "string"
                },
                "is_success": {
                    "type": "boolean"
                }
            }
        },
        "handler.VirtualMachine": {
            "type": "object",
            "properties": {
                "clusterId": {
                    "type": "string"
                },
                "dataDisks": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "datastoreId": {
                                "type": "string"
                            },
                            "format": {
                                "type": "string"
                            },
                            "id": {
                                "type": "string"
                            },
                            "key": {
                                "type": "integer"
                            },
                            "mode": {
                                "type": "string"
                            },
                            "sharing": {
                                "type": "string"
                            },
                            "size": {
                                "type": "integer"
                            }
                        }
                    }
                },
                "datacenterId": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "hostId": {
                    "type": "string"
                },
                "hostname": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "instanceUUID": {
                    "type": "string"
                },
                "ip_address": {
                    "type": "string"
                },
                "memoryMB": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "networkInterfaces": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "adapterType": {
                                "type": "string"
                            },
                            "id": {
                                "type": "string"
                            },
                            "ipInfo": {
                                "type": "array",
                                "items": {
                                    "type": "object",
                                    "properties": {
                                        "ipAddress": {
                                            "type": "string"
                                        },
                                        "state": {
                                            "type": "string"
                                        }
                                    }
                                }
                            },
                            "key": {
                                "type": "integer"
                            },
                            "macAddress": {
                                "type": "string"
                            },
                            "networkId": {
                                "type": "string"
                            }
                        }
                    }
                },
                "numCPU": {
                    "type": "integer"
                },
                "numCoresPerSocket": {
                    "type": "integer"
                },
                "osFamily": {
                    "type": "string"
                },
                "osName": {
                    "type": "string"
                },
                "power_state": {
                    "type": "string"
                },
                "sysDisk": {
                    "type": "object",
                    "properties": {
                        "datastoreId": {
                            "type": "string"
                        },
                        "format": {
                            "type": "string"
                        },
                        "id": {
                            "type": "string"
                        },
                        "key": {
                            "type": "integer"
                        },
                        "mode": {
                            "type": "string"
                        },
                        "sharing": {
                            "type": "string"
                        },
                        "size": {
                            "type": "integer"
                        }
                    }
                },
                "toolsHasInstalled": {
                    "type": "boolean"
                },
                "tools_status": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        }
    }
}
`
