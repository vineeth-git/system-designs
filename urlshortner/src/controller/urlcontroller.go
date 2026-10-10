package controller
import (
	"net/http"
	"encoding/json"
	"urlshorner/models"
	"fmt"
)

type Controller struct {
	store *models.Store
}
func NewController(store *models.Store) *Controller {
	return &Controller{store: store}
}
func (c *Controller) handleErr(w http.ResponseWriter, r *http.Request, err string, status int) { 
	http.Error(w, err, status)
}

func (c *Controller) handleOk(w http.ResponseWriter, r *http.Request, response map[string]string) {
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) HandleGenerateShortURL(w http.ResponseWriter, r *http.Request) {
	var request struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.handleErr(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	code, err := c.store.Generate(r.Context(), request.URL)
	if err != nil {
		c.handleErr(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	c.handleOk(w, r, map[string]string{
		"code": code,
		"shortUrl": fmt.Sprintf("http://localhost:8000/%s", code),
		"status": fmt.Sprint(http.StatusOK),
	})
}

func (c *Controller) HandleGetFullURL(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if longUrl, err := c.store.Get(r.Context(), code); err != nil {
		http.NotFound(w, r)
		return
	} else {
		c.handleOk(w, r, map[string]string{
			"fullUrl": longUrl,
			"status": fmt.Sprint(http.StatusOK),
		})
		// http.Redirect(w, r, longUrl, http.StatusFound)
	}
}