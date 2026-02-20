package api

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
	"github.com/k1ngalph0x/payflow/identity-service/config"
	"github.com/k1ngalph0x/payflow/identity-service/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct{
	DB  *gorm.DB
	Config *config.Config
	WalletClient *walletclient.WalletClient
}

type SignUpRequest struct{
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role"     binding:"required,oneof=user merchant"`
}

type SignInRequest struct{
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Claims struct{
	UserID string `json:"user_id"`
	Email string `json:"email"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewHandler(db *gorm.DB, config *config.Config, walletclient *walletclient.WalletClient) *Handler {
	return &Handler{DB: db, Config: config, WalletClient: walletclient}
}

func(h *Handler) generateJWT(userID, email, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "payflow-auth",
			Subject:   userID,
	},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.Config.TOKEN.JwtKey))
}

func(h *Handler) SignUp(c *gin.Context) {

	var req SignUpRequest
	err := c.ShouldBindJSON(&req)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	var existing models.User
	result := h.DB.Where("email = ?", email).First(&existing)
	if result.Error == nil{
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}else if !errors.Is(result.Error, gorm.ErrRecordNotFound){
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return 
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return 
	}
	
	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	result = h.DB.Create(&user)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, err := h.generateJWT(user.ID, user.Email, user.Role)
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return 
	}

	refreshToken, err := h.createRefreshToken(user.ID)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return 
	}

	err = h.WalletClient.CreateWallet(user.ID)
	if err != nil{
		//c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create wallet"})
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return 
	}

	h.setRefreshCookie(c, refreshToken)
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "token": token, "role": user.Role})
}

func(h *Handler) SignIn(c *gin.Context) {

	var req SignInRequest
	err := c.ShouldBindJSON(&req)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	result := h.DB.Where("email = ?", email).First(&user)
	if result.Error != nil{
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		}
		return 
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err!=nil{
		c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid credentials"})
		return 
	}

	token, err := h.generateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	refreshToken, err := h.createRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	h.setRefreshCookie(c, refreshToken)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token, "role": user.Role})
}

func(h *Handler) createRefreshToken(userID string) (string, error){

	raw := make([]byte, 32)
	_, err := rand.Read(raw)
	if err != nil{
		return "", err
	}
	
	tokenStr := base64.RawURLEncoding.EncodeToString(raw)
	refreshtoken := models.RefreshToken{
		UserID:    userID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	result := h.DB.Create(&refreshtoken)
	if result.Error != nil{
		return "", err
	}

	return tokenStr, nil
}


func (h *Handler) setRefreshCookie(c *gin.Context, token string){
	c.SetCookie(
		"refresh_token",
		token,
		30*24*60*60,
		"/", 
		"", 
		false, 
		true,
	)
}

func (h *Handler) Refresh(c *gin.Context){
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing refresh token"})
		return
	}

	var refreshToken models.RefreshToken
	result := h.DB.Preload("User").Where("token = ?", cookie).First(&refreshToken)
	if result.Error != nil || time.Now().After(refreshToken.ExpiresAt){
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	token, err := h.generateJWT(refreshToken.User.ID, refreshToken.User.Email, refreshToken.User.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}