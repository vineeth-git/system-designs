package main
import (
	"fmt"
	"time"
	"io"
	"bytes"
	"testing"
	"errors"
	"net/http"
	"encoding/json"
	"sync"
)
type Tester struct {
	client *http.Client
	host string
	mu sync.Mutex
	errorMap map[string]int
}

func (t *Tester) Generate(fullUrl string) (string, error) {
	url := fmt.Sprintf("%s/generate", t.host)
	body, _ := json.Marshal(map[string]string{"url": fullUrl})
	resp, err := t.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	data, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	var response struct {
		Code string `json:"code"`
		ShortUrl string `json:"shortUrl"`
	}
	if err := json.Unmarshal(data, &response); err != nil || response.Code == "" {
		return "", errors.New("Empty code" + string(data))
	}
	return response.Code, nil
}

func (t *Tester) GetFullUrl(code string) (string, error) {
	shortUrl := fmt.Sprintf("%s/%s", t.host, code)
	resp, err := t.client.Get(shortUrl)
	if err != nil {
		return "", err
	}
	data, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	var response struct {
		FullUrl string `json:"fullUrl"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return "", err
	}
	return response.FullUrl, nil
}
func NewTester() *Tester {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{MaxIdleConnsPerHost: 100},
	}
	return &Tester{client, "http://localhost:8080/", sync.Mutex{}, map[string]int{}}
}

func MakeOneTest(tst *Tester, dt chan<- TestStatus, wg *sync.WaitGroup) {
	defer wg.Done()
	code, err := tst.Generate("kuttyradha.com")
	if err != nil {
		tst.mu.Lock()
		tst.errorMap[err.Error()]++
		tst.mu.Unlock()
		// fmt.Println("Generate err: ", err)
		dt <- TestStatus{0, 0}
		return
	}
	_, err = tst.GetFullUrl(code)
	if err != nil {
		tst.mu.Lock()
		tst.errorMap[err.Error()]++
		tst.mu.Unlock()
		// fmt.Println("Get err: ", err)
		dt <- TestStatus{ 1, 0}
		return
	}
	dt <- TestStatus{1, 1}
}

type TestStatus struct {
	GeneratePass int
	GetPass int
}

func TestLoad(t *testing.T) {
	tst := NewTester()
	cases := 5000
	postSuccess, getSuccess := 0, 0
	fullSucces := 0
	wg := &sync.WaitGroup{}
	respa := make(chan TestStatus, cases)
	for i:= 0; i < cases; i++ {
		wg.Add(1)
		go MakeOneTest(tst, respa, wg)
	}
	go func() {
		wg.Wait()
		close(respa)
	}()
	
	for r := range respa {
		postSuccess += r.GeneratePass
		getSuccess += r.GetPass

		if r.GeneratePass + r.GetPass == 2 {
			fullSucces += 1
		}
	}
	fmt.Printf("Total Cases: %d\n TotalPostSuccess: %d\n TotalGetSuccess: %d\nTotal Success: %d\n\nSuccess P : %.2f", cases, postSuccess, getSuccess, fullSucces, float64((fullSucces*100)/cases))
	// fmt.Println("Err Map", tst.errorMap)
}
