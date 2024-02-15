package router

import (
	"github.com/Hurrinade/diplomski-backend/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(cl *mongo.Client) *gin.Engine {
	// Create a new Gin router
	r := gin.Default()

	// Register the middleware
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           86400,
	}))


	// Set up the authentication routes
	r.GET("/v1/getVrapceEvents", func(ctx *gin.Context) {	
		controller.GetEvents(
			ctx, 
			cl,
			"vrapce",
			controller.ApiDetails{
				WuURL: "https://api.weather.com/v2/pws/observations/current?stationId=IZAGRE19&format=json&units=m&apiKey=8e48f7be32604eb288f7be3260beb267",
				PljusakURL: "https://pljusak.com/1_wu/vrapce.txt",
		})
	})
	r.GET("/v1/getMlinoviEvents", func(ctx *gin.Context) {
		controller.GetEvents(
			ctx, 
			cl,
			"mlinovi", 
			controller.ApiDetails{
				WuURL: "https://api.weather.com/v2/pws/observations/current?stationId=IUNDEFIN41&format=json&units=m&apiKey=8e48f7be32604eb288f7be3260beb267",
				PljusakURL: "https://pljusak.com/1_wu/mlinovi.txt",
		})
	})
	r.GET("/v1/chartData", func(ctx *gin.Context) {
		controller.GetChartData(ctx, cl)
	})

	return r
}