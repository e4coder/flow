package flow

import (
	"errors"
	"fmt"
	"strings"
)

type ProcessHandler func(*ProcessContext, []DefinedInput) error

type FlowProcessDefinedInputsSchema struct {
	Name string
	Type string
	Meta string
}

type FlowProcessSchema struct {
	Name          string
	DefinedInputs []FlowProcessDefinedInputsSchema
}

type FlowSchema struct {
	Name               string
	ProcessDefinitions []FlowProcessSchema
}

type FlowParser struct {
	FlowList map[string]FlowSchema
	Handlers map[string]ProcessHandler
}

func NewFlowParser(handlers map[string]ProcessHandler) *FlowParser {
	flowParser := &FlowParser{}
	flowParser.Handlers = handlers
	flowParser.FlowList = make(map[string]FlowSchema)
	return flowParser
}

func (fp *FlowParser) Add(name string, fpd FlowSchema) {
	fp.FlowList[name] = fpd
}

func (fp *FlowParser) Parse(request FlowRequest) (*Flow, error) {
	err := _verifySchema(fp, request)
	if err != nil {
		return nil, errors.Join(ErrParserFailure, ErrSchemaVerificationFailure, err)
	}

	f, err := _parseRequest(fp, request)
	if err != nil {
		return nil, errors.Join(ErrParserFailure, err)
	}

	return f, nil
}

func _parseRequest(fp *FlowParser, request FlowRequest) (*Flow, error) {
	flow := &Flow{}
	flow.Name = request.Name

	for i, requestProcess := range request.FlowRequestProcesses {
		flow.Processes = append(flow.Processes, Process{
			Name: requestProcess.Name,
		})

		for _, requestInput := range requestProcess.DefinedInput {
			flow.Processes[i].DefinedInput = append(flow.Processes[i].DefinedInput, DefinedInput{
				Name:  requestInput.Name,
				Value: requestInput.Value,
				Type:  requestInput.Type,
				Meta:  requestInput.MetaType,
			})
		}

		handler, ok := fp.Handlers[requestProcess.Name]
		if !ok {
			eInfo := fmt.Errorf("handler: %s", flow.Processes[i].Name)
			return nil, errors.Join(ErrHandlerNotFound, eInfo)
		}
		flow.Processes[i].Handler = handler
	}

	return flow, nil
}

func _verifySchema(fp *FlowParser, request FlowRequest) error {
	flowSchema, ok := fp.FlowList[request.Name]
	if !ok {
		eInfo := fmt.Errorf("flow: %s", request.Name)
		return errors.Join(ErrFlowNotFound, eInfo)
	}

	if len(flowSchema.ProcessDefinitions) != len(request.FlowRequestProcesses) {
		eInfo := fmt.Errorf("len mismatch schema(%d) != request(%d)", len(flowSchema.ProcessDefinitions), len(request.FlowRequestProcesses))
		return errors.Join(ErrInvalidRequest, eInfo)
	}

	for i, schemaProcess := range flowSchema.ProcessDefinitions {
		requestProcess := request.FlowRequestProcesses[i]
		if strings.Compare(requestProcess.Name, schemaProcess.Name) != 0 {
			eInfo := fmt.Errorf("process name mismatch schema(%s) != request(%s)", schemaProcess.Name, requestProcess.Name)
			return errors.Join(ErrInvalidRequestProcess, eInfo)
		}

		if len(schemaProcess.DefinedInputs) != len(requestProcess.DefinedInput) {
			eInfo := fmt.Errorf("len mismatch schema(%d) != request(%d)", len(schemaProcess.DefinedInputs), len(requestProcess.DefinedInput))
			return errors.Join(ErrInvalidRequestInputs, eInfo)
		}
	}

	return nil
}
