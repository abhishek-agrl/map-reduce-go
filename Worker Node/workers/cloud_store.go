package workers

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
	"log"
	"time"

	gcs "cloud.google.com/go/storage"
)

func getFileData(fileLocation string, fileName string, client *gcs.Client) []byte {
	log.Printf("Requesting File Data from Google Cloud Store Bucket %v: Filename: %v", fileLocation, fileName)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucket).Object(fileLocation + "/" + fileName).NewReader(ctx)
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

func appendToFile(fileLocation string, fileName string, input string, client *gcs.Client) {
	log.Printf("Writing File Data to Google Cloud Store Bucket %v: Filename: %v", fileLocation, fileName)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	log.Println("Compressing before writing to GCS")
	file := &bytes.Buffer{}
	gob.NewEncoder(file).Encode(input)
	log.Println("Compressing complete")

	wc := client.Bucket(bucket).Object(fileLocation + "/" + fileName).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		log.Printf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		log.Printf("Writer.Close: %v", err)
	}
	log.Println("Write successful")
}

func writeToOutputFile(fileLocation string, fileName string, input string, client *gcs.Client) {
	log.Printf("Writing File Data to Google Cloud Store Bucket %v: Filename: %v", fileLocation, fileName)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	wc := client.Bucket(bucket).Object(fileLocation + "/" + fileName).NewWriter(ctx)
	if _, err := io.Copy(wc, bytes.NewBufferString(input)); err != nil {
		log.Printf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		log.Printf("Writer.Close: %v", err)
	}
	log.Println("Write successful")
}
