package controller

import (
	"awesomeProject/bvr"
	"awesomeProject/entity"
	"awesomeProject/entity/errors"
	"awesomeProject/repo"
	"awesomeProject/user"
	"awesomeProject/util"
	"github.com/gin-gonic/gin"
	"slices"
)

func SendEmail(c *gin.Context) {
	// Get User instance from middleware
	usr := user.GetUserFromContext(c)
	if usr == nil {
		// Fail if user does not exist. (Redundant)
		util.HttpFailCustom(errEntity.InvalidClientTokenId)
		return
	}

	// Bind to incoming body to params
	var params entity.SendEmailParams
	if !entity.BindRequestAndValidate(c, &params) {
		return
	}

	// If the user does not have a registered email, fail.
	if !slices.Contains(usr.GetDomains(), params.Source) {
		util.HttpFailCustom(errEntity.MailFromDomainNotVerified)
	}

	// Log Success of the API validation
	messageId := repo.LogSendEmailCallSuccess(usr.GetId())
	util.LogInfo("SendMail request received successfully with messageId:", messageId)

	// Process the mail in the background
	go bvr.AsyncSendEmail(messageId, &params)

	util.SuccessApiResponse(
		c,
		"Request received successfully.",
		gin.H{"MessageId": messageId},
	)
}

func GetEmailStatus(c *gin.Context) {
	// Ensure user
	usr := user.GetUserFromContext(c)
	if usr == nil {
		util.HttpFailCustom(errEntity.InvalidClientTokenId)
		return
	}

	var params struct {
		MessageId uint32 `xml:"MessageId"`
	}
	if !entity.BindQueryRequestAndValidate(c, &params) {
		return
	}

	util.SuccessApiResponse(
		c,
		"Data fetched successfully.",
		repo.GetEmailStatus(usr.GetId(), params.MessageId))
}

func GetQuota(c *gin.Context) {
	usr := user.GetUserFromContext(c)
	if usr == nil {
		util.HttpFailCustom(errEntity.InvalidClientTokenId)
		return
	}
	quota := usr.GetQuota()
	util.SuccessApiResponse(
		c,
		"Data fetched successfully",
		gin.H{
			"RemainingQuota": quota.RemainingQuota,
			"SentToday":      quota.SentToday,
			"ResetTime":      quota.ResetTime,
			"MaxSendLimit":   quota.MaxSendLimit,
		})
}

func GetMetrics(c *gin.Context) {
	usr := user.GetUserFromContext(c)
	if usr == nil {
		util.HttpFailCustom(errEntity.InvalidClientTokenId)
		return
	}
	util.SuccessApiResponse(c, "Data fetched successfully", repo.GetUserMetrics(usr.GetId()))
}
