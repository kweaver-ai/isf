package mapping

func GetFileCollectorMapping() string {
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
					"file_collector": {
						"properties": {
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
							"object": {
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
									"submit_by": {
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
