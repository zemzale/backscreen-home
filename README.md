# Backscreen homework

## Task

Create a microservice in Go language, which fetches currency exchange rates from https://www.bank.lv/vk/ecb_rss.xml RSS feed and shows it to users.
The microservice consists of 2 endpoints:
1. JSON Data return from database - latest currency exchange rates.
2. JSON Data return from database - History of exchange rates for a specific currency.

The microservice also has 2 console commands:
1. Command, which fetches current exchange rates from (https://www.bank.lv/vk/ecb_rss.xml) and saves them to database.
    - Create this fetching using go routines by fetching each of 10 preselected currencies in their own requests like the HTTP endpoint would only be returning the result for one currency each time.
    - So request 1 fetches info for currency GBP pulls only that currency info from the request and processes it.
    - Request 2 then again fetches the endpoint and pulls its own currency info.
    - Create the fetching of currency info into an interface so it can easily be adapted to new endpoints for fetching data.
2. Command, which starts the microservice so that the endpoints are accessible to users.

Include extensive error logging and debug logging.
Preferably use either MySQL(MariaDB) or Cassandra database.
Create a github(Or any other version control system that uses git) project with all of the code and readme so that we can run the microservice ourselves.
Readme must include all instructions on how to set up and run the microservice, preferably on Docker.
