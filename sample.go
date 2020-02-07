package telegrabotlib

type sampleTask struct {
    do Executable
}

func (d *sampleTask) Execute(input Input) *TaskResult {
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