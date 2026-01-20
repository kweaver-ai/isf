package uniqreclogeo

import (
	"AuditLog/common/utils"
)

type InputOutput struct {
	Input  interface{} `json:"input"`
	Output interface{} `json:"output"`
}

type Detail struct {
	OutCall     *InputOutput `json:"out_call,omitempty"`      // 外部调用
	PreAPICall  *InputOutput `json:"pre_api_call,omitempty"`  // 前置API调用
	LLMCall     *InputOutput `json:"llm_call,omitempty"`      // LLM调用
	PostAPICall *InputOutput `json:"post_api_call,omitempty"` // 后置API调用
}

func NewDetail() *Detail {
	return &Detail{
		OutCall:     &InputOutput{},
		PreAPICall:  &InputOutput{},
		LLMCall:     &InputOutput{},
		PostAPICall: &InputOutput{},
	}
}

func (d *Detail) LoadByInterface(i interface{}) (err error) {
	if i == nil {
		return
	}

	//    通过json来实现
	jsonStr, err := utils.JSON().Marshal(i)
	if err != nil {
		return
	}

	err = utils.JSON().Unmarshal(jsonStr, d)
	if err != nil {
		return
	}

	return
}
