package pages

import (
	"strings"

	"github.com/razzie/razlink"
)

func getIDFromRequest(r *razlink.PageRequest) (id, trailing string) {
	parts := strings.SplitN(r.RelPath, "/", 2)
	if len(parts) > 1 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}
