# Authorizations
To access "Authorization required" endpoints, you need to send a token that was issued by the server.

You can receive the token when creating the user ([`POST /user`](./users.md)), and when you access `POST /login`.

After getting a token, it must be set to the `Authorization` header with the following format.

```
Authorization: Bearer ${TOKEN}
```

For example:

```
Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjUsImV4cCI6MTYyOTE3NjA3MywiaXNzIjoibWF0Y2hpbmctYXBwIn0.0xg6jD6Qbb7MU3eCZ0s5ZBhc31WEuMS9SgwEFHs_b1d8zgS0ytYeFP99P_bbAWCQno_EYw83sQg0mqpL4DfF6nPL72ophcXdewcmfRuihB4RDZq2A7Z4yIXGu1F38IHcrjdvfpgOOGRza2RwFpfw_u75mjPETSmDPDGEj-LLTyt5tpIcg4FMV-oLPRC8UU3EWYmzR38DoYN846QOtHeLPexNTogPbU271fQW1-bu4SJcF7MEHO3-b6wZ8Ix42GywVssZixGQDgAJzfnGSYoTonJ1xg0YKhqhC15ke3K3T-VlCv_dVvSWPxFF7X8Oft0iWkId6IN69JcUHJSbIOl84Q
```

## `POST /login`
request body:
```
{
    "username": "hoge",
    "password": "hogePassword"
}
```

response:

```
{
    "token": "TOKEN COMES HERE"
}
```