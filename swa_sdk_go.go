package swa_sdk_go

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func NewStupidWebauthn(baseUrl string) *StupidWebauthn {
	return &StupidWebauthn{BaseUrl: baseUrl}
}

type StupidWebauthn struct {
	BaseUrl string
}

type SetCookieFunc func(value string)

type AuthResponse struct {
	ID           int     `json:"id"`
	Email        string  `json:"email"`
	JwtVersion   int     `json:"jwt_version"`
	GdprDeleteAt *string `json:"gdpr_delete_at"`
	CreatedAt    string  `json:"created_at"`
}

// If setCookie is nil the cookie will not be removed on 401 response
func (p *StupidWebauthn) AuthMiddleware(req *http.Request, setCookie *http.Header) (*AuthResponse, int, error) {
	return p.genericMiddleware("/auth/auth/validate")(req, setCookie)
}
func (p *StupidWebauthn) AuthCsrfMiddleware(req *http.Request, setCookie *http.Header) (*AuthResponse, int, error) {
	return p.genericMiddleware("/auth/auth/csrf/validate")(req, setCookie)
}
func (p *StupidWebauthn) AuthDoubleCheckMiddleware(req *http.Request, setCookie *http.Header) (*AuthResponse, int, error) {
	return p.genericMiddleware("/auth/auth/doublecheck/validate")(req, setCookie)
}

func (p *StupidWebauthn) genericMiddleware(uri string) func(req *http.Request, setCookie *http.Header) (*AuthResponse, int, error) {
	return func(req *http.Request, setCookie *http.Header) (*AuthResponse, int, error) {
		status := http.StatusInternalServerError
		newResp, err := p.fetch(req, http.MethodGet, uri)
		if err != nil {
			return nil, status, err
		}

		defer newResp.Body.Close()
		body, err := io.ReadAll(newResp.Body)
		if err != nil {
			return nil, status, err
		}

		switch newResp.StatusCode {
		case http.StatusUnauthorized:
			if setCookie != nil {
				(*setCookie).Set("Set-Cookie", newResp.Header.Get("Set-Cookie"))
			}
			status = newResp.StatusCode
		case http.StatusOK, http.StatusInternalServerError, http.StatusBadRequest:
			status = newResp.StatusCode
		default:
			return nil, status, fmt.Errorf("Unable to connect to auth server [%d]", &newResp.StatusCode)
		}

		if newResp.StatusCode != http.StatusOK {
			return nil, status, errors.New(string(body))
		}
		bodyJson := AuthResponse{}
		fmt.Printf("%s\n", string(body))
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		return &bodyJson, status, nil
	}
}

func (p *StupidWebauthn) Logout(req *http.Request, setCookie http.ResponseWriter) (int, error) {
	newResp, err := p.fetch(req, http.MethodGet, "/auth/logout")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if newResp.StatusCode == http.StatusCreated {
		setCookie.Header().Set("Set-Cookie", newResp.Header.Get("Set-Cookie"))
	}

	return newResp.StatusCode, nil
}

func (p *StupidWebauthn) fetch(req *http.Request, method string, uri string) (*http.Response, error) {
	newReq, err := http.NewRequest(method, p.BaseUrl+uri, nil)
	if err != nil {
		return nil, err
	}
	newReq.Header.Add("Cookie", req.Header.Get("Cookie"))
	newResp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		return nil, err
	}
	return newResp, nil
}
