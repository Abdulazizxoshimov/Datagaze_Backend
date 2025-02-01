package entity

import "time"

type Weather struct {
	ID        string    `db:"id"`
	City      string    `db:"name"`
	Country   string    `db:"country"`
	Lat       float64   `db:"lat"`
	Lon       float64   `db:"lon"`
	TempC     float64   `db:"temp_c"`
	TempColor string    `db:"temp_color"`
	WindKph   float64   `db:"wind_kph"`
	WindColor string    `db:"wind_color"`
	Cloud     int       `db:"cloud"`
	CloudColor string   `db:"cloud_color"`
	CreatedAt time.Time `db:"created_at"`
}
