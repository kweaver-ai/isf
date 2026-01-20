package mapping

import "fmt"

func GetDirVisitCdMapping() string {
	return fmt.Sprintf(`{
  "properties": {
	%s,
    "description": {
      "type": "text"
    },
	"biz_type": {
      "type": "keyword"
    },
    "detail": {
      "type": "object",
      "properties": {
        "from_object": {
          "type": "object",
          "properties": {
            "doc_lib": {
              "type": "object",
              "properties": {
                "custom_doc_lib_sub_type": {
                  "type": "keyword"
                },
                "id": {
                  "type": "keyword",
                  "fields": {
                    "text": {
                      "type": "text"
                    }
                  }
                },
                "name": {
                  "type": "text"
                },
                "type": {
                  "type": "keyword"
                }
              }
            },
            "id": {
              "type": "keyword",
              "fields": {
                "text": {
                  "type": "text"
                }
              }
            },
            "path": {
              "type": "text"
            },
            "type": {
              "type": "keyword"
            }
          }
        },
        "fav_category": {
          "type": "object",
          "properties": {
            "id": {
              "type": "keyword"
            },
            "name": {
              "type": "text"
            },
            "id_path": {
              "type": "keyword",
              "fields": {
                "text": {
                  "type": "text"
                }
              }
            },
            "name_path": {
              "type": "text"
            }
          }
        }
      }
    },
    "log_from": {
      "type": "object",
      "properties": {
        "package": {
          "type": "keyword"
        },
        "service": {
          "type": "object",
          "properties": {
            "instance": {
              "type": "object",
              "properties": {
                "id": {
                  "type": "keyword"
                }
              }
            },
            "name": {
              "type": "keyword"
            },
            "version": {
              "type": "text",
              "fields": {
                "keyword": {
                  "type": "keyword",
                  "ignore_above": 256
                }
              }
            }
          }
        }
      }
    },
    "object": {
      "type": "object",
      "properties": {
        "doc_lib": {
          "type": "object",
          "properties": {
            "custom_doc_lib_sub_type": {
              "type": "keyword"
            },
            "id": {
              "type": "keyword"
            },
            "name": {
              "type": "text"
            },
            "type": {
              "type": "keyword"
            }
          }
        },
        "id": {
          "type": "keyword",
          "fields": {
            "text": {
              "type": "text"
            }
          }
        },
        "name": {
          "type": "text"
        },
        "path": {
          "type": "text"
        },
        "size": {
          "type": "long"
        },
        "type": {
          "type": "keyword"
        }
      }
    },
    "operation": {
      "type": "keyword"
    },
    "operator": {
      "type": "object",
      "properties": {
        "agent": {
          "type": "object",
          "properties": {
            "app_type": {
              "type": "keyword"
            },
            "ip": {
              "type": "keyword"
            },
            "os_type": {
              "type": "keyword"
            },
            "type": {
              "type": "keyword"
            },
            "udid": {
              "type": "keyword"
            },
            "user_agent": {
              "type": "text"
            }
          }
        },
        "department_path": {
          "type": "nested",
          "properties": {
            "id_path": {
              "type": "keyword",
              "fields": {
                "text": {
                  "type": "text"
                }
              }
            },
            "name_path": {
              "type": "text"
            }
          }
        },
        "id": {
          "type": "keyword"
        },
        "is_system_op": {
          "type": "boolean"
        },
        "name": {
          "type": "text"
        },
        "type": {
          "type": "keyword"
        }
      }
    },
    "recorder": {
      "type": "keyword"
    },
    "referer": {
      "type": "object",
      "properties": {
        "current": {
          "type": "keyword",
          "fields": {
            "text": {
              "type": "text"
            }
          }
        },
        "previous": {
          "type": "keyword",
          "fields": {
            "text": {
              "type": "text"
            }
          }
        }
      }
    }
  }
}`, innerMappingFields())
}
