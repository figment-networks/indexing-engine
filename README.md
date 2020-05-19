# Indexing Engine

## Description
Indexing engine helps to build indexers using simple to use DSL. It's goal is to provide a logical structure to 
the process of indexing.

Every indexing pipeline has fixed stages available to hook up to:
* `Setup stage` (Chore): performs setup tasks
* `Fetcher stage` (Syncing): fetches data for indexing
* `Parser stage` (Syncing): parses and normalizes fetched data to a single structure
* `Validator stage` (Syncing): validates parsed data 
* `Syncer stage` (Syncing): saves data to datastore
* `Sequencer stage` (Indexing): Creates sequences from synced data (syncable)
* `Aggregator stage` (Indexing): Creates aggregates from synced data (syncable)
* `Cleanup stage` (Chore): Cleans up after execution

Below flow-chart depicts all the available stages that come with this package and order in which they are executed.

![indexing engine flow chart](/diagram.jpg)

Please note that all syncing phase stages are executed in sequence, whereas indexing phase stages are executed concurrently
in order to speed up the indexing process. 

## Installation

To install github.com/figment-networks/indexing-engine.git use:
```shell script
go get https://github.com/figment-networks/github.com/figment-networks/indexing-engine.git
```

## Usage

### Setting up stages
In order to set a specific stage you can use:
```shell script
p.Set[StageName]Stage(
  pipeline.SyncRunner(NewTask()),
)
```

As a parameter to `p.Set[StageName]Stage` functions you pass in StageRunner instance.
StageRunner is responsible for running individual tasks inside of the stage.
This package provides 2 types of StageRunners:
1. `SyncRunner` - executes tasks one by one
2. `AsyncRunner` - executes tasks concurrently

If you want to use your own method of running task inside of stage, you can easily
create your own implementation of StageRunner and pass it in to `p.Set[StageName]Stage`.

## Adding custom stages
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


 ### Retrying
 github.com/figment-networks/indexing-engine.git provides 2 types of retrying mechanisms:
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
Indexing engine provides you with options to run stages and individual indexing tasks selectively.
You have 4 options you can use for this purpose:
* `StagesWhitelist` - list of stages to execute
* `StagesBlacklist` - list of stages to NOT execute
* `IndexingTasksWhitelist` - list of indexing tasks to execute 
* `IndexingTasksBlacklist` - list of indexing tasks to NOT execute

In order to use above options you have to use `setOptions` method of pipeline like so:
```go
p.SetOptions(&pipeline.Options{
    IndexingTasksWhitelist: []string{"SequencerTask"},
})
```
Above example would run only `SequencerTask` during indexing process