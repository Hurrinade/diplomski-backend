/* Validate input */
package controller

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/* Send chart data also every 300 seconds */
func GetChartData(ctx *gin.Context, cl *mongo.Client) {
	chanStream := make(chan ResponseChart)
	
	go func() {
		defer close(chanStream)
		for {
			chanStream <- sendChartEvent(ctx, cl)
			time.Sleep(300 * time.Second)
		}
	}()

	gone := ctx.Stream(func(w io.Writer) bool {
		if msg, ok := <-chanStream; ok {
			ctx.SSEvent("data", msg)
			return true
		}

		return false
	})

	if gone {
		log.Println("Client disconnected")
	}
}

func GetEvents(ctx *gin.Context, cl *mongo.Client, location string, apiDetails ApiDetails) {
	db := cl.Database("meteosite")
	initial := true
	
	var collection *mongo.Collection
	if location == "vrapce"{
		collection = db.Collection("vrapce-chart")
	} else if location == "mlinovi" {
		collection = db.Collection("mlinovi-chart")
	}

	// Create channel on which we will send data from API
	chanStream := make(chan Response)
	go func() {
		defer close(chanStream)
		for {
			chanStream <- sendEvent(ctx, apiDetails, cl)
			time.Sleep(300 * time.Second)
		}
	}()

	// Create gin stream where it will wait for data that is fetched every 300 seconds and it will send it as sse event
	ctx.Stream(func(w io.Writer) bool {
		if msg, ok := <-chanStream; ok {
			ctx.SSEvent("data", msg)

			if !initial {
				/* Add data to databse */
				dbd := DbData{
					Date: msg.Data.Date,
					Temperature: msg.Data.Temperature,
					WindSpeed: msg.Data.WindSpeed,
					Precipation: msg.Data.PrecipationTotal,
					Humidity: msg.Data.Humidity,
				}
	
				_, err := collection.InsertOne(context.TODO(), dbd)

				if err != nil {
					log.Println("Failed to insert data", err)
				}
			}

			initial = false

			return true
		}
		return false
	})
}

func sendChartEvent(ctx *gin.Context, cl *mongo.Client) ResponseChart{
	db := cl.Database("meteosite")

	vrapceCollection := db.Collection("vrapce-chart")
	mlinoviCollection := db.Collection("mlinovi-chart")

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)

	// Get data for last 7 days
	filter := bson.M{"date": bson.M{"$gte": startDate, "$lt":  endDate}}

	// Get cursor to that data
	cursorVrapce, err := vrapceCollection.Find(context.TODO(), filter)
	
	if err != nil {
		log.Fatalln(err)
	}

	cursorMlinovi, err := mlinoviCollection.Find(context.TODO(), filter)

	if err != nil {
		log.Fatalln(err)
	}

	// Map data into array
	var resultsVrapce []DbData
	if err = cursorVrapce.All(context.TODO(), &resultsVrapce); err != nil {
		log.Fatalln(err)
	}

	var resultsMlinovi []DbData
	if err = cursorMlinovi.All(context.TODO(), &resultsMlinovi); err != nil {
		log.Fatalln(err)
	}

	return ResponseChart{
		Error: false,
		Notice: "",
		Data: JointDbData{
			VrapceData: resultsVrapce,
			MlinoviData: resultsMlinovi,
		},
	}
}

func sendEvent(ctx *gin.Context, apiDetails ApiDetails, cl *mongo.Client) Response {
	weatherDataResponse, err := http.Get(apiDetails.PljusakURL)

	if err != nil {
		log.Println("Error when fetching endpoint pljusak",err)
	}

	if weatherDataResponse.StatusCode == http.StatusOK {
		eventData := FormatData("pljusak", weatherDataResponse.Body)
		return eventData
	} else {
		weatherDataResponse, err = http.Get(apiDetails.WuURL)

		if err != nil {
			log.Println("Error when fetching endpoint WU",err)
		}

		eventData := FormatData("wu", weatherDataResponse.Body)
		return eventData
	}
}