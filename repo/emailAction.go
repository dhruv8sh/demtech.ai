package repo

import (
	"awesomeProject/entity"
	"awesomeProject/util"
	"github.com/gin-gonic/gin"
)

// https://github.com/sqlitebrowser/sqlitebrowser/issues/1270
// Due to the issue mentioned above, SQLite sometimes has problems based on the version you are using.
// An easy solution is to use fmt.Sprintf(), but it is not as elegant.

// UpsertEmailAction upserts the action/status of an email to 'IN_QUEUE', 'SUCCESS', 'FAILURE'
// If the action/status is 'SUCCESS', it also updates in users table
func UpsertEmailAction(messageId uint32, action entity.SendEmailAction) {
	if !action.IsValid() {
		return // Return if status is invalid
	}
	sql := `UPDATE email_records SET action = ? WHERE message_id = ?`
	if err := entity.DB.Exec(sql, string(action), messageId).Error; err != nil {
		util.LogCritical(
			"repo::LogSendEmailCallSuccess: Unhandled SQl Error ",
			err, action, messageId, sql)
		return
	}
	// Now only log the success response to the user table
	if action != entity.Success {
		return
	}
	sql = `UPDATE users SET sent_today = sent_today + 1 WHERE id = ?`
	if err := entity.DB.Exec(sql, messageId).Error; err != nil {
		util.LogCritical(
			"repo::LogSendEmailCallSuccess: Unhandled SQl Error ",
			err, messageId, sql)
	}
}

// LogSendEmailCallSuccess initializes and logs that the email is in queue
func LogSendEmailCallSuccess(userId uint32) (messageId uint32) {
	sql := `INSERT INTO email_records (action, user_id) VALUES ('IN_QUEUE', ?) RETURNING message_id`
	if err := entity.DB.Raw(sql, userId).Scan(&messageId).Error; err != nil {
		util.LogCritical(
			"repo::LogSendEmailCallSuccess: Unhandled SQl Error ",
			err, userId, sql)
	}
	return
}

// GetEmailStatus returns the status of a message.
// If the message is not found for the user, it returns FAILURE
func GetEmailStatus(usrId uint32, messageId uint32) entity.SendEmailAction {
	var out string
	sql := `SELECT action FROM email_records WHERE message_id = ? AND user_id = ?`
	if err := entity.DB.Raw(sql, messageId, usrId).Scan(&out).Error; err != nil {
		util.LogCritical(
			"repo::GetEmailStatus: Unhandled SQl Error ",
			err, messageId, usrId, sql)
		return entity.Failure
	}
	if action := entity.SendEmailAction(out); action.IsValid() {
		return action
	}
	return entity.Failure
}

// GetUserMetrics fetches metrics for one specific user
func GetUserMetrics(usrId uint32) gin.H {
	out := map[string]interface{}{}
	sql := `WITH daily_avg AS (
			SELECT user_id, COUNT(*) / COUNT(DISTINCT DATE(updated_at)) AS avg_sent_per_day
			FROM email_records WHERE user_id = ?),

			email_stats AS ( SELECT 
				user_id,
				COUNT(*) AS total_attempts,
				COUNT(CASE WHEN action = 'SUCCESS' THEN 1 END) AS success_count,
				COUNT(CASE WHEN action = 'FAILURE' THEN 1 END) AS failure_count
			FROM email_records
			WHERE user_id = ?
			GROUP BY user_id)
		SELECT 
			u.id,
			u.emails,
			u.daily_limit,
			u.sent_today,
			(u.daily_limit - u.sent_today) AS remaining_quota,
			es.total_attempts,
			es.success_count,
			es.failure_count,
			ROUND(100.0 * es.success_count / NULLIF(es.total_attempts, 0), 2) AS success_rate,
			ROUND(100.0 * es.failure_count / NULLIF(es.total_attempts, 0), 2) AS failure_rate,
			COALESCE(da.avg_sent_per_day, 0) AS avg_sent_per_day
		FROM users u
		LEFT JOIN email_stats es ON u.id = es.user_id
		LEFT JOIN daily_avg da ON u.id = da.user_id
		WHERE u.id = ?`

	if err := entity.DB.Raw(sql, usrId, usrId, usrId).Scan(&out).Error; err != nil {
		util.LogCritical(
			"repo::GetUserMetrics: Unhandled SQl Error ",
			err, usrId, sql)
		return nil
	}
	return out
}
