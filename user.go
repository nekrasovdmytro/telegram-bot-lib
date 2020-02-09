package telegrabotlib

type UserCache struct {
	s Session
}

func NewUserCache(s Session) *UserCache {
	return &UserCache{
		s: s,
	}
}

func (d *UserCache) AddUser(userID, subscribeFor string) error {
	if err :=  d.s.SetForever("storeAll", userID, userID); err != nil {
	    return err
    }

    if err :=  d.s.SetForever(subscribeFor, userID, userID); err != nil {
        return err
    }

    return nil
}

func (d *UserCache) GetAllUsers() (map[string]string, error) {
    return d.s.GetAllLike("storeAll*")
}

func (d *UserCache) GetAllSubscribedUsers(subscribeFor string) (map[string]string, error) {
    return d.s.GetAllLike("*" + subscribeFor)
}

func (d *UserCache) DeleteSubscription(subscribeFor string) error {
    return d.s.DeleteAll(subscribeFor)
}

func (d *UserCache) UnsubscribeUser(userID, subscribeFor string) error {
    return d.s.Delete(subscribeFor, userID)
}



