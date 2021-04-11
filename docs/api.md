# Fooder Server API Docs

## Endpoints

### GET /ping

Returns ping status

Response:

```
{
"status": "ok"
}
```

### POST /auth

Returns `Bearer` authorization token and sets user uniqe cookie

```
{
    "token": "fdsfdsfdsfdsfas"
}
```

### GET /history

Return whole user chat history with all events

```
[
    {
        "type": "message",
        "body": {
            "test": "Hi!"
        }
    },
    ...
]
```

### POST /send_event

Method used to set events, returns just confirmation if everything went smooth.

Request:
```
{
    "type": "event_type",
    "body": {
        // Event body
    }
}
```

Response:
```
{
    "status": "OK"
}
```

## WebSocket

WebSocket connection is used to update client with events send in chat.
Connection allows for following methods
### Login

Used to log user in

```
{
    "action": "login",
    "payload": {
        "token": "fdsdfsd"
    }
}
```

### Events

Events are send back as incoming message from bot, events are described in docs.

## Events

Events are form of thing appeared in chat, their general format is:

```
{
    "type": "event_type", // Could be on of described below
    "author_id": "bot" // could be user uuidv4 id or `bot`
    "body": {
        ...
    }
}
```

Events bodies are described below

### Message 

Type: `message`

```
{
    "message": "Hi, my name is bot"
}
```

### Card

Type: `card`

```
{
    "dish_id": "XXXX"
    "title": "Scrambled eggs",
    "description": "Well known dish",
    "img": "http://..../phoyo.png",
    "link": "http://.../dish"
}
```

### Rating Requested

Type: `rating_requested`

```
{
    "dish_id": "XXXX"
}
```

### Rating Set

Type: `rating_set`

```
{
    "dish_id": "XXXX"
    "rating": 4 // from 1-5
}
```

## Errors

Errors are send with following format:

```
{
    "error": "error description"
}
```
