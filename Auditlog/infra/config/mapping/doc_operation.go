package mapping

func GetDocOperationMapping() string {
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
					"doc_operation": {
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
									"sharer": {
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
									"target_object_status": {
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
											}
										}
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
									"dirname": {
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
									"doc": {
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
									"shared_link_id": {
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
									"sharer": {
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
									"size": {
										"type": "long"
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
									"deparment_path": {
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
