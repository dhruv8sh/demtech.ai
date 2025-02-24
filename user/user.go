package user

import (
	"awesomeProject/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// User is a struct to keep necessary information about the user
// This is initialized in the middleware and then put into the gin.Context
type User struct {
	userId            uint32
	quota             Quota
	signupTime        *time.Time
	registeredDomains []string
}

// Quota stores the available quota for the user
type Quota struct {
	MaxSendLimit   int        `gorm:"MaxSendLimit"   json:"MaxSendLimit"`
	SentToday      int        `gorm:"SentToday"      json:"SentToday"`
	RemainingQuota int        `gorm:"RemainingQuota" json:"RemainingQuota"`
	ResetTime      *time.Time `gorm:"ResetTime"      json:"ResetTime"`
}

// SetUser fetches information from a map and puts it into the User Reference
func SetUser(cu *User, data map[string]any) {
	cu.userId, _ = data["id"].(uint32)
	d, _ := data["emails"].(string)
	var regDomains []string
	if err := json.Unmarshal([]byte(d), &regDomains); err != nil {
		util.LogDebug("Error while unmarshal Register Domains, setting none", err.Error())
		cu.registeredDomains = make([]string, 0)
	} else {
		cu.registeredDomains = regDomains
	}
	cu.quota = Quota{
		MaxSendLimit:   data["daily_limit"].(int),
		SentToday:      data["sent_today"].(int),
		RemainingQuota: data["daily_limit"].(int) - data["sent_today"].(int),
		ResetTime:      data["reset_at"].(*time.Time),
	}
}

// GetUserFromContext is used like a helper function to get the pointer to the currently set user
func GetUserFromContext(c *gin.Context) *User {
	user, exists := c.Get("current_user")
	if !exists {
		util.HttpFail(http.StatusUnauthorized, "Could not find user")
		c.Abort()
		return nil
	}
	if userV, ok := user.(*User); ok {
		return userV
	}
	util.HttpFail(http.StatusInternalServerError, "Error while resolving user")
	return nil
}

// Getter functions------------

func (u *User) GetId() uint32 {
	return u.userId
}
func (u *User) GetDomains() []string {
	return u.registeredDomains
}
func (u *User) GetQuota() *Quota {
	return &u.quota
}
