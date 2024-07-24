// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/accounts/balance": {
            "get": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Get balance for account by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "account"
                ],
                "summary": "Get balance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user id",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_api_v1.balanceResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/v1/accounts/create": {
            "post": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Create account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "account"
                ],
                "summary": "Create account",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_api_v1.accountCreateInput"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/v1/accounts/deposit": {
            "patch": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Deposit on account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "account"
                ],
                "summary": "Account deposit",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_api_v1.accountDepositInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/v1/accounts/transfer": {
            "post": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Transfer from account to account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "account"
                ],
                "summary": "Account transfer",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_api_v1.accountTransferInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/v1/accounts/withdraw": {
            "patch": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Withdraw from account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "account"
                ],
                "summary": "Account withdraw",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_api_v1.accountWithdrawInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/v1/operations/history": {
            "get": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Get account transactions history",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "operation"
                ],
                "summary": "Get history",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_api_v1.operationHistoryInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/avito_intership_internal_service.HistoryOutput"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/v1/operations/report": {
            "get": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Get monthly report, ordered by products ids",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "operation"
                ],
                "summary": "Get report",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_api_v1.operationReportInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/v1/reservations/cancel": {
            "delete": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "cancel product amount reservation and return money to account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reservation"
                ],
                "summary": "cancel reservation",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_api_v1.reservationCancelInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/v1/reservations/create": {
            "post": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Create product amount reservation",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reservation"
                ],
                "summary": "create reservation",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_api_v1.reservationCreateInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_api_v1.reservationResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/v1/reservations/revenue": {
            "post": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "confirm reservation",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reservation"
                ],
                "summary": "revenue reservation",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_api_v1.reservationRevenueInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "avito_intership_internal_service.HistoryOutput": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "created_at": {
                    "type": "string"
                },
                "operation_id": {
                    "type": "integer"
                },
                "order_id": {
                    "type": "integer"
                },
                "product_id": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "echo.HTTPError": {
            "type": "object",
            "properties": {
                "message": {}
            }
        },
        "internal_api_v1.accountCreateInput": {
            "type": "object",
            "required": [
                "user_id"
            ],
            "properties": {
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "internal_api_v1.accountDepositInput": {
            "type": "object",
            "required": [
                "amount",
                "user_id"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "internal_api_v1.accountTransferInput": {
            "type": "object",
            "required": [
                "amount",
                "from",
                "to"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "from": {
                    "type": "integer"
                },
                "to": {
                    "type": "integer"
                }
            }
        },
        "internal_api_v1.accountWithdrawInput": {
            "type": "object",
            "required": [
                "amount",
                "user_id"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "internal_api_v1.balanceResponse": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number"
                }
            }
        },
        "internal_api_v1.operationHistoryInput": {
            "type": "object",
            "required": [
                "user_id"
            ],
            "properties": {
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                },
                "sort": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "internal_api_v1.operationReportInput": {
            "type": "object",
            "required": [
                "month",
                "year"
            ],
            "properties": {
                "month": {
                    "type": "integer"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "internal_api_v1.reservationCancelInput": {
            "type": "object",
            "properties": {
                "reservation_id": {
                    "type": "integer"
                }
            }
        },
        "internal_api_v1.reservationCreateInput": {
            "type": "object",
            "required": [
                "amount",
                "order_id",
                "product_id",
                "user_id"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "order_id": {
                    "type": "integer"
                },
                "product_id": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "internal_api_v1.reservationResponse": {
            "type": "object",
            "properties": {
                "reservation_id": {
                    "type": "integer"
                }
            }
        },
        "internal_api_v1.reservationRevenueInput": {
            "type": "object",
            "properties": {
                "reservation_id": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "JWT": {
            "description": "JWT token",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Api for account balance management",
	Description:      "Api for balance management. Include operations, e.g. deposit, withdraw, transfer, reservation, etc",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
