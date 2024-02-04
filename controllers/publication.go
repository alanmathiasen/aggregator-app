package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/alanmathiasen/aggregator-api/helpers"
	"github.com/alanmathiasen/aggregator-api/services"
	"github.com/alanmathiasen/aggregator-api/view/dashboard"
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

	err = publicationData.Validate()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	publicationCreated, err := publication.CreatePublication(publicationData) 
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, publicationCreated)
}

//PUT /publications/:id
func UpdatePublication(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var publicationData services.Publication
	err := json.NewDecoder(r.Body).Decode(&publicationData)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	publicationUpdated, err := publication.UpdatePublication(id, publicationData)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, publicationUpdated)
}

//DELETE /publications/:id
func DeletePublication(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := publication.DeletePublication(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"message": "succesfully deleted"})
}


func GetAllPublicationsHTML(w http.ResponseWriter, r *http.Request) {
	all, err := publication.GetAllPublications()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

    component := dashboard.DashboardPublications(all)

    err = component.Render(r.Context(), w)
    if err != nil {
				helpers.MessageLogs.ErrorLog.Println(err)
				return
    }

}

func GetPublicationHTML(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := publication.GetPublicationById(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	component := dashboard.Publication(*p)
	err = component.Render(r.Context(), w)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
}


func GetDiv(w http.ResponseWriter, r *http.Request) {	
	// id := chi.URLParam(r, "id")
	// p, err := publication.GetPublicationById(id)
	// if err != nil {
	// 	helpers.MessageLogs.ErrorLog.Println(err)
	// 	return
	// }
	// pl := services.PublicationLink{
		
	// }
	var links services.PublicationSource
	fmt.Printf(links.ID)
	tmpl, err := template.ParseFiles("templates/test.html")
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	
	err = tmpl.Execute(w, nil)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
}