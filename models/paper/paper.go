package paper

import "jim_evaluate/models"

type Paper struct {
	models.Model

	Title   string `json:"title"`
	Content string `json:"content"`
	Type    int    `json:"type"`
	Score   int    `json:"score"`

}
