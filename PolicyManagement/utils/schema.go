package utils

const (
	AccessorsSchema = `{
		"type": "array",
		"items": {
			"type": "object",
			"properties": {
				"accessor_id": {
					"type": "string"
				},
				"accessor_type": {
					"type": "string",
					"enum": ["user", "department"]
				}
			},
			"required": [
				"accessor_id",
				"accessor_type"
			]
		}
	}`

	PolicesSchema = `{
    "type": "array",
    "items": {
        "anyOf": [
            {
                "$ref": "#/definitions/password_strength_meter_schema"
            },
            {
                "$ref": "#/definitions/multi_factor_auth_schema"
            },
            {
                "$ref": "#/definitions/client_restriction_schema"
            },
            {
                "$ref": "#/definitions/user_document_sharing_schema"
            },
            {
                "$ref": "#/definitions/user_document_schema"
            },
            {
                "$ref": "#/definitions/network_resctriction_schema"
            },
            {
                "$ref": "#/definitions/no_network_policy_accessor"
            },
            {
                "$ref": "#/definitions/system_protection_levels_schema"
            }
        ]
    },
    "definitions": {
        "password_strength_meter_schema": {
            "type": "object",
            "required": [
                "name",
                "value"
            ],
            "properties": {
                "name": {
                    "enum": [
                        "password_strength_meter"
                    ]
                },
                "value": {
                    "type": "object",
                    "required": [
                        "enable",
                        "length"
                    ],
                    "properties": {
                        "enable": {
                            "type": "boolean"
                        },
                        "length": {
                            "type": "integer"
                        }
                    }
                }
            }
        },
        "multi_factor_auth_schema": {
            "type": "object",
            "required": [
                "name",
                "value"
            ],
            "properties": {
                "name": {
                    "enum": [
                        "multi_factor_auth"
                    ]
                },
                "value": {
                    "type": "object",
                    "required": [
                        "enable",
                        "image_vcode",
                        "password_error_count",
                        "sms_vcode",
                        "otp"
                    ],
                    "properties": {
                        "enable": {
                            "type": "boolean"
                        },
                        "image_vcode": {
                            "type": "boolean"
                        },
                        "password_error_count": {
                            "type": "integer",
                            "minimum": 0,
                            "maximum": 99
                        },
                        "sms_vcode": {
                            "type": "boolean"
                        },
                        "otp": {
                            "type": "boolean"
                        }
                    }
                }
            }
        },
        "client_restriction_schema": {
            "type": "object",
            "required": [
                "name",
                "value"
            ],
            "properties": {
                "name": {
                    "enum": [
                        "client_restriction"
                    ]
                },
                "value": {
                    "type": "object",
                    "required": [
                        "pc_web",
                        "mobile_web",
                        "windows",
                        "mac",
                        "android",
                        "ios",
                        "linux"
                    ],
                    "properties": {
                        "pc_web": {
                            "type": "boolean"
                        },
                        "mobile_web": {
                            "type": "boolean"
                        },
                        "windows": {
                            "type": "boolean"
                        },
                        "mac": {
                            "type": "boolean"
                        },
                        "android": {
                            "type": "boolean"
                        },
                        "ios": {
                            "type": "boolean"
                        },
                        "linux": {
                            "type": "boolean"
                        }
                    }
                }
            }
        },
        "user_document_sharing_schema": {
            "type": "object",
            "required": [
                "name",
                "value"
            ],
            "properties": {
                "name": {
                    "enum": [
                        "user_document_sharing"
                    ]
                },
                "value": {
                    "type": "object",
                    "required": [
                        "anyshare",
                        "http"
                    ],
                    "properties": {
                        "anyshare": {
                            "type": "boolean"
                        },
                        "http": {
                            "type": "boolean"
                        }
                    }
                }
            }
        },
        "user_document_schema": {
            "type": "object",
            "required": [
                "name",
                "value"
            ],
            "properties": {
                "name": {
                    "enum": [
                        "user_document"
                    ]
                },
                "value": {
                    "type": "object",
                    "required": [
                        "create",
                        "size"
                    ],
                    "properties": {
                        "create": {
                            "type": "boolean"
                        },
                        "size": {
                            "type": "number",
                            "minimum": 1
                        }
                    }
                }
            }
		},
        "network_resctriction_schema": {
            "type": "object",
            "required": [
                "name",
                "value"
            ],
            "properties": {
                "name": {
                    "enum": [
                        "network_restriction"
                    ]
                },
                "value": {
                    "type": "object",
                    "required": [
                        "is_enabled"
                    ],
                    "properties": {
                        "is_enabled": {
                            "type": "boolean"
                        }
                    }
                }
            }
        },
        "system_protection_levels_schema": {
            "type": "object",
            "required": [
                "name",
                "value"
            ],
            "properties": {
                "name": {
                    "enum": [
                        "system_protection_levels"
                    ]
                },
                "value": {
                    "type": "object",
                    "required": [
                        "level"
                    ],
                    "properties": {
                        "level": {
                            "type": "integer",
                            "minimum": 1,
                            "maximum": 3
                        }
                    }
                }
            }
        },
        "no_network_policy_accessor": {
            "type": "object",
            "required": [
                "name",
                "value"
            ],
            "properties": {
                "name": {
                    "enum": [
                        "no_network_policy_accessor"
                    ]
                },
                "value": {
                    "type": "object",
                    "required": [
                        "is_enabled"
                    ],
                    "properties": {
                        "is_enabled": {
                            "type": "boolean"
                        }
                    }
                }
            }
        }
    }
}`

	WatermarkTemplateSchema = `{
    "type": "object",
    "required": ["name","layout","user_name_config","custom_config","time_config"],
    "properties":{
        "name": {
            "type": "string",
            "minLength": 1,
            "maxLength": 128
        },
        "layout": {
            "enum": ["center","tile"]
        },
        "user_name_config": {
            "$ref": "#/definitions/common_config_schema"
        },
        "custom_config": {
            "allOf": [
                {"$ref": "#/definitions/common_config_schema"}
            ],
            "if": {
                "properties": {"enabled": {"const": false} }
            },
            "then": {
                "required": ["content"],
                "properties": {
                    "content": {
                        "type": "string",
                        "minLength": 0,
                        "maxLength": 50
                    }
                }
            },
            "else": {
                "required": ["content"],
                "properties": {
                    "content": {
                        "type": "string",
                        "minLength": 1,
                        "maxLength": 50
                    }
                }
            }
        },
        "time_config": {
            "$ref": "#/definitions/common_config_schema"
        }
    },
    "definitions": {
        "common_config_schema": {
            "type": "object",
            "required": ["color","enabled","font_size","opacity"],
            "properties": {
                "color": {
                    "type": "string",
                    "pattern": "^#[0-9a-fA-F]{6}$"
                },
                "enabled": {
                    "type": "boolean"
                },
                "font_size": {
                    "type": "integer",
                    "minimum": 5,
                    "maximum": 144
                },
                "opacity": {
                    "type": "integer",
                    "minimum": 0,
                    "maximum": 100
                }
            }
        }
    }
}`

	WatermarkPolicySchema = `{
    "type": "object",
    "required": ["docid","watermark_template_id"],
    "properties": {
        "docid": {
            "type": "string"
        },
        "watermark_template_id": {
            "type": "string"
        }
    }
}`

	BatchWatermarkPolicySchema = `{
    "type": "object",
    "required": ["docid","watermark_template_id"],
    "properties": {
        "docid": {
            "type": "array",
            "items": {
                "type": "string"
            },
            "minItems": 1
        },
        "watermark_template_id": {
            "type": "string"
        }
    }
}`
)
