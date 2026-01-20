package mapping

func GetSapMapping() string {
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
					"sap": {
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
									"archobj_tables": {
										"properties": {
											"OBJECT": {
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
											"OBJTEXT": {
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
											"SEQ_NO": {
												"type": "long"
											},
											"TABNAME": {
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
											"TABTEXT": {
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
									"data_volume": {
										"properties": {
											"DataVolumeKey": {
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
											"DataVolumeValue": {
												"type": "float"
											},
											"File": {
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
											"Size": {
												"type": "float"
											}
										}
									},
									"db_history": {
										"properties": {
											"DataUsed": {
												"type": "float"
											},
											"Date": {
												"type": "date",
												"format": "yyyy-MM-dd HH:mm:ss||strict_date_optional_time"
											},
											"LogUsed": {
												"type": "float"
											},
											"MemoryUsed": {
												"type": "float"
											},
											"Month": {
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
											"TraceUsed": {
												"type": "float"
											},
											"Year": {
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
									"db_memory": {
										"properties": {
											"DbMemoryKey": {
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
											"DbMemoryValue": {
												"type": "float"
											},
											"ProcessMemory": {
												"type": "float"
											},
											"ServiceName": {
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
											"TotalMemory": {
												"type": "float"
											}
										}
									},
									"file_size_change": {
										"type": "long"
									},
									"growing_table": {
										"properties": {
											"MemoryDelta": {
												"type": "float"
											},
											"MemoryTotal": {
												"type": "float"
											},
											"Object": {
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
											"ObjectText": {
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
											"RecordCount": {
												"type": "integer"
											},
											"SquenceNo": {
												"type": "integer"
											},
											"TabName": {
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
											"TableName": {
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
									"largest_table": {
										"properties": {
											"MemoryDelta": {
												"type": "float"
											},
											"MemoryTotal": {
												"type": "float"
											},
											"Object": {
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
											"ObjectText": {
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
											"RecordCount": {
												"type": "integer"
											},
											"SquenceNo": {
												"type": "integer"
											},
											"TabName": {
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
											"TableName": {
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
									"log_volume": {
										"properties": {
											"File": {
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
											"LogVolumeKey": {
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
											"LogVolumeValue": {
												"type": "float"
											},
											"Size": {
												"type": "float"
											}
										}
									},
									"position": {
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
									"sence_data_info": {
										"properties": {
											"edgp_service_name": {
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
									"trace_volume": {
										"properties": {
											"File": {
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
											"Size": {
												"type": "long"
											},
											"TraceKey": {
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
											"TraceValue": {
												"type": "long"
											}
										}
									}
								}
							},
							"log_from": {
								"properties": {
									"package": {
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
									"service": {
										"properties": {
											"instance": {
												"properties": {
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
													}
												}
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
											"version": {
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
							},
							"object": {
								"properties": {
									"doc_lib": {
										"properties": {
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
							"target_object": {
								"properties": {
									"doc_lib": {
										"properties": {
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
							}
						}
					}
				}
			}
		}
	}`
}
