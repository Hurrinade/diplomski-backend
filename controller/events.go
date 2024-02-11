/* Validate input */
package controller

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetEvents(ctx *gin.Context, apiDetails ApiDetails) {
	gone := ctx.Stream(func(w io.Writer) bool {
		sendEvent(ctx, apiDetails)
		time.Sleep(15 * time.Second)
		return true
	})

	if gone {
		log.Println("Client disconected")
	}
}

func sendEvent(ctx *gin.Context, apiDetails ApiDetails) {
	weatherDataResponse, err := http.Get(apiDetails.PljusakURL)

	if err != nil {
		log.Fatal(err)
	}

	if weatherDataResponse.StatusCode == http.StatusOK {
		eventData := FormatData("pljusak", weatherDataResponse.Body)
		ctx.SSEvent("data", eventData)
	} else {
		weatherDataResponse, err = http.Get(apiDetails.WuURL)

		if err != nil {
			log.Fatal(err)
		}

		eventData := FormatData("wu", weatherDataResponse.Body)
		ctx.SSEvent("data", eventData)
	}
}