# Backscreen homework

## Instructions

### Prerequisites
- Docker

### Running the sync command
```bash
docker compose run --rm sync
```

### Running the API
```bash
docker compose up -d 
```

And any other normal docker commands. The API is configured through the environment variables, 
to run inside the docker compose environment, you can use the `.env.example` file as a template.

## Task

(note: I have rewritten the task text a bit, since it was not clear enough IMO)

Create a microservice in Go language, which fetches currency exchange rates 
from https://www.bank.lv/vk/ecb_rss.xml RSS feed and shows it to users.

The micsroservice has to implement 2 endpoints that return data in JSON format:
1. Return latest currency exchange rates.
2. History of exchange rates for a specific currency.

The microservice also has 2 console commands:
1. Command, to fetch current exchange rates from (https://www.bank.lv/vk/ecb_rss.xml) and store them in database.
    - Select 10 currencies 
    - Each currency has to be fetched in their own request by pretending that
    the endpoint returns data only for one currency at a time.
    - Fetch the data using goroutines.
    - Process the data and store it in the database.
    - Abstract the currency fetching into an interface so it can be adapted to
    use new endpoints
2. Command, that starts the API server

- Include extensive error logging and debug logging.
- Preferably use either MySQL(MariaDB) or Cassandra database.
- Create a GitHub(or any other git platform) project with the source code and README so that microservice can be ran by us.

README must include all instructions on how to set up and run the microservice, preferably using Docker.

## Notes

1. What are ECB exhcnage rates? What does that exactly mean?

These are Euro foreign exchange reference rates that are published by the European Central Bank (ECB), every working day around 16:00 CET.

2. What is the format for the RSS feeds XML?

The XML is an RSS feed from the Latvijas Banka, it contains some metadata, but
the main thing is the `<description>` tag, which contains the exchange rate
under the `<![CDATA[ ]]>` tag.

    1. What is the `CDATA` tag? 


3. How to strucutre the data?

We should store the raw data from each day, in a table, so we can
restore/inspect it for any reason later.

The other table should contain columns:
- Publication date - matches the one that is the XML
- Currency - the currency code
- Exchange rate - should be stored as a string so we don't lose any precision

4. How to structure the API?

Have two simple APIs. 
One returns an object with currency->rate object with date.
Other one returns a data object that wraps the list of from the previous API, 
also has a pagination object, with default page size of 5.

We should include a filter that allows to filter to show only the requested
currency.
We should also include a filter for date range.
