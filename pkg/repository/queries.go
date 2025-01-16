package repository

import (
	"crabigateur-api/pkg/api"
	"database/sql"
	"fmt"
)

const CardSelector = `
	SELECT c.card_id, c.word, c.translation, c.word_type, c.level, c.gender,
			con.tense, con.forms, con.irregular, f.gender, f.number, f.form
	FROM cards c
	LEFT JOIN conjugations con
	ON c.card_id = con.card_id
	LEFT JOIN forms f
	ON c.card_id = f.card_id
`;

func (s* storage) LessonsQuery(userId string, numLessons int) (*sql.Rows, error) {
	limit := ""
	if numLessons > 0 {
		limit = fmt.Sprintf("LIMIT %d", numLessons)
	} 

	lessonCardsQuery := fmt.Sprintf(`
		WITH PendingLessonCardIds AS (
			SELECT DISTINCT c.card_id
			FROM cards c
			WHERE c.level = (
				SELECT u.level 
				FROM Users u 
				WHERE u.user_id = $1
			)
			AND NOT EXISTS (
				SELECT 1
				FROM UserCardStatus ucs
				WHERE ucs.user_id = $1
				AND ucs.card_id = c.card_id
			)
			%s
		)

		%s
		WHERE c.card_id IN (SELECT card_id FROM PendingLessonCardIds)
		ORDER BY c.card_id;
	`, limit, CardSelector)
	
	return s.db.Query(lessonCardsQuery, userId)
}

func (s* storage) ReviewsQuery(userId string, numReviews int, sort []api.SortOrder) (*sql.Rows, error) {
	sortOrder := ""
	for _, value := range sort {
        switch value {
		case api.DateAsc:
			sortOrder += "pr.next_review_date ASC,"
		case api.DateDesc:
			sortOrder += "pr.next_review_date DESC,"
		case api.LevelAsc:
			sortOrder += "c.level ASC,"
		case api.LevelDesc:
			sortOrder += "c.level DESC,"
		}
    }

	limit := ""
	if numReviews > 0 {
		limit = fmt.Sprintf("LIMIT %d", numReviews)
	} 

	reviewsQuery := fmt.Sprintf(`
		WITH PendingReviews AS (
			SELECT card_id, next_review_date
			FROM usercardstatus
			WHERE user_id = $1 
			AND next_review_date < NOW()
			%s
		)
		
		%s
		JOIN PendingReviews pr ON c.card_id = pr.card_id
		ORDER BY %s c.card_id;
	`, limit, CardSelector, sortOrder)

	return s.db.Query(reviewsQuery, userId)
}

func (s* storage) ReviewsInsert(userId string, review api.Review) (sql.Result, error) {
	reviewsInsert := `
		INSERT INTO Reviews (user_id, card_id, review_date, success, previous_stage)
		VALUES (
			$1::VARCHAR, 
			$2, 
			$3,
			$4,
			(
				SELECT s.stage_id
				FROM SRSStages s
				JOIN UserCardStatus ucs ON s.stage_id = ucs.stage_id
				WHERE ucs.user_id = $1::VARCHAR AND ucs.card_id = $2
			)
		)
	`

	return s.db.Exec(reviewsInsert, userId, review.CardId, review.ReviewDate, review.Success)
}

func (s* storage) UserCardStatusInsert(userId string, review api.Review) (*sql.Row) {
	insertQuery := `
		WITH inserted AS (
			INSERT INTO UserCardStatus (user_id, card_id, stage_id, next_review_date)
			VALUES (
				$1, 
				$2, 
				1, 
				(
					SELECT $3::TIMESTAMPTZ + s.stage_interval
					FROM SRSStages s
					WHERE s.stage_id = 1
				)
			)
			RETURNING card_id
		)
		SELECT i.card_id, c.word, true AS success, 1 AS stage_id
		FROM inserted i
		JOIN Cards c ON i.card_id = c.card_id;
	`

	return s.db.QueryRow(insertQuery, userId, review.CardId, review.ReviewDate)

}

func (s* storage) UserCardStatusUpdate(userId string, review api.Review) (*sql.Row) {
	stageUpdate := `ucs.stage_id + 1`
	if(!*review.Success){
		stageUpdate = fmt.Sprintf(`
			GREATEST(
				1,
				ucs.stage_id - CEIL(
					ROUND(%d / 2.0, 1)
				) * (
					SELECT s.stage_penalty
					FROM SRSStages s
					WHERE s.stage_id = ucs.stage_id
				)
			)`, *review.IncorrectCount)
	}
	
	updateQuery := fmt.Sprintf(`
		WITH calculated_stage AS (
			SELECT %s 
				AS new_stage_id,
				ucs.card_id,
				ucs.user_id
			FROM UserCardStatus ucs
			WHERE ucs.user_id = $1 AND ucs.card_id = $2
		),
		updated AS (
			UPDATE UserCardStatus ucs
			SET stage_id = cs.new_stage_id,
				next_review_date = $3::TIMESTAMPTZ + s.stage_interval
			FROM calculated_stage cs
			JOIN SRSStages s
				ON cs.new_stage_id = s.stage_id
			WHERE ucs.user_id = cs.user_id AND ucs.card_id = cs.card_id
			RETURNING ucs.card_id, ucs.stage_id
		)

		SELECT u.card_id, c.word, $4 AS success, u.stage_id
		FROM updated u
		JOIN Cards c ON u.card_id = c.card_id;
	`, stageUpdate)

	return s.db.QueryRow(updateQuery, userId, review.CardId, review.ReviewDate, *review.Success)
}

func (s* storage) MostRecentReviewQuery(userId string, cardId int) (*sql.Row) {
	mostRecentReview := `
		SELECT r.card_id, c.word, r.success, ucs.stage_id
		FROM Reviews r
		JOIN UserCardStatus ucs
		ON r.user_id = ucs.user_id AND r.card_id = ucs.card_id
		JOIN Cards c
		ON r.card_id = c.card_id
		WHERE r.user_id = $1 AND r.card_id = $2
		ORDER BY r.review_date DESC
		LIMIT 1;
	`

	return s.db.QueryRow(mostRecentReview, userId, cardId)
}

func (s* storage) CardQuery(id int) (*sql.Rows, error) {
	cardQuery := fmt.Sprintf(`
		%s
		WHERE c.card_id = $1;
	`, CardSelector)
	
	return s.db.Query(cardQuery, id)
}
