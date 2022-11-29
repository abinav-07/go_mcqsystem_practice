package services

import (
	"context"
	"fmt"
	"github/abinav-07/mcq-test/infrastructure"
	"strconv"

	"cloud.google.com/go/firestore"
)

type FireStoreService struct {
	fireStoreClient *firestore.Client
	env             infrastructure.Env
}

func NewFireStoreClient(fireStoreClient *firestore.Client, env infrastructure.Env) FireStoreService {
	return FireStoreService{
		fireStoreClient: fireStoreClient,
		env:             env,
	}
}

func (fsc *FireStoreService) GetFireStoreClient() *firestore.Client {
	fireStoreClient, err := firestore.NewClient(context.Background(), fsc.env.FireStoreProject)

	if err != nil {
		fmt.Println("FATAL: Failed to create firestore client", err)

	}

	return fireStoreClient
}

func (fsc *FireStoreService) SaveOrUpdateEntityWithId(entityName string, id uint, docData map[string]interface{}) (string, error) {
	fireStoreClient := fsc.GetFireStoreClient()
	defer fireStoreClient.Close()

	stringId := strconv.FormatUint(uint64(id), 10)
	_, err := fireStoreClient.Collection(entityName).Doc(stringId).Set(context.Background(), docData)

	if err != nil {
		return "", err
	}
	return "Document Data Created!", nil
}

func (fsc *FireStoreService) GetEntityWithId(entityName string, id uint) (*firestore.DocumentSnapshot, error) {
	fireStoreClient := fsc.GetFireStoreClient()
	defer fireStoreClient.Close()

	stringId := strconv.FormatUint(uint64(id), 10)
	docData, err := fireStoreClient.Collection(entityName).Doc(stringId).Get(context.Background())
	return docData, err
}
