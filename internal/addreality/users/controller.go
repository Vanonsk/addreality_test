package users

import (
	"time"

	"addreality_t/internal/addreality/users/withdb"

	"github.com/paulbellamy/ratecounter"
)

type Controller struct {
	uw UserReadWriter
}

type UserReadWriter interface {
	SaveRateUser(user *withdb.User)
	GetAllUsers() ([]*withdb.User, error)
	GetUser(userId string) (*withdb.User, error)
}

func NewController(uw UserReadWriter) *Controller {
	return &Controller{
		uw: uw,
	}
}

func (c *Controller) AddCounterUser(userId string) {
	user, _ := c.uw.GetUser(userId)
	if user == nil {
		counter := ratecounter.NewRateCounter(60 * time.Second)
		counter.Incr(1)
		c.uw.SaveRateUser(&withdb.User{ID: userId, Counter: counter})
	} else {
		user.Counter.Incr(1)
		c.uw.SaveRateUser(user)
	}
}

func (c *Controller) GetRobotUserCount() int {
	robotCountUser := 0
	users, _ := c.uw.GetAllUsers()
	for _, user := range users {
		if user.Counter.Rate() >= 100 {
			robotCountUser++
		}
	}
	return robotCountUser
}
