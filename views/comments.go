// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gorilla/csrf"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/templates"
)

func editComment(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)

	user, _ := r.Context().Value(CtxUserKey).(*models.User)
	domain := r.Context().Value(ctxDomain).(*models.Domain)

	topicIDStr := chi.URLParam(r, "topicID")
	categoryIDStr := chi.URLParam(r, "categoryID")
	commentIDStr := chi.URLParam(r, "commentID")
	quoteIDStr := r.URL.Query().Get("quote")

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	category := models.GetCategoryByID(categoryID)
	if category == nil || category.DomainID != domain.DomainID || category.ArchivedAt.Valid {
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

	var comment *models.Comment
	if commentIDStr != "" {
		// Edit comment
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		comment = models.GetCommentByID(commentID)
		if comment == nil || comment.TopicID != topic.TopicID {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
	}

	var quote string
	if comment == nil {
		if quoteIDStr == "topic" {
			quote = "> " + strings.ReplaceAll(topic.Content, "\n", "\n> ")
		} else {
			if quoteCommentID, err := strconv.Atoi(quoteIDStr); err == nil {
				if quoteComment := models.GetCommentByID(quoteCommentID); quoteComment != nil {
					if quoteComment.TopicID == topic.TopicID {
						quote = "> " + strings.ReplaceAll(quoteComment.Content, "\n", "\n> ")
					}
				}
			}
		}
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		action := r.PostFormValue("action")
		content := r.PostFormValue("content")
		isSticky := r.PostFormValue("is_sticky") == "1"
		if !(user.IsSuperAdmin || user.IsSuperMod) {
			isSticky = false
			if topic != nil {
				isSticky = topic.IsSticky
			}
		}
		if action == "Submit" {
			if topic.IsReadOnly {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			newCommentID := models.CreateComment(topic.TopicID, user.UserID, content, isSticky)
			if newCommentID >= 0 {
				http.Redirect(w, r, basePath+"categories/"+strconv.Itoa(category.CategoryID)+
					"/topics/"+strconv.Itoa(topic.TopicID)+
					"#comment-"+strconv.Itoa(newCommentID),
					http.StatusSeeOther)
				return
			}
		} else if action == "Update" {
			models.UpdateCommentByID(comment.CommentID, content, isSticky)
			http.Redirect(w, r, basePath+"categories/"+strconv.Itoa(category.CategoryID)+
				"/topics/"+strconv.Itoa(topic.TopicID)+
				"#comment-"+strconv.Itoa(comment.CommentID),
				http.StatusSeeOther)
			return
		} else if action == "Delete" {
			models.DeleteCommentByID(comment.CommentID, comment.UserID, topic.TopicID)
			http.Redirect(w, r, basePath+"categories/"+strconv.Itoa(category.CategoryID)+
				"/topics/"+strconv.Itoa(topic.TopicID),
				http.StatusSeeOther)
			return
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	templates.CommentEdit.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		DomainField:      domain,
		UserField:        user,
		"Category":       category,
		"Topic":          topic,
		"Comment":        comment,
		"Quote":          quote,
	})

}
