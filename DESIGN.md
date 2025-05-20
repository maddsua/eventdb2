## Basic idea

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

## Overall composition

### Microservices?

The short answer is - nah.

If you need the wildly flexible architecture that supports sharding and replication - just stick with Grafana. This project tries to achieve quite the opposite. The point here is to deploy it to a VPS or something like <railway.com> and just forget about it existing.

I do not want to give it a massive scalability potential that will inevitably be wasted. Again, if your apps generates literal gigabytes of logs - there are other tools that are made specifically to handle that. Go use Lokie or something similar.

Realistically, the use case for this thing is to keep logs of your Vercel, Netlify or Cloudflare Workers deployments and if you manage to get any of them to send more than a few hundreds logs per hour - you have a much bigger problems than the vertical scaling that eventdb might require.


### The DB?

Since we've established that we do not need to handle huge amounts of data AND we also want to keep things in a single container, SQLite seems like an obvoius choise.

But isn't it slow? Yeah sorta, but only if you don't use WAL and try to do crap like CIDR matching using a hand crafted query that literally does the bigint math.

Otherwise, SQLite is more than capable to provide the needed base storage for this application.

I should probably make a storage interface so that it would be possible to use an external DB like timescale in the future, but that's not something I see doing in the forseable future.


### Log labeling or filters

Only the basic data such as the timestamp and stream id would be written as separate SQL columns, while labels and tags should be stored as serialized filers.

Even though that would increase the data transfer between the DB and the dashboard backend - this should be sufficient enough considering the requirements.

This SQL table could be used as a reference:
```sql
------------------------------------------------------------------------------------------------
|    id    |    timestamp    |    stream_id    |    message    |    labels    |    metadata    |
------------------------------------------------------------------------------------------------
|  integer |     integer     |       blob      |      text     |     blob     |      blob      |
------------------------------------------------------------------------------------------------
|  serial  |    unix epoch   |       uuid      |  raw message  | binary array |    json map    |
------------------------------------------------------------------------------------------------
```
