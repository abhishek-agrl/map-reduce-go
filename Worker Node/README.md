# Worker Node - Operational Node for MapReduce

## Setup
1) Complete basic setup of Google Cloud Storage
2) Make sure you have Go v1.21
3) Run command - ```go mod tidy```
4) Copy your google_cloud_credentials.json to the folder
5) Make sure the worker node is functional, I recommend running the code in VS Code in the run and debug section (Make sure you run and debug while on the main.go file) and looking for any errors in the console. If the worker was able to successfully boot, the debug console should print -
   ```
   YYYY/MM/DD HH:MM:SS New Worker Node Initiated
   YYYY/MM/DD HH:MM:SS Google Cloud Store Client Ready
   YYYY/MM/DD HH:MM:SS Starting New HTTP Server at Port 8080
    ```

To run in multiple concurrent mode, we can use docker -
1) In the folder with the docker file run - ```docker image build -t worker-node:0.0.1 .```
2) It takes a while but once done, you should have a docker image of the worker node, now, simply just run the image 8 times, to mimic 8 worker nodes. Recommended way -
   ```
   docker run -p 8001:8080 worker-node:0.0.1
   docker run -p 8002:8080 worker-node:0.0.1
   docker run -p 8003:8080 worker-node:0.0.1
   docker run -p 8004:8080 worker-node:0.0.1
   docker run -p 8005:8080 worker-node:0.0.1
   docker run -p 8006:8080 worker-node:0.0.1
   docker run -p 8007:8080 worker-node:0.0.1
   docker run -p 8008:8080 worker-node:0.0.1
   ```
By using these port configurations, no changes in the master node would be required.
(Note: Due to hardcoding in master for some values, the system won't function properly until 8 workers are running)
