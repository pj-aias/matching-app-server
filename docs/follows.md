# Follows API
An user can follow another user.

## `POST /follow/:userid`
Authorization required.

Follow a user (`target`) specified by the `userid` newly from the authenticated user.

response:

```json
{
    "target": 123,
    "following": true,
    "followed": true
}
```

## `GET /follow/:userid`
Authorization required.

Show whether the authenticated user follows a user specified by the `userid` or not and the reverse.

`"following"`: `true` if the authenticated user is following `"target"`, or else `false`.

`"followed"`: `true` if the authenticated user is followed by `"target"`, or else `false`.

```json
{
    "target": 123,
    "following": true,
    "followed": true
}
```

## `DELETE /follow/:userid`
Authorization required.

Un-follow a user specified by the `userid`.

response:

```json
{
    "target": 123,
    "following": false,
    "followed": true
}
```

## `GET /followers`
Authorization required.

Get followers of the authenticated user.

```json
{
    "followers": [
        {
            "id": "123",
            "username": "hoge",
            "bio": "hoge bio",
            "avatar": "https://example.com/hoge.png",
        },
        {
            "id": 124,
            "username": "fuga",
            "bio": "fuga bio",
            "avatar": "https://example.com/fuga.png",
        },
    ]
}
```

## `GET /followees`
Authorization required.

Get followees of the authenticated user.

```json
{
    "followees": [
        {
            "id": 124,
            "username": "fuga",
            "bio": "fuga bio",
            "avatar": "https://example.com/fuga.png",
        },
    ]
}
```

