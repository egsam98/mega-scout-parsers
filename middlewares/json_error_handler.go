package middlewares

import (
	"errors"
	"fmt"
	"github.com/egsam98/MegaScout/models"
	errors2 "github.com/egsam98/MegaScout/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var FetchHtmlError = &errors2.FetchHtmlError{}
var ClientError = &errors2.ClientError{}

func JSONErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}

		if errors.As(err.Err, &FetchHtmlError) {
			c.JSON(408, models.ErrorJSON{
				Code:    408,
				Error:   "Request Timeout",
				Message: err.Error(),
			})
			return
		}

		if errors.As(err.Err, &ClientError) {
			err := err.Err.(*errors2.ClientError)
			c.JSON(err.Code, models.ErrorJSON{
				Code:    err.Code,
				Error:   http.StatusText(err.Code),
				Message: err.Error(),
			})
			return
		}

		c.JSON(500, models.ErrorJSON{
			Code:    500,
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		fmt.Printf("%+v\n", err.Err)
	}
}
