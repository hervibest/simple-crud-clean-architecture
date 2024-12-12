package repository

import (
	"simple-crud-clean-architecture/internal/entity"

	"github.com/sirupsen/logrus"
)

type EmployeeRepository struct {
	Repository[entity.Employee]
	Log *logrus.Logger
}

func NewEmployeeRepository(log *logrus.Logger) *EmployeeRepository {
	return &EmployeeRepository{
		Log: log,
	}
}
