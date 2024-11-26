package flow_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/e4coder/flow"
)

const FLOW_REQUEST = `{
    "name": "Stake",
    "process_definitions": [
        {
            "name": "prepareCallData",
            "defined_input": [
                {
                    "name": "stake_amount",
                    "value": "100000",
                    "type": "string",
                    "metaType": "bigint"
                }
            ]
        },
        {
            "name": "processUserOp",
            "defined_input": [
            ]
        }
    ]
}`

func TestMain(t *testing.T) {
	handlers := map[string]flow.ProcessHandler{}
	handlers["prepareCallData"] = prepareCallData
	handlers["processUserOp"] = processUserOp

	flowParser := flow.NewFlowParser(handlers)
	stakeSchema := _stake_schema()
	flowParser.Add(stakeSchema.Name, stakeSchema)

	flowRequest := flow.FlowRequest{}
	err := json.Unmarshal([]byte(FLOW_REQUEST), &flowRequest)
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}
	parsedFlow, err := flowParser.Parse(flowRequest)
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}
	c := &flow.ProcessContext{}
	c.Vals = make(map[interface{}]interface{})

	parsedFlow.Process(c)

	val, ok := c.Vals["response"].(map[string]string)
	if !ok {
		t.Logf("%v", err)
		t.Fail()
	}

	data := val["data"]
	if strings.Compare(data, "to mars") != 0 {
		t.Log("invalid response")
		t.Fail()
	}
}

func TestProcessError(t *testing.T) {
	handlers := map[string]flow.ProcessHandler{}
	handlers["prepareCallData"] = signUserOp
	handlers["processUserOp"] = processUserOp

	flowParser := flow.NewFlowParser(handlers)
	stakeSchema := _stake_schema()
	flowParser.Add(stakeSchema.Name, stakeSchema)

	flowRequest := flow.FlowRequest{}
	err := json.Unmarshal([]byte(FLOW_REQUEST), &flowRequest)
	if err != nil {
		t.Logf("%v", err)
		t.FailNow()
	}
	parsedFlow, err := flowParser.Parse(flowRequest)
	if err != nil {
		t.Logf("%v", err)
		t.FailNow()
	}
	c := &flow.ProcessContext{}
	c.Vals = make(map[interface{}]interface{})

	err = parsedFlow.Process(c)
	if !errors.Is(err, flow.ErrProcessFailure) {
		t.Logf("%v\n", err)
		t.FailNow()
	}
}

func TestFlowNotFound(t *testing.T) {
	handlers := map[string]flow.ProcessHandler{}
	handlers["prepareCallData"] = prepareCallData
	handlers["processUserOp"] = processUserOp

	flowParser := flow.NewFlowParser(handlers)
	stakeSchema := _unstake_schema()
	flowParser.Add(stakeSchema.Name, stakeSchema)

	flowRequest := flow.FlowRequest{}
	err := json.Unmarshal([]byte(FLOW_REQUEST), &flowRequest)
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}
	_, err = flowParser.Parse(flowRequest)
	if !errors.Is(err, flow.ErrSchemaVerificationFailure) {
		t.Logf("%v", err)
		t.Fail()
	}

	if !errors.Is(err, flow.ErrFlowNotFound) {
		t.Logf("%v", err)
		t.Fail()
	}
}

func TestInvalidRequestProcessLen(t *testing.T) {
	handlers := map[string]flow.ProcessHandler{}
	handlers["prepareCallData"] = prepareCallData
	handlers["processUserOp"] = processUserOp

	flowParser := flow.NewFlowParser(handlers)
	stakeSchema := _stake_prepare_only()
	flowParser.Add(stakeSchema.Name, stakeSchema)

	flowRequest := flow.FlowRequest{}
	err := json.Unmarshal([]byte(FLOW_REQUEST), &flowRequest)
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}
	_, err = flowParser.Parse(flowRequest)
	if !errors.Is(err, flow.ErrSchemaVerificationFailure) {
		t.Logf("%v", err)
		t.Fail()
	}

	if !errors.Is(err, flow.ErrInvalidRequest) {
		t.Logf("%v", err)
		t.Fail()
	}
}

func TestInvalidRequestProcess(t *testing.T) {
	handlers := map[string]flow.ProcessHandler{}
	handlers["prepareCallData"] = prepareCallData
	handlers["processUserOp"] = processUserOp

	flowParser := flow.NewFlowParser(handlers)
	stakeSchema := _stake_prepare_only()
	flowParser.Add(stakeSchema.Name, stakeSchema)

	flowRequest := flow.FlowRequest{}
	err := json.Unmarshal([]byte(FLOW_REQUEST), &flowRequest)
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}
	_, err = flowParser.Parse(flowRequest)
	if !errors.Is(err, flow.ErrSchemaVerificationFailure) {
		t.Logf("%v", err)
		t.Fail()
	}

	if !errors.Is(err, flow.ErrInvalidRequest) {
		t.Logf("%v", err)
		t.Fail()
	}
}

func TestInvalidRequestProcessName(t *testing.T) {
	handlers := map[string]flow.ProcessHandler{}
	handlers["prepareCallData"] = prepareCallData
	handlers["processUserOp"] = processUserOp

	flowParser := flow.NewFlowParser(handlers)
	stakeSchema := _stake_sign_process()
	flowParser.Add(stakeSchema.Name, stakeSchema)

	flowRequest := flow.FlowRequest{}
	err := json.Unmarshal([]byte(FLOW_REQUEST), &flowRequest)
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}
	_, err = flowParser.Parse(flowRequest)
	if !errors.Is(err, flow.ErrSchemaVerificationFailure) {
		t.Logf("%v", err)
		t.Fail()
	}

	if !errors.Is(err, flow.ErrInvalidRequestProcess) {
		t.Logf("%v", err)
		t.Fail()
	}
}

func TestInvalidRequestProcessInputs(t *testing.T) {
	handlers := map[string]flow.ProcessHandler{}
	handlers["prepareCallData"] = prepareCallData
	handlers["processUserOp"] = processUserOp

	flowParser := flow.NewFlowParser(handlers)
	stakeSchema := _stake_no_data()
	flowParser.Add(stakeSchema.Name, stakeSchema)

	flowRequest := flow.FlowRequest{}
	err := json.Unmarshal([]byte(FLOW_REQUEST), &flowRequest)
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}
	_, err = flowParser.Parse(flowRequest)
	if !errors.Is(err, flow.ErrSchemaVerificationFailure) {
		t.Logf("%v", err)
		t.Fail()
	}

	if !errors.Is(err, flow.ErrInvalidRequestInputs) {
		t.Logf("%v", err)
		t.Fail()
	}
}

func TestHandlerNotFound(t *testing.T) {
	handlers := map[string]flow.ProcessHandler{}
	handlers["processUserOp"] = processUserOp

	flowParser := flow.NewFlowParser(handlers)
	stakeSchema := _stake_schema()
	flowParser.Add(stakeSchema.Name, stakeSchema)

	flowRequest := flow.FlowRequest{}
	err := json.Unmarshal([]byte(FLOW_REQUEST), &flowRequest)
	if err != nil {
		t.Logf("%v", err)
		t.Fail()
	}
	_, err = flowParser.Parse(flowRequest)
	if !errors.Is(err, flow.ErrParserFailure) {
		t.Logf("%v", err)
		t.Fail()
	}

	if !errors.Is(err, flow.ErrHandlerNotFound) {
		t.Logf("%v", err)
		t.Fail()
	}
}

func prepareCallData(c *flow.ProcessContext, inputs []flow.DefinedInput) error {
	c.Vals["calldata"] = "to mars"

	return nil
}

func processUserOp(c *flow.ProcessContext, input []flow.DefinedInput) error {
	val, ok := c.Vals["calldata"].(string)
	if !ok {
		return fmt.Errorf("invalid calldata")
	}

	response := map[string]string{}
	response["data"] = val

	c.Vals["response"] = response
	return nil
}

func signUserOp(_ *flow.ProcessContext, _ []flow.DefinedInput) error {
	return fmt.Errorf("failed to sign userOp")
}

func _stake_schema() flow.FlowSchema {
	return flow.FlowSchema{
		Name: "Stake",
		ProcessDefinitions: []flow.FlowProcessSchema{
			{
				Name: "prepareCallData",
				DefinedInputs: []flow.FlowProcessDefinedInputsSchema{
					{
						Name: "stake_amount",
						Type: "string",
						Meta: "bigint",
					},
				},
			},
			{
				Name:          "processUserOp",
				DefinedInputs: []flow.FlowProcessDefinedInputsSchema{},
			},
		},
	}
}

func _unstake_schema() flow.FlowSchema {
	return flow.FlowSchema{
		Name: "Unstake",
		ProcessDefinitions: []flow.FlowProcessSchema{
			{
				Name: "prepareCallData",
				DefinedInputs: []flow.FlowProcessDefinedInputsSchema{
					{
						Name: "stake_amount",
						Type: "string",
						Meta: "bigint",
					},
				},
			},
			{
				Name:          "processUserOp",
				DefinedInputs: []flow.FlowProcessDefinedInputsSchema{},
			},
		},
	}
}

func _stake_prepare_only() flow.FlowSchema {
	return flow.FlowSchema{
		Name: "Stake",
		ProcessDefinitions: []flow.FlowProcessSchema{
			{
				Name: "prepareCallData",
				DefinedInputs: []flow.FlowProcessDefinedInputsSchema{
					{
						Name: "stake_amount",
						Type: "string",
						Meta: "bigint",
					},
				},
			},
		},
	}
}

func _stake_sign_process() flow.FlowSchema {
	return flow.FlowSchema{
		Name: "Stake",
		ProcessDefinitions: []flow.FlowProcessSchema{
			{
				Name: "signUserOp",
				DefinedInputs: []flow.FlowProcessDefinedInputsSchema{
					{
						Name: "stake_amount",
						Type: "string",
						Meta: "bigint",
					},
				},
			},
			{
				Name:          "processUserOp",
				DefinedInputs: []flow.FlowProcessDefinedInputsSchema{},
			},
		},
	}
}

func _stake_no_data() flow.FlowSchema {
	return flow.FlowSchema{
		Name: "Stake",
		ProcessDefinitions: []flow.FlowProcessSchema{
			{
				Name:          "prepareCallData",
				DefinedInputs: []flow.FlowProcessDefinedInputsSchema{},
			},
			{
				Name:          "processUserOp",
				DefinedInputs: []flow.FlowProcessDefinedInputsSchema{},
			},
		},
	}
}
