# Distributed Systems - MapReduce in Go

## Setup
The individual setup for worker and master nodes is in the respective folders!

Common Setup needed - Google Cloud Storage (GCS)

As the system is heavily dependent on Google Cloud Storage it is recommended that you setup GCS to make the system usable out of the box.

### GCS Setup
1) Setup a new Google Cloud Store Bucket with name - ```track_1_files```
2) Upload the input files to path ```customer_trends/``` and make sure the files are of the following naming convention - ```##_customer_trends.txt``` ## being a positive integer.
(Note: Due to hardcoding the number of files read, the system only processes ten files - ```1_customer_trends.txt``` to ```10_customer_trends.txt```

3) Copy your Google Cloud Application Default Credentials to file ```google_cloud_credentials.json``` and place it in both worker-node and master-node folder - This is needed to authorize and actually be able to connect with the GCloud Storage
(Follow this guide - https://cloud.google.com/docs/authentication/application-default-credentials#personal)

This rounds up all the setup needed!

## Conclusion
I am grateful to work on this project and put theories into practice! This assignment/project helped me get up to speed with a lot of technology that I was unfamiliar with - Go, Google Cloud, Networking, Docker, and much more. Through trial and error, I have learned that the design of a distributed system can easily get complex with an amalgamation of many technologies and often lead to more problems than the original, therefore keeping it simple is the real challenge.

As some tech guru once said, "Often simple solutions are the hardest to come up with." 
