package telegrabotlib

type SampleTask struct {
    do Executable
}

func (d *SampleTask) Execute(input Input) *TaskResult {
    if d == nil {
        return InvalidTaskResult
    }

    return d.do(input)
}

type TextInput struct {
    Text string
    Location
}

func (d *TextInput) InputData() interface{} {
    return d.Text
}

type Location struct {
    Lat float64
    Lng float64
}