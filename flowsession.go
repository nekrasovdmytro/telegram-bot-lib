package telegrabotlib

import (
    "log"
)

const (
    taskFlowSessionKey = "task_flow"
    CurrentFlowTaskKey = "currentFlowTask"
    CurrentTaskKey     = "currentTask"
)

func NewUserFlowSession(session Session) *UserFlowSession {
    return &UserFlowSession{session: session}
}

type UserFlowSession struct {
    session Session
}

func (d UserFlowSession) GetFlowID(userId string) (string, error) {
    f, err := d.session.Get(userId+taskFlowSessionKey, CurrentFlowTaskKey)
    if err != nil {
        return "", err
    }

    return f, nil
}

func (d UserFlowSession) GetTaskID(userId string) (string, error) {
    f, err := d.session.Get(userId+taskFlowSessionKey, CurrentTaskKey)
    if err != nil {
        return "", err
    }

    return f, nil
}

func (d UserFlowSession) WithFlowWaiting(userId string) error {
    err := d.session.Set(userId+taskFlowSessionKey, "flowWaits", "wait")
    if err != nil {
        return err
    }

    return nil
}

func (d UserFlowSession) FlowWaiting(userId string) bool {
    f, err := d.session.Get(userId+taskFlowSessionKey, "flowWaits")
    if err != nil {
        return false
    }

    return f != ""
}

func (d UserFlowSession) StopFlowWaiting(userId string) error {
    err := d.session.Delete(userId+taskFlowSessionKey, "flowWaits")
    if err != nil {
        return err
    }

    return nil
}

func (d UserFlowSession) SaveFlowTask(userId, flowTaskId, taskId string) error {
    log.Printf("Saved flow task %s %s %s", userId, flowTaskId, taskId)
    if err := d.session.Set(userId+taskFlowSessionKey, CurrentFlowTaskKey, flowTaskId); err != nil {
        return err
    }

    if err := d.session.Set(userId+taskFlowSessionKey, CurrentTaskKey, taskId); err != nil {
        return err
    }

    return nil
}

func (d UserFlowSession) FinishFlow(userId string) error {
    if err := d.session.DeleteAll(userId+ taskFlowSessionKey); err != nil {
        return err
    }

    return nil
}

