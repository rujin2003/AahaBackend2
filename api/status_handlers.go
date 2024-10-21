package api

import (
	"net/http"
)

func (s *ApiServer) handleGetStatus(w http.ResponseWriter, r *http.Request) error {
	err := s.store.AddSellStatus("TRUE")
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, "Status updated")

}

func (s *ApiServer) handleUpdateStatus(w http.ResponseWriter, r *http.Request) error {
	err := s.store.UpdateSellStatus("FALSE")
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, "Status updated")
}
