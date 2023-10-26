package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const TokenExpirDuation = time.Hour * 2

// define Secret
var mySecret = []byte("2401pt")

// 内嵌一个jwt.StandardClaims
type MyClaims struct {
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

// define exp(expiration time)
func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}

// summon access token and refresh token
func GenToken(username string) (aToken, rToken string, err error) {
	//Create a custom definition tailored(量身定制的) to our needs
	c := MyClaims{
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			//ExpiresAt: time.Now().Add(TokenExpirDuation).Unix(),
			Issuer: "2401pt", // Issuer
		},
	}
	// encrypt and get the string token encode
	// jwt.NewWithClaims is used to create a new JWT (JSON Web Token) object and associate it with the specified claims.
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)
	//and the function need two args.signing method and Claims(声明) .Claims include two members ,ExpiresAt (过期时间) iss(发行者)
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * 30).Unix(),
		Issuer:    "2401pt",
	}).SignedString(mySecret)
	return
}

// parse token and return MyClaims's point
func ParseToken(tokenString string) (claims *MyClaims, err error) {
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		err = errors.New("invalid token")
	}
	return
}

// refresh accessToken
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return "", "", err
	}
	var claims MyClaims
	// 解析访问Token，获取Claims信息
	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
	v, _ := err.(*jwt.ValidationError)
	//if not expired ,this will get new token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserName)
	}
	return
}
