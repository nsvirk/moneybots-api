# Moneybots API

Documentation for the Moneybots API

## API's

| Sr  | API         | Route Prefix   | Description                       | Authentication                           | Status      |
| --- | ----------- | -------------- | --------------------------------- | ---------------------------------------- | ----------- |
| 1.  | User        | `/user`        | For managing Kite API users       | Request Form Body: `user_id`, `password` | Implemented |
| 2.  | Instruments | `/instruments` | For managing Kite API instruments | Authorization Header: `Bearer token`     | Implemented |
| 3.  | Ticker      | `/ticker`      | For managing Kite Ticker API data | Authorization Header: `Bearer token`     | Implemented |
| 4.  | Orders      | `/orders`      | For managing Kite Orders          | Authorization Header: `Bearer token`     | ToDo        |
| 5.  | Positions   | `/positions`   | For managing Kite Positions       | Authorization Header: `Bearer token`     | ToDo        |

### 1. User API

| Method  | Route           | Description             | Request Form Payload                  | Response Data                                        |
| ------- | --------------- | ----------------------- | ------------------------------------- | ---------------------------------------------------- |
| `POST`  | `/user/create`  | Create a new user       | `user_id`, `password`                 | `user_id`, `created_at`                              |
| `PATCH` | `/user/update`  | Update password of user | `user_id`, `password`, `new_password` | `user_id`, `updated_at`                              |
| `POST`  | `/user/delete`  | Delete a user           | `user_id`, `password`                 | `user_id`, `deleted_at`                              |
| `POST`  | `/user/login`   | Login a user            | `user_id`, `password`, `totp_secret`  | `user_id`, `enctoken`, `token`, `login_time`         |
| `POST`  | `/user/logout`  | Logout a user           | `user_id`, `password`                 | `user_id`, `logout_time`                             |
| `POST`  | `/user/profile` | Profile of a user       | `user_id`, `password`                 | `user_id`, `username`, `user_shortname`, `email` ... |

### 2. Instruments API

| Method | Route                  | Request Querystring     | Description                                       | Response Data                            |
| ------ | ---------------------- | ----------------------- | ------------------------------------------------- | ---------------------------------------- |
| `GET`  | `/instruments/update`  |                         | Update instruments in database from Kite API      | `exchanges`, `instruments`, `updated_at` |
| `GET`  | `/instruments/details` | `?i=NSE:INFY&i=NSE:IOC` | Get instrument detail from exchange:tradingsymbol | `instrument_map`                         |
| `GET`  | `/instruments/details` | `?t=26565&t=264969`     | Get instrument detail from instrument_token       | `instrument_token_map`                   |

### 3. Ticks API

| Method | Route          | Description     | Request Body  | Response Data                                       |
| ------ | -------------- | --------------- | ------------- | --------------------------------------------------- |
| `POST` | `/ticks/start` | Start the ticks | `instruments` | `user_id`, `channel`, `instruments_ct`,`started_at` |
| `GET`  | `/ticks/stop`  | Stop the ticks  |               | `user_id`, `channel`, `published_ct`, `stopped_at`  |

---

## Authentication

### Uses token returned on login

```curl
curl https://www.moneybots.app/instruments/update \
    -H "Authorization: Bearer token" \
```

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
