package http

import (
	"net/http"

	"github.com/RakibulBh/AI-pr-reviewer/internal/utils/json"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) Health(w http.ResponseWriter, r *http.Request) {
	json.WriteSuccessJSON(w, http.StatusOK, "OK", nil)
}
