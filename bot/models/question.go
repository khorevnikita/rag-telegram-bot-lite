package models

import (
	"github.com/google/uuid"
	"time"
)

type QuestionType string

const (
	QuestionTypeText   QuestionType = "text"
	QuestionTypeNumber QuestionType = "number"
	QuestionTypeSelect QuestionType = "select"
	QuestionTypeEmail  QuestionType = "email"
	/*QuestionTypeDate     QuestionType = "date"
	QuestionTypeCity     QuestionType = "city"
	QuestionTypeBirthday QuestionType = "birthday"
	QuestionTypeGender   QuestionType = "gender"*/
)

type Question struct {
	ID                     uuid.UUID        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	IsPublished            bool             `gorm:"default:false"`
	IsRequired             bool             `gorm:"default:false"`
	Text                   string           `gorm:""`
	Hint                   *string          `gorm:""`
	Type                   QuestionType     `gorm:""`
	SelectableOptionsCount int              `gorm:""`
	Order                  int              `gorm:""`
	CreatedAt              time.Time        `gorm:"autoCreateTime"`
	UpdatedAt              time.Time        `gorm:"autoUpdateTime"`
	DeletedAt              *time.Time       `gorm:"index"` // Для soft delete
	QuestionOptions        []QuestionOption // Аннотация для указания внешнего ключа
	Answers                []Answer
}
