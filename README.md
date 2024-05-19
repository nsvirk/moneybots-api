# Moneybots API

Documentation for the Moneybots API

## Routes

| Main Routes | Route Prefix | Description                                      |
| ----------- | ------------ | ------------------------------------------------ |
| Users       | `/user`      | Register, Update, Delete, Login and Logout Users |

### User API - CRUD

| Method  | Route           | Description             | Request                               | Response                            |
| ------- | --------------- | ----------------------- | ------------------------------------- | ----------------------------------- |
| `POST`  | `/user/create`  | Create a new user       | `user_id`, `password`                 | `user_id`, `created_at`             |
| `PATCH` | `/user/update`  | Update password of user | `user_id`, `password`, `new_password` | `user_id`, `updated_at`             |
| `POST`  | `/user/delete`  | Delete a user           | `user_id`, `password`                 | `user_id`, `deleted_at`             |
| `POST`  | `/user/login`   | Login a user            | `user_id`, `password`, `totp_secret`  | `user_id`, `enctoken`, `login_time` |
| `POST`  | `/user/logout`  | Logout a user           | `user_id`, `password`                 | `user_id`, `logout_time`            |
| `POST`  | `/user/profile` | Profile of a user       | `user_id`, `password`                 | `user_id`, `username`, ...          |

---

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

### Status Codes

| Status Code | Description                                         |
| ----------- | --------------------------------------------------- |
| `200`       | Ok                                                  |
| `400`       | Bad Request: Input(s) missing                       |
| `401`       | Unauthorized Access: Unregistered user              |
| `403`       | Forbidden: User not allowed to access this resource |
| `422`       | Unprocessable Entity: Input(s) malformed            |
| `500`       | Internal Server Error: Server side error            |

---
