# Authorization Server (Draft)

[![License](https://img.shields.io/badge/license-BSD--3--Clause--Clear-brightgreen)](LICENSE)

The authorization server gives out upload tokens to health authorities / doctors / ...

They can then give these tokens to an infected person to upload his/her contact numbers.

Each upload token is valid once for max. ~5 days.

![Diagram](/diagram.svg)

## Upload token

A token is a random sequence of ~10 numbers and letters. Numbers and letters should be chosen to be well distinguishable in text and speech.

## API for TCN Backend

HTTP API with pinned public keys (peer auth)

```
POST /check_and_invalidate_token 
params: token=XXXXXXXXX
returns: "valid" or "invalid"
```

This has to be somehow transaction safe, i.e. token invalidation and upload should either both fail or none of them. I'm not sure how to best do that right now.

## Interface for Health Authorities/Doctors

Webinterface to generate tokens

Authentication using smartcard e.g. https://uziregister.nl/ for the Netherlands, Praxisausweis (SMC-B) for Germany

## Database schema

```
TABLE tokens:
	token 
	created_at
```
