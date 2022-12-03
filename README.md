# Manager service
## Endpoints

```
GET  localhost:8080/v1/users/:id
GET  localhost:8080/v1/users?offset=0&limit=10
POST localhost:8080/v1/users
```

User JSON:
```
{
  "name": "User01",
  "email": "admin@test.com",
  "password": "test",
  "activated": true
}
```