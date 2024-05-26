# Moneybots API

Documentation for the Moneybots API

## API's

| API         | Route Prefix   | Description                       | Authentication                           |
| ----------- | -------------- | --------------------------------- | ---------------------------------------- |
| User        | `/user`        | For managing Kite API users       | Request Form Body: `user_id`, `password` |
| Instruments | `/instruments` | For managing Kite API instruments | Authorization Header: `Bearer token`     |
| Ticker      | `/ticker`      | For managing Kite Ticker API data | Authorization Header: `Bearer token`     |

### User API

| Method  | Route           | Description             | Request Form Payload                  | Response Data                                        |
| ------- | --------------- | ----------------------- | ------------------------------------- | ---------------------------------------------------- |
| `POST`  | `/user/create`  | Create a new user       | `user_id`, `password`                 | `user_id`, `created_at`                              |
| `PATCH` | `/user/update`  | Update password of user | `user_id`, `password`, `new_password` | `user_id`, `updated_at`                              |
| `POST`  | `/user/delete`  | Delete a user           | `user_id`, `password`                 | `user_id`, `deleted_at`                              |
| `POST`  | `/user/login`   | Login a user            | `user_id`, `password`, `totp_secret`  | `user_id`, `enctoken`, `token`, `login_time`         |
| `POST`  | `/user/logout`  | Logout a user           | `user_id`, `password`                 | `user_id`, `logout_time`                             |
| `POST`  | `/user/profile` | Profile of a user       | `user_id`, `password`                 | `user_id`, `username`, `user_shortname`, `email` ... |

### Instruments API

| Method | Route                                       | Description                                       | Response Data                                       |
| ------ | ------------------------------------------- | ------------------------------------------------- | --------------------------------------------------- |
| `GET`  | `/instruments/update`                       | Update instruments in database from Kite API      | `exchanges`, `instruments`, `updated_at`            |
| `GET`  | `/instruments/details?i=NSE:INFY&i=NSE:IOC` | Get instrument detail from exchange:tradingsymbol | `exchange`, `trading_symbol`, `expiry`, `strike`... |
| `GET`  | `/instruments/details?t=26565&t=264969`     | Get instrument detail from instrument_token       | `exchange`, `trading_symbol`, `expiry`, `strike`... |

### Ticks API

| Method | Route          | Description     | Request Body  | Response Data                                    |
| ------ | -------------- | --------------- | ------------- | ------------------------------------------------ |
| `GET`  | `/ticks/start` | Start the ticks | `instruments` | `user_id`, `instruments`,`channel`, `started_at` |

---

## Authentication

```curl
curl https://api.moneybots.app/instruments/update \
    -H "Authorization: Bearer jwt" \
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
