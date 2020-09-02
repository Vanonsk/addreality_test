package restapi

import (
	"fmt"
	"net/http"

	"addreality_t/internal/addreality/users"
	"addreality_t/internal/metrics"
)

type ServiceHandlers struct {
	UserController *users.Controller
}

func (h *ServiceHandlers) CountUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics.UserRequestCount.Inc()

		userId := r.URL.Query().Get("user_id")
		h.UserController.AddCounterUser(userId)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *ServiceHandlers) GetRobotUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics.UserStatusRequestCount.Inc()

		robotCountUsers := h.UserController.GetRobotUserCount()
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, robotCountUsers)
	}
}
