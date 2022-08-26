package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type prometheusMessage struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []struct {
		Status string `json:"status"`
		Labels struct {
			Alertname   string `json:"alertname"`
			Instance    string `json:"instance"`
			Job         string `json:"job"`
			Monitor     string `json:"monitor"`
			ServiceID   string `json:"serviceId"`
			ServiceName string `json:"serviceName"`
			Severity    string `json:"severity"`
		} `json:"labels"`
		Annotations struct {
			Description string `json:"description"`
			Summary     string `json:"summary"`
			Title       string `json:"title"`
		} `json:"annotations"`
		StartsAt     time.Time `json:"startsAt"`
		EndsAt       time.Time `json:"endsAt"`
		GeneratorURL string    `json:"generatorURL"`
		Fingerprint  string    `json:"fingerprint"`
	} `json:"alerts"`
	GroupLabels struct {
	} `json:"groupLabels"`
	CommonLabels struct {
		Alertname   string `json:"alertname"`
		Instance    string `json:"instance"`
		Job         string `json:"job"`
		Monitor     string `json:"monitor"`
		ServiceID   string `json:"serviceId"`
		ServiceName string `json:"serviceName"`
		Severity    string `json:"severity"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
		Title       string `json:"title"`
	} `json:"commonAnnotations"`
	ExternalURL     string `json:"externalURL"`
	Version         string `json:"version"`
	GroupKey        string `json:"groupKey"`
	TruncatedAlerts int    `json:"truncatedAlerts"`
}

func main() {
	receiveMessage()
}

func getMessage(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	var p prometheusMessage
	err := json.Unmarshal(body, &p)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sendMessage(p)
	w.WriteHeader(http.StatusOK)
}

func receiveMessage() {

	http.HandleFunc("/prometheus-chatwork", getMessage)

	log.Fatal(http.ListenAndServe(":8001", nil))

}

func sendMessage(p prometheusMessage) {

	url := fmt.Sprintf("https://api.chatwork.com/v2/rooms/%s/messages", os.Getenv("chatwork_roomID"))
	info := fmt.Sprintf(`
	[info][title]🚧 %s %s[/title]
	🛸 job: %s
	🤖 serviceName: %s
	[hr]
	📋 %s
	[hr]
	💻 %s
	[/info]
	`, p.CommonLabels.Instance,
		p.CommonAnnotations.Title,
		p.CommonLabels.Job,
		p.CommonLabels.ServiceName,
		p.CommonAnnotations.Description,
		p.CommonAnnotations.Summary)
	payload := strings.NewReader(fmt.Sprintf("self_unread=0&body=%s", info))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-ChatWorkToken", os.Getenv("chatwork_token"))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res)
	fmt.Println(string(body))

}
