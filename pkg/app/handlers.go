package app

import (
	"crabigateur-api/pkg/api"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (s *Server) ApiStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		response := map[string]string{
			"data": "crabigateur API running smoothly",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) GetUserLessons() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var pathParams api.UserPath
		if err := c.ShouldBindUri(&pathParams); err != nil {
			log.Printf("handler error: invalid uri params: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}

		var queryParams api.QueryParams
		if err := c.ShouldBindWith(&queryParams, binding.Query); err != nil {
			log.Printf("handler error: invalid query params: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
			return
		}

		cards, ids, err := s.userService.LessonCards(pathParams.UserId, queryParams.NumCards)
		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		response := map[string]interface{}{
			"data": map[string]interface{}{
				"cards":    cards,
				"card_ids": ids,
				"total":    len(cards),
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) GetUserReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var pathParams api.UserPath
		if err := c.ShouldBindUri(&pathParams); err != nil {
			log.Printf("handler error: invalid uri params: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}

		var queryParams api.QueryParams
		if err := c.ShouldBindWith(&queryParams, binding.Query); err != nil {
			log.Printf("handler error: invalid query params: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
			return
		}

		review, err := s.userService.ReviewCard(pathParams.UserId, queryParams.FirstReview, queryParams.Sort)
		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		} else if review.IsEmpty() {
			c.JSON(http.StatusNoContent, gin.H{"data": "No cards to review"})
			return
		}

		response := map[string]interface{}{
			"data": review,
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) PostUserReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var pathParams api.UserPath
		if err := c.ShouldBindUri(&pathParams); err != nil {
			log.Printf("handler error: invalid uri params: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}

		var list api.QuizList
		err := c.ShouldBindJSON(&list)
		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review format"})
			return
		}

		result, err := s.userService.AddReviews(pathParams.UserId, list.CardIds)
		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		response := map[string]interface{}{
			"data": result,
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) PutUserReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var pathParams api.UserPath
		if err := c.ShouldBindUri(&pathParams); err != nil {
			log.Printf("handler error: invalid uri params: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}

		var review api.Review
		err := c.ShouldBindJSON(&review)
		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review format"})
			return
		}

		result, err := s.userService.UpdateReview(pathParams.UserId, review)
		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		response := map[string]interface{}{
			"data": result,
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) GetUserQuizSummary() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var pathParams api.UserPath
		if err := c.ShouldBindUri(&pathParams); err != nil {
			log.Printf("handler error: invalid uri params: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}

		var queryParams api.QueryParams
		if err := c.ShouldBindWith(&queryParams, binding.Query); err != nil {
			log.Printf("handler error: invalid query params: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
			return
		}

		summary, err := s.userService.GetQuizSummary(pathParams.UserId, queryParams.NumCards)
		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		response := map[string]interface{}{
			"data": summary,
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) GetCardById() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var pathParams api.CardPath
		if err := c.ShouldBindUri(&pathParams); err != nil {
			log.Printf("handler error: invalid uri params: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card_id"})
			return
		}

		card, err := s.cardService.GetCardById(pathParams.CardId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		} else if card.IsEmpty() {
			c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
			return
		}

		response := map[string]interface{}{
			"data": card,
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) CreateCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		var card api.Card
		err := c.ShouldBindJSON(&card)
		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review format"})
			return
		}

		result, err := s.cardService.CreateCard(card)
		if err != nil {
			log.Printf("service error: %v", err)
			if strings.Contains(err.Error(), "duplicate card") {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "A card with this word already exists"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}

		response := map[string]interface{}{
			"data": result,
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) UpdateCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pathParams api.CardPath
		if err := c.ShouldBindUri(&pathParams); err != nil {
			log.Printf("handler error: invalid uri params: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card_id"})
			return
		}

		var card api.Card
		err := c.ShouldBindJSON(&card)
		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card format"})
			return
		}

		result, err := s.cardService.UpdateCard(pathParams.CardId, card)
		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		response := map[string]interface{}{
			"data": result,
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) DeleteCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var pathParams api.CardPath
		if err := c.ShouldBindUri(&pathParams); err != nil {
			log.Printf("handler error: invalid uri params: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card_id"})
			return
		}

		err := s.cardService.DeleteCard(pathParams.CardId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		} // TODO: especify errors ??

		response := map[string]interface{}{
			"data": "Card deleted successfully",
		}

		c.JSON(http.StatusOK, response)
	}
}
