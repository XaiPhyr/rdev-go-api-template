package dto_test

import (
	"testing"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/dto"
)

func TestQuery(t *testing.T) {
	t.Run("test query", func(t *testing.T) {
		q := &dto.BaseFilters{}

		q.SanitizeQuery([]string{})
	})
}
