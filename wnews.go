package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//CustomResponse ...
type CustomResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Body    interface{} `json:"body,omitempty"`
	Error   string      `json:"error,omitempty"`
}

const (
	dataStoreProjectID = "inapp-infrastructure-190215"
)

//GetWeather ...
func GetWeather(w http.ResponseWriter, r *http.Request) {
	var latitude, longitude string
	latitude = r.URL.Query().Get("latitude")
	longitude = r.URL.Query().Get("longitude")

	var response CustomResponse
	if latitude == "" || longitude == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = CustomResponse{
			Status:  false,
			Message: "Missing query params",
		}
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	url := fmt.Sprintf("https://weatherport.co/hapi/getWeatherData?latitude=%s&longitude=%s", latitude, longitude)
	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	body, _ := ioutil.ReadAll(resp.Body)
	var res map[string]interface{}
	json.Unmarshal(body, &res)
	//body, _ := json.Marshal(resp.Body)
	response = CustomResponse{
		Status:  true,
		Message: "Success",
		Body:    res,
	}
	_ = json.NewEncoder(w).Encode(response)
}

//GetNews ..
func GetNews(w http.ResponseWriter, r *http.Request) {

	var response CustomResponse

	url := fmt.Sprintf("https://hapi.newsprompt.co/HomePageExtension/getApiArticles?v=1.0&sec=usnews&ldesc=128&actno=104&origin=extension&d=newsprompt.co&maxno=100")
	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response = CustomResponse{
			Status:  false,
			Message: "Something went wrong",
		}
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	body, _ := ioutil.ReadAll(resp.Body)
	var res map[string]interface{}
	json.Unmarshal(body, &res)
	//body, _ := json.Marshal(resp.Body)
	response = CustomResponse{
		Status:  true,
		Message: "Success",
		Body:    res,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func makeHTTPGETRequest(url string) error {
	log.Println("Calling url: ", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("GET request failed, err: %+v", err)
		return err
	}

	defer resp.Body.Close()

	log.Println("status: ", resp.StatusCode)
	return nil
}

func main() {
	log.Print("wnews service started")

	http.HandleFunc("/weather", GetWeather)
	http.HandleFunc("/news", GetNews)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
