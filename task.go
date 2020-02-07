package telegrabotlib

import (
	"strconv"
)

const (
	FirstTask   = 0    //first task must have this index
	OnQueryTask = 9898 //onQuery task must have this index
)

var (
	LastStepIndex     = -1
	LastStep          = Step{Index: LastStepIndex}
	ErrorStep         = Step{Index: -2}
	InvalidTaskResult = &TaskResult{res: []*SingleResult{{Type: INVALID}}, Step: ErrorStep}

	StopExecutionText = "<b>To stop click</b> /stop"
)

type Input interface {
	InputData() interface{}
}

type Executable func(input Input) *TaskResult

type Task interface {
	Execute(input Input) *TaskResult
}

type TaskFlow map[int]Task

type TaskResultType int

const (
	INVALID TaskResultType = iota
	TEXT
	LOCATION
)

type TaskResult struct {
	Step Step

	res []*SingleResult
}

type SingleResult struct {
	Type TaskResultType

	res interface{}
}

func (s TaskResult) Result() []*SingleResult {
	return s.res
}

func (s SingleResult) Result() interface{} {
	return s.res
}

func (s TaskResult) LastStep() bool {
	return s.Step.Index == LastStep.Index
}

func NewTaskResult(value []*SingleResult, step Step) *TaskResult {
	return &TaskResult{res: value, Step: step}
}

func NewSingleTextTaskResult(value string, step Step) *TaskResult {
	return NewTaskResult([]*SingleResult{NewSingleResult(TEXT, value)}, step)
}

func NewSingleResult(t TaskResultType, value interface{}) *SingleResult {
	return &SingleResult{Type: t, res: value}
}

type Step struct {
	Index   int
	Routine func(input Input) *TaskResult
}

func (s Step) String() string {
	return strconv.Itoa(s.Index)
}
