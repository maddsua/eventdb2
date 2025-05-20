# Basic idea

This project is planned as a complete rewrite of <https://gitlab.com/maddsua/eventdb>.

What's eventdb? It was my attempt a few years ago to get logs from all the apps that I maintained back then in one place.

It was quite annoying having to jump over all the different dashboards just to tell what the heck is wrong with that "serverless" function. This is how the idea was born.

Over time I tried to embrace the industry standard (Grafana and it's stack) but got unsatisfactory results:

1. The entire setup process is too cumbersome for such use as it involves setting up at least 3 services just to get some meaningful logs recorded.
2. The performance overhead is a bit too large for a use case where you get at most a few hundred log entrie per day.
3. SQL data source support (pg, timescale) is quite crap for logs and, to be frank, just any general data points
4. There is just way too much shit going on in grafana itself

And keep in mind that I have to run that stack just to monitor a few old friends-of-friends apps (which are technically prod, yes) or just some mess-around deployments, which makes running the Grafana stack seem like a huge overhead.

So, the basic requirements that we get in the end here are:

1. Log aggregation for volumes of under 100,000 entries per day
2. Support for arbitrary label filtering
3. Structured metadata
4. http service uptime checks
5. Event notifications
6. Simple metrics recording
7. User GraphQL API
8. Public data ingest REST API
9. TypeScript client library

Doesn't look like Grafana killa, does it? Well because it isn't! It would be much more close to the analitics inside of DigitalOcean's dashboard.
