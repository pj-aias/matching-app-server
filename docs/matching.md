# Matching API Specifications
## `POST /matching`
Authorization required.

Pick a user randomly from all users in the service.

After this request, the client can optionally create a chatroom with the returned user.

response:

```
{
    "user": {
        "id": 124
        "username": "fuga",
        "biio": "fugafuga bio",
        "avatar": "https://example.com/fuga.png"
    }
}
```