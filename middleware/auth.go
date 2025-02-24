package middleware

import (
	"awesomeProject/entity"
	errEntity "awesomeProject/entity/errors"
	"awesomeProject/user"
	"awesomeProject/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

/*
	The JWT was implemented as a Bearer Token inside Header->Authorization
	Due to time constraints, the token is non-obfuscated and the token is the user id itself.
	For example,
		For user_id = 1
		Auth token = "Bearer 1"
*/

// Authenticate middleware checks for API Authentication except for /api/auth
func Authenticate(c *gin.Context) {
	//if strings.HasPrefix(c.Request.RequestURI, "/api/auth/") {
	//	c.Next()
	//	return
	//}

	// Get Bearer JWT Token - Validate to get user data and set token in the request context
	bearerToken := ""
	reqToken := c.Request.Header.Get("Authorization")
	if reqToken != "" {
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) == 2 {
			bearerToken = splitToken[1]
		}
	}

	// Get and then set the current user
	c.Set("current_user", getUserFromToken(bearerToken))
	c.Next()
}

// Fetches user.User object given a Token
// If the reset time has passed, it sets the reset time to one day from now.
// It also warms up email to a max of 1000 emails per day incrementing by 100 for the first 10 days.
func getUserFromToken(token string) *user.User {
	var dbRow struct {
		Id         uint32     `gorm:"id"`
		Emails     string     `gorm:"emails"`
		DailyLimit int        `gorm:"daily_limit"`
		SendLimit  int        `gorm:"send_limit"`
		SentToday  int        `gorm:"sent_today"`
		CreatedAt  *time.Time `gorm:"created_at"`
		ResetAt    *time.Time `gorm:"reset_at"`
	}

	id := decodeObfuscatedIdFromToken(token)
	sql := `UPDATE users
			SET reset_at =
				CASE 
                  WHEN reset_at < CURRENT_TIMESTAMP 
                  THEN DATETIME(CURRENT_TIMESTAMP, '+1 day')
                  ELSE reset_at
                END,
    		sent_today = 
				CASE 
                   WHEN reset_at < CURRENT_TIMESTAMP 
                   THEN 0
                   ELSE sent_today
                END,
    		daily_limit = CAST(
				CASE 
                    WHEN reset_at < CURRENT_TIMESTAMP 
                         AND (JULIANDAY('now') - JULIANDAY(created_at)) < 10 
                    THEN (JULIANDAY('now') - JULIANDAY(created_at)) * 100
                    ELSE daily_limit
                END AS INTEGER)
			WHERE id = ?
			RETURNING *`
	if err := entity.DB.Raw(sql, id).Scan(&dbRow).Error; err != nil || dbRow.Id == 0 {
		util.HttpFail(http.StatusUnauthorized, "Could not resolve user")
		return nil
	}
	usr := user.User{}
	user.SetUser(&usr, map[string]any{
		"id":          dbRow.Id,
		"emails":      dbRow.Emails,
		"created_at":  dbRow.CreatedAt,
		"daily_limit": dbRow.DailyLimit,
		"sent_today":  dbRow.SentToday,
		"reset_at":    dbRow.ResetAt,
	})
	return &usr
}

func decodeObfuscatedIdFromToken(token string) string {
	// panic("Currently unimplemented")
	return token
}

// AuthorizeAccess middleware checks if sending mail is paused for a user
func AuthorizeAccess(c *gin.Context) {
	// Do access authorization here
	c.Next()
}

// SendMailQuotaVerify checks if the user can send any e-mails at the moment
func SendMailQuotaVerify(c *gin.Context) {
	usr := user.GetUserFromContext(c)
	if usr == nil {
		util.HttpFailCustom(errEntity.InvalidClientTokenId)
		return
	}
	quota := usr.GetQuota()
	if quota.RemainingQuota <= 0 {
		util.HttpFailCustom(errEntity.AccountSendingPaused)
		return
	}
	c.Next()
}
