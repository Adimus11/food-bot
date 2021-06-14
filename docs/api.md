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
    "event_id": "..." // uuid v4 id, required for user sending `select` or `rating_event`
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



### Rating Event

Type: `rating_event`

```
{
    "dish_id": "XXXX",
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

# Example chat scenario

## User sends: 

```
{
    "type": "message",
    "body": {
        "message": "Hi man"
    }
}
```

## Bot Responded:

```
[
    {
        "event_id": "148db137-a878-4996-81ae-1558a6c48ec0",
        "type": "message",
        "author_id": "bot",
        "body": {
            "message": "Hi! From which ingredients do you want to eat?"
        },
        "creation_timestamp": 891700234
    }
]
```

## User sends: 

```
{
    "type": "message",
    "body": {
        "message": "I don't know, how about some eggs?"
    }
}
```

## Bot Responded:

```
[
    {
        "event_id": "84359e5c-5dd8-4043-bfea-36315ab8bc37",
        "type": "message",
        "author_id": "bot",
        "body": {
            "message": "This is what I found for you!"
        },
        "creation_timestamp": 132278530
    },
    {
        "event_id": "0b5c589b-99fa-483f-b28d-90e10497a9bd",
        "type": "select",
        "author_id": "bot",
        "body": {
            "options": [
                {
                    "option_id": 0,
                    "option_text": "Scrambled eggs 1",
                    "dish_id": "6ee22dc6-9e6a-462a-bf83-8ca98a93b325"
                },
                {
                    "option_id": 1,
                    "option_text": "Scrambled eggs 2",
                    "dish_id": "2b41e888-63cc-4470-ba5f-c719a0b9609b"
                },
                {
                    "option_id": 2,
                    "option_text": "Scrambled eggs 3",
                    "dish_id": "3f3c2e4e-c575-444d-bd45-9f1c435185ab"
                }
            ]
        },
        "creation_timestamp": 148007661
    }
]
```

## User sends: 

```
{
    "event_id": "0b5c589b-99fa-483f-b28d-90e10497a9bd",
    "type": "select",
    "body": {
        "selected_option_id": 1
    }
}
```

## Bot Responded:

```
[
    {
        "event_id": "d9573629-1db2-40bb-a8c5-6f767d03a128",
        "type": "message",
        "author_id": "9f11638d-d221-4be3-b2de-1110f1fdf73f",
        "body": {
            "message": "Hi I would like `Scrambled eggs 2` for my dish"
        },
        "creation_timestamp": 945171778
    },
    {
        "event_id": "fd682819-a30a-4db7-b3c5-b3636a81cc0d",
        "type": "card",
        "author_id": "bot",
        "body": {
            "dish_id": "2b41e888-63cc-4470-ba5f-c719a0b9609b",
            "title_id": "Scrambled eggs 2",
            "description": "Well known dish",
            "img": "https://images.immediate.co.uk/production/volatile/sites/30/2020/08/recipe-image-legacy-id-1201452_12-7f7a0fa.jpg?quality=90&webp=true&resize=440,400",
            "link": "https://www.bbcgoodfood.com/recipes/perfect-scrambled-eggs-recipe"
        },
        "creation_timestamp": 947284004
    },
    {
        "event_id": "a1095210-4d82-4110-8771-f3fae1961716",
        "type": "message",
        "author_id": "bot",
        "body": {
            "message": "How you liked your meal?"
        },
        "creation_timestamp": 949595289
    },
    {
        "event_id": "1f06049d-fd50-4048-bd8e-4a316c4f5497",
        "type": "rating_event",
        "author_id": "bot",
        "body": {
            "dish_id": "2b41e888-63cc-4470-ba5f-c719a0b9609b"
        },
        "creation_timestamp": 951647711
    }
]
```

## User sends: 

```
{
    "event_id": "1f06049d-fd50-4048-bd8e-4a316c4f5497",
    "type": "rating_event",
    "body": {
        "rating": 4,
        "dish_id": "2b41e888-63cc-4470-ba5f-c719a0b9609b"
    }
}
```

## Bot Responded:

```
[
    {
        "event_id": "3d037bb2-dc3a-4a73-8123-8569d7f61250",
        "type": "message",
        "author_id": "9f11638d-d221-4be3-b2de-1110f1fdf73f",
        "body": {
            "message": "I think it's solid `4`"
        },
        "creation_timestamp": 716543618
    },
    {
        "event_id": "7c7fc45c-4f83-4b10-bbe1-e10c15960c73",
        "type": "message",
        "author_id": "bot",
        "body": {
            "message": "Thanks for you opinion, remember I'm always here for you to help!"
        },
        "creation_timestamp": 718234877
    }
]
```

# And all over again

## User sends: 

```
{
    "type": "message",
    "body": {
        "message": "Hi man, whats up?"
    }
}
```

## Bot Responded:

```
[
    {
        "event_id": "b4a5086c-e5ff-4951-82e6-e4d0b0f9e9bb",
        "type": "message",
        "author_id": "bot",
        "body": {
            "message": "Hi! From which ingredients do you want to eat?"
        },
        "creation_timestamp": 641794255
    }
]
```

