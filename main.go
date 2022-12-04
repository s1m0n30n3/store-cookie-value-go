package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", serveHtml)
	http.HandleFunc("/submit", submitInfo)
	http.ListenAndServe(":8080", nil)
}

func getCode(message string) string {
	h := hmac.New(sha256.New, []byte("some strings"))
	h.Write([]byte(message))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func submitInfo(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Redirect(response, request, "/", http.StatusSeeOther)
		return
	}

	email := request.FormValue("email")
	if email == "" {
		http.Redirect(response, request, "/", http.StatusSeeOther)
		return
	}

	code := getCode(email)

	cookie := http.Cookie{
		Name:  "session",
		Value: code + "|" + email,
	}

	http.SetCookie(response, &cookie)
	http.Redirect(response, request, "/", http.StatusSeeOther)
}

func serveHtml(response http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("session")
	if err != nil {
		cookie = &http.Cookie{}
	}

	isEqual := true
	sliceOfStrings := strings.SplitN(cookie.Value, "|", 2)
	if len(sliceOfStrings) == 2 {
		cookieCode := sliceOfStrings[0]
		cookieEmail := sliceOfStrings[1]

		code := getCode(cookieEmail)

		isEqual = hmac.Equal([]byte(cookieCode), []byte(code))
	}

	message := "Not logged in"
	if isEqual {
		message = "Logged in"
	}

	html := `
	<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<title>Doc Cookie</title>
		</head>
		<body>
      <p>Cookie value: ` + cookie.Value + `</p>
      <p>` + message + `</p>
			<form action="/submit" method="post">
				<input type="email" name="email" />
				<input type="submit" />
			</form>
		</body>
	</html>`

	io.WriteString(response, html)
}
