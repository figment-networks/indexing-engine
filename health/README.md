# Health

Health package is a a source of all the default metrics and health measurements we may like to get from the services.
It's written in a way so it would be easily extensible. Any structure fulfilling interface `Prober`, would automatically attach all the necessary information into process.


## HTTP Interface

- `/health` - Endpoint that always return 200 for dumb health check
- `/readiness` - Endpoint meant to return information from different probes and return meaningful information about is service ready for connecting it to the traffic.

## Metrics
Service produces a set of metrics

### Database
- `health_database_ping` - Duration how long it takes to ping the database (float seconds)
    - Tags:
        - `database_type` - The name od database type (eg. `postgres`)
- `health_database_size` - Current size of database
    - Tags:
        - `database_type` - The name od database type (eg. `postgres`)
