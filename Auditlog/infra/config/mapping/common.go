package mapping

func innerMappingFields() string {
	return `
    "__rec_log_created_time__": {
		"type": "long"
	},
    "__rec_log_created_time_txt__": {
		"type": "text",
        "index": false
	},
	"__rec_log_ulid__": {
		"type": "keyword"
	}
`
}
