package flow

type FlowRequestProcessDefinedInputs struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Type     string `json:"type"`
	MetaType string `json:"metaType"`
}

type FlowRequestProcess struct {
	Name         string                            `json:"name"`
	DefinedInput []FlowRequestProcessDefinedInputs `json:"defined_input"`
}

type FlowRequest struct {
	Name                 string               `json:"name"`
	FlowRequestProcesses []FlowRequestProcess `json:"process_definitions"`
}
