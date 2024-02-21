package workers

import (
	"fmt"
	"log"
	"os"
	"strings"

	gcs "cloud.google.com/go/storage"
)

type WordMap struct {
	Word  string
	Count uint
}

const (
	number_of_workers  = 5
	number_of_mappers  = 3
	number_of_reducers = 2
	bucket             = "track_1_files"
)

var CloudStoreClient *gcs.Client

func Mapper(fileName string) {
	inputStr := getFileData("customer_trends", fileName, CloudStoreClient)
	mapOutput := startMap(string(inputStr))
	log.Println("Unique words frequency calculation complete, writing to GCS")
	writeToIntermediate(mapOutput, CloudStoreClient)
}

func Reducer(bucketNames []string) {

	log.Printf("Getting all intermediate files for bucket %v", bucketNames[0][0])
	var bucketGobs [][]byte
	for _, bucketName := range bucketNames {
		bucketGobs = append(bucketGobs, getFileData("intermediate_files", bucketName, CloudStoreClient))
	}

	log.Printf("Starting Reduce Computation for %v files", len(bucketGobs))
	reducedOutput := startReducer(bucketGobs)

	var outputFileName strings.Builder
	fmt.Fprintf(&outputFileName, "%v_output.txt", string(bucketNames[0][0]))
	log.Printf("Reducer Computation Complete, writing to GCS: output_files/%v", outputFileName.String())

	writeToOutputFile("output_files", outputFileName.String(), reducedOutput, CloudStoreClient)
}

func InitWorker() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./google_cloud_credentials.json")
	CloudStoreClient = initGoogleCloudStore()
	log.Println("Google Cloud Store Client Ready")
}
