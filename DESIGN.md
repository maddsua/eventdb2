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
10. Single tenant + multiuser

Doesn't look like Grafana killa, does it? Well because it isn't! It would be much more close to the analitics inside of DigitalOcean's dashboard.

## Main draft

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

Here, `binary array` indicates a binary data structure that is used instead of JSON in order to improve parsing speed. It contains of one or more messages written to the data stream sequentially.
```
---> stream
---------------------------------------------
|      message     |      message      |   ...
---------------------------------------------
```

Messages consist of two integers indicating content size followed by the raw contents.
```
---------------------------------------------------------------------
|                                message                            |
---------------------------------------------------------------------
|  key_size  |  data_size  |      key raw      |      data raw      |
---------------------------------------------------------------------
|  uint8_le  |  uint16_le  |       bytes       |        bytes       |
---------------------------------------------------------------------
```

So we do waste entire 3 bytes on data size indication, which is nothing comparing to the overheard of having JSON or url encoding, which looks justifiable to me. Oh an also, this WILL handle all the weird cases of unicode and whatnot (YES I AM LOOKING AT YOU LOKI)

Even tho the used int sizes would limit label key and label content sizes to 256 and 65536 bytes respecively, at actual maximal allowed size should be limited to 200 bytes for the key and a 1000 bytes for the value. There's no need to allow anyone to just dump huge amounts of data here. Labels should only be used for filtering, for everythyng else there's the metadata field.


### Structued log metadata

With this one it's dead simple - it's literally just a plain JSON object of the following format:
```jsonc
{
  "key": "value",
  // ...
  "client_ip": "127.0.0.1",
}
```

There aren't any technical limitations here, expect for both keys and values having to be strings. Key size should be restricted to around 100 symbols though to avoid having absurdly long keys that would break the UI. Values could be much longer but it still makes sense to limit them to let's say a 1000 characters (not bytes, unlike the labels). These limits can be freely adjusted at a later stage.


### http uptime checks

Pretty much just joink it from <https://github.com/maddsua/pulse> but instead of having a config file just pull the options from the database. To avoid excessive writes the minimal poll interval should be set to a sensible. 15 seconds should be sufficient. 10 probes writing data each 15 seconds would generate (24 * 60 * 4) * 28 = 161,280 entries monthly which is well inside the SQlite cpabilities.

On the UI side of things, I want to display a simple response time and http code graphs, average uptime percentage (which was an absolute pain in the ass to achieve with Grafana btw) and probably some other calculated values.


### Event notifications

Singe everything is sitting in the same process anyway it is trivial to monitor data changes and dispatch notifications when certain conditions are met. This feature can be pretty much just copied directly from v1 minus maybe adding message templates or something similar.


### Dashboard API

Using REST for it the previous time was a huge mistake. It took too much time and effort to write it. Using GraphQL would make this task much easier since it will generate all the boilerplace code for you. Oh, and not just on the server but also no the client!


### Ingest REST API

However, when it comes to the ingest API there's no reason to overcomplicate things. We won't have many endpoints here and the ones that we'll do would be pretty simple.


### TypeScript client library

To make it simpler to connect an already existing app to eventdb using a client library is preferred. The interfaces similar to <https://github.com/maddsua/logpush> can be used to connect a `console.log`-like class to the REST API of the backend.


### Single tenant + multiuser

Simply put: I want to have a way to add readonly or otherwise limited users to the dashboard but at the same time I am not building it to scale to the billions, so if you want to use it for, let's say, different customers - just deploy multiple services instead. I can't be bothered with writing the overly comples system to handle different resoucres access and sharing.
