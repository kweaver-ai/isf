package persconsts

const (
	UserStatic = `{
		"id": {
			"type": "string",
			"title": "用户id",
			"description": "用户id",
			"required": true
		},
		"name": {
			"type": "string",
			"title": "名称",
			"description": "用户名称",
			"required": true
		},
		"email": {
			"type": "string",
			"title": "邮箱",
			"description": "用户邮箱",
			"required": true
		},
		"telephone": {
			"type": "string",
			"title": "电话",
			"description": "用户电话",
			"required": true
		},
		"level": {
			"type": "string",
			"title": "密级",
			"description": "用户密级",
			"required": true
		},
		"roles": {
			"type": "array",
			"title": "角色",
			"description": "用户角色",
			"required": true,
			"items": {
				"type": "object",
				"properties": {
					"role_id": {
						"type": "string",
						"title": "角色id",
						"description": "角色id"
					},
					"role_name": {
						"type": "string",
						"title": "角色名",
						"description": "角色名"
					}
				}
			}
		},
		"parent_deps": {
			"type": "array",
			"title": "父部门信息",
			"description": "父部门信息，描述多个父部门的层级关系信息，每个父部门层级数组内第一个对象是根部门，最后一个对象是直接父部门",
			"required": true,
			"items": {
				"type": "object",
				"properties": {
					"dept_id": {
						"type": "string",
						"title": "部门id",
						"description": "部门id"
					},
					"dept_name": {
						"type": "string",
						"title": "部门名称",
						"description": "部门名称"
					}
				}
			}
		},
		"groups": {
			"type": "array",
			"title": "用户组列表",
			"description": "用户组列表",
			"required": true,
			"items": {
                "type": "object",
                "properties": {
                    "id": {
                        "type": "string",
                        "title": "用户组id",
                        "description": "用户组id"
                    },
                    "name": {
                        "type": "string",
                        "title": "用户组名称",
                        "desciption": "用户组名称"
                    },
                    "type": {
                        "type": "string",
                        "title": "用户组类别",
						"description": "用户组类别"
					}
				}
			}
		},
		"enabled": {
			"type": "boolean",
			"title": "用户状态",
			"description": "是否可用",
			"required": true
		},
		"frozen": {
			"type": "boolean",
			"title": "是否冻结",
			"description": "是否冻结",
			"required": true
		},
		"positions": {
			"type": "array",
			"title": "岗位职务",
			"description": "岗位职务",
			"items": {
				"type": "string"
			}
		},
		"graduated": {
			"type": "array",
			"title": "毕业院校",
			"desciption": "毕业院校",
			"items": {
				"type": "string"
			}
		},
		"native": {
			"type": "string",
			"title": "籍贯",
			"desciption": "籍贯"
		},
		"birthday": {
			"type": "string",
			"title": "出生日期",
			"desciption": "出生日期"
		},
		"certificate": {
			"type": "array",
			"title": "认证信息",
			"desciption": "认证证书名称",
			"items": {
				"type": "string"
			}
		},
		"address": {
			"type": "string",
			"title": "工作地点",
			"desciption": "地方代码"
		}
	}`

	DeptStatic = `{
		"id": {
			"type": "string",
			"title": "部门id",
			"description": "部门id"
		},
		"name": {
			"type": "string",
			"title": "部门名称",
			"description": "部门名称"
		},
		"level": {
			"type": "integer",
			"title": "部门层级",
			"description": "部门层级"
		},
		"managers": {
			"type": "array",
			"title": "部门管理员列表",
			"description": "部门管理员列表",
			"items": {
				"type": "object",
				"properties": {
					"id": {
						"type": "string",
						"title": "管理员id",
						"description": "管理员id"
					},
					"name": {
						"type": "string",
						"title": "管理员名称",
						"description": "管理员名称"
					}
				}
			}
		},
		"parent_deps": {
			"type": "array",
			"title": "父部门列表",
			"description": "父部门列表，描述多个父部门的层级关系信息，每个父部门层级数组内第一个对象是根部门，最后一个对象是直接父部门",
			"items": {
				"type": "object",
				"properties": {
					"id": {
						"type": "string",
						"title": "部门id",
						"description": "部门id"
					},
					"name": {
						"type": "string",
						"title": "部门名称",
						"description": "部门名称"
					}
				}
			}
		}
	}`

	DocCenterDynamic = `{
		"{obj_type}": {
			"type": "object",
			"title": "文档对象类型",
			"description": "文档对象类型变量，文档类型：file|folder等",
			"properties": {
				"operation_preferences": {
					"type": "array",
					"title": "文档操作偏好特征",
					"description": "文档操作偏好特征，取Top",
					"items": {
						"type": "object",
						"properties": {
							"{operation}": {
								"type": "integer",
								"title": "文档操作以及对应次数",
								"description": "operation为变量，例如 create, edit, delete等"
							}
						}
					}
				},
				"obj_preferences": {
					"type": "array",
					"title": "文档内容偏好特征",
					"description": "文档内容偏好特征，取Top",
					"items": {
						"type": "object",
						"properties": {
							"{obj_id}": {
								"type": "array",
								"title": "文档id",
								"description": "文档id",
								"items": {
									"type": "object",
									"properties": {
										"obj_id": {
											"type": "string",
											"title": "文档id",
											"description": "文档id"
										},
										"obj_name": {
											"type": "string",
											"title": "名称",
											"description": "名称"
										},
										"path": {
											"type": "string",
											"title": "路径",
											"description": "路径"
										},
										"last_op_time": {
											"type": "string",
											"title": "最后一次操作时间",
											"description": "最后一次操作时间"
										},
										"count": {
											"type": "integer",
											"title": "累计操作次数",
											"description": "累计操作次数"
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}`

	KCCenterDynamic = `{
		"{obj_type}": {
			"type": "object",
			"title": "知识中心对象类型",
			"description": "知识中心对象变量，对象类型：wikidoc|topic|circle|file|space|qa等",
			"properties": {
				"operation_preferences": {
					"type": "array",
					"title": "操作偏好特征",
					"description": "操作偏好特征，取Top",
					"items": {
						"type": "object",
						"properties": {
							"{operation}": {
								"type": "integer",
								"title": "操作以及对应次数",
								"description": "operation为变量，例如 create, edit, delete等"
							}
						}
					}
				},
				"obj_preferences": {
					"type": "array",
					"title": "内容偏好特征",
					"description": "内容偏好特征，取Top",
					"items": {
						"type": "object",
						"properties": {
							"{obj_id}": {
								"type": "array",
								"title": "对象id",
								"description": "对象id",
								"items": {
									"type": "object",
									"properties": {
										"obj_id": {
											"type": "string",
											"title": "对象id",
											"description": "对象id"
										},
										"obj_name": {
											"type": "string",
											"title": "对象名称",
											"description": "对象名称"
										},
										"path": {
											"type": "string",
											"title": "对象路径",
											"description": "对象路径"
										},
										"last_op_time": {
											"type": "string",
											"title": "最后一次操作时间",
											"description": "最后一次操作时间"
										},
										"count": {
											"type": "integer",
											"title": "累计操作次数",
											"description": "累计操作次数"
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}`
)
