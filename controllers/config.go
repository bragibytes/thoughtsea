package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	POST          string = "POST"
	GET           string = "GET"
	PUT           string = "PUT"
	DELETE        string = "DELETE"
	AUTHORIZATION string = "Authorization"
)

// Define the payload for the JWT token
type Claims struct {
	ID primitive.ObjectID `json:"id"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token for the given user
func GenerateJWT(id primitive.ObjectID) (string, error) {
	// Define the expiration time for the token (in this example, 24 hours)
	expireTime := time.Now().Add(time.Hour * 24).Unix()

	// Create the claims object for the JWT token
	claims := Claims{
		id,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			IssuedAt:  time.Now().Unix(),
			Subject:   "auth_token",
		},
	}

	// Generate the JWT token using the claims object and the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// AuthMiddleware makes sure there is a valid token in the header before allowing access to a route
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user is authenticated
		if userIsAuthenticated(c) {
			// If the user is authenticated, allow them to access the protected route
			c.Next()
		} else {
			// If the user is not authenticated, return a 401 Unauthorized status
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

// usersMatch checks to make sure the current user is authorized to manipulate some data
func usersMatch(c *gin.Context, id primitive.ObjectID) bool {
	headerValue, exists := c.Get(AUTHORIZATION)
	if !exists {
		// handle case where header does not exist
		return false
	}
	oid, ok := headerValue.(primitive.ObjectID)
	if !ok {
		// handle case where header value is not a valid ObjectID
		return false
	}
	if oid == id {
		return true
	}
	return false
}

// userIsAuthenticated checks to make sure the token is valid, and if it is replaces the Authorization header with the user id
func userIsAuthenticated(c *gin.Context) bool {
	tokenString := c.GetHeader(AUTHORIZATION)
	if tokenString == "" {
		return false
	}
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the signing method used to sign the token is HMAC.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key used to sign the token.
		return []byte(os.Getenv("TOKEN_KEY")), nil
	})
	if err != nil {
		fmt.Println("error parsing token ", err.Error())
		return false
	}

	// Get the claims from the token.
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		id := claims["id"].(string)
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			fmt.Println("error converting hex string to oid")
			return false
		}
		c.Set(AUTHORIZATION, oid)
		fmt.Println("Token is valid")
		return true
	} else {
		fmt.Println("Invalid token.")
		return false
	}
}

// response is the standard response object for all requests
type response struct {
	Success bool        `json:"success"`
	Payload interface{} `json:"payload"`
	Message string      `json:"message"`
}
