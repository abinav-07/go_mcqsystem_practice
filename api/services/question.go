package services

import (
	"github/abinav-07/mcq-test/database/models"
	"github/abinav-07/mcq-test/infrastructure"
)

type QuestionService struct {
	repository infrastructure.Database
}

// contructor
func NewQuestionService(database infrastructure.Database) QuestionService {
	return QuestionService{
		repository: database,
	}
}

// Get  Question By Id
func (c QuestionService) GetTestQuestionsById(TestID uint) (*models.Question, error) {
	question := models.Question{}

	return &question, c.repository.DB.Where("test_id = ?", TestID).Preload("Option").Find(&question).Error
}

// Create question and answer
func (c QuestionService) CreateQuestion(question models.Question) (*models.Question, error) {
	return &question, c.repository.DB.Preload("Option").Create(&question).Error
}

// Get Correct answers for a test: Returns Questions
func (c QuestionService) GetCorrectAnswers(TestID uint) ([]models.Option, error) {
	testQuestions := []models.Question{}
	testQA := []models.Option{}

	//Get only question for a test
	if testErr := c.repository.DB.Where("test_id = ?", TestID).Find(&testQuestions).Error; testErr != nil {
		return nil, testErr
	}

	//Question Id to array
	questionIds := []uint{}
	for _, question := range testQuestions {
		questionIds = append(questionIds, question.ID)
	}

	//Get answers of test question
	return testQA, c.repository.DB.Model(&models.Option{}).Where("is_correct = ? AND question_id IN (?)", true, questionIds).Select("QuestionID", "ID as option_id").Find(&testQA).Error

}
