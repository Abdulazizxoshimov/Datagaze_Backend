package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Abdulazizxoshimov/Datagaze_Backend/config"
	"github.com/Abdulazizxoshimov/Datagaze_Backend/entity"
	repo "github.com/Abdulazizxoshimov/Datagaze_Backend/internal/repo/postgres"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
)

type WeatherService struct {
	weatherRepo *repo.WeatherRepo
	apiKey      string
}

func NewWeatherService(weatherRepo *repo.WeatherRepo, cfg config.Config) *WeatherService {
	return &WeatherService{
		weatherRepo: weatherRepo,
		apiKey:      cfg.WeatherAPIKey, // API kalitni configdan olish
	}
}

//Ob-havo ma'lumotlarini API dan olib, bazaga yozish
func (s *WeatherService) FetchAndStoreWeather(city string) error {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", s.apiKey, city)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch weather: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("weather API returned status: %s", resp.Status)
	}

	var data struct {
		Location struct {
			Name    string  `json:"name"`
			Country string  `json:"country"`
			Lat     float64 `json:"lat"`
			Lon     float64 `json:"lon"`
		} `json:"location"`
		Current struct {
			TempC   float64 `json:"temp_c"`
			WindKph float64 `json:"wind_kph"`
			Cloud   int     `json:"cloud"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode weather response: %v", err)
	}

	weather := entity.Weather{
		ID:         uuid.NewString(),
		City:       data.Location.Name,
		Country:    data.Location.Country,
		Lat:        data.Location.Lat,
		Lon:        data.Location.Lon,
		TempC:      data.Current.TempC,
		TempColor:  GetTempColor(data.Current.TempC),  // Temperaturaga qarab rang belgilash
		WindKph:    data.Current.WindKph,
		WindColor:  GetWindColor(data.Current.WindKph), // Shamol tezligiga qarab rang berish
		Cloud:      data.Current.Cloud,
		CloudColor: GetCloudColor(data.Current.Cloud), // Bulutlilik uchun rang
		CreatedAt:  time.Now(),
	}
	pp.Println(weather)

	//  Ma'lumotni bazaga saqlash
	return s.weatherRepo.SaveWeather(&weather)
}

//  Bazadan ob-havo ma'lumotlarini olish
func (s *WeatherService) GetWeather(city string) (*entity.Weather, error) {
	if s.weatherRepo == nil{
		pp.Println("salom mana shu joyda errror")
	}
	return s.weatherRepo.GetWeatherByCity(city)
}
