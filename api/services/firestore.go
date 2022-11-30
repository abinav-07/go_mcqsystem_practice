package services

import (
	"context"
	"fmt"
	"github/abinav-07/mcq-test/infrastructure"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
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
	_, err := fireStoreClient.Collection(entityName).Doc(stringId).Set(context.Background(), docData, firestore.MergeAll)

	if err != nil {
		return "", err
	}
	return "Document Data Created!", nil
}

func (fsc *FireStoreService) UpdateEntityFieldWithId(entityName string, id uint, docData map[string]interface{}) (string, error) {
	fireStoreClient := fsc.GetFireStoreClient()
	defer fireStoreClient.Close()

	stringId := strconv.FormatUint(uint64(id), 10)
	_, err := fireStoreClient.Collection(entityName).Doc(stringId).Set(context.Background(), docData, firestore.MergeAll)

	if err != nil {
		return "", err
	}
	return "Document Data Updated!", nil
}

func (fsc *FireStoreService) GetEntityWithId(entityName string, id uint) (*firestore.DocumentSnapshot, error) {
	fireStoreClient := fsc.GetFireStoreClient()
	defer fireStoreClient.Close()

	stringId := strconv.FormatUint(uint64(id), 10)
	docData, err := fireStoreClient.Collection(entityName).Doc(stringId).Get(context.Background())
	return docData, err
}

func (fsc *FireStoreService) DeleteCollectionWithId(entityName string, id uint) error {
	fireStoreClient := fsc.GetFireStoreClient()
	defer fireStoreClient.Close()
	stringId := strconv.FormatUint(uint64(id), 10)
	for {
		iter := fireStoreClient.Collection(entityName).Documents(context.Background())

		numDeleted := 0

		batch := fireStoreClient.Batch()

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			if doc.Ref.ID == stringId {
				batch.Delete(doc.Ref)
				numDeleted++
			}

		}

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(context.Background())
		if err != nil {
			return err
		}

	}

}
