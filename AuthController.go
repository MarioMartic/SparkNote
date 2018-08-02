package main

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"log"
	"io/ioutil"
	"time"
	"github.com/dgrijalva/jwt-go/request"
	"crypto/rsa"
	"golang.org/x/crypto/bcrypt"
	"database/sql"
)

const (
	privKeyPath = "keys/app.rsa"
	pubKeyPath  = "keys/app.rsa.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)
var verifyBytes, signBytes []byte

func initKeys() {
	var err error

	signBytes, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatalf("Error reading private key: %v", err)
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatalf("Error parsing private key: %v", err)
	}
	verifyBytes, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatalf("Error reading public key: %v", err)
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatalf("Error parsing public key: %v", err)
	}
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	user := User{Username:username}

	err := user.findByUsername(a.DB)

	if err != nil {
		fmt.Println("Invalid username.")
		fmt.Println(err.Error())

		respondWithError(w, http.StatusForbidden, "Wrong username/password.")

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		fmt.Println("Invalid password.")
		fmt.Println(err.Error())

		respondWithError(w, http.StatusForbidden, "Wrong username/password.")

		return
	}

	//create a rsa 256 signer and set claims
	signer := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "admin",
		"exp": time.Now().Add(time.Minute * 60).Unix(),
		"id": user.ID,
		"username": user.Username,
	})


	tokenString, err := signer.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		log.Printf("Error signing token: %v\n", err)
	}

	//create a token instance using the token string
	response := Token{tokenString}
	respondWithJSON(w, http.StatusOK, response)

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")
	name := r.FormValue("name")

	log.Print(username, password, name)

	var usernameFromDB string
	err := a.DB.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&usernameFromDB)

	if usernameFromDB != "" {
		http.Error(w, "Username already exist.", http.StatusForbidden)
		respondWithError(w, http.StatusForbidden, "Username already exist!")
		log.Println("Username in use!")
		return
	}

	switch {

	case err == sql.ErrNoRows:

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if err != nil {
			http.Error(w, "Server error, unable to create your account.", 500)
			respondWithError(w, http.StatusForbidden, "Server error.")
			log.Println(err.Error())
			return
		}

		user := User{Username:username, Password:hashedPassword, Name:name}

		err = user.createUser(a.DB)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Server error, unable to create your account.")
			log.Println(err.Error())
			return
		}

		w.Write([]byte("User created!"))
		return
	case err != nil:
		respondWithError(w, http.StatusInternalServerError, "Server error, unable to create your account.")
		log.Println(err.Error())
		return
	default:
		respondWithError(w, http.StatusInternalServerError,"Unknown error.")
		log.Println(err.Error())
	}
}

//AUTH TOKEN VALIDATION

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//validate token
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err == nil {

		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorised access to this resource")
	}

}

func getUserFromToken(w http.ResponseWriter, r *http.Request) (User, error){

	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return User{}, err
	}

	tokenString := token.Raw

	claims, ok := extractClaims(tokenString)

	if ok == false {
		return User{}, fmt.Errorf("token error")
	}

	id:= int(claims["id"].(float64))
	var u = User { ID: id, Username: claims["username"].(string) }

	return u, nil

}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil {
		log.Printf(err.Error())
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Printf("Returned claims")
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

