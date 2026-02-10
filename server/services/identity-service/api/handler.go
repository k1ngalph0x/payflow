package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"

	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
	"github.com/k1ngalph0x/payflow/identity-service/config"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct{
	DB *sql.DB
	Config *config.Config
	WalletClient *walletclient.WalletClient
}

type SignUpRequest struct{
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role 	  string 	`json:"role"`
}

type SignInRequest struct{
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type Claims struct{
	UserID string `json:"user_id"`
	Email string `json:"email"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewHandler(db *sql.DB, cfg *config.Config, walletclient *walletclient.WalletClient) *Handler {
	return &Handler{DB: db, Config: cfg, WalletClient: walletclient}
}

func(h *Handler) GenerateJWT(userId, email, role string)(string, error){

	expiration := time.Now().Add(24 * time.Hour)
	//expiration := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		UserID: userId,
		Role: role,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: "payflow-auth",
			Subject: userId,
		},
	}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(h.Config.TOKEN.JwtKey))

	if err != nil{
		fmt.Println(err)
		return "", err
	}

	return tokenString, nil
}

func(h *Handler) SignUp(c *gin.Context) {
	
	//var user models.User
	var req SignUpRequest
	var userId string 

	err := c.ShouldBindJSON(&req)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	if req.Email == "" || req.Password == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Email and password are required"})
		return 
	}

	if req.Role != "user" && req.Role != "merchant"{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Role must be user or merchant"})
		return 
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := strings.TrimSpace(req.Password)

	var existingEmail string
	query := `SELECT email FROM payflow_auth WHERE email = $1`
	err = h.DB.QueryRow(query, email).Scan(&existingEmail)
	if err == nil{
		c.JSON(http.StatusConflict, gin.H{"error":"Email already exists"})
		return 
	}

	if err != sql.ErrNoRows{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create user"})
		return 
	}


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create a new user"})
		return 
	}
	
	//created_at := time.Now().UTC()
	//updated_at := time.Now().UTC()
	
	//insert the new user
	// insertQuery := `
	// INSERT INTO payflow_auth (email, password, created_at, updated_at) 
	// VALUES ($1, $2, $3, $4)
	// `

	insertQuery := `
	INSERT INTO payflow_auth (email, password, role) 
	VALUES ($1, $2, $3)
	`
	//_, err = h.DB.Exec(insertQuery, email, string(hashedPassword), created_at, updated_at)
	_, err = h.DB.Exec(insertQuery, email, string(hashedPassword), req.Role)
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create a new user"})
		return 
	}

	selectQuery := `SELECT id FROM payflow_auth
					WHERE email = $1`
	err = h.DB.QueryRow(selectQuery, email).Scan(&userId)
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Something went wrong"})
		return 
	}


	token, err := h.GenerateJWT(userId, email, req.Role)
	if err!=nil{
		//c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create a new user"})
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Something went wrong"})
		return 
	}

	refreshToken, err := h.GenerateRefreshToken()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return 
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	insertQuery = `INSERT INTO payflow_refresh_tokens (user_id, token, expires_at)
    				VALUES ($1, $2, $3)`
	_, err = h.DB.Exec(insertQuery, userId, refreshToken, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	err = h.WalletClient.CreateWallet(userId)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create wallet"})
		return 
	}
	c.SetCookie(
		"refresh_token",
		refreshToken,
		30*24*60*60,
		"/",
		"",
		false, 
		true,  
	)


	c.JSON(http.StatusCreated, gin.H{"message":"User created successfully", "token":token, "role": req.Role})
}

func(h *Handler) SignIn(c *gin.Context) {
	var role string
	var userId string
	var userEmail string
	var req SignInRequest
	var hashedPassword string


	err := c.ShouldBindJSON(&req)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	if req.Email == "" || req.Password == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Email and password are required"})
		return 
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	//existing mail chech
	selectQuery := `SELECT id, email, password, role FROM payflow_auth 
					WHERE email = $1`
	
	err = h.DB.QueryRow(selectQuery, email).Scan(&userId, &userEmail, &hashedPassword, &role)

	if err == sql.ErrNoRows{
		c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid credentials"})
		return 
	}

	if err!=nil{
		//c.JSON(http.StatusInternalServerError, gin.H{"error":"Internal server error"})
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Something went wrong"})
		return 
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err!=nil{
		c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid credentials"})
		return 
	}

	token, err := h.GenerateJWT(userId, email, role)

	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create a new user"})
		return 
	}

	refreshToken, err := h.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	insertQuery := `INSERT INTO payflow_refresh_tokens (user_id, token, expires_at)
    				VALUES ($1, $2, $3)`
	_, err = h.DB.Exec(insertQuery, userId, refreshToken, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	c.SetCookie("refresh_token",
		refreshToken,
		30*24*60*60,
		"/",
		"",
		false,  
		true,  
	)

	c.JSON(http.StatusOK, gin.H{"message":"Login successful", "token":token, "role":role})
}

func(h *Handler) GenerateRefreshToken() (string, error){
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil{
		return "", err
	}

	    return base64.RawURLEncoding.EncodeToString(token), nil
}

func(h *Handler) Refresh(c *gin.Context){
	refreshToken, err := c.Cookie("refresh_token")
	var userId, role, email string
	var expiresAt time.Time
	if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing refresh token"})
        return
    }

	query := `
        SELECT rt.user_id, a.email, a.role, rt.expires_at
        FROM payflow_refresh_tokens rt
        JOIN payflow_auth a ON a.id = rt.user_id
        WHERE rt.token = $1
    `
	err = h.DB.QueryRow(query, refreshToken).Scan(&userId, &email, &role, &expiresAt)
	if err != nil || time.Now().After(expiresAt) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
        return
    }

	token, err := h.GenerateJWT(userId, email, role)
	if err != nil{
		 c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
        return
	}

	 c.JSON(http.StatusOK, gin.H{"token": token})
}