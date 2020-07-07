# Indexing Engine

## Description
Indexing engine helps to build indexers using simple to use DSL. It's goal is to provide a logical structure to 
the process of indexing.

Every indexing pipeline has fixed stages available to hook in to:
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

![indexing engine flow chart](/diagram.jpg)

Please note that all syncing phase stages are executed in sequence, whereas indexing phase stages are executed concurrently
in order to speed up the indexing process. 

## Installation

To install github.com/figment-networks/indexing-engine use:
```shell script
go get https://github.com/figment-networks/github.com/figment-networks/indexing-engine
```

## Usage

### Setting up stages
In order to set a specific stage you can use:
```shell script
p.SetStage(
  [Name of the stage],
  pipeline.SyncRunner(NewTask()),
)
```

As a parameter to `p.SetStage` function you pass in stage name and StageRunner instance.
StageRunner is responsible for running individual tasks inside of the stage.
This package provides 2 types of StageRunners:
1. `SyncRunner` - executes tasks one by one
2. `AsyncRunner` - executes tasks concurrently

If you want to use your own method of running task inside of stage, you can easily
create your own implementation of StageRunner and pass it in to `p.Set[StageName]Stage`.

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

p.AddStageBefore(pipeline.StageFetcher, CustomStageName, afterFetcherFunc)
```

### Running pipeline stages in custom order

If you want to run stages in a different order than outlined in the flowchart,
then create a new custom pipeline

```go
p := pipeline.NewCustom(payloadFactory)
```

This creates a pipeline where the stage dependencies are empty. You must set
a stage dependency with every stage. If a stage has no dependencies, call `SetDependency` with an empty list

```go

p.SetStage(
  pipeline.StageSetup,
  pipeline.SyncRunner(NewTask()),
)

p.SetDependency(pipeline.StageSetup, []pipeline.StageName{})
```

Before a stage runs, all its dependencies must run first. In the example below, `StageSetup` is set as a dependency of `StageFetcher`. This means `StageFetcher` will
run only after `StageSetup` has completed.

```go
p.SetStage(
  pipeline.StageFetcher,
  pipeline.SyncRunner(NewTask()),
)

p.SetDependency(pipeline.StageFetcher, []pipeline.StageName{pipeline.StageSetup})
```


If multiple stages are able to run, then they will run in parallel. In the example below, `StageAggregator` and `StageSequencer` both depend on `StageFetcher`. Once `StageFetcher` is complete, both `StageAggregator` and `StageSequencer` will run concurrently.


```go
p.SetStage(
  pipeline.StageAggregator,
  pipeline.SyncRunner(NewTask()),
)

p.SetDependency(pipeline.StageAggregator, []pipeline.StageName{pipeline.StageFetcher})


p.SetStage(
  pipeline.StageSequencer,
  pipeline.SyncRunner(NewTask()),
)

p.SetDependency(pipeline.StageSequencer, []pipeline.StageName{pipeline.StageFetcher})


```

All running stages must finish before starting any new stages. In the example below, although `StageCleanup` only depends on `StageSequencer`, both `StageAggregator` and `StageSequencer` need to finish (since they have both started because they are running in parallel) before `StageCleanup` is run.


```go
p.SetStage(
  pipeline.StageCleanup,
  pipeline.SyncRunner(NewTask()),
)

p.SetDependency(pipeline.StageCleanup, []pipeline.StageName{pipeline.StageSequencer})

```


 ### Retrying
 github.com/figment-networks/indexing-engine provides 2 types of retrying mechanisms:
 * `RetryingStageRunner` - which is responsible for retrying the entire stage if error occurred
 * `RetryingTask` - which is responsible for retrying individual tasks if it return error

 In order to implement retrying mechanism you need to wrap stage or task with above functions.
 Here is an example of use of `RetryingTask`:
 ```go
p.SetFetcherStage(
    pipeline.AsyncRunner(
        pipeline.RetryingTask(NewFetcherTask(), func(err error) bool {
            // Make error always transient for simplicity
            return true
        }, 3),
    ),
)
``` 

### Selective execution
Indexing engine provides you with options to run stages and individual tasks selectively.
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

## Examples
In `/examples` folder you can find an example of a pipeline. To run it use:
```shell script
go run example/main.go
```

## To-Dos:
* Collect stats for every stage and every task
* Add context cancellation