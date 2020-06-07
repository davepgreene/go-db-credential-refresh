# Custom Store Example

This example is shows the Connector in operation. The only dependency required for this example is a running PostgreSQL instance. You can easily set one up by running:

```shell
docker run --name pg -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres
```

Once the database is ready, you can run the example:

```shell
go run main.go
```

Creating a custom store and integrating it is very simple. See `store.go` in this directory for an example and
`main.go` for how it would be used.

Stores must implement the `driver.Store` interface.

``` go
type Store interface {
    Get() (driver.Credentials, error)
    Refresh() (driver.Credentials, error)
}
```

The distinction between `Get()` and `Refresh()` is one that allows caching or other early retry mechanism before a
more complex operation is done. An example of this would be if you had a token-based auth to a system that held
credentials, you might want to use the token you have to try to get the creds via the `Get()` method. If, however, the token had expired or otherwise been invalidated, you may want to go through a more complex and lengthy process of renewing the token and then retrieving credentials by using `Refresh()`.
