package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"iss-model-backend/internal/services"
	"iss-model-backend/internal/utils"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// @Summary Get All Posts
// @Tags Blog
// @Produce json
// @Success 200 {array} models.Post
// @Router /blog/posts [get]
func (h *PostHandler) HandleGetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get posts", err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, posts)
}

// @Summary Get Post by ID
// @Tags Blog
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} models.Post
// @Router /blog/posts/{id} [get]
func (h *PostHandler) HandleGetPostByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		// Obsługa błędu "nie znaleziono"
		utils.SendErrorResponse(w, http.StatusNotFound, "Post not found", err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, post)
}

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// @Summary Create New Post
// @Security ApiKeyAuth
// @Tags Blog (Admin)
// @Accept json
// @Produce json
// @Param request body CreatePostRequest true "Post data"
// @Success 201 {object} models.Post
// @Router /admin/blog/posts [post]
func (h *PostHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	post, err := h.postService.CreatePost(req.Title, req.Content)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create post", err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusCreated, post)
}

// @Summary Update Post
// @Security ApiKeyAuth
// @Tags Blog (Admin)
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param request body CreatePostRequest true "Post data"
// @Success 200 {object} models.Post
// @Router /admin/blog/posts/{id} [put]
func (h *PostHandler) HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	post, err := h.postService.UpdatePost(uint(id), req.Title, req.Content)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to update post", err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, post)
}

// @Summary Delete Post
// @Security ApiKeyAuth
// @Tags Blog (Admin)
// @Param id path int true "Post ID"
// @Success 204 "No Content"
// @Router /admin/blog/posts/{id} [delete]
func (h *PostHandler) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	if err := h.postService.DeletePost(uint(id)); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to delete post", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
