package app

import (
	"crabigateur-api/pkg/api"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (s *Server) ApiStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		response := map[string]string{
			"data":   "crabigateur API running smoothly",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) GetUserLessons() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var pathParams api.PathParams
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

		cards, err := s.userService.LessonCards(pathParams.UserId, queryParams.NumLessons)
		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		response := map[string]interface{}{
			"data": map[string]interface{}{
				"cards": cards,
				"total": len(cards),
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) GetUserReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var pathParams api.PathParams
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

		reviews, err := s.userService.ReviewCards(pathParams.UserId, queryParams.NumReviews, queryParams.Sort)
		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		response := map[string]interface{}{
			"data":   map[string]interface{}{
				"cards": reviews,
				"total": len(reviews),
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) PostUserReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var pathParams api.PathParams
		if err := c.ShouldBindUri(&pathParams); err != nil {
			log.Printf("handler error: invalid uri params: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}

		var reviews api.Reviews
		err := c.ShouldBindJSON(&reviews)
		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review format"})
			return
		}

		results, err := s.userService.AddReviews(pathParams.UserId, reviews.Reviews)
		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		response := map[string]interface{}{
			"data":   results,
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) PutUserReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var pathParams api.PathParams
		if err := c.ShouldBindUri(&pathParams); err != nil {
			log.Printf("handler error: invalid uri params: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}

		var reviews api.Reviews
		err := c.ShouldBindJSON(&reviews)
		if err != nil {
			log.Printf("handler error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review format"})
			return
		}

		results, err := s.userService.UpdateReviews(pathParams.UserId, reviews.Reviews)
		if err != nil {
			log.Printf("service error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		response := map[string]interface{}{
			"data":   results,
		}

		c.JSON(http.StatusOK, response)
	}
}