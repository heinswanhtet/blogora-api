package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/heinswanhtet/blogora-api/controllers"
)

func TestAuthorServiceHandler(t *testing.T) {
	controller := controllers.NewAuthorController(s)
	t.Run("should handle get authors", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/author", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("GET /author", controller.HandleGetAuthors)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fail creating an author if the payload is missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/author", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("POST /author", controller.HandleCreateAuthor)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}
