package middlewares

import (
	"fmt"
	"github.com/egsam98/MegaScout/models"
	errors2 "github.com/egsam98/MegaScout/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"net/http"
)

func JSONErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		ginErr := c.Errors.Last()
		if ginErr == nil {
			return
		}

		switch err := ginErr.Err.(type) {
		case *errors.Error:
			handleStackTraceError(c, err)
		default:
			if err, ok := ginErr.Err.(*errors2.ClientError); ok {
				c.JSON(err.Code, models.ErrorJSON{
					Code:    err.Code,
					Error:   http.StatusText(err.Code),
					Message: err.Error(),
				})
				return
			}
			panic(fmt.Errorf("unhandled error type %v", err))
		}
	}
}

func handleStackTraceError(c *gin.Context, err *errors.Error) {
	switch err.Err.(type) {
	case *errors2.FetchHtmlError:
		c.JSON(408, models.ErrorJSON{
			Code:    408,
			Error:   "Request Timeout",
			Message: err.Error(),
		})
	default:
		c.JSON(500, models.ErrorJSON{
			Code:    500,
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		fmt.Println(err.ErrorStack())
	}
	//if _, ok := err.Err.(*errors2.FetchHtmlError); ok {
	//	c.JSON(408, models.ErrorJSON{
	//		Code:    408,
	//		Error:   "Request Timeout",
	//		Message: err.Error(),
	//	})
	//	return
	//}
	//
	//c.JSON(500, models.ErrorJSON{
	//	Code:    500,
	//	Error:   "Internal Server Error",
	//	Message: err.Error(),
	//})
	//fmt.Println(err.ErrorStack())
}
