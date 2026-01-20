package common

const (
	// 新增审计日志
	PostAuditLog = `{
		"type": "object",
		"required": ["user_id", "user_type", "level", "date", "msg", "op_type", "out_biz_id"],
		"properties": {
			"user_id": {
				"type": "string",
				"minLength": 0,
				"maxLength": 40
			},
			"user_name": {
				"type": ["string", "null"],
				"minLength": 0,
				"maxLength": 128
			},
			"user_type": {
				"type": "string",
				"enum": ["authenticated_user", "anonymous_user", "app", "internal_service"]
			},
			"level": {
				"type": "integer"
			},
			"date": {
				"type": "integer"
			},
			"ip": {
				"type": ["string", "null"],
				"minLength": 0,
				"maxLength": 40
			},
			"mac": {
				"type": ["string", "null"],
				"minLength": 0,
				"maxLength": 40
			},
			"msg": {
				"type": "string"
			},
			"ex_msg": {
				"type": ["string", "null"]
			},
			"user_agent": {
				"type": ["string", "null"],
				"maxLength": 1024
			},
			"additional_info": {
				"type": ["string", "null"]
			},
			"op_type": {
				"type": "integer"
			},
			"out_biz_id": {
				"type": "string",
				"minLength": 1,
				"maxLength": 128
			},
			"dept_paths": {
				"type": ["string", "null"]
			}
		}
	}`

	PutDumpStrategy = `{
		"type": "object",
		"properties": {
			"retention_period": {
				"type": "integer",
				"minimum": 1,
				"maximum": 999999
			},
			"retention_period_unit": {
				"type": "string",
				"enum": ["day", "week", "month", "year"]
			},
			"dump_time": {
				"type": "string",
                "pattern": "^(2[0-3]|[01][0-9]):[0-5][0-9]:[0-5][0-9]$"
			},
			"dump_format": {
				"type": "string",
				"enum": ["xml", "csv"]
			}
		}
	}`

	ScopeStrategy = `{
        "type": "object",
		"required": ["category", "type", "role", "scope"],
		"properties": {
			"category": {
				"type": "integer",
				"enum": [1, 2]
			},
			"type": {
				"type": "integer",
				"enum": [10, 11, 12]
			},
			"role": {
				"type": "string",
				"enum": ["sys_admin", "sec_admin", "audit_admin"]
			},
			"scope": {
				"type": "array",
				"items": {
					"type": "string",
					"enum": ["sys_admin", "sec_admin", "audit_admin", "normal_user"]
				}
			}
		}
    }`

	PutHistoryPwdStatus = `{
		"type": "object",
		"required": ["status"],
		"properties": {
			"status": {
				"type": "boolean"
			}
		}
	}`

	PostHistoryTask = `{
		"type": "object",
		"required": ["obj_id"],
		"properties": {
			"pwd": {
				"type": "string"
			},
			"obj_id": {
				"type": "string"
			}
		}
	}`

	// 获取对象个性化特征值
	GetPersFeature = `{
		"type": "object",
		"required": ["obj_id", "obj_type", "key"],
		"properties": {
			"obj_id": {
				"type": "string",
				"minLength": 1,
				"maxLength": 128
			},
			"obj_type": {
				"type": string,
				"enum": ["user", "dept"]
			},
			"key": {
				"type": "string",
				"minLength": 1,
				"maxLength": 64
			}
		}
	}`
)
