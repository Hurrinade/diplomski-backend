package controller

import (
	"time"
)

type ApiDetails struct {
	WuURL      string
	PljusakURL string
}

type Metric struct {
	Temp       float64     `json:"temp"`
	HeatIndex  float64     `json:"heatIndex"`
	Dewpt      float64     `json:"dewpt"`
	WindChill  float64     `json:"windChill"`
	WindSpeed  float64     `json:"windSpeed"`
	WindGust   float64     `json:"windGust"`
	Pressure   float64 `json:"pressure"`
	PrecipRate float64 `json:"precipRate"`
	PrecipTotal float64 `json:"precipTotal"`
	Elev       int     `json:"elev"`
}

type Observation struct {
	StationID          string    `json:"stationID"`
	ObsTimeUtc         time.Time `json:"obsTimeUtc"`
	ObsTimeLocal       string    `json:"obsTimeLocal"`
	Neighborhood       string    `json:"neighborhood"`
	SoftwareType       string    `json:"softwareType"`
	Country            string    `json:"country"`
	SolarRadiation     float64   `json:"solarRadiation"`
	Lon                float64   `json:"lon"`
	RealtimeFrequency  interface{} `json:"realtimeFrequency"`
	Epoch              int64     `json:"epoch"`
	Lat                float64   `json:"lat"`
	UV                 float64   `json:"uv"`
	Winddir            float64       `json:"winddir"`
	Humidity           int       `json:"humidity"`
	QCStatus           float64       `json:"qcStatus"`
	Metric             Metric    `json:"metric"`
}

type WuWeatherData struct {
	Observations []Observation `json:"observations"`
}

type ResponseWeatherData struct {
	Temperature      string
	HeatIndex        string
	Pressure         string
	SolarRadiation  string
	UV               string
	WindSpeed        string
	WindGust         string
	WindChill        string
	WindDirection    string
	WindSide         string
	PrecipationRate  string
	PrecipationTotal string
	Humidity         string
	DewPoint         string
	Elevation        string
}

type PljusakWeatherData struct {
	Data             string    `json:"data"`
	Date             string    `json:"datum"`
	Time             string    `json:"vrijeme"`
	Temperature      string   `json:"temperatura"`
	Humidity         string       `json:"vlaga"`
	Pressure         string   `json:"tlak"`
	WindDirection    string    `json:"smjer"`
	WindSpeed        string   `json:"brzina"`
	WindGust         string   `json:"udar"`
	Precipitation    string   `json:"oborined"`
	IndoorTemperature string  `json:"temperatura_IN"`
	IndoorHumidity   string       `json:"vlaga_IN"`
	SolarRadiation   string   `json:"solar"`
	UV               string       `json:"uv"`
	WindChill        string   `json:"wind_chill"`
	DewPoint         string   `json:"dew_point"`
	LocationID       string    `json:"all_ID"`
	LocationPassword string    `json:"all_PASSWORD"`
	IndoorTempF      string   `json:"all_indoortempf"`
	TempF            string   `json:"all_tempf"`
	DewptF           string   `json:"all_dewptf"`
	WindchillF       string   `json:"all_windchillf"`
	IndoorHumidityF  string       `json:"all_indoorhumidity"`
	HumidityF        string       `json:"all_humidity"`
	WindspeedMph     string   `json:"all_windspeedmph"`
	WindgustMph      string   `json:"all_windgustmph"`
	WindDirectionF   string       `json:"all_winddir"`
	AbsBaromin       string   `json:"all_absbaromin"`
	Baromin          string   `json:"all_baromin"`
	Rainin           string   `json:"all_rainin"`
	DailyRainin      string   `json:"all_dailyrainin"`
	WeeklyRainin     string   `json:"all_weeklyrainin"`
	MonthlyRainin    string   `json:"all_monthlyrainin"`
	Solarradiation   string   `json:"all_solarradiation"`
	AllUV            string       `json:"all_UV"`
	DateUTC          string `json:"all_dateutc"`
	Realtime         string       `json:"all_realtime"`
	Rtfreq           string       `json:"all_rtfreq"`
}

type Response struct {
	Error bool 					`json:"error"`
	Notice string				`json:"notice"`
	Data ResponseWeatherData 	`json:"data"`
}

type HeatIndexInput struct {
	Temperature float64
	Humidity    float64
	Fahrenheit  bool
}

