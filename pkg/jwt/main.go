package jwt

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(ttl time.Duration, payload interface{}, privateKey string, publicKey string) (string, error) {
	privateKeyData, err := b64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("createToken: decode: private key: %w", err)
	}
	publicKeyData, err := b64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("createToken: decode: public key: %w", err)
	}

	parsePrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", fmt.Errorf("createToken: parse: private key: %w", err)
	}
	parsePublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return "", fmt.Errorf("createToken: parse: public key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["sub"] = payload
	claims["ebf"] = now.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	token.Raw, err = token.SignedString(parsePrivateKey)
	if err != nil {
		return "", fmt.Errorf("createToken: signing string: %w", err)
	}

	_, err = jwt.ParseWithClaims(token.Raw, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return parsePublicKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("createToken: parse with claims: %w", err)
	}

	return token.Raw, nil
}

func ValidateStringToken(token string, publicKey string) (string, error) {
	publicKeyData, err := b64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("validateToken: decode: public key: %w", err)
	}

	parsedPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return "", fmt.Errorf("validateToken: parse: public key: %w", err)
	}

	parsedToken, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return parsedPublicKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("validateToken: parse token: %w", err)
	}

	claims, ok := parsedToken.Claims.(*jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return "", fmt.Errorf("validateToken: claims: %w", err)
	}

	return fmt.Sprintf("%v", (*claims)["sub"]), nil
}

func ValidateToken(token string, publicKey string) (uint, error) {
	publicKeyData, err := b64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return 0, fmt.Errorf("validateToken: decode: public key: %w", err)
	}

	parsedPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return 0, fmt.Errorf("validateToken: parse: public key: %w", err)
	}

	parsedToken, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return parsedPublicKey, nil
	})
	if err != nil {
		return 0, fmt.Errorf("validateToken: parse token: %w", err)
	}

	claims, ok := parsedToken.Claims.(*jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return 0, fmt.Errorf("validateToken: claims: %w", err)
	}

	return uint((*claims)["sub"].(float64)), nil
}

func CheckCookie(r *http.Request) (uint, error) {
	cookie, err := r.Cookie("vote")
	// check if err is ErrNoCookie
	if err != nil {
		if err == http.ErrNoCookie {
			return 0, nil
		} else {
			fmt.Println("Error at cookie: ", err)
			return 0, err
		}
	}
	token := cookie.Value

	if token == "" {
		return 0, nil
	}

	sub, err := ValidateToken(token, os.Getenv("PUBLIC_KEY"))
	if err != nil {
		fmt.Println("Error validating token: ", err)
		return 0, err
	}

	return sub, nil
}
