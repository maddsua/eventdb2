package sqlite

import (
	"strings"

	"github.com/maddsua/eventdb2/storage/model"
)

func Matchlabels(labels model.StringMap, filters *model.LogLabelFilter) bool {

	if labels == nil || filters == nil || filters.Key == "" {
		return false
	}

	val, has := labels[filters.Key]
	if filters.IsEmpty.Valid {
		if (filters.IsEmpty.Bool && has) || !has {
			return false
		}
	}

	if filters.Equal.Valid && !strings.EqualFold(filters.Equal.String, val) {
		return false
	}

	if filters.NotEqual.Valid && strings.EqualFold(filters.NotEqual.String, val) {
		return false
	}

	if filters.Contains.Valid && !strings.Contains(strings.ToLower(val), strings.ToLower(filters.Contains.String)) {
		return false
	}

	if filters.NotContains.Valid && strings.Contains(strings.ToLower(val), strings.ToLower(filters.NotContains.String)) {
		return false
	}

	return true
}
