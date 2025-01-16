package app

import "github.com/gin-gonic/gin"

func (s *Server) Routes() *gin.Engine {
	router := s.router

	// group all routes under /v1/api
	v1 := router.Group("/v1/api")
	{
		v1.GET("/status", s.ApiStatus())
		v1.GET("/lessons/:user_id", s.GetUserLessons())
		v1.GET("/reviews/:user_id", s.GetUserReviews())
		v1.POST("/lessons/:user_id", s.PostUserReviews())
		v1.PUT("/reviews/:user_id", s.PutUserReviews())
		v1.GET("/quiz_summary/:user_id", s.GetUserQuizSummary())
		v1.GET("/card/:card_id", s.GetCardById())
	}

	return router

}