# user-service
This is a small microservice to handle Users in Golang

## Running the service
To run the service call ```go run main.go``` in the root directory. It should start the service on port 8080, but you can change it using a configuration file.

There are 6 routes for this microservice
```
GET /
GET /user/{id}
PUT /user/{id}
POST /user
DELETE /user/{id}
GET /search/{criteria}/{search}
```

The responses are all JSON, including errors. The input for the POST consumes application/x-www-form-urlencoded data with the following values
- first_name
- last_name
- nickname
- password
- email
- country

## Tests
The client has a test suite. To run it go into the client folder and run ```go test```

## Design choices

I structed the code like this:
- Load in a configuration file. This configuration file allows you to change the host address of the service as well as what database type the service will use to hold data. This also included a field for a AMQP brooker address for sending events to other services. I did not implement this feature fully because it was out of scope and I don't have that much free time.
- Create the database layer of the service. This is an interface that contains 5 functions. This is to enable quick changing of the database. For this example I mocked the database so the database layer is simply a slice of User structs. The 5 functions are:
	- AddUser
	- FindUserByID
	- DeleteUser
	- FindUserByCriteria
	- UpdateUser
- Pass in the database layer to create the client. The client is an interface as well that implents the 6 routes for the service. Again this is to allow much quicker implementation of different clients if you so wish. 

## Other liberties I took due to time
- I was going to make the search route user query parameters instead of path parameters, but I wasn't sure if you wanted to be able to filter by multiple items.
- I did not write tests for the MockDB because it is a mock database
- I did not validate the input for creating users
- No security on the client
- I started commenting my code, but stopped half way through. If you really want me to I will do that.
