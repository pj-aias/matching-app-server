# API Specifications (Draft)
## `POST /message`
Authorization required.

Create a new chatroom with a specified `target` user.

request body:

```
{
    "target": 123
}
```

response:

```
{
    "chatroom": {
        "id": 1,
        "users": [
            {
                "id": 123,
                "username": "hoge",
                "bio": "hogehoge bio",
                "avatar": "https://examlpe.com/hoge.png"
            },
            {
                "id": 124,
                "username": "fuga",
                "bio": "fugafuga bio",
                "avatar": "https://examlpe.com/fuga.png"
            }
        ]
    }
}
```

## `POST /message/:roomid`
Authorization required.

Post a new message to a chatroom.

request body:

```
{
    "content": "hoge message"
}
```

response:

```json
{
    "message": {
        "id": 321,
        "chatroom_id": 432,
        "user": {
            "id": 123,
            "username": "hoge",
            "bio": "hogehoge bio",
            "avatar": "https://example.com/hoge.png"
        },
        "content": "hogehoge message"
    },
    "chatroom": {
        "id": 432,
        "users": [
            {
                "id": 123,
                "username": "hoge",
                "bio": "hogehoge bio",
                "avatar": "https://examlpe.com/hoge.png"
            },
            {
                "id": 124,
                "username": "fuga",
                "bio": "fugafuga bio",
                "avatar": "https://examlpe.com/fuga.png"
            }
        ]
    },
    "messages": [
        {
             "id": 321,
             "chatroom_id": 432,
             "user": {
                 "id": 123,
                 "username": "hoge",
                 "bio": "hogehoge bio",
                 "avatar": "https://example.com/hoge.png"
             },
            "content": "hogehoge message"
        },
        {
             "id": 322,
             "chatroom_id": 432,
             "user": {
                 "id": 124,
                 "username": "fuga",
                 "bio": "fugafuga bio",
                 "avatar": "https://example.com/fuga.png"
             },
            "content": "fugafuga message"
        }
    },
}
```

## `GET /message/:roomid`
Authorization requied.

Get messages in the chatroom.

response:

```
{
    "chatroom": {
        "id": 432,
        "users": [
            {
                "id": 123,
                "username": "hoge",
                "bio": "hogehoge bio",
                "avatar": "https://examlpe.com/hoge.png"
            },
            {
                "id": 124,
                "username": "fuga",
                "bio": "fugafuga bio",
                "avatar": "https://examlpe.com/fuga.png"
            }
        ]
    },
    "messages": [
        {
             "id": 321,
             "chatroom_id": 432,
             "user": {
                 "id": 123,
                 "username": "hoge",
                 "bio": "hogehoge bio",
                 "avatar": "https://example.com/hoge.png"
             },
            "content": "hogehoge message"
        },
        {
             "id": 322,
             "chatroom_id": 432,
             "user": {
                 "id": 124,
                 "username": "fuga",
                 "bio": "fugafuga bio",
                 "avatar": "https://example.com/fuga.png"
             },
            "content": "fugafuga message"
        }
    },
}
```

## `GET /message/rooms`
Authorization required.

Get all chatrooms which the user is in now.

response: 

```
{
    "chatrooms": [
        {
            "id": 432,
            "users": [
                {
                    "id": 123,
                    "username": "hoge",
                    "bio": "hogehoge bio",
                    "avatar": "https://examlpe.com/hoge.png"
                },
                {
                    "id": 124,
                    "username": "fuga",
                    "bio": "fugafuga bio",
                    "avatar": "https://examlpe.com/fuga.png"
                }
            ],
        },
        {
            "id": 433,
            "users": [
                {
                    "id": 123,
                    "username": "hoge",
                    "bio": "hogehoge bio",
                    "avatar": "https://examlpe.com/hoge.png"
                },
                {
                    "id": 124,
                    "username": "fuga",
                    "bio": "fugafuga bio",
                    "avatar": "https://examlpe.com/fuga.png"
                }
            ]
        }
    ]
}
```