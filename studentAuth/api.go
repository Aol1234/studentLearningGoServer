package studentAuth

import (
	"golang.org/x/net/context"

	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"log"
)

func InitializeAppWithServiceAccount(ctx context.Context) *firebase.App {
	// [START initialize_app_service_account_golang]
	opt := option.WithCredentialsFile("credential.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	// [END initialize_app_service_account_golang]
	return app
}

func VerifyUser(ctx context.Context, idToken string) (*auth.Token, error) {
	app := InitializeAppWithServiceAccount(ctx) // FIXME: Stop initializing service account each time this function called
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
		return nil, err
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
		return nil, err
	}
	return token, nil
}
