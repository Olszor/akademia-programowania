package academy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGradeStudent(t *testing.T) {
	t.Run("Test student not found", func(t *testing.T) {
		repoMock := NewRepository(t)
		repoMock.On("Get", "Adam").Return(nil, ErrStudentNotFound)
		var want error = nil
		got := GradeStudent(repoMock, "Adam")
		assert.Equal(t, want, got)
	})
	t.Run("Test invalid grade", func(t *testing.T) {
		repoMock := NewRepository(t)
		studentMock := NewStudent(t)
		repoMock.On("Get", "Adam").Return(studentMock, nil)
		studentMock.On("FinalGrade").Return(6)
		want := ErrInvalidGrade
		got := GradeStudent(repoMock, "Adam")
		assert.Equal(t, want, got)
	})
	t.Run("Test student does not get promoted to next year", func(t *testing.T) {
		repoMock := NewRepository(t)
		studentMock := NewStudent(t)
		repoMock.On("Get", "Adam").Return(studentMock, nil)
		studentMock.On("FinalGrade").Return(1)
		studentMock.On("Year").Return(uint8(1))
		studentMock.On("Name").Return("Adam")
		repoMock.On("Save", "Adam", studentMock.Year()).Return(nil)
		var want error = nil
		got := GradeStudent(repoMock, "Adam")
		assert.Equal(t, want, got)
		repoMock.AssertCalled(t, "Save", "Adam", studentMock.Year())
	})
	t.Run("Test student graduates", func(t *testing.T) {
		repoMock := NewRepository(t)
		studentMock := NewStudent(t)
		repoMock.On("Get", "Adam").Return(studentMock, nil)
		studentMock.On("FinalGrade").Return(5)
		studentMock.On("Year").Return(uint8(3))
		studentMock.On("Name").Return("Adam")
		repoMock.On("Graduate", "Adam").Return(nil)
		var want error = nil
		got := GradeStudent(repoMock, "Adam")
		assert.Equal(t, want, got)
		repoMock.AssertCalled(t, "Graduate", "Adam")
	})
	t.Run("Test student gest promoted to next year", func(t *testing.T) {
		repoMock := NewRepository(t)
		studentMock := NewStudent(t)
		repoMock.On("Get", "Adam").Return(studentMock, nil)
		studentMock.On("FinalGrade").Return(3)
		studentMock.On("Year").Return(uint8(1))
		studentMock.On("Name").Return("Adam")
		repoMock.On("Save", "Adam", studentMock.Year()+1).Return(nil)
		var want error = nil
		got := GradeStudent(repoMock, "Adam")
		assert.Equal(t, want, got)
		repoMock.AssertCalled(t, "Save", "Adam", studentMock.Year()+1)
	})
}

func TestGradeYear(t *testing.T) {
	t.Run("Test no students", func(t *testing.T) {
		repoMock := NewRepository(t)
		repoMock.On("List", uint8(3)).Return([]string{}, nil)
		var want error = nil
		got := GradeYear(repoMock, uint8(3))
		assert.Equal(t, want, got)
	})
	t.Run("Test student not found", func(t *testing.T) {
		repoMock := NewRepository(t)
		repoMock.On("List", uint8(3)).Return([]string{"Adam"}, nil)
		repoMock.On("Get", "Adam").Return(nil, ErrStudentNotFound)
		var want error = nil
		got := GradeYear(repoMock, uint8(3))
		assert.Equal(t, want, got)
	})
}
