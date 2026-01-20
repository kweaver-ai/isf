package oprlogmq

//nolint:unused
func getMockMsgSingleDocumentDomainSync() []byte {
	return []byte(`{
  "config": {
    "max_file_count": -1,
    "max_file_size": 0
  },
  "description": "文档域同步上传源端文件“m1/math.txt”到目标端“10.4.134.83:443/m2/math.txt”",
  "detail": {
    "plan": {
      "id": "4b5c108a-e918-4d6e-9d65-cdc96e159328",
      "source": {
        "id": "gns://B1194865B64646C2A327B1577973341B",
        "type": "user_doc_lib"
      },
      "target": {
        "domain": {
          "id": "40f90d80-0276-4311-bdbf-740469600e65",
          "name": "10.4.134.83:443"
        },
        "library": "m2"
      }
    },
    "result": "200"
  },
  "log_from": {
    "package": "AnyShareMainModule",
    "service": {
      "name": "document-sync-scheduling"
    }
  },
  "object": {
    "id": "gns://B1194865B64646C2A327B1577973341B/91337341EFDD45968F3BE8FE8BE0869E",
    "name": "math.txt",
    "path": "m1/math.txt",
    "size": 105,
    "type": "file"
  },
  "operation": "document_domain_upload",
  "operator": {
    "id": "ff73e81a-c768-4fce-bd15-3084127af924",
    "name": "document-sync-scheduling",
    "type": "internal_service"
  },
  "recorder": "Anyshare",
  "target_object": {
    "document_domain": "10.4.134.83:443",
    "path": "m2/math.txt",
    "size": 105,
    "type": "file"
  }
}
`)
}

//nolint:unused
func getMockMsgSingleDirVisit() []byte {
	return []byte(`{
  "operation": "cd",
  "recorder": "Anyshare",
  "description": "用户“张三”从“部门文档库1/a”进入到“部门文档库1/a/b”。",
  "log_from": {
    "package": "AnyShareMainModule",
    "service": {
      "instance": {
        "id": "docset-6475f48ff6-cdbkf"
      },
      "version": "0.0.0-20240913153111-e71d017e",
      "name": "docset"
    }
  },
  "operator": {
    "id": "f49e9bee-167b-11ef-a1fa-3e35e19e5cab",
    "name": "李宇（Aaron）",
    "type": "authenticated_user",
    "department_path": [
      {
        "id_path": "4e8bfbda-d99c-11eb-35b9-24e8e0506805/4bfdae8b-d9c9-1eb1-5b39-5068024e8e05/e8bfbda4-d31c-12ab-34c9-50680524e8e0",
        "name_path": "爱数/数据智能产品BG/AnyShare研发线/智能搜索研发部"
      }
    ],
    "agent": {
      "type": "windows",
      "os_type": "windows",
      "app_type": "sync_disk",
      "ip": "192.168.50.100",
      "udid": "",
      "user_agent": ""
    }
  },
  "object": {
    "id": "gns://D42F2729C56E489A948985D4E75C5813/4e8bfbda-d99c-11eb-35b9-24e8e050xxx5",
    "path": "部门文档库1/a",
    "name": "a",
    "size": 0,
    "type": "folder",
    "doc_lib": {
      "id": "gns://D42F2729C56E489A948985D4E75C5813",
      "type": "department_doc_lib",
      "name": "mock_doc_lib_name"
    }
  },
  "rec": {
    "ext_info_json": "{\"k1\":{\"k2\":\"v2\"}}",
    "not_use_for_rec": false
  },
  "detail": {
    "from_object": {
      "id": "gns://D42F2729C56E489A948985D4E75C4813",
      "path": "部门文档库1/a",
      "type": "normal_dir",
      "doc_lib": {
        "id": "gns://D42F2729C56E489A948985D4E75C4813",
        "type": "department_doc_lib",
        "name": "mock_doc_lib_name"
      }
    },
    "obj1": {
      "key1": "val1",
      "key2": "val2"
    },
    "key3": "val3"
  },
  "referer": {
    "current": "preview",
    "previous": "dir"
  }
}
`)
}

//nolint:unused
func getMockMsgSingleInvalidButtonClick() []byte {
	return []byte(`
    {
        "description": "收藏IMG_20210310_132259.jpg",
        "detail": {
            "object": "{\"name\":\"IMG_20210310_132259.jpg\",\"operation_position\":null}",
            "agent": "{\"app_type\":\"android_app\",\"background\":0,\"client_version\":\"7.0.6.1\",\"client_version_code\":1783,\"device_model\":\"TAS-AL00\",\"device_name\":\"HUAWEI\",\"device_system_version\":\"10\",\"host\":\"10.4.134.83\",\"ip\":\"\",\"network_type\":\"NETWORK_WIFI\",\"os_type\":\"android\",\"release\":0,\"server_version\":\"7.0.6.1.20241107\",\"type\":\"android\",\"udid\":\"\",\"user_agent\":\"cn.aishu.anyshare/7.0.6.1 (build:1783;versionCode:147;HUAWEI/TAS-AL00;isAndroid:true;10)\"}"
        },
        "object": {
            "doc_lib": {
                "id": "gns://0FE0A88137294709B4A5C7B18E7BDDAD",
                "name": "",
                "type": "user_doc_lib"
            },
            "id": "gns://0FE0A88137294709B4A5C7B18E7BDDAD/EA00BD819DA84BB8B522D72C0582102C",
            "name": "IMG_20210310_132259.jpg",
            "path": "",
            "size": 1458598,
            "tags": [
                
            ],
            "type": "file"
        },
        "operation": "starred",
        "operator": {
            "agent": {
                "app_type": "android_app",
                "ip": "",
                "os_type": "android",
                "type": "android",
                "udid": "7B36D80778A53EBDB9FFB2A275029DC31B6A889B",
                "user_agent": "cn.aishu.anyshare/7.0.6.1 (build:1783;versionCode:147;HUAWEI/TAS-AL00;isAndroid:true;10)"
            },
            "department_path": [
                {
                    "id_path": "151bcb65-48ce-4b62-973f-0bb6685f9cb8",
                    "name_path": "组织结构"
                }
            ],
            "id": "e0d6a6a0-a15a-11ef-a5e7-0e2760d4301d",
            "is_system_op": false,
            "name": "app1",
            "type": "authenticated_user"
        },
        "rec": {
            "ext_info_json": "",
            "not_use_for_rec": false
        },
        "recorder": "AnyShare",
        "referer": null,
        "targetObject": null
    }
`)
}

//nolint:unused
func getMockMsgManyMenuButtonClick() []byte {
	return []byte(`[
  {
    "operation": "rename",
    "recorder": "Anyshare",
    "description": "用户“李宇（Aaron）”将文档“年度报告.pdf”重命名为“2023年度总结.pdf”。",
    "operator": {
      "id": "8b085b72-567c-11ed-aecc-063c8a32c7bf",
      "name": "李宇（Aaron）",
      "type": "authenticated_user",
      "is_system_op": false,
      "department_path": [
        {
          "id_path": "4e8bfbda-d99c-11eb-35b9-24e8e0506805/4bfdae8b-d9c9-1eb1-5b39-5068024e8e05/e8bfbda4-d31c-12ab-34c9-50680524e8e0",
          "name_path": "爱数/数据智能产品BG/AnyShare研发线/智能搜索研发部"
        }
      ],
      "agent": {
        "type": "web",
        "os_type": "windows",
        "app_type": "app",
        "ip": "192.168.50.100"
      }
    },
    "object": {
      "id": "gns://D42F2729C56E489A948985D4E75C5813/4e8bfbda-d99c-11eb-35b9-24e8e050xxx5",
      "path": "/documents/年度报告.pdf",
      "name": "年度报告.pdf",
      "size": 2048,
      "type": "file"
    },
    "log_from": {
      "package": "AnyShareMainModule",
      "service": {
        "instance": {
          "id": "docset-6475f48ff6-cdbkf"
        },
        "version": "0.0.0-20240913153111-e71d017e",
        "name": "docset"
      }
    },
    "rec": {
      "not_use_for_rec": false,
      "ext_info_json": "{\"k1\":{\"k2\":\"v2\"}}"
    },
    "detail": {
      "op_name": "重命名",
      "op_type": "op_configed_built-in",
      "position": {
        "type": "dir",
        "path_id": "gns://D42F2729C56E489A948985D4E75C4813",
        "path_name": "部门文档库1/a",
        "area_type": 2
      }
    },
    "referer": {
      "current": "xx",
      "previous": "xx"
    }
  },
  {
    "operation": "rename",
    "recorder": "Anyshare",
    "description": "用户“李宇（Aaron）”将文档“年度报告2.pdf”重命名为“2023年度总结.pdf”。",
    "operator": {
      "id": "8b085b72-567c-11ed-aecc-063c8a32c7bf",
      "name": "李宇（Aaron）",
      "type": "authenticated_user",
      "is_system_op": false,
      "department_path": [
        {
          "id_path": "4e8bfbda-d99c-11eb-35b9-24e8e0506805/4bfdae8b-d9c9-1eb1-5b39-5068024e8e05/e8bfbda4-d31c-12ab-34c9-50680524e8e0",
          "name_path": "爱数/数据智能产品BG/AnyShare研发线/智能搜索研发部"
        }
      ],
      "agent": {
        "type": "web",
        "os_type": "windows",
        "app_type": "app",
        "ip": "192.168.50.100"
      }
    },
    "object": {
      "id": "gns://D42F2729C56E489A948985D4E75C5813/4e8bfbda-d99c-11eb-35b9-24e8e050xxx6",
      "path": "/documents/年度报告2.pdf",
      "name": "年度报告2.pdf",
      "size": 2048,
      "type": "file"
    },
    "log_from": {
      "package": "AnyShareMainModule",
      "service": {
        "instance": {
          "id": "docset-6475f48ff6-cdbkf"
        },
        "version": "0.0.0-20240913153111-e71d017e",
        "name": "docset"
      }
    },
    "rec": {
      "not_use_for_rec": false,
      "ext_info_json": "{\"k1\":{\"k2\":\"v2\"}}"
    },
    "detail": {
      "op_name": "重命名",
      "op_type": "op_configed_built-in",
      "position": {
        "type": "dir",
        "path_id": "gns://D42F2729C56E489A948985D4E75C4813",
        "path_name": "部门文档库1/a",
        "area_type": 2
      }
    },
    "referer": {
      "current": "xx",
      "previous": "xx"
    }
  },
  {
    "operation": "rename",
    "recorder": "Anyshare",
    "description": "用户“李宇（Aaron）”将文档“年度报告2.pdf”重命名为“2023年度总结.pdf”。",
    "operator": {
      "id": "8b085b72-567c-11ed-aecc-063c8a32c7bf",
      "name": "李宇（Aaron）",
      "type": "authenticated_user",
      "is_system_op": false,
      "department_path": [
        {
          "id_path": "4e8bfbda-d99c-11eb-35b9-24e8e0506805/4bfdae8b-d9c9-1eb1-5b39-5068024e8e05/e8bfbda4-d31c-12ab-34c9-50680524e8e0",
          "name_path": "爱数/数据智能产品BG/AnyShare研发线/智能搜索研发部"
        }
      ],
      "agent": {
        "type": "web",
        "os_type": "windows",
        "app_type": "app",
        "ip": "192.168.50.100"
      }
    },
    "object": {
      "id": "gns://D42F2729C56E489A948985D4E75C5813/4e8bfbda-d99c-11eb-35b9-24e8e050xxx6",
      "path": "/documents/年度报告2.pdf",
      "name": "年度报告2.pdf",
      "size": 2048,
      "type": "file"
    },
    "log_from": {
      "package": "AnyShareMainModule",
      "service": {
        "instance": {
          "id": "docset-6475f48ff6-cdbkf"
        },
        "version": "0.0.0-20240913153111-e71d017e",
        "name": "docset"
      }
    },
    "rec": {
      "not_use_for_rec": false,
      "ext_info_json": "{\"k1\":{\"k2\":\"v2\"}}"
    },
    "detail": {
      "op_name": "重命名",
      "op_type": "op_configed_built-in",
      "position": {
        "type": "dir",
        "path_id": "gns://D42F2729C56E489A948985D4E75C4813",
        "path_name": "部门文档库1/a",
        "area_type": 2
      }
    },
    "referer": {
      "current": "xx",
      "previous": "xx"
    }
  }
]
`)
}

func getMockMsgDocOperation() []byte {
	return []byte(`{
    "detail": {
        "target_object_status": ""
    },
    "operation": "create",
    "biz_type": "doc_operation",
    "description": "用户“于涛（Tanya）”新建文件“财务管理线/财务部/09资金类/5、对账单电子档/2024年对账单/厦门爱数/农行5933/对账单/账户明细查询列表 (11).pdf”。",
    "object": {
        "path": "财务管理线/财务部/09资金类/5、对账单电子档/2024年对账单/厦门爱数/农行5933/对账单/账户明细查询列表 (11).pdf",
        "size": 101175,
        "type": "file",
        "id": "gns://F8743A83A660432EBA316367DE5B3C40/13F422B002704D4B9A26C4CADDBDD7F7/B5C59B6A6E6B4C63B25AA9CBC7B072C9/0CB45ABE1FC446CC924E0C4AC336CA8B/B053C43FC2604D4C862B4F60F8099FEC/B5704AC814B74E5CBA65212035AD845E/0846B6D4C179473D9C22F3235A367209/8B2E20CEF2DC4137A9E9EC99E7371FD9/8F71CC66784D4211B1F79DE1843EE038",
        "doc_lib_name": "财务管理线",
        "name": "账户明细查询列表 (11).pdf",
        "basename": "账户明细查询列表 (11).pdf",
        "doc_lib": {
            "id": "gns://F8743A83A660432EBA316367DE5B3C40",
            "type": "department_doc_lib"
        },
        "extension": "pdf"
    },
    "recorder": "AnyShare",
    "log_from": {
        "service": {
            "name": "efast"
        },
        "package": "IdentifyAndAuthentication"
    },
    "operator": {
        "agent": {
            "ip": "58.247.3.114",
            "type": "windows",
            "user_agent": "AnyShare-WinClient",
            "udid": "2C-33-58-39-22-DB"
        },
        "department_path": [
            {
                "id_path": "e9af4dba-5e16-11e3-b9fe-dcd2fc061e41/e0fd1826-bf84-11e3-a95e-dcd2fc061e41/e10c2dde-bf84-11e3-a95e-dcd2fc061e41/e137d754-bf84-11e3-a95e-dcd2fc061e41/e1711b54-bf84-11e3-a95e-dcd2fc061e41",
                "name_path": "上海爱数信息技术股份有限公司/aishu/财务管理线/财务部/资金组"
            }
        ],
        "name": "于涛（Tanya）",
        "type": "authenticated_user",
        "id": "24a87614-3749-11ef-8d90-5a7df6eefa12"
    }
}`)
}

func getMockMsgKcOperation() []byte {
	return []byte(`{
    "ar_types": [
        "kc-statistics"
    ],
    "detail": null,
    "event_name": "",
    "event_user": {
        "display_name": "李宇（Aaron）",
        "is_expert": 0,
        "department_ids": [
            "7a3f97c0-b706-11ed-9a71-56a992d5db74"
        ],
        "address": "310000",
        "is_knowledge": 0,
        "as_id": "8b085b72-567c-11ed-aecc-063c8a32c7bf"
    },
    "operation": "create",
    "recorder": "KnowledgeCenter",
    "log_from": {
        "service": {
            "name": "kc-mc",
            "version": "",
            "instance": {
                "id": "kc-mc-fb7dc7c8f-tfwwf"
            }
        },
        "package": "KnowledgeCenter"
    },
    "object_id": "292371997501505537",
    "operator": {
        "agent": null,
        "department_path": [
            {
                "id_path": "e9af4dba-5e16-11e3-b9fe-dcd2fc061e41/e0fd1826-bf84-11e3-a95e-dcd2fc061e41/4a060486-e050-11ed-a5dc-363580fb9098/eec14dc8-9c7a-11e9-b135-dcd2fc061e41/0fd2c5b2-f8ed-11e5-8a84-d8490bc791a2/7a3f97c0-b706-11ed-9a71-56a992d5db74",
                "name_path": "上海爱数信息技术股份有限公司/aishu/数据智能产品BG/AnyShare研发线/应用场景研发部/SAP数据资产管理后端开发组"
            }
        ],
        "name": "李宇（Aaron）",
        "type": "authenticated_user",
        "id": "8b085b72-567c-11ed-aecc-063c8a32c7bf"
    },
    "others": {
        "space_type": "2",
        "space_name": "李宇（Aaron）",
        "user_device": "nonmobile",
        "wikidoc_source": 0,
        "space_id": "223030506168448032",
        "wikidoc_suffix": ""
    },
    "object_name": "aa",
    "biz_type": "kc_operation",
    "description": "用户“李宇（Aaron）”&{create 创建}。",
    "object_type": "space_article",
    "create_time": 1736146351,
    "action_type": "create"
}`)
}

func getMockMsgClientOperation() []byte {
	return []byte(`{
    "detail": {
        "url": "https://10.4.134.243/anyshare/zh-cn/microapps/fun_250108xext1g65dmnfnkd50/knowledge-center/wikidoc/space?article_id=293698625900462082&space_id=293698625900462081&status=detail",
        "title": "222_AnyShare",
        "language": "zh-CN",
        "page_name": "space",
        "event_name": "点击导航栏入口[首页]",
        "operation_position": "navigation"
    },
    "referer": {
        "previous": "",
        "current": "knowledge-center"
    },
    "operation": "navigation_home",
    "recorder": "Anyshare",
    "description": "用户niko点击了导航栏入口[首页]",
    "operator": {
        "type": "authenticated_user",
        "department_path": [
            {
                "id_path": "151bcb65-48ce-4b62-973f-0bb6685f9cb8",
                "name_path": "组织结构"
            },
            {
                "id_path": "151bcb65-48ce-4b62-973f-0bb6685f9cb8/e8aff8de-c064-11f0-8e0e-52116f620329/ee67e110-c064-11f0-b10a-52116f620329",
                "name_path": "组织结构/测试部门/子部门"
            },
            {
                "id_path": "151bcb65-48ce-4b62-973f-0bb6685f9cb8/c9291990-cfd9-11ef-bbe6-52116f620329",
                "name_path": "组织结构/研发部"
            }
        ],
        "agent": {
            "release": 1,
            "sdk_version": "0.1.0",
            "browser_type": "Chrome",
            "type": "web_portal",
            "userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
            "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
            "device_system_version": "10",
            "client_version": "7.0.6.2",
            "ip": "",
            "os_type": "windows",
            "host": "10.4.134.243",
            "browser_engine": "Blink",
            "device_name": "Windows",
            "browser_version": "131.0.0.0",
            "server_version": "7.0.6.2.20241218",
            "app_type": "web",
            "udid": ""
        },
        "name": "niko",
        "is_system_op": false,
        "id": "c40542b2-cf3a-11ef-ae53-52116f620329"
    },
    "biz_type": "client_operation"
}`)
}
