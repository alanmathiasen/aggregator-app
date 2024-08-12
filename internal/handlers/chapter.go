package handlers

import (
	"net/http"

	"github.com/alanmathiasen/aggregator-api/internal/services"
	"github.com/alanmathiasen/aggregator-api/pkg/utils"
	"github.com/go-chi/chi"
)

var chapter services.Chapter

func GetAllChaptersByPublicationID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	chapters, err := chapter.GetAllChaptersByPublicationID(r.Context(), id)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		utils.ErrorJSON(w, err, http.StatusBadRequest)
	}
	utils.WriteJSON(w, http.StatusOK, chapters)
}

func CreateChapterForPublication(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var chapterData services.Chapter

	err := utils.ReadJSON(w, r, &chapterData)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = chapterData.Validate()
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = chapter.CreateChapterForPublication(r.Context(), id, &chapterData)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, chapterData)
}
