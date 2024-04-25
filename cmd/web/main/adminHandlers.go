//cmd/web/main/adminHandlers.go

package main

import (
	"io/ioutil"
	"net/http"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := checkAdminAuthentication(r)

	var jsonData string
	if isAuthenticated {
		data, err := loadJSONData("phonebook.json")
		if err != nil {
			http.Error(w, "Error loading JSON data", http.StatusInternalServerError)
			return
		}
		jsonData = data
	}

	tmpl.ExecuteTemplate(w, "admin.tmpl", struct {
		Authenticated bool
		JSONData      string
	}{
		Authenticated: isAuthenticated,
		JSONData:      jsonData,
	})
}

func adminLoginHandler(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")
	isAuthenticated := authenticateAdmin(username, password)

	if isAuthenticated {
		http.SetCookie(w, &http.Cookie{
			Name:  "admin_authenticated",
			Value: "true",
		})
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func adminSaveJSONHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAdminAuthentication(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	jsonData := r.FormValue("json_data")
	err := saveJSONData("phonebook.json", jsonData)
	if err != nil {
		http.Error(w, "Error saving JSON data", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func adminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "admin_authenticated",
		Value:  "",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func checkAdminAuthentication(r *http.Request) bool {
	cookie, err := r.Cookie("admin_authenticated")
	if err != nil {
		return false
	}
	return cookie.Value == "true"
}

func authenticateAdmin(username, password string) bool {
	return username == "admin" && password == "0000"
}

func loadJSONData(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func saveJSONData(filename, data string) error {
	return ioutil.WriteFile(filename, []byte(data), 0644)
}
