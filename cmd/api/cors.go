package main

import (
	"github.com/rs/cors"
	"net/http"
)

func CorsSettings() *cors.Cors {
	c := cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodPost, http.MethodGet, http.MethodDelete, http.MethodPatch,
		},
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:8081",
			"exp://192.168.0.16:8081",
			"http://192.168.0.16:8081",
			"https://192.168.0.16:8081",
			"192.168.0.16:8081",
			"http://127.0.0.1:5500",
			"http://10.0.114.182:8081",
			"http://172.20.10.2:8081",
			"http://172.20.10.2",
			"http://127.0.0.1:5501",
			"http://192.168.154.15",
			"http://192.168.154.15:8081",
			"http://192.168.2.15:8081",
			"http://192.168.0.10:8081",
			"http://192.168.40.15:8081",
			"http://192.168.0.15:8081",
			"http://172.20.10.2:8081",
			"http://172.20.10.2",
			"http://172.20.10.2",
		},
		AllowCredentials: true,
		AllowedHeaders: []string{
			"Access-Control-Allow-Origin",
			"Content-Type",
			"Authorization",
		},
		OptionsPassthrough: true,
		ExposedHeaders: []string{
			"Access-Control-Allow-Origin",
			"Content-Type",
		},
		Debug: false,
	})
	return c
}
