# user-microservice-go
test commitment for F.I. 

Test SQLite database `users.db` is available and populated.

## data model
User type:
- `id`
- `first_name`
- `last_name`
- `nickname`
- `password`
- `email`
- `country`
- `created_at`
- `updated_at`

## Starting the server
```
go run main.go
```
The server will run and listen localhost on the port `8080`.


## API for _adding_ a new User
Creating a new user:
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

## API for _modifying_ User
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


## API for _removing_ User
Deleting user with id 1:
```
curl --location --request POST 'http://localhost:8080/user/1' 
```

## API return _all_ Users
Get all Users:
```
curl --location --request GET 'http://localhost:8080/users'
```

## API, return _paginated list_ of Users,:
Paginate all list of Users, 3 items on page, get page #2:

```
curl --location --request GET 'http://localhost:8080/users/3/2'
```

## API _filter_ Users:

```
curl --location --request GET 'http://localhost:8080/users'
```
