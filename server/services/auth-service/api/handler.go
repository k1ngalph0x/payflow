package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/k1ngalph0x/payflow/auth-service/config"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct{
	DB *sql.DB
	Config *config.Config
}

type SignUpRequest struct{
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type SignInRequest struct{
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type Claims struct{
	UserID string `json:"user_id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewHandler(db *sql.DB, cfg *config.Config) *Handler {
	return &Handler{DB: db, Config: cfg}
}

func(h *Handler) GenerateJWT(userId, email string)(string, error){

	expiration := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userId,
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

	// check the inputs
	if req.Email == "" || req.Password == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Email and password are required"})
		return 
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := strings.TrimSpace(req.Password)

	//existing mail check
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

	//hash the password
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
	INSERT INTO payflow_auth (email, password) 
	VALUES ($1, $2)
	`
	//_, err = h.DB.Exec(insertQuery, email, string(hashedPassword), created_at, updated_at)
	_, err = h.DB.Exec(insertQuery, email, string(hashedPassword))
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


	token, err := h.GenerateJWT(userId, email)
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create a new user"})
		return 
	}

	c.JSON(http.StatusCreated, gin.H{"message":"User created successfully", "token":token})
}

func(h *Handler) SignIn(c *gin.Context) {
	var req SignInRequest
	var userId string
	var userEmail string
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
	selectQuery := `SELECT id, email, password FROM payflow_auth 
					WHERE email = $1`
	
	err = h.DB.QueryRow(selectQuery, email).Scan(&userId, &userEmail, &hashedPassword)

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

	token, err := h.GenerateJWT(userId, email)

	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create a new user"})
		return 
	}


	c.JSON(http.StatusOK, gin.H{"message":"Login successful", "token":token})
}