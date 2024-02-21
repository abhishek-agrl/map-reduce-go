package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	gcs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

var CloudStoreClient *gcs.Client

const number_of_reducers = 2
const bucket = "track_1_files"
const localDev = true

var worker_addr []string

func initGoogleCloudStore() *gcs.Client {
	ctx := context.Background()
	client, err := gcs.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	return client
}

func getIntermediateFileNames(bucketNumber int) string {
	log.Printf("Requesting all intermediate file names starting with bucket number %v from GCS", bucketNumber)
	var returnStr string

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	it := CloudStoreClient.Bucket("track_1_files").Objects(ctx, &gcs.Query{
		Prefix: "intermediate_files/" + strconv.Itoa(bucketNumber) + "_bucket",
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println("Bucket(track_1_files).Objects(): %w", err)
		}
		returnStr = returnStr + " " + strings.Split(attrs.Name, "/")[1]
	}

	return returnStr
}

func sendToWorker(fileName string, workType string, workerAddr string, workerAddrChan chan string) {
	log.Printf("%v is %v-ing %v", workerAddr, workType, fileName)
	reqBody := bytes.NewBuffer([]byte(fileName))
	resp, err := http.Post(workerAddr+"/"+workType, "application/text", reqBody)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Response From %v: %v", workerAddr, string(responseBody))

	workerAddrChan <- workerAddr

}

func startMapping() {
	log.Printf("Mapping Stage Started")
	workers := make(chan string, 8)
	var wg sync.WaitGroup

	for _, addr := range worker_addr {
		workers <- addr
	}

	for i := 1; i < 11; i++ {
		fileName := strconv.Itoa(i) + "_customer_trends.txt"
		wg.Add(1)
		go func() {
			defer wg.Done()
			sendToWorker(fileName, "map", <-workers, workers)
		}()
	}
	wg.Wait()
}

func startReducing() {
	log.Printf("Reducing Stage Started")
	workers := make(chan string, number_of_reducers)
	var wg sync.WaitGroup

	for i := 0; i < number_of_reducers; i++ {
		workers <- worker_addr[i]
	}

	for i := 0; i < 2; i++ {
		fileNames := getIntermediateFileNames(i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			sendToWorker(fileNames, "reduce", <-workers, workers)
		}()
	}
	wg.Wait()
}

func cleanup() {
	log.Println("Cleanup Started")
	err := deleteFile()
	if err != nil {
		log.Fatalf("Error while cleanup: %v", err)
	}
	log.Println("Cleanup complete")
}

func Home(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received to Home Page")
	fmt.Fprint(w, "Welcome to My Implementation of Map Reduce. To start, call endpoint /mapreduce")
}

func StartMapReduce(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received to Start Map Reduce")
	fmt.Fprintln(w, "Map-Reduce Started")

	cleanup()
	log.Println("Old Intermediate Files Clean-up Complete")
	fmt.Fprintln(w, "1) Old Intermediate Files Clean-up Complete")

	startMapping()
	log.Println("Mapping Stage Complete")
	fmt.Fprintln(w, "2) Mapping Stage Complete")

	startReducing()
	log.Printf("Reducing Stage Complete, Final output generated.")
	fmt.Fprintln(w, "3) Reducing Stage Complete, Final output available at /output")
}

func ManualCleanup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Manual clean up of old intermediate files started")
	cleanup()
	fmt.Fprintln(w, "Manual clean up complete")
}

func Output(w http.ResponseWriter, r *http.Request) {
	log.Println("Requesting all output files from GCS")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	it := CloudStoreClient.Bucket("track_1_files").Objects(ctx, &gcs.Query{
		Prefix: "output_files/",
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println("Bucket(track_1_files).Objects(): %w", err)
		}
		fileName := strings.Split(attrs.Name, "/")[1]
		if fileName != "" {
			fmt.Fprintf(w, "[%v]\n%v\n", fileName, string(getFileData("output_files", fileName)))
		}
	}
}

func Evaluation(w http.ResponseWriter, r *http.Request) {
	//TBD
}

func main() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./google_cloud_credentials.json")

	log.Println("Initiatiing Worker Node Addresses")
	worker_addr = getWorkerNodes(localDev)
	CloudStoreClient = initGoogleCloudStore()
	log.Println("Google Cloud Store Client Ready")
	defer CloudStoreClient.Close()

	http.HandleFunc("/", Home)
	http.HandleFunc("/mapreduce", StartMapReduce)
	http.HandleFunc("/cleanup", ManualCleanup)
	http.HandleFunc("/output", Output)
	http.HandleFunc("/evaluation", Evaluation)

	log.Println("Starting New HTTP Server at Port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
