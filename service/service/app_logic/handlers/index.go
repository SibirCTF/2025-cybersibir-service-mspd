package handlers

import (
	"html/template"
	"net/http"
	"sibir2025/service/app_logic/utils"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles("templates/index.html") // Парсим HTML-шаблон
	if err != nil {
		utils.DropError(w, r, err, http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil) // Выполняем шаблон и передаем его в ResponseWriter
	if err != nil {
		utils.DropError(w, r, err, http.StatusInternalServerError)
		return
	}
}
