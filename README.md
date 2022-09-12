# awesomego
This is a Golang webapplication, which include four requests:
1. "/name": get a random name from https://names.mcquay.me/api/v0/.    
The response data should look like ```{"first_name":"Patti","last_name":"Flamm"}```
2. "/joke": get a random joke from "http://api.icndb.com/jokes/random?firstName=John&lastName=Doe&limitTo=[nerdy]".    
The response data should look like ```{ "type": "success", "value": { "id": 543, "joke": "John Doe'ss programs can pass the Turing Test by staring at the interrogator.", "categories": ["nerdy"] } }```
3. "/": combine the results from step1 and step2, in which a joke from step2 will be returned to users with "John Doe" be replaced with the name from step1.    
The response should look like ```Carmine Wildfong's programs never exit, they terminate.```   
**Note: If the first_name or last_name we get from step1 is empty, replace will not happen, and server return a status code of 500 with a message**
4. "/status" : for liveness check purpose only   

## To run the project locally without docker
1. Navigate to the src directory
2. Run ```go run main.go```
3. Application will start and listen on port 8080, to change the port, please edit main.go
4. Visit the application with http://localhost:8080 by default   

## To run the project locally with docker 
1. Build a docker image with the Dockerfile. Run ```docker build -t {image_tag} ., example: docker build -t zxcarrot/awesomego:v2 .```
2. Run the application in docker container. Run ```docker run -i -t -p {port}:8080 {image_tag}, exmaple: docker run -i -t -p 5000:8080 zxcarrot/awesomego:v2```
3. Feel free to come up with your own image tag and port.
4. Visit the application with http://localhost:5000 by default   

## To run the project on k8s cluster
1. Push the image to docker hub. Run ```docker push {image_tag}, example: docker push zxcarrot/awesomego:v2```
2. Navigate to "deploy" directory   
2.1. Change the image name to the one you pushed.   
2.2. Run ``` kubectl apply -f deployment.yaml```. This file will help create a deployment with two replicas, and a hpa(HorizontalPodAutoscaler) with maximum of 10 replicas based on CPU and memory resource.
2.3 Run ```kubectl apply -f service.yaml```. This file will help create a LoadBalancer service, which help visit service
3. Visit the application with http://localhost:5000 on k8s cluster.   
**Note: the autoscaling was not tested due to my cluster limitation**


## Goals/Evaluation
* How will the application handle large spikes in traffic?   
autoscaling should take care of the traffic spike
* What kind of observability exists,and is it enough to debug production issues and aid incapacity planning?   
 No observability metric was added for now, but we should add:    
(1)CPU/memory/number of file descriptor in container   
(2)number of 5XX/4XX requests    
(3) number of container restart count    
* What steps have been taken to secure the service? Are there mitigations to common application security issues?   
This is a simple AOU service without considering security issue. For any public APIs, we should definitely add    
(1) Rate limiter in case of DDoS attack   
(2) If we want to limit anyone who is using the APIs, we can add another OAuth service to help authenticate

* How does the application respond if external service dependencies misbehave?    
For now, my API return with status code of 500 with information given tentatively. But it always open to discussion in this case   
* How will the application be deployed?
The application can be deployed within a k8s cluster by running the above command. In the ideal case, it should be managed by helm chars.


## Places can be improved
1. Log: currently using a simple log package without log level
2. Make it more configurable: replica numbers, images, port etc. are now hard codes, they should be easy to manage without changing source code
3. Rate limit: can be added to prevent DoS attack 
4. helm char: can be added to help manage deployment more automatically