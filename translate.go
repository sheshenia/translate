package translate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func Text(langFrom, langTo, text string) (string, error) {
	return Do(nil, langFrom, langTo, "text", text)
}

func TextWithContext(ctx context.Context, langFrom, langTo, text string) (string, error) {
	return Do(ctx, langFrom, langTo, "text", text)
}

func Html(langFrom, langTo, html string) (string, error) {
	return Do(nil, langFrom, langTo, "html", html)
}

func HtmlWithContext(ctx context.Context, langFrom, langTo, html string) (string, error) {
	return Do(ctx, langFrom, langTo, "html", html)
}

var translationAPIs = []string{
	`https://translate.argosopentech.com/translate`,
	`https://libretranslate.de/translate`,
	`https://translate.mentality.rip/translate`,
	`https://trans.zillyhuhn.com/translate`,
	//`https://translate.astian.org/translate`,
}

type Request struct {
	Q      string `json:"q"`
	Source string `json:"source"` //source: "ja",
	Target string `json:"target"` //target: "en",
	Format string `json:"format"` //format: "text"
}

type Response struct {
	TranslatedText string `json:"translatedText,omitempty"` //"translatedText"
	Error          string `json:"error,omitempty"`
}

var customHttp = http.Client{Timeout: time.Second * 5}

func Do(ctx context.Context, langFrom, langTo, format, text string) (string, error) {
	var (
		buf bytes.Buffer
		req *http.Request
		err error
	)

	tr := Request{
		Q:      text,
		Source: langFrom,
		Target: langTo,
		Format: format,
	}

	if err := json.NewEncoder(&buf).Encode(&tr); err != nil {
		return "", err
	}

	api := translationAPIs[rand.Intn(len(translationAPIs))]

	if ctx != nil {
		req, err = http.NewRequestWithContext(ctx, "POST", api, &buf)
	} else {
		req, err = http.NewRequest("POST", api, &buf)
	}
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := customHttp.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.New("response status: " + strconv.Itoa(res.StatusCode))
	}

	var trResp Response
	if err := json.NewDecoder(res.Body).Decode(&trResp); err != nil {
		return "", err
	}

	if trResp.Error != "" {
		return "", errors.New(trResp.Error)
	}

	return trResp.TranslatedText, nil
}

func init() {
	rand.Seed(time.Now().Unix())
}
