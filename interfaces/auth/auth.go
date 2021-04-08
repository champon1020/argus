package auth

import (
	"crypto/rsa"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// VerifyJWTToken verifies jwt token.
func VerifyJWTToken(token string) error {
	user, err := readUserInfo()
	if err != nil {
		return err
	}

	tok, err := decodeJwt(token)
	if err != nil {
		return err
	}

	if err := validate(user, tok); err != nil {
		return err
	}

	return nil
}

func validate(user *authUser, tok jwt.MapClaims) error {
	if user.Iss != tok["iss"] ||
		user.Azp != tok["azp"] ||
		user.Aud != tok["aud"] ||
		user.Sub != tok["sub"] ||
		user.Email != tok["email"] {
		return errors.New("failed to verify token")
	}
	return nil
}

type authUser struct {
	Iss   string `json:"iss"`
	Azp   string `json:"azp"`
	Aud   string `json:"aud"`
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

type secret struct {
	Web struct {
		CertURL string `json:"auth_provider_x509_cert_url"`
	} `json:"web"`
}

func readUserInfo() (*authUser, error) {
	file := os.Getenv("ARGUS_USER_PATH")
	body, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	user := new(authUser)
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return user, nil
}

func decodeJwt(token string) (jwt.MapClaims, error) {
	pubKey, err := readPublicKey(token)
	if err != nil {
		return nil, err
	}

	tok, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !tok.Valid {
		return nil, errors.New("invalid token")
	}

	claim, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return claim, nil
}

func readPublicKey(token string) (*rsa.PublicKey, error) {
	file := os.Getenv("ARGUS_PUBLIC_KEY_PATH")

	if _, err := os.Stat(file); err != nil {
		if err := fetchPublicKey(token); err != nil {
			return nil, err
		}
	}

	body, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(body)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func fetchPublicKey(token string) error {
	// Read secret file.
	secret, err := readSecret()
	if err != nil {
		return err
	}

	// Convert token to map.
	tok, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return err
	}
	if tok == nil {
		return errors.New("invalid token")
	}

	// Download public key via web.
	if err := downloadPublicKey(secret.Web.CertURL, tok.Header["kid"].(string)); err != nil {
		return err
	}

	return nil
}

func readSecret() (*secret, error) {
	file := os.Getenv("ARGUS_CLIENT_SECRET_PATH")
	body, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	sec := &secret{}
	if err := json.Unmarshal(body, sec); err != nil {
		return nil, err
	}

	return sec, nil
}

func downloadPublicKey(url string, kid string) error {
	client := new(http.Client)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	keys := make(map[string]string)
	if err := json.Unmarshal(body, &keys); err != nil {
		return err
	}

	if err := saveKeyAsPem(keys[kid]); err != nil {
		return err
	}

	return nil
}

func saveKeyAsPem(key string) error {
	block, _ := pem.Decode([]byte(key))
	file := os.Getenv("ARGUS_PUBLIC_KEY_PATH")

	k, err := os.Create(file)
	if err != nil {
		return err
	}

	if err := pem.Encode(k, block); err != nil {
		return err
	}

	return nil
}
