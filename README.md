# awesomego
This is a Golang webapplication, which include four requests:
1. "/name": get a random name from https://names.mcquay.me/api/v0/. The result should look like
   ```{"first_name":"Patti","last_name":"Flamm"}```
2. "/joke": get a random joke from "http://api.icndb.com/jokes/random?firstName=John&lastName=Doe&limitTo=[nerdy]". 
The result should look like ```{ "type": "success", "value": { "id": 543, "joke": "John Doe'ss programs can pass the Turing Test by staring at the interrogator.", "categories": ["nerdy"] } }```
3. "/": combine the results from step1 and step2, in which a joke from step2 will be returned to users with "John Doe" be replaced with the name from step1
4. "/status" : for liveness check purpose

## To run the project locally without docker
1. Navigate to the src directory
2. Run ```go run main.go```
3. Application will start and listen on port 8080, to change the port, please edit main.go
## To run the project locally with docker 
1. Build a docker image with the Dockerfile, run ```docker build -t zxcarrot/awesomego:v2 .```
2. Run the application in docker container, run ```docker run -i -t -p 5000:8080 zxcarrot/awesomego:v2```
   

2. 