package postgres

import (

	"github.com/Abdulazizxoshimov/Datagaze_Backend/entity"
	"github.com/jmoiron/sqlx"
)

type WeatherRepo struct {
	db *sqlx.DB
}

func NewWeatherRepo(db *sqlx.DB) *WeatherRepo {
	return &WeatherRepo{db: db}
}

func (w *WeatherRepo) GetWeatherByCity(city string) (*entity.Weather, error) {
	var weather entity.Weather
	query := `
		SELECT id, name, country, lat, lon, temp_c, temp_color, wind_kph, wind_color, cloud, cloud_color, created_at
		FROM weather
		WHERE name = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := w.db.QueryRow(query, city).Scan(
		&weather.ID, &weather.City, &weather.Country, &weather.Lat, &weather.Lon,
		&weather.TempC, &weather.TempColor, &weather.WindKph, &weather.WindColor,
		&weather.Cloud, &weather.CloudColor, &weather.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}

func (w *WeatherRepo) SaveWeather(weather *entity.Weather) error {
	query := `
		INSERT INTO weather (id, name, country, lat, lon, temp_c, temp_color, wind_kph, wind_color, cloud, cloud_color, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := w.db.Exec(query,
		weather.ID, weather.City, weather.Country, weather.Lat, weather.Lon,
		weather.TempC, weather.TempColor, weather.WindKph, weather.WindColor,
		weather.Cloud, weather.CloudColor, weather.CreatedAt,
	)
	return err
}
