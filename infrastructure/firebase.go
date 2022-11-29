package infrastructure

import (
	"context"
	"path/filepath"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func NewFireBaseApp() *firebase.App {
	//Get the underlying context which is built with gin for our case
	ctx := context.Background()

	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")

	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}

	fbOption := option.WithCredentialsFile(serviceAccountKeyFilePath)

	//Firebase App initialization
	app, err := firebase.NewApp(ctx, nil, fbOption)
	if err != nil {
		panic("Firebase load error")
	}

	fsClient, fsErr := app.Firestore(ctx)
	if fsErr != nil {
		panic("Firestore error" + fsErr.Error())
	}
	defer fsClient.Close()

	return app
}

// NewFBAuth creates new
func NewFirebaseAuth(app *firebase.App) *auth.Client {

	//New copy of parent context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//New Firebase Auth client
	firebaseAuthClient, err := app.Auth(ctx)

	if err != nil {
		panic("Firebase Auth client err")
	}

	return firebaseAuthClient
}

// NewFirestoreClient creates new firestore client
func NewFirestoreClient(app *firebase.App) *firestore.Client {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//New Firebase Store client
	fireStoreClient, err := app.Firestore(ctx)
	if err != nil {
		panic("Fire Store Client err")
	}

	return fireStoreClient
}
