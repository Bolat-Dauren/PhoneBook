//cmd/web/main/authHandlers.go

package main

import (
	"PhoneBook_AP/pkg/drivers"
	"PhoneBook_AP/pkg/models"
	_ "fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Создание экземпляра логгера
var log = logrus.New()

func init() {
	// Настройка форматтера для логгера
	log.SetFormatter(&logrus.JSONFormatter{})
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		hashedPassword, err := models.GetHashedPassword(username)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"user":  username,
			}).Error("Error getting hashed password")
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				log.WithFields(logrus.Fields{
					"error": err,
					"user":  username,
				}).Warning("Invalid password")
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}
			log.WithFields(logrus.Fields{
				"error": err,
				"user":  username,
			}).Error("Error comparing passwords")
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Логирование успешного входа
		log.WithFields(logrus.Fields{
			"user": username,
		}).Info("Login successful")

		http.Redirect(w, r, "/application", http.StatusSeeOther)
		return
	}
	tmpl.ExecuteTemplate(w, "login.tmpl", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		var exists bool
		err := drivers.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error":    err,
				"username": username,
			}).Error("Error checking for existing username")
			http.Error(w, "Error checking for existing username", http.StatusInternalServerError)
			return
		}
		if exists {
			log.WithFields(logrus.Fields{
				"username": username,
			}).Warning("Username already exists")
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		err = drivers.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"email": email,
			}).Error("Error checking for existing email")
			http.Error(w, "Error checking for existing email", http.StatusInternalServerError)
			return
		}
		if exists {
			log.WithFields(logrus.Fields{
				"email": email,
			}).Warning("Email already exists")
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		user := models.User{Username: username, Password: password, Email: email}
		err = models.CreateUser(user)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"user":  user,
			}).Error("Error creating user")
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		// Логирование успешной регистрации
		log.WithFields(logrus.Fields{
			"user":  username,
			"email": email,
		}).Info("Registration successful")

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl.ExecuteTemplate(w, "register.tmpl", nil)
}
