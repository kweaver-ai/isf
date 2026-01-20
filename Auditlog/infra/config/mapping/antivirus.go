package mapping

func GetAntivirusMapping() string {
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
					"antivirus": {
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
							"detail": {
								"properties": {
									"result": {
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
									"basename": {
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
									"doc_lib_name": {
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
									"extension": {
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
							}
						}
					}
				}
			}
		}
	}`
}
