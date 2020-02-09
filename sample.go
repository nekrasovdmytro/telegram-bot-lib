package telegrabotlib

type SampleTask struct {
    Do Executable
}

func (d *SampleTask) Execute(input Input) *TaskResult {
    if d == nil {
        return InvalidTaskResult
    }

    return d.Do(input)
}

type TextInput struct {
    UserId int
    Username string
    Text string
    Location
}

func (d *TextInput) InputData() interface{} {
    return d.Text
}

type Location struct {
    Lat float32
    Lng float32
}