//cmd/web/main/main.go

package main

import (
	"PhoneBook_AP/pkg/drivers"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"path/filepath"
	_ "time"
)

var limiter = rate.NewLimiter(1, 3)

func main() {
	drivers.InitDB("user=postgres dbname=postgres password=0000 sslmode=disable\n")
	mux := http.NewServeMux()

	// Используем Rate Limiting Middleware для всех обработчиков
	http.Handle("/", rateLimitMiddleware(mux))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/register", register)
	mux.HandleFunc("/login", login)

	mux.HandleFunc("/application", application)

	mux.HandleFunc("/city/", searchPageHandler)
	mux.HandleFunc("/search/", searchHandler)

	mux.HandleFunc("/city/hospitals/", searchHospitalsPageHandler)
	mux.HandleFunc("/search/hospitals/", searchHospitalsHandler)

	mux.HandleFunc("/city/schools/", searchSchoolsPageHandler)
	mux.HandleFunc("/search/schools/", searchSchoolsHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.HandleFunc("/admin", adminHandler)
	mux.HandleFunc("/admin/login", adminLoginHandler)
	mux.HandleFunc("/admin/logout", adminLogoutHandler)
	mux.HandleFunc("/admin/save-json", adminSaveJSONHandler)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	log.Println("Starting server on :5000")
	err := http.ListenAndServe(":5000", nil)
	log.Fatal(err)
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func init() {
	logFilePath := "logs/application.log"
	logDir := filepath.Dir(logFilePath)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			panic(err)
		}
	}
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(logFile)
}
