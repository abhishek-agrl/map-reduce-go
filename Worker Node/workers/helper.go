package workers

import (
	"context"
	"fmt"
	"log"
	"strings"

	gcs "cloud.google.com/go/storage"
	"github.com/google/uuid"
)

func initGoogleCloudStore() *gcs.Client {
	ctx := context.Background()
	client, err := gcs.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	return client
}

func writeToIntermediate(buckets [number_of_reducers][]WordMap, cloudStoreClient *gcs.Client) {
	for i, bucket := range buckets {
		// bucketsStr[i] = convertMaptoText(bucket)
		var fileNameBuilder strings.Builder
		uniqueId := uuid.New().String()
		fmt.Fprintf(&fileNameBuilder, "%v_bucket_%v.gob", i, uniqueId)
		appendToFile("intermediate_files", fileNameBuilder.String(), convertMaptoText(bucket), cloudStoreClient)
	}
}
func convertMaptoText(inputArr []WordMap) string {
	log.Println("Converting Array of Custom Objects to String Output")
	var sb strings.Builder

	for _, wordMap := range inputArr {
		fmt.Fprintf(&sb, "%v %v\n", wordMap.Word, wordMap.Count)
	}
	log.Println("String Output Ready")
	return sb.String()
}
