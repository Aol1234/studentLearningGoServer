package studentAuth

import (
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
)

func VerifyUser(ctx context.Context, idToken string) (*auth.Token, error) {
	// Verify users token
	app := InitializeAppWithServiceAccount(ctx) // Connect to Firebase
	client, err := app.Auth(ctx)
	if err != nil {
		log.Printf("error getting Auth client: %v\n", err)
		return nil, err
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Printf("error verifying ID token: %v\n", err)
		return nil, err
	}
	return token, nil
}

func InitializeAppWithServiceAccount(ctx context.Context) *firebase.App {
	// Open communication channel with Firebase
	opt := option.WithCredentialsFile("credential.json")
	app, err := firebase.NewApp(ctx, nil, opt) // Open Connection
	if err != nil {
		log.Printf("error initializing app: %v\n", err)
	}
	return app
}
