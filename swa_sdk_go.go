package swa_sdk_go

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewStupidWebauthn(baseUrl string) *StupidWebauthn {
	return &StupidWebauthn{Url: baseUrl + "/auth/auth/validate"}
}

type StupidWebauthn struct {
	Url string
}

type AuthResponse struct {
	UserID        int    `json:"user_id"`
	UserEmail     string `json:"user_email"`
	UserCreatedAt string `json:"user_created_at"`
}

func (p *StupidWebauthn) Middleware(req *http.Request) (*AuthResponse, int, error) {
	newReq, err := http.NewRequest("POST", p.Url, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	newReq.Header.Add("Cookie", req.Header.Get("Cookie"))
	newResp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	defer newResp.Body.Close()
	body, err := io.ReadAll(newResp.Body)
	if err != nil {
		return nil, newResp.StatusCode, err
	}

	if newResp.StatusCode != http.StatusOK {
		return nil, newResp.StatusCode, errors.New(string(body))
	}
	bodyJson := AuthResponse{}
	fmt.Printf("%s\n", string(body))
	err = json.Unmarshal(body, &bodyJson)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &bodyJson, newResp.StatusCode, nil
}

func (*StupidWebauthn) RemoveAuthCookie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     "swa_auth",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}
