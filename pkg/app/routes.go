package app

import "github.com/gin-gonic/gin"

func (s *Server) Routes() *gin.Engine {
	router := s.router

	// group all routes under /v1/api
	v1 := router.Group("/v1/api")
	{
		v1.GET("/status", s.ApiStatus())
		v1.GET("/lessons/:user_id", s.GetUserLessons())

		reviews := v1.Group("/reviews")
		{
			reviews.GET("/:user_id", s.GetUserReviews())
			reviews.POST("/:user_id", s.PostUserReviews())
			reviews.PUT("/:user_id", s.PutUserReviews())
		}

		v1.GET("/quiz_summary/:user_id", s.GetUserQuizSummary())
		v1.GET("/stats/:user_id", s.GetUserStats())

		card := v1.Group("/card")
		{
			card.GET("/:card_id", s.GetCardById())
			card.POST("", s.CreateCard())
			card.PUT("/:card_id", s.UpdateCard())
			card.DELETE("/:card_id", s.DeleteCard())
			card.GET("/search", s.SearchCards())
		}
	}

	return router

}
