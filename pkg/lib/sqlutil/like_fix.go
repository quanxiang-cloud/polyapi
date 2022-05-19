package sqlutil

import (
	"strings"
)

// LikeEscape convert key as safe in like statement.
// eg "foo_bar"->"foo\_bar"
// TODO: don't convert \_ to \\_
func LikeEscape(key string) string {
	return strings.Replace(key, `_`, `\_`, -1)
}
