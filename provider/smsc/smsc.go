package smsc

import (
	"fmt"
	"github.com/messagebird/sachet"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}

const SmscRequestTimeout = time.Second * 60

type Smsc struct {
	Login    string
	Password string
}

func NewSmsc(config Config) *Smsc {
	Smsc := &Smsc{Login: config.Login, Password: config.Password}
	return Smsc
}

func (c *Smsc) Send(message sachet.Message) (err error) {
	for _, number := range message.To {
		err = c.SendOne(message, number)
		if err != nil {
			return fmt.Errorf("Failed to make API call to smsc:%s", err)
		}
	}
	return
}

func (c *Smsc) SendOne(message sachet.Message, PhoneNumber string) (err error) {
	encoded_message := url.QueryEscape(message.Text)
	smsURL := fmt.Sprintf("https://smsc.ru/sys/send.php?login=%s&psw=%s&phones=%s&sender=%s&fmt=0&mes=%s", c.Login, c.Password, PhoneNumber, message.From, encoded_message)
	var request *http.Request
	var resp *http.Response
	request, err = http.NewRequest("GET", smsURL, nil)
	if err != nil {
		return
	}
	httpClient := &http.Client{}
	httpClient.Timeout = SmscRequestTimeout
	resp, err = httpClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var body []byte
	resp.Body.Read(body)
	if resp.StatusCode == http.StatusOK && err == nil {
		return
	}
	return fmt.Errorf("Failed sending sms:Reason: %s, StatusCode : %d", string(body), resp.StatusCode)
}
