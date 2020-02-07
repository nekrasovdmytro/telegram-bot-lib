package telegrabotlib

type Ability struct {
	Short       string
	Name        string
	Flow        TaskFlow
	Description string
}

type AbilityMap map[string]*Ability

func (d AbilityMap) GetOnQueryTask(s string) (Task, bool) {
    for _ , a := range d {
        if a.Short == s {

            if t, ok := a.Flow[OnQueryTask]; ok {
                return t, true
            }
        }
    }

    return nil, false
}
