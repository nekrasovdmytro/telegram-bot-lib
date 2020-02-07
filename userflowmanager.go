package telegrabotlib

import (
    "log"
    "strconv"
)

type UserFlowManager struct {
    uSession *UserFlowSession
}

func NewUserFlowManager(us *UserFlowSession) *UserFlowManager {
    return &UserFlowManager{
        uSession: us,
    }
}

func (d *UserFlowManager) CurrentTask(abilities AbilityMap, userID string) (flowName string, currentTask Task, ok bool) {
    flowID, err := d.uSession.GetFlowID(userID)
    if err != nil {
        return "", nil, false
    }

    taskID, err := d.uSession.GetTaskID(userID)
    if err != nil {
        d.uSession.FinishFlow(userID) //delete if no current task - flow is broken
        return "", nil, false
    }

    log.Printf("Got flow %s and task %s", flowID, taskID)
    if tf, flowFound := abilities[flowID]; flowFound {
        i, _ := strconv.Atoi(taskID)
        t, taskFound := tf.Flow[i]
        if !taskFound {
            log.Printf("Delete if task no found")
            d.uSession.FinishFlow(userID) //delete if task no found
            return "", nil, false
        }

        return flowID, t, true
    }

    log.Printf("Delete all - nothing found")
    d.uSession.FinishFlow(userID) //delete if taskFlow doesn't exist

    return "", nil, false
}

func (d *UserFlowManager) ExecuteFlow(bot Bot, userID string, r Recipient, input Input)  {
    if d.uSession.FlowWaiting(userID) {
        bot.SendMessage(r, "please wait...")
        return
    }

    ability, abilityRequest := d.GetAbility(bot, input)

    c, t, ok := d.CurrentTask(bot.Abilities(), userID)
    if ok && !abilityRequest {
        log.Print("Found in progress")
        res := bot.Execute(userID, r, t.Execute, input)

        if res.Step.Routine != nil {
            log.Print("Execute routine")
            d.uSession.WithFlowWaiting(userID)
            defer d.uSession.StopFlowWaiting(userID)
            bot.Execute(userID, r, res.Step.Routine, input)
        }

        if res.LastStep() {
            d.uSession.FinishFlow(userID)
            return
        }

        if err := d.uSession.SaveFlowTask(userID, c, res.Step.String()); err != nil {
            log.Fatal(err)
        }

        return
    }

    if abilityRequest {
        log.Print("Found ability")
        firstTask := ability.Flow[FirstTask]
        res := bot.Execute(userID, r, firstTask.Execute, input)

        d.uSession.SaveFlowTask(userID, input.InputData().(string), res.Step.String())
    } else {
        bot.SendMessage(r, "This what I can do:")
        for k, a := range bot.Abilities() {
            bot.SendMessage(r, k + " - " + a.Description)
        }
    }
}

func (d *UserFlowManager) GetAbility(bot Bot, input Input) (*Ability, bool) {
    a, ok :=  bot.Abilities()[input.InputData().(string)]
    return a, ok
}
