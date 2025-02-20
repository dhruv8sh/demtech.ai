package middleware

import (
	"awesomeProject/user"
	"awesomeProject/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Authenticate middleware checks for API Authentication except for /api/auth
func Authenticate(c *gin.Context) {
	if strings.HasPrefix(c.Request.RequestURI, "/api/auth/") {
		c.Next()
		return
	}

	// Get Bearer JWT Token - Validate to get user data and set token in the request context
	bearerToken := ""
	reqToken := c.Request.Header.Get("Authorization")
	if reqToken != "" {
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) == 2 {
			bearerToken = splitToken[1]
		}
	}

	u, ok := getUserFromToken(bearerToken)
	if !ok {

	}
	c.Set("token", u)
	c.Next()
}

// Fetches user.User object given a Token
func getUserFromToken(token string) (user.User, bool) {
	dbUser := user.User{
		UserId:     0,
		PlanActive: false,
	}
	/*

		id := decodeObfuscatedIdFromToken(token)

		sql := `SELECT id, is_plan_active FROM users WHERE id = ?`

		err := dbConnection.getResult(sql, id).Scan(dbUser)
		if err != nil {
			LogCritical("Unhandled SQL Error", err, id, sq)
			return dbUser, false
		}

	*/
	return dbUser, true
}

// AuthorizeAccess middleware checks if the user has an active plan
func AuthorizeAccess(c *gin.Context) {
	if user.GetUserFromContext(c).HasActivePlan() {
		c.Next()
		return
	}

	util.HttpFail(http.StatusForbidden, "Your plan has expired")
}
