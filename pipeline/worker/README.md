# Worker

The `worker` package implements a simple load balancing algorithm based on the [master-worker scheme](https://en.wikipedia.org/wiki/Load_balancing_(computing)#Master-Worker_Scheme).
In this algorithm, the workload is distributed among multiple worker processes by a single manager process.

## Overview

The diagram below shows the available components and how they relate to each other:

![Worker components](/assets/worker-components.svg)

### Components

*Click the name of a component to see its description.*

<details>
  <summary><strong>Pool</strong></summary>

  Distributes the workload among registered pool workers.
</details>

<details>
  <summary><strong>PoolWorker</strong></summary>

  Represents a worker on the manager side.
  Uses a client to communicate with the worker.
</details>

<details>
  <summary><strong>Client</strong></summary>

  Sends requests to a server and receives responses from it.
  Used on the manager side.
</details>

<details>
  <summary><strong>Server</strong></summary>

  Receives requests from a client and sends back responses.
  Used on the worker side.
</details>

<details>
  <summary><strong>Loop</strong></summary>

  Represents the processing loop of a worker.
  Uses a server to communicate with a manager.
</details>

## Communication

Although the communication between a worker and a manager is protocol-independent, the package comes with a default implementation based on the [WebSocket protocol](https://en.wikipedia.org/wiki/WebSocket).
You can find it in the `websocket.go` file.

## Backoff algorithm

The package includes an implementation of the [exponential backoff](https://en.wikipedia.org/wiki/Exponential_backoff) algorithm.
This mechanism can be used to decrease the rate of requests in case of repeated failures.
You can find the implementation of the algorithm in the `backoff.go` file.
