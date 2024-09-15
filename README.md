<div align=center><img src="https://i.postimg.cc/FR1MRJXv/box-icon128.png" width="128" alt="box-icon"/><h1>-ChatBox-</h1></div>

This project is a fullstack web application designed to achieve communication between users in real-time through websockets.

## Tech Stack

**Client:** Vue, Vuetify

**Server:** Go, Chi, MongoDB

**Infra:** Docker, Render

## API Reference

### Base URL
>     https://chatbox-back.onrender.com
#### Websocket
>     wss://chatbox-back.onrender.com

### Routes
##### Public:
- [SignUp](#user-signup)
- [SignIn](#user-signin)
- [Logout](#user-logout)
- [Renew Token](#renew-token)
- [Revoke Token](#revoke-token)
##### Authorization required:
- [User Data](#user-data)
- [Save Personal Data](#save-personal-data)
- [Add Chat](#add-chat)
- [Load Chat](#load-chat)
- [Load Messages](#load-messages)
- Websockets ↓
	- [Send Message](#send-message)
	-  [Edit Message](#edit-message)
	- [Delete Message](#delete-messages)


I opted to store values such as access tokens, user id's and session id's in `localStorage` for the sake of simplicity, easy debugging and straightforward management of session data during development of frontend. I am fully aware of the potential security risks associated with this type of data management.
<hr>

### User SignUp
```http
POST /signup
```
#### Request body
```javascript
{
  "name": "John Doe"
  "password": "1234"
}
```
| Headers	       | Value              |
| :------------- | :----------------- |
| `Content-Type` | `application/json` |
#### Response
```javascript
{
  "access_exp": "2024-09-14T02:53:10Z",
  "access_token": "eyJhbGciOiJ...",
  "message": "User added successfully",
  "refresh_exp": "2024-09-15T02:38:10Z",
  "session_id": "ce8306f9-d614-4439-8337-8f7efdeff1b1",
  "user": {
    "id": "66e4f71243d2b68b40ca71ac",
    "name": "John Doe"
  }
}
```
<hr>

### User SignIn
```http
POST /signin
```
| Headers	       | Value              |
| :------------- | :----------------- |
| `Content-Type` | `application/json` |
#### Request body
```javascript
{
  "name": "John Doe"
  "password": "1234"
}
```

#### Response
```javascript
{
  "access_exp": "2024-09-14T02:56:28Z",
  "access_token": "eyJhbGciOiJ...",
  "message": "User logged successfully",
  "refresh_exp": "2024-09-15T02:41:28Z",
  "session_id": "5aeb2038-2e16-4ba3-85c7-672ebd6a6a5b",
  "user": {
    "id": "66e4f71243d2b68b40ca71ac",
    "name": "John Doe"
  }
}
```
<hr>

### User Logout
```http
DELETE /logout/{session-id}
```
#### Response
```javascript
{
  "message": "User logged out"
}
```
<hr>

### Renew Token
```http
POST /token/renew/{session-id}
```
#### Response
```javascript
{
  "access_exp":  "2024-09-14T03:07:52Z",
  "access_token":  "eyJhbGciOiJ...",
  "message":  "Access token renewed"
}
```
<hr>

### Revoke Token
```http
POST /token/revoke/{session-id}
```
#### Response
```javascript
204 No Content (Revoked successfully)
```
<hr>

### User Data
```http
GET /user/{id}
```
| Headers	       | Value              |
| :------------- | :----------------- |
| `Authorization`| `access-token`     |
#### Response
```javascript
{
  "claims": {
    "expires_at": "Sat Sep 14 04:19:44 2024",
    "subject": "312a36d6-178e-4f14-b671-bd6d0e656c4f"
  },
  "message": "User data retrieved",
    "user_data": {
      "id": "66e4f71243d2b68b40ca71ac",
      "name": "John Doe",
      "personal": {
      "Email": "john@gmail.com",
      "Country": "United States",
      "Age": 32
    }
  }
}
```
<hr>

### Save Personal Data
```http
PUT /user/save-personal/{id}
```
| Headers	       | Value              |
| :------------- | :----------------- |
| `Content-Type` | `application/json` |
| `Authorization`| `access-token`     |
#### Request body
```javascript
{
  "email": "johndoe@gmail.com",
  "age": 32,
  "country": "United States"
}
```
#### Response
```javascript
{
  "message":  "Personal data saved successfully"
}
```
<hr>

### Add Chat
```http
POST /chat
```
| Headers	       | Value              |
| :------------- | :----------------- |
| `Content-Type` | `application/json` |
| `Authorization`| `access-token`     |
#### Request body
```javascript
{
  "username": "New Chat",
  "petitioner_id": "66e4f71243d2b68b40ca71ac",
  "petitioner": "John Doe"
}
```
#### Response
```javascript
{
  "chat": {
    "id": "66e5ab4cdff129912cd0e80b",
    "participants": [
      {
        "id": "66e5ab40dff129912cd0e80a",
        "name": "Jane Doe"
      },
      {
        "id": "66e4f71243d2b68b40ca71ac",
        "name": "John Doe"
      }
    ]
  },
  "contact": {
    "username": "Jane Doe",
    "petitioner_id": "66e4f71243d2b68b40ca71ac",
    "petitioner": "John Doe"
  },
  "message": "Chat added successfully"
}
```
<hr>

### Load Chats
```http
GET /chat/{user-id}
```
| Headers	       | Value              |
| :------------- | :----------------- |
| `Content-Type` | `application/json` |
| `Authorization`| `access-token`     |
#### Response
```javascript
{
  "chats": [
    {
      "id": "66e5ab4cdff129912cd0e80b",
      "participants":  [
        {
          "id": "66e4f71243d2b68b40ca71ac",
          "name": "John Doe"
        },
        {
          "id": "66e35284b75d98c6ce1a58c7",
          "name": "Jane Doe"
        }
      ]
    },
    [...]
  ],
  "message": "User chats retrieved"
}
```
<hr>

### Delete Chat
```http
DELETE /chat/{chat-id}
```
| Headers	     | Value              |
| :------------- | :----------------- |
| `Content-Type` | `application/json` |
| `Authorization`| `access-token`     |
#### Response
```javascript
{
  "message":  "Chat deleted successfully"
}
```
<hr>

### Load Messages
```http
GET /chat/{chat-id}/messages
```
| Headers	       | Value              |
| :------------- | :----------------- |
| `Content-Type` | `application/json` |
| `Authorization`| `access-token`     |
#### Response
```javascript
{
  "message": "Messages loaded",
  "messages": [
    {
      "id": "66e65747c19386daf16f545e",
      "chat_id": "66e5ab4cdff129912cd0e80b",
      "sender_id": "66e4f71243d2b68b40ca71ac",
      "content": "Hello",
      "sent_at": "2024-09-15T03:40:55.063Z"
    },
    [...]
  ]
}
```
<hr>

### Send Message
```plaintext
WEBSOCKET /ws/send-msg
```
| Parameters  | Value          |
| :---------- | :------------- |
| `wsauth`    | `access-token` |
#### Request body
```javascript
{
  "chat_id": "66e5ab4cdff129912cd0e80b",
  "sender_id": "66e4f71243d2b68b40ca71ac",
  "content": "Hello"
}
```
#### Response
```javascript
{
  "id": "66e65747c19386daf16f545e",
  "chat_id": "66e5ab4cdff129912cd0e80b",
  "sender_id": "66e4f71243d2b68b40ca71ac",
  "content": "Hello",
  "sent_at": "2024-09-15T03:40:55.063Z"
}
```
<hr>

### Edit Message
```plaintext
WEBSOCKET /ws/edit-msg
```
| Parameters  | Value          |
| :---------- | :------------- |
| `wsauth`    | `access-token` |
#### Request body
```javascript
{
  "message_id": "66e65747c19386daf16f545e",
  "new_message": "Bye"
}
```
#### Response
```javascript
{
  "id": "66e65747c19386daf16f545e",
  "chat_id": "66e5ab4cdff129912cd0e80b"
  "sender_id": "66e4f71243d2b68b40ca71ac",
  "content": "Bye",
  "sent_at": "2024-09-15T03:40:55.063Z"
}
```
<hr>

### Delete Message
```plaintext
WEBSOCKET /ws/del-msg
```
| Parameters  | Value          |
| :---------- | :------------- |
| `wsauth`    | `access-token` |
#### Request body
Message ID in Text format
```plaintext
66e65747c19386daf16f545e
```
#### Response
```javascript
{
  "message": "Message deleted successfully"
  "message_id": "66e65747c19386daf16f545e"
}
```
