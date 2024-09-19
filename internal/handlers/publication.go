package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alanmathiasen/aggregator-api/pkg/utils"

	"github.com/alanmathiasen/aggregator-api/internal/auth"
	"github.com/alanmathiasen/aggregator-api/internal/services"
	"github.com/alanmathiasen/aggregator-api/internal/views/dashboard"
	"github.com/alanmathiasen/aggregator-api/internal/views/discover"

	pub "github.com/alanmathiasen/aggregator-api/internal/views/publication"
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
	userID, ok := session.Values["userID"].(uint)
	if !ok {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	all, err := publication.GetAllPublications(r.Context(), userID)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"publications": all})
}

// POST /publications
func CreatePublication(w http.ResponseWriter, r *http.Request) {
	var publicationData services.Publication
	err := json.NewDecoder(r.Body).Decode(&publicationData)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}

	err = publicationData.Validate()
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	publicationCreated, err := publication.CreatePublication(r.Context(), publicationData)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, publicationCreated)
}

// PUT /publications/:id
func UpdatePublication(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var publicationData services.Publication
	err := json.NewDecoder(r.Body).Decode(&publicationData)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}

	publicationUpdated, err := publication.UpdatePublication(id, publicationData)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, publicationUpdated)
}

// DELETE /publications/:id
func DeletePublication(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := publication.DeletePublication(id)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "succesfully deleted"})
}

func GetAllPublicationsHTML(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.SessionKey).(*sessions.Session)
	userID, ok := session.Values["userID"].(uint)
	if userID == 0 || !ok {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	all, err := publication.GetAllPublications(r.Context(), userID)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}

	component := discover.DiscoverPage(all)

	err = component.Render(r.Context(), w)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}
}

func UpsertPublicationFollowHTML(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	publicationIDUint, err := utils.StringToUint(id)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}

	status := r.FormValue("status")

	chapterID := r.FormValue("chapter_id")
	chapterIDUint, err := utils.StringToUint(chapterID)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}

	session := r.Context().Value(auth.SessionKey).(*sessions.Session)
	userID, ok := session.Values["userID"].(uint)
	if !ok {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	publicationFollow := &services.UserPublicationFollows{
		PublicationID: publicationIDUint,
		ChapterID:     chapterIDUint,
		Status:        status,
		UserID:        userID,
	}

	_, err = userPublicationFollows.UpsertUserPublicationFollows(r.Context(), *publicationFollow)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}
	p, err := publication.GetPublicationById(r.Context(), id, userID)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return

	}
	component := pub.DashboardPublication(p)
	err = component.Render(r.Context(), w)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}
}

func DeletePublicationFollowHTML(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	publicationIDUint, err := utils.StringToUint(id)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}

	session := r.Context().Value(auth.SessionKey).(*sessions.Session)
	userID, ok := session.Values["userID"].(uint)
	if !ok {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	err = userPublicationFollows.DeleteUserPublicationFollow(r.Context(), publicationIDUint, userID)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}

	component := pub.Publication(&publication)
	err = component.Render(r.Context(), w)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}
}

func DashboardHTML(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(auth.SessionKey).(*sessions.Session)
	userID, ok := session.Values["userID"].(uint)
	if userID == 0 || !ok {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	publications, err := publication.GetAllPublications(r.Context(), userID)
	if err != nil {
		fmt.Println("Error", err)
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}

	component := dashboard.Page(publications)
	err = component.Render(r.Context(), w)
	if err != nil {
		utils.MessageLogs.ErrorLog.Println(err)
		return
	}
}
