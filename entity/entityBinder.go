package entity

import (
	errEntity "awesomeProject/entity/errors"
	"awesomeProject/util"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
	"net/http"
	"strings"
)

func init() {
	validate.AddValidator("emails", validateEmails)
}

// Validates a string of emails using gookit email validation
func validateEmails(emails []string) bool {
	for _, email := range emails {
		email = strings.TrimSpace(email)
		if !validate.IsEmail(email) {
			return false
		}
	}
	return true
}

// BindRequestAndValidate binds to body in POST, PUT, etc.
func BindRequestAndValidate(c *gin.Context, s any) bool {
	err := c.ShouldBindXML(s)
	if err != nil {
		util.HttpFailCustom(errEntity.MissingParameter)
		return false
	}

	// Perform validations using gookit
	v := validate.Struct(s)
	if !v.Validate() {
		util.HttpFailCustom(errEntity.InvalidParameterValue)
		return false
	}
	return true
}

// BindQueryRequestAndValidate binds to GET request params
func BindQueryRequestAndValidate(c *gin.Context, s interface{}) bool {
	err := c.ShouldBindQuery(s)
	if err != nil {
		util.FailApiResponse(c, http.StatusUnprocessableEntity, err)
		return false
	}

	// Perform validations using gookit
	v := validate.Struct(s)
	if !v.Validate() {
		util.FailApiResponse(c, http.StatusUnprocessableEntity, v.Errors.String())
		return false
	}

	return true
}
