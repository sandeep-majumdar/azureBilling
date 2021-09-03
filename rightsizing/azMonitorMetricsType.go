package rightsizing

type azMonitorMetricsType struct {
	Cost           int           `json:"cost"`
	Interval       string        `json:"interval"`
	Namespace      string        `json:"namespace"`
	Resourceregion string        `json:"resourceregion"`
	Timespan       string        `json:"timespan"`
	Value          []metricValue `json:"value"`
}

type metricValue struct {
	DisplayDescription string             `json:"displayDescription"`
	ErrorCode          string             `json:"errorCode"`
	Id                 string             `json:"id"`
	Name               metricName         `json:"name"`
	ResourceGroup      string             `json:"resourceGroup"`
	Timeseries         []metricTimeSeries `json:"timeseries"`
	MetricType         string             `json:"type"`
	Unit               string             `json:"unit"`
}

type metricName struct {
	LocalizedValue string `json:"localizedValue"`
	Value          string `json:"value"`
}

type metricTimeSeries struct {
	Data           []metricObservation `json:"data"`
	MetadataValues []string            `json:"metadataValues"`
}

type metricObservation struct {
	TimeStamp string  `json:"timeStamp"`
	Average   float64 `json:"average"`
	Count     float64 `json:"count"`
	Maximum   float64 `json:"maximum"`
	Minimum   float64 `json:"minimum"`
	Total     float64 `json:"total"`
}

func (mo *metricObservation) getValuesArray() (a []float64) {
	a = append(a, mo.Average)
	a = append(a, mo.Maximum)
	a = append(a, mo.Minimum)
	a = append(a, mo.Count)
	a = append(a, mo.Total)
	return a
}

/*
{
	"cost": 2879,
	"interval": "1 day, 0:00:00",
	"namespace": "Microsoft.Compute/virtualMachines",
	"resourceregion": "westeurope",
	"timespan": "2021-08-24T00:00:00Z/2021-08-26T00:00:00Z",
	"value": [
	  {
		"displayDescription": "The percentage of allocated compute units that are currently in use by the Virtual Machine(s)",
		"errorCode": "Success",
		"id": "/subscriptions/ecd48036-a17a-47c1-90f9-5a13975853a3/resourceGroups/rgpazewnmlit001prodpohdb/providers/Microsoft.Compute/virtualMachines/anmlp1podhdb01/providers/Microsoft.Insights/metrics/Percentage CPU",
		"name": {
		  "localizedValue": "Percentage CPU",
		  "value": "Percentage CPU"
		},
		"resourceGroup": "rgpazewnmlit001prodpohdb",
		"timeseries": [
		  {
			"data": [
			  {
				"average": 1.5980469226081657,
				"count": null,
				"maximum": null,
				"minimum": null,
				"timeStamp": "2021-08-24T00:00:00+00:00",
				"total": null
			  },
			  {
				"average": 0.9500577528876444,
				"count": null,
				"maximum": null,
				"minimum": null,
				"timeStamp": "2021-08-25T00:00:00+00:00",
				"total": null
			  }
			],
			"metadatavalues": []
		  }
		],
		"type": "Microsoft.Insights/metrics",
		"unit": "Percent"
	  }
	]
  }
*/
