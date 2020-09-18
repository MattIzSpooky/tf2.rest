## tf2.rest
A Team Fortress 2 Responses API built in Go.

Made by me and [Captain Hippie](https://github.com/Captain-Hippie)

### Endpoints
[The root endpoint](https://api.tf2.rest/) returns a random response.
```json
{
    "class": "medic",
    "response": "Ze healing leaves little time for ze hurting.",
    "audioFile": "https://wiki.teamfortress.com/w/images/4/47/Medic_specialcompleted03.wav",
    "type": "Event-related Responses",
    "context": "Under the effects of an ÜberCharge"
}
```

[The by class endpoint (scout example)](https://api.tf2.rest/by-class/scout) returns a random response by given class.
```json
{
    "class": "scout",
    "response": "Hey good job there, hardhat!",
    "audioFile": "https://wiki.teamfortress.com/w/images/7/79/Scout_thanksfortheteleporter02.wav",
    "type": "Event-related",
    "context": "Teleportation"
}
```