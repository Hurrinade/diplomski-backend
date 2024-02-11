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
	// Create channel on which we will send data from API
	chanStream := make(chan Response)
	go func() {
		defer close(chanStream)
		for {
			chanStream <- sendEvent(ctx, apiDetails)
			time.Sleep(300 * time.Second)
		}
	}()

	// Create gin stream where it will wait for data that is fetched every 15 seconds and it will send it as sse event
	ctx.Stream(func(w io.Writer) bool {
		if msg, ok := <-chanStream; ok {
			ctx.SSEvent("data", msg)
			return true
		}
		return false
	})
}

func sendEvent(ctx *gin.Context, apiDetails ApiDetails) Response {
	weatherDataResponse, err := http.Get(apiDetails.PljusakURL)

	if err != nil {
		log.Fatal(err)
	}

	if weatherDataResponse.StatusCode == http.StatusOK {
		eventData := FormatData("pljusak", weatherDataResponse.Body)
		return eventData
	} else {
		weatherDataResponse, err = http.Get(apiDetails.WuURL)

		if err != nil {
			log.Fatal(err)
		}

		eventData := FormatData("wu", weatherDataResponse.Body)
		return eventData
	}
}