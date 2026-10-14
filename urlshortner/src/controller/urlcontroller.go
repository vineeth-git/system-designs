package controller
import (
	"net/http"
	"encoding/json"
	"urlshorner/models"
	"fmt"
	"sync/atomic"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

)

type Controller struct {
	store *models.Store
	totalRequests uint64
	req *prometheus.CounterVec
}
func NewController(store *models.Store) *Controller {
	var reqs = promauto.NewCounterVec(
		prometheus.CounterOpts{Name: "shortener_requests_total"},
		[]string{"route"},
	)	
	return &Controller{store: store, totalRequests: 0, req: reqs}
}
func (c *Controller) handleErr(w http.ResponseWriter, r *http.Request, err string, status int) { 
	http.Error(w, err, status)
}

func (c *Controller) handleOk(w http.ResponseWriter, r *http.Request, response map[string]string) {
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) HandleGenerateShortURL(w http.ResponseWriter, r *http.Request) {
	c.req.WithLabelValues("shorten").Inc()

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
	c.req.WithLabelValues("get").Inc()
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

func (c *Controller) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	count := atomic.LoadUint64(&c.totalRequests)
	fmt.Fprintf(w, "# TYPE http_requests_total counter\n")
	fmt.Fprintf(w, "http_requests_total %d\n", count)
}

func (c *Controller) CounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/metrics" {
			atomic.AddUint64(&c.totalRequests, 1)
		}
		next.ServeHTTP(w, r)
	})
}