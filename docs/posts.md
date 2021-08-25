# Posts API Documentation (Deprecated)
Deprecated.

Posts API were altered to Messages, and no longer used.

## `POST /post`
Authorization required.

Create a new post.

request body:

```
{
    "content": "An amazing content."
}
```

response: 
```
{
    "post": {
        "id": 1235",
        "user": {
            "id": 123,
            "username": "hoge",
            "avatar": "https://example.com/hoge.png",
            "bio": "hogehoge bio"
        },
        "content": "An amazing content."
    }
}
```

## `GET /post/:postid`
Authorization required.

Get a post with ID.

response:

```
{
    "post": {
        "id": 1234",
        "user": {
            "id": 123,
            "username": "hoge",
            "avatar": "https://example.com/hoge.png",
            "bio": "hogehoge bio"
        },
        "content": "post's content"
    }
}
```

## `PATCH /post/:postid`
Authorization required.

Update the content of the given post. If the post is not the user's, it fails.

request body:

```
{
    "content": "Another amazing content."
}
```

response: 
```
{
    "post": {
        "id": 1235",
        "user": {
            "id": 123,
            "username": "hoge",
            "avatar": "https://example.com/hoge.png",
            "bio": "hogehoge bio"
        },
        "content": "Another amazing content."
}
```

## `DELETE /post/:postid`
Authorization required.

Delete the given post.

response: 
```
{}
```

## `GET /post/recent`
Authorization required.

Get recent posts.

The number of posts to get can be specified with `"count"`.

request body:

```
{
    "count": 10
}
```

response:

```
{
    "posts": [
        {
            "id": 1234",
            "user": {
                "id": 123,
                "username": "hoge",
                "avatar": "https://example.com/hoge.png",
                "bio": "hogehoge bio"
            },
            "content": "post's content"
        },
        {
            "id": 1233",
            "user": {
                "id": 124,
                "username": "fuga",
                "avatar": "https://example.com/fuga.png",
                "bio": "fugafuag bio"
            },
            "content": "tihs is post's content"
        }
    ]
}
```
