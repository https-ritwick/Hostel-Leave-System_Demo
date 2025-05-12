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
	Role     string `json:"role"`
}

type WebPageData struct {
	WebsiteTitle                 string
	H1Heading                    string
	BodyParagraphText            string
	PostResponseMessage          string
	PosrResponseHTTPResponseCode string
}

/* type LeaveEntry struct {
	Reason   string
	FromDate string
	ToDate   string
	Status   string
} */

type StudentProfile struct {
	ID         int
	Name       string
	RoomNumber string
	Contact    string
	Gender     string
}

type StudentDashboardData struct {
	WebsiteTitle  string
	StudentName   string
	RoomNumber    string
	ContactNumber string
	Gender        string
	LeaveHistory  []db.LeaveEntry
}

// -------------------- TEMPLATES --------------------

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

// -------------------- REGISTER --------------------

func RegisterHandler(httpWriter http.ResponseWriter, r *http.Request) {
	webPageDataLogin := WebPageData{
		WebsiteTitle:        "User Registration",
		H1Heading:           "Enter Your Registration Details",
		BodyParagraphText:   "",
		PostResponseMessage: "",
	}

	if r.Method == http.MethodPost {
		applicantname := r.FormValue("applicantname")
		applicantemail := r.FormValue("applicantemail")
		username := r.FormValue("username")
		password := r.FormValue("password")
		roomnumber := r.FormValue("roomnumber")
		address := r.FormValue("address")
		contact := r.FormValue("contact")
		role := r.FormValue("role")
		gender := r.FormValue("gender")

		id, err := db.CreateUserWithAllDetails(applicantname, applicantemail, username, password, roomnumber, address, contact, role, gender)
		if err != nil {
			webPageDataLogin.PostResponseMessage = "DB Error Occurred. Please contact administrator."
			webPageDataLogin.PosrResponseHTTPResponseCode = "500"
		} else {
			webPageDataLogin.PostResponseMessage = "Registration successful for username: " + username + " with userid: " + strconv.Itoa(int(id))
			webPageDataLogin.PosrResponseHTTPResponseCode = "200"
		}
	}

	templateRender(httpWriter, webPageDataLogin, "register")
}

// -------------------- LOGIN --------------------

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
		http.Error(httpWriter, "Login Attempt Failed", http.StatusForbidden)
		return
	}

	role, err := db.FetchUserRoleByUsername(username)
	if err != nil {
		http.Error(httpWriter, "Unable to fetch user role", http.StatusInternalServerError)
		return
	}

	session.Values["authenticatedUser"] = true
	session.Values["role"] = role
	session.Values["username"] = username
	session.Save(r, httpWriter)

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   creds.Username,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		http.Error(httpWriter, "Error creating token", http.StatusInternalServerError)
		return
	}

	httpWriter.Header().Set("Content-Type", "application/json")

	redirectPath := "/student_dashboard"
	if role == "warden" {
		redirectPath = "/warden_dashboard"
	}

	json.NewEncoder(httpWriter).Encode(map[string]string{
		"token":         tokenString,
		"redirect_path": redirectPath,
	})

}

// -------------------- RESTRICTED WARDEN PAGE --------------------

func LeaveHandleHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "ajaxjwtdemo.com")
	role, ok := session.Values["role"].(string)
	if !ok || role != "warden" {
		http.Error(w, "Unauthorized: Only wardens can access this page.", http.StatusForbidden)
		return
	}
	http.ServeFile(w, r, "templates/leave_handle.html")
}

// -------------------- COMMON --------------------

func HomePageHandler(httpWriter http.ResponseWriter, r *http.Request) {
	mapData := map[string]string{
		"WebsiteTitle":    "Go Lang Home Portal",
		"HomePageHeading": "Welcome to the home page",
	}
	if r.Method == http.MethodGet {
		nameReceived := strings.TrimSpace(r.URL.Query().Get("name"))
		if nameReceived != "" {
			mapData["WelcomeMessage"] = fmt.Sprintf("Welcome, %s!", nameReceived)
		}
	}
	templateRenderWithCustomeDataMap(httpWriter, mapData, "home")
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "randomerror123"
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

func StudentDashboardHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "ajaxjwtdemo.com")

	role, ok := session.Values["role"].(string)
	if !ok || role != "student" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	username, ok := session.Values["username"].(string)
	if !ok {
		http.Error(w, "Session expired", http.StatusUnauthorized)
		return
	}

	// Fetch profile
	profile, err := db.FetchStudentProfile(username)
	if err != nil {
		http.Error(w, "Failed to load student profile", http.StatusInternalServerError)
		return
	}

	// Fetch leave history
	leaveHistory, err := db.FetchLeaveHistory(profile.ID)
	if err != nil {
		http.Error(w, "Failed to load leave history", http.StatusInternalServerError)
		return
	}
	for _, l := range leaveHistory {
		fmt.Println("LEAVE STATUS:", l.Status)
	}

	data := StudentDashboardData{
		WebsiteTitle:  "Student Dashboard",
		StudentName:   profile.Name,
		RoomNumber:    profile.RoomNumber,
		ContactNumber: profile.Contact,
		Gender:        profile.Gender,
		LeaveHistory:  leaveHistory,
	}

	tmpl, err := template.ParseFiles("templates/student_dashboard.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "ajaxjwtdemo.com")

	// Set flash message
	session.AddFlash("Logged out successfully")

	// Clear session values
	session.Values["authenticatedUser"] = false
	delete(session.Values, "role")
	delete(session.Values, "username")
	delete(session.Values, "nonce")

	// Clear token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "ajaxjwtdemo.com")
	isAuthenticated, ok := session.Values["authenticatedUser"].(bool)
	if ok && isAuthenticated {
		http.Redirect(w, r, "/student_dashboard", http.StatusSeeOther)
		return
	}

	// Get nonce
	nonce := generateRandomString(16)
	session.Values["nonce"] = nonce
	session.Values["authenticatedUser"] = false

	// Get flash messages
	var flashMsg string
	if flashes := session.Flashes(); len(flashes) > 0 {
		flashMsg = fmt.Sprintf("%v", flashes[0])
	}
	session.Save(r, w)

	mapData := map[string]string{
		"WebsiteTitle":      "Login Page",
		"H1Heading":         "Login Page",
		"BodyParagraphText": flashMsg,
		"Nonce":             nonce,
	}
	templateRenderWithCustomeDataMap(w, mapData, "login")
}

// -------------------- WARDEN DASHBOARD HANDLER --------------------

func WardenDashboardHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "ajaxjwtdemo.com")
	role, ok := session.Values["role"].(string)
	if !ok || role != "warden" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	tmpl, err := template.ParseFiles("templates/warden_dashboard.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
