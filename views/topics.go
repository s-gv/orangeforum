// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/csrf"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/templates"
)

func getTopicList(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)

	user, _ := r.Context().Value(CtxUserKey).(*models.User)
	domain, _ := r.Context().Value(ctxDomain).(*models.Domain)

	categoryIDStr := chi.URLParam(r, "categoryID")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	category := models.GetCategoryByID(categoryID)
	if category == nil || category.DomainID != domain.DomainID {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	beforeStr := r.URL.Query().Get("before")
	before, err := strconv.ParseInt(beforeStr, 10, 64)
	if err != nil {
		before = time.Now().UnixNano()
	}

	topics := models.GetTopicsByCategoryID(category.CategoryID, time.Unix(0, before))

	moreDt := time.Now()
	if len(topics) > 0 {
		moreDt = topics[len(topics)-1].ActivityAt
	}
	templates.TopicList.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		UserField:        user,
		"Category":       category,
		"Topics":         topics,
		"MoreDt":         moreDt.UnixNano(),
	})
}

func editTopic(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)

	user, _ := r.Context().Value(CtxUserKey).(*models.User)
	domain, _ := r.Context().Value(ctxDomain).(*models.Domain)

	topicIDStr := chi.URLParam(r, "topicID")
	categoryIDStr := chi.URLParam(r, "categoryID")

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	category := models.GetCategoryByID(categoryID)
	if category == nil || category.DomainID != domain.DomainID {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var topic *models.Topic
	if topicIDStr != "" {
		topicID, err := strconv.Atoi(topicIDStr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		topic = models.GetTopicByID(topicID)
		if topic.CategoryID != category.CategoryID {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		action := r.PostFormValue("action")
		title := r.PostFormValue("title")
		content := r.PostFormValue("content")
		if len(title) < 2 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if action == "Submit" {
			newTopicID := models.CreateTopic(category.CategoryID, user.UserID, title, content)
			if newTopicID >= 0 {
				http.Redirect(w, r, basePath+
					"/categories/"+strconv.Itoa(category.CategoryID)+
					"/topics/"+strconv.Itoa(newTopicID),
					http.StatusSeeOther)
				return
			}
		} else if action == "Update" {
			models.UpdateTopicByID(topic.TopicID, title, content)
			http.Redirect(w, r, basePath+
				"/categories/"+strconv.Itoa(category.CategoryID)+
				"/topics/"+strconv.Itoa(topic.TopicID),
				http.StatusSeeOther)
			return
		} else if action == "Delete" {
			models.DeleteTopicByID(topic.TopicID)
			http.Redirect(w, r, basePath+"/categories/"+strconv.Itoa(category.CategoryID), http.StatusSeeOther)
			return
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	templates.TopicEdit.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		UserField:        user,
		"Category":       category,
		"Topic":          topic,
	})
}

func getTopic(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)

	user, _ := r.Context().Value(CtxUserKey).(*models.User)
	domain, _ := r.Context().Value(ctxDomain).(*models.Domain)

	topicIDStr := chi.URLParam(r, "topicID")
	categoryIDStr := chi.URLParam(r, "categoryID")

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	category := models.GetCategoryByID(categoryID)
	if category == nil || category.DomainID != domain.DomainID {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	topic := models.GetTopicByID(topicID)
	if topic.CategoryID != category.CategoryID {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	comments := models.GetCommentsByTopicID(topic.TopicID)

	templates.Topic.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		UserField:        user,
		"Category":       category,
		"Topic":          topic,
		"Comments":       comments,
	})
}
