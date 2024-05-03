package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/alanmathiasen/aggregator-api/auth"
	"github.com/alanmathiasen/aggregator-api/helpers"
	"github.com/alanmathiasen/aggregator-api/services"
	"github.com/alanmathiasen/aggregator-api/view/dashboard"
	"github.com/alanmathiasen/aggregator-api/view/discover"

	pub "github.com/alanmathiasen/aggregator-api/view/publication"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
)

var (
	publication            services.Publication
	userPublicationFollows services.UserPublicationFollows
)

// GET /publications
func GetAllPublications(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.SessionKey).(*sessions.Session)
	user, ok := session.Values["user"].(*services.User);
	if user == nil || !ok {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	all, err := publication.GetAllPublications(r.Context(), user.ID)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"publications": all})
}

// GET /publications/:id
// func GetPublicationById(w http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")
// 	publication, err := publication.GetPublicationById(id)
// 	if err != nil {
// 		helpers.MessageLogs.ErrorLog.Println(err)
// 		return
// 	}

// 	helpers.WriteJSON(w, http.StatusOK, publication)
// }

// POST /publications
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

	publicationCreated, err := publication.CreatePublication(r.Context(), publicationData)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, publicationCreated)
}

// PUT /publications/:id
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

// DELETE /publications/:id
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
	session := r.Context().Value(auth.SessionKey).(*sessions.Session)
	user, ok := session.Values["user"].(*services.User);
	if user == nil || !ok {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	all, err := publication.GetAllPublications(r.Context(), user.ID)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	component := discover.DiscoverPage(all)

	err = component.Render(r.Context(), w)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
}

func GetUserPublicationsHTML(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.SessionKey).(*sessions.Session)
	user, ok := session.Values["user"].(*services.User);
	if user == nil || !ok {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	all, err := publication.GetAllPublications(r.Context(), user.ID)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	component := discover.DiscoverPage(all)
	err = component.Render(r.Context(), w)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
}

// func GetPublicationHTML(w http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")
// 	p, err := publication.GetPublicationById(id)
// 	if err != nil {
// 		helpers.MessageLogs.ErrorLog.Println(err)
// 		return
// 	}
// 	component := dashboard.Publication(*p)
// 	err = component.Render(r.Context(), w)
// 	if err != nil {
// 		helpers.MessageLogs.ErrorLog.Println(err)
// 		return
// 	}
// }

func UpsertPublicationFollowHTML(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	publicationIDUint, err := helpers.StringToUint(id) 
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	
	queryParams := r.URL.Query()
	status := queryParams.Get("status")
	
	chapterID := queryParams.Get("chapter_id")
	chapterIDUint, err := helpers.StringToUint(chapterID)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}


	session := r.Context().Value(auth.SessionKey).(*sessions.Session)
	user, ok := session.Values["user"].(*services.User);
	if user == nil || !ok {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	publicationFollow := &services.UserPublicationFollows{
		PublicationID: publicationIDUint,
		ChapterID: chapterIDUint,
		Status: status,
		UserID: user.ID,
	}
	
	_, err = userPublicationFollows.UpsertUserPublicationFollows(r.Context(), *publicationFollow)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	p, err := publication.GetPublicationById(id, user.ID)
	
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return

	}
	component := pub.Publication(p)
	err = component.Render(r.Context(), w)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
}

func DeletePublicationFollowHTML(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	publicationIDUint, err := helpers.StringToUint(id) 
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	session := r.Context().Value(auth.SessionKey).(*sessions.Session)
	user, ok := session.Values["user"].(*services.User);
	if user == nil || !ok {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	err = userPublicationFollows.DeleteUserPublicationFollow(r.Context(), publicationIDUint, user.ID)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}


	component := pub.DeletedPublication()
	err = component.Render(r.Context(), w)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
}

func DashboardHTML(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.SessionKey).(*sessions.Session)
	user, ok := session.Values["user"].(*services.User);
	if user == nil || !ok {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	publications, err := publication.GetAllPublications(r.Context(), user.ID)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	component := dashboard.Page(publications)
	err = component.Render(r.Context(), w)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

}