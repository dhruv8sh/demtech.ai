package user

import (
	"awesomeProject/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	UserId     uint32 `json:"user_id"`
	PlanActive bool   `json:"is_plan_active"`
}

func GetUserFromContext(c *gin.Context) *User {
	user, exists := c.Get("user")
	if !exists {
		util.FailApiResponse(c, http.StatusUnauthorized, "Could not find user")
	}
	if userV, ok := user.(User); ok {
		return &userV
	}
	util.FailApiResponse(c, http.StatusInternalServerError, "Error while resolving user")
	return nil
}

func (u User) GetId() uint32 {
	return u.UserId
}
func (u User) HasActivePlan() bool {
	return u.PlanActive
}
