package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
)

// @Security      BearerAuth
// @Summary       Get Weather
// @Description   API for getting weather data by city
// @Tags          weather
// @Accept        json
// @Produce       json
// @Param         city query string true "City name"
// @Success       200 {object} entity.Weather
// @Failure       400 {object} entity.Error
// @Failure       500 {object} entity.Error
// @Router        /weather [get]
func (h *HandlerV1) GetWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City is required"})
		return
	}
	pp.Println(city)
	if err := h.WeatherService.FetchAndStoreWeather(city); err != nil {
		// Xato yuz beradigan bo'lsa, HTTP status 500 (Internal Server Error) va xatolikni qaytarish
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("could not fetch and store weather for city %s: %v", city, err),
		})
		h.Logger.Error("error while fetch weather", zap.Error(err))
		return
	}
	

	weather, err := h.WeatherService.GetWeather(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.Logger.Error("error while get weather", zap.Error(err))

		return
	}

	c.JSON(http.StatusOK, weather)
}
