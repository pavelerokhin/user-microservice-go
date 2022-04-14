# user-microservice-go
Test commitment for F.I. 

Test SQLite example database `user.db` is available in the repo and populated.

## Data model
User type:
- `id`: type`uint` (autoincremental, provided by `GORM` library)
- `first_name`: type`string`, required
- `last_name`: type`string`, required
- `nickname`: type`string`, required
- `password`: type`string`, required  
- `email`: type`string`, required
- `country`: type`string`, required
- `created_at`: type`time.Time` (provided by `GORM` library)
- `updated_at`: type`time.Time` (provided by `GORM` library)

## Start the server
```
go run main.go [-port PORT]
```
The server will run and listen localhost on the port, by default it is `8080`.

## Run tests
```
go test ./...
```

## Run with Docker
First you need to build the Docker image:
```
docker build -t golang-user-api
```
Check if the image is present in your local repo:
```
docker images
```
Run the container:
```
docker run -p 8080:8080 golang-user-api
```
NB the port 8080 is the default port which microservice exposes, you can cutomize it by changing the last directive of the `Dockerfile`, e.g.:
```
CMD ["./user-microservice-go","-port", "9000"]
```
Then you need to rebuild the image and run the container once more exposing the correct port:
```
docker run -p 8080:9000 golang-user-api
```

## Microservice APIs

### Adding a new User
You can add a new user by sending `POST` request with user data in the request body. All user data are
required and microservice return `InternalServerError` if any of the required fields is empty 
The `id` field can be set explicitly, otherwise microservice assigns it automatically.
`created_at` and `updated_at` fields must be empty in the request.

Example of creating a new user:
```
curl --location --request POST 'http://localhost:8080/user' \
--header 'Content-Type: application/json' \
--data-raw '{
    "first_name": "name",
    "last_name": "surname",
    "nickname": "nick",
    "password": "12345",
    "email": "mail@mail.com",
    "country": "Israel"
}'
```

### Modifying an existent User
User can be modified via sending a `POST` request with URI `/user/<user_id>`. The request body may contain 
the data to be modified.
Modifying user with id 1:

```
curl --location --request POST 'http://localhost:8080/user/1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "first_name": "new name",
    "last_name": "new surname",
    "nickname": "new nick",
    "email": "mail@supermail.com",
}'
```

### Remove a User
You can delete a user by its `id`, sending a `DELETE` request to the URI `/user/<id>`.

Deleting user with id 1:
```
curl --location --request DELETE 'http://localhost:8080/user/1' 
```

### Return all Users
`GET` request to 
the URI `/users` returns the list of users available in the database.

Get all Users:
```
curl --location --request GET 'http://localhost:8080/users'
```

### Return paginated list of Users:
Paginate all list of Users, 3 items on page, get page #2:

```
curl --location --request GET 'http://localhost:8080/users/3/2'
```
You can combine pagination with filtering (API below),

### Return filtered list of Users:
You can get filtered l;list of users by sending a `GET` request to `/users`. The request body must 
contain a json with all filtering request.
Get all users from Israel:
```
curl --location --request GET 'http://localhost:9000/users' \
--header 'Content-Type: application/json' \
--data-raw '    {
        "country": "Israel"
    }'
```
You can combine pagination with filtering (API above),

### Return User by `id`
It wasn't requested by the test commitment, but it is handy to have this API available. 
You can get a single user by its `id` sending a `GET` request to `/user/<id>`.

E.g.:
```
curl --location --request GET 'http://localhost:9000/user/1' \
--header 'Content-Type: application/json' \
--data-raw '    {
        "country": "Israel"
    }'
```
