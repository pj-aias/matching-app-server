# Users API

## GET `/user/:id`
Authorization required.
```json
{
    "id": 123,
    "username": "hoge",
    "avatar": "https://example.com/some-image.png",
    "bio": "I am hoge."
}
```

## POST `/user`
You need AIAS Signature to crate a user.

If the creation was successful, server returns the user's information and a token. (For more information about authorization, see [authorizations documentation](./authorizations.md).)

request body:

```json
{
    "username": "fuga",
    "password": "hogehoge-password",
    "signature": "<AIAS SIGNATURE>"
}
```

response:

```json
{
    "user": {
        "id": 124,
        "username": "fuga",
        "avatar": "",
        "bio": "",
    },
    "token": "TOKEN"
}
```

## PATCH `/user`
Authorization requeired.

`avatar` and `bio` can be updated.

In the request body, fields can be omitted if it's not intended to update.

request body:

```json
{
    "avatar": "https://examlpe.com/another-image.png",
    "bio": "I am fuga.",
}
```

response

```json
{
    "id": 124,
    "username": "fuga",
    "avatar": "https://examlpe.com/another-image.png",
    "bio": "I am fuga."
}
```