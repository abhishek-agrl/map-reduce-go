package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	gcs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func getFileData(fileLocation string, fileName string) []byte {
	log.Printf("Requesting File Data from Google Cloud Store Bucket %v: Filename: %v", fileLocation, fileName)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	rc, err := CloudStoreClient.Bucket(bucket).Object(fileLocation + "/" + fileName).NewReader(ctx)
	if err != nil {
		log.Printf("File Not Found: %v/%v", fileLocation, fileName)
		return nil
	}
	defer rc.Close()

	buf := bytes.NewBuffer([]byte{})
	if _, err := io.Copy(buf, rc); err != nil {
		log.Fatal(err)
	}
	log.Printf("File Data received from Google Cloud Store Bucket %v: Filename: %v", fileLocation, fileName)
	return buf.Bytes()

}

func deleteFile() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	cloudBucket := CloudStoreClient.Bucket(bucket)
	it := cloudBucket.Objects(ctx, &gcs.Query{
		Prefix: "intermediate_files/",
	})
	var fileNames []string
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Bucket(track_1_files).Objects(): %w", err)
		}
		fileName := strings.Split(attrs.Name, "/")[1]
		if fileName != "" {
			fileNames = append(fileNames, fileName)
		}
	}

	for _, fileName := range fileNames {
		object := "intermediate_files/" + fileName
		o := cloudBucket.Object(object)

		if err := o.Delete(ctx); err != nil {
			return fmt.Errorf("Object(%q).Delete: %w", object, err)
		}
		log.Printf("Deleted File: %v", object)
	}

	return nil
}

func getWorkerNodes(localDev bool) []string {
	if localDev {
		return []string{
			"http://localhost:8001",
			"http://localhost:8002",
			"http://localhost:8003",
			"http://localhost:8004",
			"http://localhost:8005",
			"http://localhost:8006",
			"http://localhost:8007",
			"http://localhost:8008",
		}

	}
	return []string{
		"http://worker-node-1.europe-west3-a.c.distributed-systems-412017.internal:8080",
		"http://worker-node-2.europe-west3-a.c.distributed-systems-412017.internal:8080",
		"http://worker-node-3.europe-west3-a.c.distributed-systems-412017.internal:8080",
		"http://worker-node-4.europe-west3-a.c.distributed-systems-412017.internal:8080",
		"http://worker-node-5.europe-west3-a.c.distributed-systems-412017.internal:8080",
		"http://worker-node-6.europe-west3-a.c.distributed-systems-412017.internal:8080",
		"http://worker-node-7.europe-west3-a.c.distributed-systems-412017.internal:8080",
	}
}
