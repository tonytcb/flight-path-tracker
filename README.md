[![CI](https://github.com/tonytcb/flight-path-tracker/actions/workflows/makefile.yml/badge.svg)](https://github.com/tonytcb/flight-path-tracker/actions/workflows/makefile.yml)

# Flight Path Tracker App

The goal of this project is design a Golang application serving an HTTP API to return the person's original flight path, given a list of flights.

## Design Solution

The application architecture follows the principles of the **Clean Architecture**, originally described by Robert C. Martin. The foundation of this kind of architecture is the dependency injection, producing systems that are independent of external agents, highly testable and easier to maintain.

You can read more about Clean Architecture [here](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

### Flight Path algorithm

The algorithm used to find out the original flight path is quite simple. It consists in search the **origin** that never appears as a destination, and the opposite as well, search the **destination** that never appears as a origin. Within these values in hands, we have the first origin and last destination.

As a performant algorithm, it runs on a linear time complexity, O(n).

## Tools

- [Golang 1.21](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Docker-compose](https://docs.docker.com/compose/)

## HTTP API

### Calculate

- Method: `POST`
- Path: `/calculate`
- Headers
- - Content-Type: application/json
- Payload:
```json
[
    {
        "source": "IND",
        "destination": "EWR"
    },
    {
        "source": "SFO",
        "destination": "ATL"
    },
    {
        "source": "GSO",
        "destination": "IND"
    },
    {
        "source": "ATL",
        "destination": "GSO"
    }
]
```
- Response:
```json
{
    "source": "SFO",
    "destination": "EWR"
}
```

## Commands

- `make help` to see all commands;
- `make up` to starts the app serving the http api;
- `make test` to run all tests.

## TODO Improvements

- [ ] validate corner cases: today the api validation is pretty simple, validating only non-empty and duplicated flights
- [ ] flights data generator: implement a function returning a huge list of flights
- [ ] benchmark: would be a good improvement running a benchmark to evaluate the current algorithm, and compare with future changes.

 