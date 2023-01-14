package handlers

import (
	"encoding/json"
	"net/http"
	articledto "server/dto/article"
	dto "server/dto/result"
	"server/models"
	"server/repositories"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

var path_file = "http://localhost:5000/uploads/"

type handlerArticle struct {
	ArticleRepository repositories.ArticleRepository
}

func HandlerArticle(ArticleRepository repositories.ArticleRepository) *handlerArticle {
	return &handlerArticle{ArticleRepository}
}

func (h *handlerArticle) FindArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	article, err := h.ArticleRepository.Finduser()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	for i, d := range article {
		article[i].Image = path_file + d.Image
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: article}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerArticle) GetArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	article, err := h.ArticleRepository.GetArticle(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	article.Image = path_file + article.Image
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: article}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerArticle) AddArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	dataContex := r.Context().Value("dataFile") // add this code
	filename := dataContex.(string)             // add this code
	request := articledto.ArticleRequest{
		Title:       r.FormValue("title"),
		Image:       filename,
		Description: r.FormValue("description"),
		UserID:      userId,
	}
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	article := models.Article{
		Title:       request.Title,
		Image:       request.Image,
		Description: request.Description,
		UserID:      request.UserID,
	}

	data, err := h.ArticleRepository.AddArticle(article)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	datas, err := h.ArticleRepository.GetArticle(data.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: datas}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerArticle) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	Article, err := h.ArticleRepository.GetArticle(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.ArticleRepository.DeleteArticle(Article)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)
}

func articleConvertResponse(u models.Article) articledto.ArticleResponse {
	return articledto.ArticleResponse{
		Title:       u.Title,
		Image:       u.Image,
		Description: u.Description,
		UserID:      u.UserID,
	}
}
