package contract_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/raafvargas/wrapit/contract"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := contract.NewError(http.StatusOK, "ok")
	assert.EqualError(t, err, "Code: 200 - Messages: ok")
}

func TestFromValidationError(t *testing.T) {
	err := errors.New("invalid")

	e := contract.FromValidationError(err)
	assert.Equal(t, http.StatusInternalServerError, e.Code)

	s := &struct {
		Value int `validate:"gt=0"`
	}{
		Value: -1,
	}

	validate := validator.New()

	err = validate.Struct(s)
	assert.Error(t, err)

	e = contract.FromValidationError(err)
	assert.Equal(t, http.StatusUnprocessableEntity, e.Code)
}

func TestBusinessError(t *testing.T) {
	err := contract.BusinessErrorf("invalid arg %s", ":(")
	assert.Equal(t, http.StatusConflict, err.Code)
	assert.Contains(t, err.Messages, "invalid arg :(")
}
