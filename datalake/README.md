# Data Lake

The `datalake` package is responsible for storing and retrieving raw data.
By adding a layer of abstraction on top of various storage providers, the package enables you to access different storages in the exact same manner, as well as easily switch between them.

The package supports the following storage providers:

- File system
- Amazon S3
- Redis

## Usage

Let's assume we are in the process of indexing the Oasis blockchain, we fetched a list of validators from a node and we want to store that data in an S3 bucket.

### Setting up a data lake

Before we can configure a data lake, we need to initialize the storage provider of our choice.
As we want to store data in an S3 bucket, we need to pass a region code and the bucket name to the `NewS3Storage` function.

```go
storage := datalake.NewS3Storage(
  os.Getenv("AWS_S3_REGION"),
  os.Getenv("AWS_S3_BUCKET"),
)
```

Now we can execute the `NewDataLake` function, passing the storage object as the last parameter (along with the network name and the chain name).

```go
dl := datalake.NewDataLake("oasis", "mainnet", storage)
```

### Serializing a resource

The next step is creating a resource object by serializing our list of validators into the JSON format.

```go
res, err := datalake.NewJSONResource("validators.json", validators)
if err != nil {
  log.Fatal(err)
}
```

In the example above, the `validators` variable points to a slice of validators and `validators.json` is an arbitrary name used to reference the resource.

The package supports the following serialization formats:

- JSON
- Binary
- Base64

### Storing a resource

Once the resource object is created, we can store it in the data lake with the following code:

```go
err := dl.StoreResource(res)
if err != nil {
  log.Fatal(err)
}
```

In case of the S3 storage, the resource is stored in the bucket under the `oasis/mainnet/validators.json` key.

### Retrieving a resource

To retrieve a resource from the data lake we need to run the following code:

```go
res, err = dl.RetrieveResource("validators.json")
if err != nil {
  log.Fatal(err)
}
```

Once again, `validators.json` is a name used to reference the resource.

### Parsing the resource data

Now we need to parse the resource data using the same format it has been serialized with.
In our case it's JSON, so we need to use the `ScanJSON` method.

```go
var validators []Validator

err := res.ScanJSON(&validators)
if err != nil {
  log.Fatal(err)
}
```

As a result, the `validators` slice contains the validator data retrieved from the data lake.

### Checking if a resource is stored

If we want to check if a resource has been stored in the data lake, we can pass its name to the `IsResourceStored` method.

```go
stored, err := dl.IsResourceStored("validators.json")
if err != nil {
  log.Fatal(err)
}
```

The method returns `true` if the resource exists in the data lake and `false` otherwise.

### Storing a resource at height

In case we want to store data associated with a specific height (such as transactions), we need to use the `StoreResourceAtHeight` method and pass the height number as the second parameter.
For example:

```go
res, err := datalake.NewJSONResource("transactions.json", transactions)
if err != nil {
  log.Fatal(err)
}

err = dl.StoreResourceAtHeight(res, 3000000)
if err != nil {
  log.Fatal(err)
}
```

This time, the resource is stored under the `oasis/mainnet/height/3000000/transactions.json` key.

### Retrieving a resource at height

To retrieve a resource stored in such a way, we need to use the `RetrieveResourceAtHeight` method and pass the same height number as the second parameter.

```go
res, err := dl.RetrieveResourceAtHeight("transactions.json", 3000000)
if err != nil {
  log.Fatal(err)
}
```
