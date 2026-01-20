package mapping

func GetClientOperationMapping() string {
	return `{
		"properties": {
			"@timestamp": {
				"type": "date"
			},
			"Body": {
				"properties": {
					"Type": {
						"type": "text",
						"fields": {
							"keyword": {
								"type": "keyword",
								"ignore_above": 256
							}
						}
					},
					"client_operation": {
						"properties": {
							"biz_type": {
								"type": "text",
								"fields": {
									"keyword": {
										"type": "keyword",
										"ignore_above": 1024
									}
								},
								"norms": false,
								"analyzer": "ik_max_word"
							},
							"description": {
								"type": "text",
								"fields": {
									"keyword": {
										"type": "keyword",
										"ignore_above": 1024
									}
								},
								"norms": false,
								"analyzer": "ik_max_word"
							},
							"detail": {
								"properties": {
									"button_source": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"child_event_id": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"click_info": {
										"properties": {
											"active_position": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"child_component_name": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											}
										}
									},
									"client_microwidget_id": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"client_microwidget_name": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"client_microwidget_number": {
										"type": "long"
									},
									"client_operate": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"dir_child_description": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"dir_contentType": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"dir_doc_contentType": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"dir_location": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"dir_naviTreeType": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"doc_name": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"doc_path": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"event_name": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"event_type": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"fcp": {
										"type": "long"
									},
									"fp": {
										"type": "long"
									},
									"language": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"operate_index": {
										"type": "long"
									},
									"operation_position": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"operation_type": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"page_name": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"state": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"title": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"url": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									}
								}
							},
							"object": {
								"properties": {
									"contentType": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"docid": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"icon": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"id": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"name": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"open_mode": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"operation_position": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"path": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"size": {
										"type": "long"
									},
									"type": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									}
								}
							},
							"operation": {
								"type": "text",
								"fields": {
									"keyword": {
										"type": "keyword",
										"ignore_above": 1024
									}
								},
								"norms": false,
								"analyzer": "ik_max_word"
							},
							"operator": {
								"properties": {
									"agent": {
										"properties": {
											"app_type": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"browser_engine": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"browser_type": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"browser_version": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"client_version": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"device_name": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"device_system_version": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"host": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"ip": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"os_type": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"release": {
												"type": "long"
											},
											"sdk_version": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"server_version": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"type": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"udid": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"userAgent": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"user_agent": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											}
										}
									},
									"department_path": {
										"properties": {
											"id_path": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											},
											"name_path": {
												"type": "text",
												"fields": {
													"keyword": {
														"type": "keyword",
														"ignore_above": 1024
													}
												},
												"norms": false,
												"analyzer": "ik_max_word"
											}
										}
									},
									"id": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"is_system_op": {
										"type": "boolean"
									},
									"name": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"type": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									}
								}
							},
							"recorder": {
								"type": "text",
								"fields": {
									"keyword": {
										"type": "keyword",
										"ignore_above": 1024
									}
								},
								"norms": false,
								"analyzer": "ik_max_word"
							},
							"referer": {
								"properties": {
									"current": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									},
									"previous": {
										"type": "text",
										"fields": {
											"keyword": {
												"type": "keyword",
												"ignore_above": 1024
											}
										},
										"norms": false,
										"analyzer": "ik_max_word"
									}
								}
							}
						}
					}
				}
			}
		}
	}`
}
