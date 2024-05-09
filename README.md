# Moneybots API

Documentation for the Moneybots API

## Routes

| Main Routes | Route Prefix | Description                      |
| ----------- | ------------ | -------------------------------- |
| User        | `/user/`     | Register, Login and Logout Users |

## Responses

### Success

```sh
HTTP/1.1 200 Ok
{
  "status": "ok",
  "data": {Any}
}
```

### Error

```sh
HTTP/1.1 400 Bad Request
Content-Type: application/json
{
  "status": "error",
  "error_type": "string",
  "message": "string",
}
```

---

## User API

| Method   | Route           | Description             |
| -------- | --------------- | ----------------------- |
| `POST`   | `/user/login`   | Create new user session |
| `DELETE` | `/user/logout`  | Delete the user session |
| `GET`    | `/user/profile` | Return the user profile |

### `POST /user/login`

Create new user session.

```sh
curl -X POST https://api.moneybots.app/user/login \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "user_id=`user_id`" \
        -d "password=`password`" \
        -d "totp_secret=`totp_secret`"
```

#### Request Parameters

| Form Body     | Type   | Description      |
| ------------- | ------ | ---------------- |
| `user_id`     | string | User Id          |
| `password`    | string | User Password    |
| `totp_secret` | string | User TOTP secret |

#### Response Attributes

| JSON Body    | Type   | Description        |
| ------------ | ------ | ------------------ |
| `user_id`    | string | User Id            |
| `enctoken`   | string | Session token      |
| `login_time` | string | Session Login Time |

#### Error Responses

| Status Code | Description                                         |
| ----------- | --------------------------------------------------- |
| `400`       | Bad Request: Input(s) missing                       |
| `401`       | Unauthorized Access: Unregistered user              |
| `403`       | Forbidden: User not allowed to access this resource |
| `422`       | Unprocessable Entity: Input(s) malformed            |

---

### `DELETE /user/logout`

Delete the user session.

```sh
curl -X DELETE https://api.moneybots.app/user/logout \
        -H "`Authorization`: `user_id`:`enctoken`"
```

#### Error Responses

| Status Code | Description                            |
| ----------- | -------------------------------------- |
| `401`       | Unauthorized Access: Unauthorized user |

---

### `GET /user/profile`

Return the user profile.

```sh
curl -X GET https://api.moneybots.app/user/profile \
        -H "`Authorization`: `user_id`:`enctoken`"
```

#### Response Attributes

| Body             | Type     | Description                          |
| ---------------- | -------- | ------------------------------------ |
| `user_id`        | string   | Unique user id of the user           |
| `user_name`      | string   | Real name of the user                |
| `user_shortname` | string   | Short version of the real name       |
| `email`          | string   | Email of the user                    |
| `broker`         | string   | Broker used by the user              |
| `exchanges`      | string[] | Exchanges enabled for the user       |
| `products`       | string[] | Margin products enabled for the user |
| `order_types`    | string[] | Order types enabled for the user     |

#### Error Responses

| Status Code | Description                            |
| ----------- | -------------------------------------- |
| `401`       | Unauthorized Access: Unauthorized user |

---
