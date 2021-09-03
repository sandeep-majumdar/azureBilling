# Observability

## Logging

## Logging controls

```bash
export LOG_ENABLED=true
# DEBUG, INFO, WARN, ERROR
export LOG_LEVEL=INFO
```

### Logging Usage

- "Exit" - Exits
- "Fatal" - Panics
- "Debug"
- "Info"
- "Warn"
- "Error"

```go

observability.SetAppName("myapp")
observability.Logger("Info", "my string"))

```

### Correlation and Causation

Set or generate causation and correlationids to allow tracking of log messages across technical domains (apis & events)

```go
observability.GenCorrId()
observability.SetCausationId(observability.GetCorrId())
```

Example output from main.go

```txt
D 2020-10-28 14:52:38.2849 [NoCausationId => NoCorrId]  main.go:9       Debug message
I 2020-10-28 14:52:38.2849 [NoCausationId => NoCorrId]  main.go:10      Info message
W 2020-10-28 14:52:38.2850 [NoCausationId => NoCorrId]  main.go:11      Warn message
E 2020-10-28 14:52:38.2850 [NoCausationId => NoCorrId]  main.go:12      Error message
I 2020-10-28 14:52:38.2850 [NoCausationId => f0d7a028-abe8-4b84-8004-15873744b289]  main.go:15  Test message with CorrId
I 2020-10-28 14:52:38.2851 [f0d7a028-abe8-4b84-8004-15873744b289 => 322ca6d3-f9ee-459b-9e48-77fdfc2b0b05]  main.go:18   Test message with CorrId and CausationId
I 2020-10-28 14:52:38.2851 [f0d7a028-abe8-4b84-8004-15873744b289 => 322ca6d3-f9ee-459b-9e48-77fdfc2b0b05]  main.go:22   Timing Stage 1 - should print completed in 0.00 ms
I 2020-10-28 14:52:38.2853 [f0d7a028-abe8-4b84-8004-15873744b289 => f643203a-b282-40a5-bb18-11d4b96ceb7b]  main.go:27   MEMORY Alloc = 0 MiB    TotalAlloc = 0 MiB      Sys = 68 MiB    NumGC = 0
I 2020-10-28 14:52:38.2853 [f0d7a028-abe8-4b84-8004-15873744b289 => de11535f-6d4d-4b37-ba37-11deeb478e72]  Metrics.go:64        metric count 1.000000
I 2020-10-28 14:52:38.2854 [f0d7a028-abe8-4b84-8004-15873744b289 => de11535f-6d4d-4b37-ba37-11deeb478e72]  Metrics.go:64        metric decimal 0.100000
F 2020-10-28 14:52:38.2854 [f0d7a028-abe8-4b84-8004-15873744b289 => 2500cd7b-d4a2-409c-b386-38533beacba9]  main.go:37   Fatal message
exit status 3
```

## Timing

```go
func (timer *Timer) Start(timing bool, str string) {
func (timer *Timer) EndAndPrint(timing bool) {
func (timer *Timer) EndAndPrintStderr(timing bool) {
```

```go
timingOn := true // toggle, eg set via environment variable
t1 := observability.Timer{}
t1.Start(timingOn, fmt.Sprintf("loadImageFromFile=%s", path))
// do stuff
t1.EndAndPrintStderr(timingOn)
```

## Write memory consumption to log output

```go
observability.LogMemory("Info")
```

## Metrics

```go
func (ms *Metrics) setKeyValue(key string, m Metric)
func (ms *Metrics) SetDuration(key string, d time.Duration)
func (ms *Metrics) SetInteger(key string, i int)
func (ms *Metrics) SetFloat(key string, f float64)
func (ms *Metrics) Dump()
```

```go
type a {
    metrics     observability.Metrics
}

a.metrics.Init()

t1 = time.Now()
// Do load
count++
t2 = time.Now()
ela = t2.Sub(t1)
a.metrics.SetInteger("LOAD_COUNT", count)
a.metrics.SetDuration("LOAD_TOTAL_TIME_S", ela)
```
