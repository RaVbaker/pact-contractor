package hooks

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpHook struct {
	Method  string
	Url     string
	Body    string
	Headers map[string]string
}

func (l *HttpHook) Run(path string) error {
	url := templateString(path, l.Url)
	body := strings.NewReader(templateString(path, l.Body))

	fmt.Printf("%s Request to: %q\n", l.Method, url)

	request, err := http.NewRequest(l.Method, url, body)
	if err != nil {
		return err
	}

	for key, value := range l.Headers {
		request.Header.Set(key, templateString(path, value))
	}

	resp, err := http.DefaultClient.Do(request)

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Response (status code: %d) received\n", resp.StatusCode)
	bodyString := string(bodyBytes)
	if len(bodyString) > 0 {
		fmt.Println("Body:")
		fmt.Println(bodyString)
	}

	return err
}
