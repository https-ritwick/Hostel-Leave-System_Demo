package handlers

import (
	"DynamicWebsiteProject/db"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
)

var JWTKey []byte

var store = sessions.NewCookieStore([]byte(generateRandomString(256)))

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nonce    string `json:"nonce"`
}

type WebPageData struct {
	WebsiteTitle                 string
	H1Heading                    string
	BodyParagraphText            string
	PostResponseMessage          string
	PosrResponseHTTPResponseCode string
}

func templateRender(httpWriter http.ResponseWriter, wd WebPageData, templateName string) {
	templateFilePath := "templates/" + templateName + ".html"
	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		http.Error(httpWriter, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(httpWriter, wd)
}
func templateRenderWithCustomeDataMap(httpWriter http.ResponseWriter, mapData map[string]string, templateName string) {
	templateFilePath := "templates/" + templateName + ".html"
	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		http.Error(httpWriter, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(httpWriter, mapData)
}

func RegisterHandler(httpWriter http.ResponseWriter, r *http.Request) {
	webPageDataLogin := WebPageData{"User Registration", "Enter Your Registration Details", "", "", ""}

	if r.Method == http.MethodPost {
		applicantname := r.FormValue("applicantname")
		applicantemail := r.FormValue("applicantemail")
		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Println(applicantname, applicantemail, username, password)
		id, err := db.CreateUserWithCredentials(applicantname, applicantemail, username, password)
		if err != nil {
			webPageDataLogin.PostResponseMessage = "DB Error Occured, Please contact administrator"
			webPageDataLogin.PosrResponseHTTPResponseCode = "200"
		} else {
			webPageDataLogin.PostResponseMessage = "Registration successful for username:" + username + " with userid:" + strconv.Itoa(int(id))
			fmt.Println("Registration successful for user: ", username, " with userid: ", id)
			webPageDataLogin.PosrResponseHTTPResponseCode = "200"
		}
	}
	templateRender(httpWriter, webPageDataLogin, "register")
}
func LoginCheckHandler(httpWriter http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(httpWriter, "Invalid request", http.StatusBadRequest)
		return
	}

	username := creds.Username
	password := creds.Password
	nonce := creds.Nonce
	// Get nonce from session
	session, _ := store.Get(r, "ajaxjwtdemo.com")
	savedNonce := session.Values["nonce"]
	if nonce != savedNonce {
		http.Error(httpWriter, "Invalid session or token", http.StatusForbidden)
		return
	}
	result, err := db.ValidateUserCredentials(username, password, nonce)
	if err != nil {
		http.Error(httpWriter, "DB Error, Contact DBA", http.StatusForbidden)
		return
	} else if !result {
		http.Error(httpWriter, "Login Attempt Failed, Unable to match username/password", http.StatusForbidden)
		return
	}
	// Create JWT token
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   creds.Username,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		fmt.Println("Error creating token", err)
		http.Error(httpWriter, "Error creating token", http.StatusInternalServerError)
		return
	}

	httpWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(httpWriter).Encode(map[string]string{
		"token": tokenString,
	})

}
func LoginHandler(httpWriter http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "ajaxjwtdemo.com")
	isAuthenticated, ok := session.Values["authenticatedUser"].(bool)
	if ok && isAuthenticated {
		http.Redirect(httpWriter, r, "/home", http.StatusSeeOther)
		return
	}
	nonce := generateRandomString(16) // Save this in session or map with a short expiry
	// Save to session or temp store for later verification
	session, _ = store.Get(r, "ajaxjwtdemo.com")
	session.Values["nonce"] = nonce
	session.Save(r, httpWriter)
	session.Values["authenticatedUser"] = false
	mapData := map[string]string{
		"WebsiteTitle":      "Login Page",
		"H1Heading":         "Login Page",
		"BodyParagraphText": "Enter login Details",
		"Nonce":             nonce}
	templateRenderWithCustomeDataMap(httpWriter, mapData, "login")
}

func HomePageHandler(httpWriter http.ResponseWriter, r *http.Request) {
	mapData := map[string]string{
		"WebsiteTitle":    "Go Lang Home Portal",
		"HomePageHeading": "Welcome to the home page",
	}
	if r.Method == http.MethodGet {
		nameReceived := strings.TrimSpace(r.URL.Query().Get("name"))
		if nameReceived != "" {
			mapData["WelcomeMessage"] = fmt.Sprintf("Welcome, %s!", nameReceived)
			fmt.Println(mapData["WelcomeMessage"])
		}
	}

	templateRenderWithCustomeDataMap(httpWriter, mapData, "home")
}
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "randomerror123" // fallback or handle error better
	}
	return hex.EncodeToString(bytes)
}

func AuthMiddlewareAPI(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		} else {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		username, err := validateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		fmt.Println("Authenticated user:", username)
		next.ServeHTTP(w, r)
	}
}
func AuthMiddlewareCookie(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil || cookie == nil {
			//http.Error(w, "Missing token", http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}
		tokenStr := cookie.Value

		username, err := validateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		fmt.Println("Authenticated user:", username)
		next.ServeHTTP(w, r)
	}
}
func validateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return JWTKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["sub"].(string)
		return username, nil
	}
	return "", jwt.ErrTokenInvalidClaims
}
