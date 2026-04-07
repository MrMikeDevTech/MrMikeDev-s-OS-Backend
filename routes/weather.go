package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MrMikeDevTech/mrmikedevs-os/database"
	"github.com/MrMikeDevTech/mrmikedevs-os/models"
	"github.com/MrMikeDevTech/mrmikedevs-os/utils"
	"github.com/gofiber/fiber/v2"
)

func WeatherRoutes(app *fiber.App) {
	weather := app.Group("/weather")

	weather.Get("/", func(c *fiber.Ctx) error {
		cacheKey := "weather"

		cachedWeather, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()
		if err == nil {
			var weather models.Weather
			if err := json.Unmarshal([]byte(cachedWeather), &weather); err == nil {
				return c.JSON(fiber.Map{
					"status": "success",
					"data":   weather,
				})
			}
		}

		lat := utils.GetEnv("LAT")
		lon := utils.GetEnv("LON")
		apiKey := utils.GetEnv("WEATHER_API_KEY")

		resWeather, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&units=metric&appid=%s", lat, lon, apiKey))

		if err != nil {
			return c.JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}
		defer resWeather.Body.Close()

		type OpenWeatherResponse struct {
			Main struct {
				Temp     float64 `json:"temp"`
				Humidity float64 `json:"humidity"`
			} `json:"main"`
			Wind struct {
				Speed float64 `json:"speed"`
			} `json:"wind"`
			Rain map[string]float64 `json:"rain"`
		}

		var apiResponse OpenWeatherResponse
		if err := json.NewDecoder(resWeather.Body).Decode(&apiResponse); err != nil {
			return c.JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		rainProb := 0.0
		if val, ok := apiResponse.Rain["1h"]; ok {
			rainProb = val
		}

		weatherData := models.Weather{
			Temp:     apiResponse.Main.Temp,
			Humidity: apiResponse.Main.Humidity,
			Wind:     apiResponse.Wind.Speed,
			RainProb: rainProb,
		}

		if weatherJSON, err := json.Marshal(weatherData); err == nil {
			database.RedisClient.Set(database.Ctx, cacheKey, weatherJSON, 10*time.Minute)
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"data":   weatherData,
		})
	})

}
