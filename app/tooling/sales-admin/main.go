package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	if err := gentoken(); err != nil {
		log.Fatal(err)
	}
}

// func genkey() error {
// 	// ================================
// 	// Create a private key file.
// 	// ================================

// 	// Generate a new private key.
// 	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
// 	if err != nil {
// 		return fmt.Errorf("generating key: %w", err)
// 	}

// 	// Create a file for the private key information in PEM form.
// 	privateFile, err := os.Create("private.pem")
// 	if err != nil {
// 		return fmt.Errorf("creating private file: %w", err)
// 	}
// 	defer privateFile.Close()

// 	// Construct a PEM block for the private key.
// 	privateBlock := pem.Block{
// 		Type:  "PRIVATE KEY",
// 		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
// 	}

// 	// Write the private key to the private key file.
// 	if err := pem.Encode(privateFile, &privateBlock); err != nil {
// 		return fmt.Errorf("encoding to private file: %w", err)
// 	}

// 	// ================================
// 	// Create a public key file.
// 	// ================================

// 	// Create a file for the public key information in PEM form.
// 	publicFile, err := os.Create("public.pem")
// 	if err != nil {
// 		return fmt.Errorf("creating public file: %w", err)
// 	}
// 	defer publicFile.Close()

// 	// Marshal the public key from the private key to PKIX.
// 	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
// 	if err != nil {
// 		return fmt.Errorf("marshaling public key: %w", err)
// 	}

// 	// Construct a PEM block for the public key.
// 	publicBlock := pem.Block{
// 		Type:  "PUBLIC KEY",
// 		Bytes: asn1Bytes,
// 	}

// 	// Write the public key to the public key file.
// 	if err := pem.Encode(publicFile, &publicBlock); err != nil {
// 		return fmt.Errorf("encoding to public file: %w", err)
// 	}

// 	return nil
// }

func gentoken() error {

	file, err := os.Open("infra/keys/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1.pem")
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	// limit PEM file size to 1 megabyte. This should be reasonable for
	// almost any PEM file and prevents shenanigans like linking the file
	// to /dev/random or something like that.
	pem, err := io.ReadAll(io.LimitReader(file, 1024*1024))
	if err != nil {
		return fmt.Errorf("reading auth private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		return fmt.Errorf("parsing private key: %w", err)
	}

	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
	claims := struct {
		jwt.RegisteredClaims
		Roles []string `json:"roles"`
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1234567890",
			Issuer:    "sales-api",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}

	method := jwt.GetSigningMethod(jwt.SigningMethodRS256.Name)
	// kid := "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"
	// Create a new JWT token with the claims and the signing method

	token := jwt.NewWithClaims(method, claims)
	token.Header["kid"] = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"

	// if err != nil {
	// 	return fmt.Errorf("parsing private pem: %w", err)
	// }

	// Sign the token with the private key.
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("signing token: %w", err)
	}
	fmt.Printf("-----BEGIN TOKEN-----\n%s\n-----END TOKEN-----\n\n", tokenString)

	// ------------------------------------------------------------

	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}))

	claims2 := struct {
		jwt.RegisteredClaims
		Roles []string `json:"roles"`
	}{}

	parsedToken, err := parser.ParseWithClaims(tokenString, &claims2, func(token *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})
	if err != nil {
		return fmt.Errorf("parsing token: %w", err)
	}

	if !parsedToken.Valid {
		return fmt.Errorf("INVALID TOKEN")
	}

	fmt.Printf("-----BEGIN CLAIMS-----\n%s\n-----END CLAIMS-----\n\n", claims2)

	// // ------------------------------------------------------------

	// // ------------------------------------------------------------

	// // claims3 := struct {
	// // 	jwt.RegisteredClaims
	// // 	Roles []string `json:"roles"`
	// // }{}

	// // parsedToken2, _, err := parser.ParseUnverified(tokenString, &claims3)
	// // if err != nil {
	// // 	return fmt.Errorf("error parsing token unverified: %w", err)
	// // }

	// // Marshal the public key from the private key to PKIX.
	// asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	// if err != nil {
	// 	return fmt.Errorf("marshaling public key: %w", err)
	// }

	// // Construct a PEM block for the public key.
	// publicBlock := pem.Block{
	// 	Type:  "PUBLIC KEY",
	// 	Bytes: asn1Bytes,
	// }

	// var b bytes.Buffer

	// // Write the public key to the public key file.
	// if err := pem.Encode(&b, &publicBlock); err != nil {
	// 	return fmt.Errorf("encoding to public file: %w", err)
	// }

	// input := map[string]any{
	// 	"Key":   b.String(),
	// 	"Token": tokenString,
	// }

	// if err := opaPolicyEvaluation(context.Background(), opaAuthentication, input); err != nil {
	// 	return fmt.Errorf("authentication failed : %w", err)
	// }

	// fmt.Println("Authentication successful")
	return nil
}

// // opaPolicyEvaluation asks opa to evaluate the token against the specified token
// // policy and public key.
// func opaPolicyEvaluation(ctx context.Context, regoScript string, input any) error {
// 	query := fmt.Sprintf("x = data.%s.%s", "foodflow.rego", "auth")

// 	q, err := rego.New(
// 		rego.Query(query),
// 		rego.Module("policy.rego", regoScript),
// 	).PrepareForEval(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	results, err := q.Eval(ctx, rego.EvalInput(input))
// 	if err != nil {
// 		return fmt.Errorf("query: %w", err)
// 	}

// 	if len(results) == 0 {
// 		return errors.New("no results")
// 	}

// 	result, ok := results[0].Bindings["x"].(bool)
// 	if !ok || !result {
// 		return fmt.Errorf("bindings results[%v] ok[%v]", results, ok)
// 	}

// 	return nil
// }
