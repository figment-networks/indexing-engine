# Pipeline

The `pipeline` package helps to build indexers using simple to use DSL. Its goal is to provide a logical structure to the process of indexing.

## Default pipeline

Every default indexing pipeline has fixed stages available to hook in to:
* `Setup stage`: performs setup tasks
* `Syncer stage`: creates syncable
* `Fetcher stage`: fetches data for indexing
* `Parser stage`: parses and normalizes fetched data to a single structure
* `Validator stage`: validates parsed data
* `Sequencer stage`: Creates sequences from fetched or/and parsed data
* `Aggregator stage`: Creates aggregates from fetched or/and parsed data
* `Persistor stage`: Saves data to data store
* `Cleanup stage`: Cleans up after execution

Besides that there are 2 additional components: Source and Sink.
`Source` is responsible for providing height iterator for the pipeline and `Sink` is gathering output data which can be used after the pipeline is done processing.

Below flow-chart depicts all the available stages that come with this package and order in which they are executed.

![Pipeline flow](/assets/pipeline-flow.jpg)

Please note that all syncing phase stages are executed in sequence, whereas indexing phase stages are executed concurrently in order to speed up the indexing process.

### Creating new pipeline

To create a new default pipeline use:
```go
NewDefault([payloadFactory])
```

### Setting up tasks in stages

In order to set tasks for a specific stage you can use:
```go
p.SetTasks(
  [Name of the stage],
  NewTask(),
  NewTask(),
)
```

To set tasks that will run at the same time, use:
```go
p.SetAsyncTasks(
  [Name of the stage],
  NewTask(),
  NewTask(),
)
```

If you want to use your own method of running task inside of a stage, you can easily create your own implementation of a `StageRunnerFunc` and pass it to `SetCustomStage`.

```go
p.SetCustomStage(
  [Name of the stage],
  [custom Stagerunnerfunc]
)
```

### Starting pipeline

Once stages are setup, we can run our pipeline
```go
options := &pipeline.Options{}
if err := p.Start(ctx, NewSource(), NewSink(), options); err != nil {
    return err
}
```
This will execute all the tasks for every iteration of all the items in the source created with `NewSource()`
If you want to run one-off iteration of pipeline for specific height you can use `Run()`
```go
height := 100
payload, err := p.Run(ctx, height, options)
```
It will return a payload collected for that one iteration of the source.

### Adding custom stages

If you want to perform some action on but provided stages are not good logic fit for it, you can always add
custom stages BEFORE or AFTER existing ones. In order to do that you can use:
* `AddStageBefore` - adds stage before provided existing stage
* `AddStageAfter` - adds stage after provided existing stage
Below is an example showing how you can add custom stage (as a func) after Fetcher stage
```go
const (
    CustomStageName = "AfterFetcher"
)

afterFetcherFunc := pipeline.StageRunnerFunc(func(ctx context.Context, p pipeline.Payload, f pipeline.TaskValidator) error {
    //...
    return nil
})

p.AddStageAfter(
  pipeline.StageFetcher,
  pipeline.NewCustomStage(CustomStageName, afterFetcherFunc)
)
```

 ### Retrying

 The indexing pipeline provides 2 types of retrying mechanisms:
 * `RetryStage` - which is responsible for retrying the entire stage if error occurred
 * `RetryingTask` - which is responsible for retrying individual tasks if it return error

 In order to implement retrying mechanism you need to wrap stage or task with above functions.
 Here is an example of use of `RetryingTask`:
 ```go
p.SetTasks(
  pipeline.StageFetcher,
  pipeline.RetryingTask(NewFetcherTask(), func(err error) bool {
      // Make error always transient for simplicity
      return true
  }, 3),
)
```

### Selective execution

Indexing pipeline provides you with options to run stages and individual tasks selectively.
You have 2 options you can use for this purpose:
* `StagesBlacklist` - list of stages to NOT execute
* `TasksWhitelist` - list of indexing tasks to execute

In order to use above options you have to use `setOptions` method of pipeline like so:
```go
p.SetOptions(&pipeline.Options{
    TasksWhitelist: []string{"SequencerTask"},
})
```
Above example would run only `SequencerTask` during indexing process. It is useful if you want to reindex the data but you only care about specific set of data.

## Custom pipeline

If the default pipeline stages and run order does not suit your needs, then you can create an empty pipeline where you must add each stage individually.

```go
p := pipeline.NewCustom(payloadFactory)
```

This creates a pipeline with no set stages. To add a stage, run:

```go
p.AddStage(
  pipeline.NewStage([Name of the stage],  NewTask())
)
```

To add a stage where tasks will run concurrently, run:
```go
p.AddStage(
  pipeline.NewAsyncStage([Name of the stage],  NewTask1(), NewTask2())
)
```

The order in which the stages will run is determined by the order in which they are added.

If you want to run stages concurrently, add them together using `AddConcurrentStages`

```go
p.AddConcurrentStages(
  pipeline.NewStage(pipeline.StageAggregator, NewTask())
  pipeline.NewStage(pipeline.StageSequencer, NewTask())
)
```

## Built-in metrics

The indexing pipeline comes with a set of built-in metrics:

| Name                               | Description                                       |
|------------------------------------|---------------------------------------------------|
| `indexer_pipeline_task_duration`   | The total time spent processing an indexing task  |
| `indexer_pipeline_stage_duration`  | The total time spent processing an indexing stage |
| `indexer_pipeline_height_duration` | The total time spent indexing a height            |
| `indexer_pipeline_heights_total`   | The total number of successfully indexed heights  |
| `indexer_pipeline_errors_total`    | The total number of indexing errors               |

For more information about metrics, see the documentation of the [`metrics`](/metrics) package.

## Examples

In `example` folder you can find an example of a pipeline. To run it use:
```shell
go run example/default/main.go

go run example/custom/main.go
```
