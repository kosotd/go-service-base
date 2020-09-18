package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kosotd/go-service-base/utils"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"strings"
)

var TokenChecker func(username string, token string) error
var tokenPattern = regexp.MustCompile(bearerTokenPattern)

func CheckTokenHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			utils.LogErrorAndSetStatus(c.Writer, http.StatusExpectationFailed, errors.New(fmt.Sprint(r)))
			c.Abort()
		}
	}()

	if TokenChecker == nil {
		c.Next()
		return
	}

	username := c.GetHeader(usernameHeader)
	if utils.Empty(username) {
		utils.LogErrorAndSetStatus(c.Writer, http.StatusForbidden, errors.New("Username header not found"))
		c.Abort()
		return
	}

	auth := c.GetHeader(authorizationHeader)
	subm := tokenPattern.FindStringSubmatch(auth)
	if subm == nil || len(subm) < 1 {
		utils.LogErrorAndSetStatus(c.Writer, http.StatusForbidden, errors.New("Authorization header not found"))
		c.Abort()
		return
	}

	if err := TokenChecker(username, strings.Trim(subm[1], " ")); err == nil {
		c.Next()
	} else {
		utils.LogErrorAndSetStatus(c.Writer, http.StatusForbidden, errors.Wrapf(err, "handler.checkTokenHandler -> service.CheckToken"))
		c.Abort()
	}
}
