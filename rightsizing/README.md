# Rightsizing

## Overview

This is an experimental capability still under development.

The following files are produced:

- vmDetails - list of VMs with high level info
- vmDayValues - list of VMs by P+Q per day during the month
- vmMonitorMetrics - azure monitor output per hour for the month

The next step will be to do something useful with them

Output naming is described by the config file

```txt
"outputVmDetailsCSVFile": "vmDetails.csv",
"outputVmDayValuesCSVFile": "vmDayValues.csv",
"outputVmMonitorMetricsFile": "vmMonitorMetrics.json",
"rightsizingMaxThreads": 20
```

## Algorithm

1. Parse the azure bill file for VMs
2. Write the summary output to config.OutputVmDetailsCSVFile
3. Write the daily costs and quantity to config.OutputVmDayValuesCSVFile
4. For each VM in the summary output, query azure monitor for %age CPU and Available Mem
5. Write the metrics to config.OutputVmMonitorMetricsFile
6. If any of the above files exist, the processing is skipped
7. *Processing to be written*


## HowTo

```bash
go run cmd/main.go azure rightsizing vm -c ./config.json
```

### Useful command to give each VM its own newline in the file

```bash
cat vmMonitorMetricsOrig.json | sed 's/},\"\//},\n\"\//g' > vmMonitorMetricsParsed.json
head -100 vmMonitorMetricsParsed.json > vmMonitorMetrics.json 
```

## Sample output

Rightsizing processing of 100 VMs takes approximately 1m20s with MAX_THREADS=20 on a core i7

```txt
I 2021-09-03 15:19:51.8872 [NoCausationId => NoCorrId]  platformMapLookup.go:77 MEMORY Alloc = 6 MiB    TotalAlloc = 9 MiB      Sys = 70 MiB    NumGC = 2
I 2021-09-03 15:19:51.8873 [NoCausationId => NoCorrId]  platformMapLookup.go:27 PlatformMapLookup has 9820 records
I 2021-09-03 15:19:51.8918 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmRightsizing.go:46     Loading from output files
I 2021-09-03 15:19:52.0126 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmDetails.go:93 Loaded 7887 vmDayValues from file /mnt/c/Users/adria/Downloads//vmDetails.csv
I 2021-09-03 15:19:52.0127 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmDetails.go:94 MEMORY Alloc = 9 MiB    TotalAlloc = 15 MiB     Sys = 71 MiB    NumGC = 3
I 2021-09-03 15:19:53.4110 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmDayValues.go:111      Loaded 98149 vmDayValues from file /mnt/c/Users/adria/Downloads//vmDayValues.csv
I 2021-09-03 15:19:53.4111 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmDayValues.go:112      MEMORY Alloc = 91 MiB   TotalAlloc = 113 MiB    Sys = 138 MiB   NumGC = 7
I 2021-09-03 15:19:53.4120 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmRightsizing.go:154    Processing Metrics
I 2021-09-03 15:19:53.4122 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  AzureCli.go:33  Initialised azcli
I 2021-09-03 15:19:53.4122 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmRightsizing.go:162    Created azure cli 0 : len(azClis)=1
I 2021-09-03 15:19:54.8911 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  AzureCli.go:33  Initialised azcli
I 2021-09-03 15:19:54.8912 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmRightsizing.go:162    Created azure cli 1 : len(azClis)=2
I 2021-09-03 15:19:55.8772 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  AzureCli.go:33  Initialised azcli
I 2021-09-03 15:19:55.8773 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmRightsizing.go:162    Created azure cli 2 : len(azClis)=3
I 2021-09-03 15:19:56.8873 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  AzureCli.go:33  Initialised azcli
....
I 2021-09-03 15:20:15.1875 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmRightsizing.go:162    Created azure cli 18 : len(azClis)=19
I 2021-09-03 15:20:16.0789 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  AzureCli.go:33  Initialised azcli
I 2021-09-03 15:20:16.0791 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmRightsizing.go:162    Created azure cli 19 : len(azClis)=20
I 2021-09-03 15:20:16.9568 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmRightsizing.go:166    Azure cli creation complete
I 2021-09-03 15:20:41.9332 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmRightsizing.go:201    Processed 100 virtual machines
I 2021-09-03 15:20:41.9332 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmRightsizing.go:202    Have 100 vm metric records
I 2021-09-03 15:20:41.9334 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmRightsizing.go:203    MEMORY Alloc = 167 MiB  TotalAlloc = 228 MiB    Sys = 205 MiB   NumGC = 8
I 2021-09-03 15:21:11.3446 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmRightsizing.go:201    Processed 200 virtual machines
I 2021-09-03 15:21:11.3447 [1 => c96f96f6-6be4-47d7-b0e2-cd0ae7dda3ec]  vmRightsizing.go:202    Have 200 vm metric records
I 2021-09-03 16:55:42.2363 [1 => 837f2e0b-e33e-42d1-8ba5-4f82a2a97f81]  vmRightsizing.go:210    Writing vm day values file to /mnt/c/Users/adria/Downloads/vmMonitorMetrics.json
```
