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
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiIzNzJlMDM2ZS00OWEzLTQyMDctYTYyOS0wN2Y5NzlmM2ZhNTgifQ.bpx8Ry1QjhDB_q4Am6n0hYGJ4QwXV9HIp_hKkAp-_OU",
    "user_id": "372e036e-49a3-4207-a629-07f979f3fa58"
}
```

### GET /history

Return whole user chat history with all events

```
[
    {
        "type": "message",
        "author_id": "372e036e-49a3-4207-a629-07f979f3fa58",
        "body": {
            "message": "Helo elo 2"
        },
        "creation_timestamp": 656265007
    },
    {
        "type": "card",
        "author_id": "bot",
        "body": {
            "img": "https://images.immediate.co.uk/production/volatile/sites/30/2020/08/recipe-image-legacy-id-1201452_12-7f7a0fa.jpg?quality=90&webp=true&resize=440,400",
            "link": "https://www.bbcgoodfood.com/recipes/perfect-scrambled-eggs-recipe",
            "dish_id": "tmp",
            "title_id": "Scrambled Eggs",
            "description": "Simple but nutritious dish"
        },
        "creation_timestamp": 665149383
    },
    {
        "type": "message",
        "author_id": "372e036e-49a3-4207-a629-07f979f3fa58",
        "body": {
            "message": "Helo elo 2"
        },
        "creation_timestamp": 317486182
    },
    {
        "type": "message",
        "author_id": "bot",
        "body": {
            "message": "Thanks for you opinion, remember I'm always here for you to help"
        },
        "creation_timestamp": 354581260
    }
]
```

### POST /send_event

Method used to set events, returns bot response in form of array of events event.

Request:
```
{
    "type": "message",
    "body": {
        "message": "Helo elo 2"
    }
}
```

Response:
```
[
    {
        "type": "message",
        "author_id": "bot",
        "body": {
            "message": "Thanks for you opinion, remember I'm always here for you to help"
        },
        "creation_timestamp": 354581260
    }
]
```

### Events

Events are send back as incoming message from bot, events are described in docs.

## Events

Events are form of thing appeared in chat, their general format is:

```
{
    "type": "event_type", // Could be on of described below
    "author_id": "bot" // could be user uuidv4 id or `bot`
    "event_id": "..." // uuid v4 id
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

### Select

Type: `select`

```
{
    "options": [
        {
            "option_id": 1,
            "option_text": "Dish Name",
            "dish_id": 123
        },
        ...
    ],
    "selected_option_id": 1 // Should be sent only when dish is selected
}
```

If user clicks selected type then only "selected_option_id" should be sent



### Rating Requested [not implemented]

Type: `rating_requested`

```
{
    "dish_id": "XXXX"
}
```

### Rating Set [not implemented]

Type: `rating_set`

```
{
    "dish_id": "XXXX"
    "rating": 4 // from 1-5
}
```

### System Status [not implemented]

Type: `status_changed`

```
{
    "information": "Chat went inactive"
}
```

## Errors

Errors are send with following format:

```
{
    "error": "error description"
}
```
