# Matching API Specifications
## `POST /matching`
Authorization required.

Pick a user randomly from all users in the service, and create a chatroom with the matched user.

response:

```
{
    "matched_user": {
        "id": 124
        "username": "fuga",
        "biio": "fugafuga bio",
        "avatar": "https://example.com/fuga.png"
    },
    "chatroom": {
        "id": 4,
        "users": [123, 124]
    }
}
```