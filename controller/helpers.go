package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

func FormatData(source string, data io.ReadCloser) Response {
	// Response body
	var eventData ResponseWeatherData
	bodyBytes, err := io.ReadAll(data)
	if err != nil {
		return Response{
			Error:  true,
			Notice: err.Error(),
			Data:   ResponseWeatherData{},
		}
	}

	defer data.Close()
	
	if source == "pljusak" {
		var respBody PljusakWeatherData
		
		err = json.Unmarshal([]byte(formatPljusak(string(bodyBytes))), &respBody)

		if err != nil {
			return Response{
				Error:  true,
				Notice: err.Error(),
				Data:   ResponseWeatherData{},
			}
		}

		temperature, err := strconv.ParseFloat(respBody.Temperature, 64)

		if err != nil {
			return Response{
				Error:  true,
				Notice: err.Error(),
				Data:   ResponseWeatherData{},
			}
		}

		humidity, err := strconv.ParseFloat(respBody.Humidity, 64)

		if err != nil {
			return Response{
				Error:  true,
				Notice: err.Error(),
				Data:   ResponseWeatherData{},
			}
		}

		eventData = ResponseWeatherData{
			Temperature:      fmt.Sprintf("%.1f", temperature),
			HeatIndex:        fmt.Sprintf("%.2f", heatIndex(HeatIndexInput{
				Temperature: temperature,
				Humidity: humidity,
				Fahrenheit: false,
			})),
			Pressure:         respBody.Pressure,
			SolarRadiation:   respBody.SolarRadiation,
			UV:               respBody.UV,
			WindSpeed:        respBody.WindSpeed,
			WindGust:         respBody.WindGust,
			WindChill:        respBody.WindChill,
			WindDirection:    respBody.WindDirectionF,
			WindSide:         respBody.WindDirection,
			PrecipationRate:  respBody.DailyRainin,
			PrecipationTotal: respBody.Precipitation,
			Humidity:         fmt.Sprintf("%d", int(humidity)),
			DewPoint:         respBody.DewPoint,
		}
	} else if source == "wu" {
		var respBody WuWeatherData

		err = json.Unmarshal(bodyBytes, &respBody)

		if err != nil {
			return Response{
				Error:  true,
				Notice: err.Error(),
				Data:   ResponseWeatherData{},
			}
		}

		observation := respBody.Observations[0]

		eventData = ResponseWeatherData{
			Temperature:      fmt.Sprintf("%.1f", observation.Metric.Temp),
			HeatIndex:        fmt.Sprintf("%.2f", observation.Metric.HeatIndex),
			Pressure:         fmt.Sprintf("%.2f", observation.Metric.Pressure),
			SolarRadiation:   fmt.Sprintf("%.2f", observation.SolarRadiation),
			UV:               fmt.Sprintf("%.2f", observation.UV),
			WindSpeed:        fmt.Sprintf("%.2f", observation.Metric.WindSpeed),
			WindGust:         fmt.Sprintf("%.2f", observation.Metric.WindGust),
			WindChill:        fmt.Sprintf("%.2f", observation.Metric.WindChill),
			WindDirection:    fmt.Sprintf("%.2f", observation.Winddir),
			WindSide:         getWindSide(observation.Winddir),
			PrecipationRate:  fmt.Sprintf("%.2f", observation.Metric.PrecipRate),
			PrecipationTotal: fmt.Sprintf("%.2f", observation.Metric.PrecipTotal),
			Humidity:         fmt.Sprintf("%d", observation.Humidity),
			DewPoint:         fmt.Sprintf("%.2f", observation.Metric.Dewpt),
			Elevation:        fmt.Sprintf("%d", observation.Metric.Elev),
		}
	}

	return Response{
		Error:  false,
		Notice: "Successfull",
		Data:   eventData,
	}
}

func getWindSide(windDirection float64) string {
	if windDirection < 0 || windDirection != windDirection {
		return "N/A"
	}

	// Keep within the range: 0 <= d < 360
	windDirection = float64(int(windDirection) % 360)

	switch {
	case windDirection >= 11.25 && windDirection < 33.75:
		return "NNE"
	case windDirection >= 33.75 && windDirection < 56.25:
		return "NE"
	case windDirection >= 56.25 && windDirection < 78.75:
		return "ENE"
	case windDirection >= 78.75 && windDirection < 101.25:
		return "E"
	case windDirection >= 101.25 && windDirection < 123.75:
		return "ESE"
	case windDirection >= 123.75 && windDirection < 146.25:
		return "SE"
	case windDirection >= 146.25 && windDirection < 168.75:
		return "SSE"
	case windDirection >= 168.75 && windDirection < 191.25:
		return "S"
	case windDirection >= 191.25 && windDirection < 213.75:
		return "SSW"
	case windDirection >= 213.75 && windDirection < 236.25:
		return "SW"
	case windDirection >= 236.25 && windDirection < 258.75:
		return "WSW"
	case windDirection >= 258.75 && windDirection < 281.25:
		return "W"
	case windDirection >= 281.25 && windDirection < 303.75:
		return "WNW"
	case windDirection >= 303.75 && windDirection < 326.25:
		return "NW"
	case windDirection >= 326.25 && windDirection < 348.75:
		return "NNW"
	default:
		return "N"
	}
}

func formatPljusak (pljusak string) string {
    data := strings.Replace(pljusak, "<!--\n", "{\"", 1);
    data = strings.ReplaceAll(data, "=", "\":\"");
    data = strings.ReplaceAll(data, ";", "");
    data = strings.Replace(data, "-->\n", "\"}", 1);
    data = strings.ReplaceAll(data, "\n", "\",\"");
    return data
}

func heatIndex(input HeatIndexInput) (float64) {
	if input.Temperature == 0 && input.Humidity == 0 {
		return 0
	}

	if !input.Fahrenheit {
		// Convert temperature to Fahrenheit if not already in Fahrenheit
		input.Temperature = toFahrenheit(input.Temperature)
	}

	if input.Humidity < 0 || input.Humidity > 100 {
		return 0
	}

	// Steadman's result
	heatIndex := 0.5 * (input.Temperature + 61 + (input.Temperature - 68) * 1.2 + input.Humidity * 0.094)

	// Regression equation of Rothfusz is appropriate
	if input.Temperature >= 80 {
		heatIndexBase := -42.379 +
			2.04901523*input.Temperature +
			10.14333127*input.Humidity +
			-0.22475541*input.Temperature*input.Humidity +
			-0.00683783*input.Temperature*input.Temperature +
			-0.05481717*input.Humidity*input.Humidity +
			0.00122874*input.Temperature*input.Temperature*input.Humidity +
			0.00085282*input.Temperature*input.Humidity*input.Humidity +
			-0.00000199*input.Temperature*input.Temperature*input.Humidity*input.Humidity

		// Adjustment
		if input.Humidity < 13 && input.Temperature <= 112 {
			heatIndex = heatIndexBase - (13-input.Humidity)/4*math.Sqrt((17-math.Abs(input.Temperature-95))/17)
		} else if input.Humidity > 85 && input.Temperature <= 87 {
			heatIndex = heatIndexBase + ((input.Humidity-85)/10)*((87-input.Temperature)/5)
		} else {
			heatIndex = heatIndexBase
		}
	}

	return toCelsius(heatIndex)
}

func toFahrenheit(celsius float64) float64 {
	return celsius*9/5 + 32
}

func toCelsius(fahrenheit float64) float64 {
	return (fahrenheit - 32) * 5 / 9
}