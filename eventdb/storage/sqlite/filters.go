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
	if filters.IsEmpty != nil {
		if ((*filters.IsEmpty) && has) || !has {
			return false
		}
	}

	if filters.Equal != nil && !strings.EqualFold(*filters.Equal, val) {
		return false
	}

	if filters.NotEqual != nil && strings.EqualFold(*filters.NotEqual, val) {
		return false
	}

	if filters.Contains != nil && !strings.Contains(strings.ToLower(val), strings.ToLower(*filters.Contains)) {
		return false
	}

	if filters.NotContains != nil && strings.Contains(strings.ToLower(val), strings.ToLower(*filters.NotContains)) {
		return false
	}

	return true
}
