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

func (s* storage) ReviewQuery(userId string, firstReview bool, sort []api.SortOrder) (*sql.Rows, error) {
	sortOrder := ""
	for _, value := range sort {
        switch value {
		case api.DateAsc:
			sortOrder += "ucs.next_review_date ASC,"
		case api.DateDesc:
			sortOrder += "ucs.next_review_date DESC,"
		case api.LevelAsc:
			sortOrder += "c.level ASC,"
		case api.LevelDesc:
			sortOrder += "c.level DESC,"
		}
    }

	// Remove trailing comma and space
	if len(sortOrder) > 0 {
		sortOrder = sortOrder[:len(sortOrder)-2]
	} else {
		sortOrder = "ucs.next_review_date ASC" // Default sorting
	}

	reviewsQuery := fmt.Sprintf(`
		WITH PendingReviews AS (
			SELECT ucs.card_id, ucs.stage_id, ucs.next_review_date, c.level
			FROM UserCardStatus ucs
			JOIN Cards c ON ucs.card_id = c.card_id
			WHERE ucs.user_id = $1 AND ucs.stage_id < 9
			AND (
				$2::boolean = true AND ucs.stage_id = 0
				OR
				$2::boolean = false AND ucs.next_review_date < NOW()
			)
			ORDER BY %s
			LIMIT 1
		)
		
		%s
		JOIN PendingReviews pr ON c.card_id = pr.card_id;
	`, sortOrder, CardSelector)

	return s.db.Query(reviewsQuery, userId, firstReview)
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
				SELECT COALESCE(
					(SELECT stage_id FROM UserCardStatus WHERE user_id = $1::VARCHAR AND card_id = $2),
					0
				) AS stage_id
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
				0, 
				$3::TIMESTAMPTZ
			)
			RETURNING card_id
		)

		SELECT i.card_id, c.word, $4 AS success, 1 AS stage_id
		FROM inserted i
		JOIN Cards c ON i.card_id = c.card_id;
	`

	return s.db.QueryRow(insertQuery, userId, review.CardId, review.ReviewDate, review.Success)

}

func (s* storage) UserCardStatusUpdate(userId string, review api.Review) (*sql.Row) {
	stageUpdate := fmt.Sprintf(`
    CASE 
        WHEN ucs.stage_id = 0 THEN 1  -- If stage_id is 0, always move to 1
        WHEN %t THEN ucs.stage_id + 1  -- If success, move to next stage
        ELSE GREATEST(
            1,
            ucs.stage_id - CEIL(
                ROUND(%d / 2.0, 1)
            ) * s.stage_penalty
        )
    END`, *review.Success, *review.IncorrectCount)
	
	updateQuery := fmt.Sprintf(`
		UPDATE UserCardStatus ucs
		SET stage_id = cs.new_stage_id,
			next_review_date = $3::TIMESTAMPTZ + s.stage_interval
		FROM (
			SELECT ucs.user_id, ucs.card_id, %s AS new_stage_id
			FROM UserCardStatus ucs
			LEFT JOIN SRSStages s ON ucs.stage_id = s.stage_id
			WHERE ucs.user_id = $1 AND ucs.card_id = $2
		) AS cs
		JOIN SRSStages s ON cs.new_stage_id = s.stage_id
		WHERE ucs.user_id = cs.user_id AND ucs.card_id = cs.card_id
		RETURNING ucs.card_id, ucs.stage_id`, stageUpdate)

	finalQuery := fmt.Sprintf(`
		WITH updated AS (
			%s
		)
		SELECT u.card_id, c.word AS card_word, %t AS success, u.stage_id
		FROM updated u
		JOIN Cards c ON u.card_id = c.card_id;`, updateQuery, *review.Success)

	return s.db.QueryRow(finalQuery, userId, review.CardId, review.ReviewDate)
}

func (s* storage) MostRecentReviewsQuery(userId string, numCards int) (*sql.Rows, error) {
	fmt.Print(numCards, userId)
	mostRecentReview := `
		SELECT r.card_id, c.word AS card_word, r.success, ucs.stage_id
		FROM Reviews r
		JOIN UserCardStatus ucs
			ON r.user_id = ucs.user_id AND r.card_id = ucs.card_id
		JOIN Cards c
			ON r.card_id = c.card_id
		WHERE r.user_id = $1
		ORDER BY ABS(EXTRACT(EPOCH FROM (r.review_date - NOW()))) ASC
		LIMIT $2;
	`

	return s.db.Query(mostRecentReview, userId, numCards)
}

func (s* storage) CardQuery(id int) (*sql.Rows, error) {
	cardQuery := fmt.Sprintf(`
		%s
		WHERE c.card_id = $1;
	`, CardSelector)
	
	return s.db.Query(cardQuery, id)
}
