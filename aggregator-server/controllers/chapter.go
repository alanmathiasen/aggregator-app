package controllers

import (
	"fmt"
	"net/http"

	"github.com/alanmathiasen/aggregator-api/helpers"
	"github.com/alanmathiasen/aggregator-api/services"
	"github.com/go-chi/chi"
)

var chapter services.Chapter

func GetAllChaptersByPublicationID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	chapters, err := chapter.GetAllChaptersByPublicationID(r.Context(), id)
	if len(chapters) == 0 {
		fmt.Printf("no chapters found")

	}
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
	}
	helpers.WriteJSON(w, http.StatusOK, chapters)
}

func CreateChapterForPublication(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var chapterData services.Chapter
	err := helpers.ReadJSON(w, r, &chapterData)
	err = chapterData.Validate()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	err = chapterData.Validate()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	err = chapter.CreateChapterForPublication(r.Context(), id, &chapterData)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
	}
	helpers.WriteJSON(w, http.StatusCreated, chapterData)
}
