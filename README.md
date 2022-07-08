Project - Build REST Service 

Overview 

For this project, you will build a REST service that meets the defined acceptance criteria for the user stories. Each user story is a separate task  within the project. This project is representative of functionality within FoodLogiQ Connect, but it not directly based on the implementations within  FoodLogiQ Connect. 

By the end of the project, you will have created REST APIs for create, read, delete and list commands utilizing events as well as providing an  easy way to build and deploy this within docker container (or docker compose). There is no required language or technology to implement this  project; use what you are comfortable with so we can evaluate your capability. 

In this project, you will build REST services around an event data structure as defined by the following JSON document. 
{ 
 "id": "ljadfj", // derived. Internal ID of  event 
 "createdAt": "2020-12-29T14:41:31.123Z", // derived. datetime the  event is created 
 "createdBy": "<userid>", // derived. id of the user  who created the event 
 "isDeleted": false, // derived. False when  created. True when deleted. 
 "type": "shipping", // required. valid entries  are shipping and receiving 
 "contents": [ 
 { 
 "gtin": "1234", // required. Global Trade  Item Number. 14-digit number. 
 "lot": "adffda", // required. any value. GTIN  + Lot are a compound identifier 
 "bestByDate": "2021-01-13", // optional. date value  "expirationDate": "2021-01-17", // optional. date value  }, 
 ... 
 ] 
}

Authentication 

REST APIs should be protected to avoid having an external party create issues with the data. For this project, you will use rudimentary  authentication system. Authentication will be performed using a bearer token in the Authorization header. 

There are only two potential users in  the system as defined below. 
Bearer Token 
UserID 
Business
74edf612f393b4eb01fbc2c29dd96671 
12345 
Acme
d88b4b1e77c70ba780b56032db1c259b 
98765 
Ajax



If any bearer token other than the above are received, return an authentication error status code. 

If you need more background on how bearer tokens are utilized, please refer to https://swagger.io/docs/specification/authentication/bearer authentication/.

Project Delivery 
In order to submit your project, please assemble all of your needed files into a github private repository. Please include instructions on how to  run and call each API within your project. As we are evaluating your code, please ensure the source code is in the repository and not just the  binaries needed to run the project. 
Please provide access to the github repository to the github user bonczj (Josh Bonczkowski). 

User Stories 
As an API consumer, I want a REST API to create an event 
As an API consumer, I want a REST API that I can pass a JSON document representing an event to. The API will validate the data structure and  data provided. If the data represents a well formed event, then the event must be persisted and success is returned. If the data is not a well  formed event, then the event is not persisted and an error is returned. 
Acceptance Criteria 
If the user authentication fails, return an authentication failure status code 
REST API is built that accepts JSON document 
Required values must be provided and are validated 
Derived values are applied before the event is stored 
If the optional date values are provided within the contents, the date values are validated as dates 
If all data is provided successfully, store the event internally and return an appropriate success status code 
If there are any errors, return an error message about the error and an appropriate failure status code 
As an API consumer, I want a REST API to delete a specific event 
Build a REST API that will attempt to remove access to an event. Due to audit constraints, we cannot just remove data from the platform. Instead,  if a user wants to remove data, it is soft deleted so that the information can be filtered out in the other APIs. 
Acceptance Criteria 
If the user authentication fails, return an authentication failure status code 
If the supplied event ID does not exist, return an appropriate failure status code 
If the supplied event ID does exist, but not accessible by the user, return an appropriate failure status code 
If the supplied event ID does exist and is accessible by the user, mark the event as deleted and return an appropriate success status  code 
If the event has already been deleted, consider it a success 
As an API consumer, I want a REST API to retrieve a specific event 
Build a REST API that will list a specific event given an ID value. You are building a multi-tenant REST API, so only return the data if the user  who is requesting the data is also the user who created the data. 
Acceptance Criteria 
If the user authentication fails, return an authentication failure status code 
If the ID provided is not found, return an appropriate failure status code 
If the ID provided is found, but was not created by the user who is making the REST API call, return the same failure as if the ID was not  found 
If the event of the ID has been deleted, return the same failure as if the ID was not found 
If the ID provided is found and the user making the REST API call also created the event, then return the event. 
As an API consumer, I want a REST API to list all of my events 
Build a REST API that will retrieve and list all events that the current user has access to. The API shall return a list of events (or an empty list) for  the user. The events should be sorted based on the createdAt datetime with the newest events first.
Acceptance Criteria 
If the user authentication fails, return an authentication failure status code 
If the user has no events, return an empty list 
If the user has events, return the list of events sorted based on createdAt with the newest events first 
If any of the events have been deleted, do not return them in the results 
As an DevOps engineer, I want to deploy this project utilizing docker / docker-compose 
At FoodLogiQ, out platform utilizes docker containers to run. This ensures the code is portable across servers and operating systems. For this  project, create a docker container that will assemble your REST API and any required resources. If you are utilizing additional services such as a  database or web server, then also a docker-compose file that will start up and connect all of the dependent services with each other. 
Acceptance Criteria 
Create a Dockerfile that will build and assemble your REST service 
If additional services are utilized, also include a docker-compose.yml file to ensure all services are wired up correctly Create a Makefile that can be used to simplify how to assemble your project
