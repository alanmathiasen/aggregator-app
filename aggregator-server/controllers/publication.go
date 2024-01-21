package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/alanmathiasen/aggregator-api/helpers"
	"github.com/alanmathiasen/aggregator-api/services"
	"github.com/go-chi/chi"
)

var publication services.Publication 

//GET /publications
func GetAllPublications(w http.ResponseWriter, r *http.Request) {
	all, err := publication.GetAllPublications()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"publications": all})
}

//GET /publications/:id
func GetPublicationById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	publication, err := publication.GetPublicationById(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, publication)
}

//POST /publications
func CreatePublication(w http.ResponseWriter, r *http.Request) {
	var publicationData services.Publication
	err := json.NewDecoder(r.Body).Decode(&publicationData)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	publicationCreated, err := publication.CreatePublication(publicationData) 
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusCreated, publicationCreated)
}

func DeletePublication(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := publication.DeletePublication(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusNoContent, "")
}