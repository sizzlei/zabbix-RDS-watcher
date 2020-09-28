package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"flag"
	"time"
	"fmt"
	"strings"
)

func main() {
	IninstanceId		:= flag.String("identifier","","DBInstanceIdentifier")
	InRegion			:= flag.String("region","ap-northeast-2","AWS Region") 
	InMetric			:= flag.String("metric","","AWS CloudWatch Database Metric")
	InClass				:= flag.String("class","RDS","RDS or REDSHIFT ")

	flag.Parse()
	
	istId	:= *IninstanceId
	region	:= *InRegion
	metric	:= *InMetric
	class	:= *InClass

	// AWS KEY
	accKey 	:= ""
	secKey	:= ""

	// Redshift Support Metric
	rdMetring := []string{
		"CommitQueueLength","ConcurrencyScalingActiveClusters","CPUUtilization","DatabaseConnections",
		"HealthStatus","MaintenanceMode","MaxConfiguredConcurrencyScalingClusters","NetworkReceiveThroughput",
		"NetworkTransmitThroughput","PercentageDiskSpaceUsed","ReadIOPS","ReadLatency","ReadThroughput",
		"TotalTableCount","WriteIOPS","WriteLatency","WriteThroughput",
	}

	// Validate Opt
	if len(istId) == 0 || len(accKey) == 0 || len(secKey) == 0 {
		panic("No Request Instance Name Or Key")
	}
	
	// Metric Collect Time
	now 	:= time.Now()
	data 	:= now.UTC().Format(time.RFC3339)
	end, _ 	:= time.Parse(time.RFC3339,data)
	start 	:= end.Add(time.Minute * -5)

	// Create CloudWatch Session
	cw := cloudwatch.New(session.New(), &aws.Config{
        Region: aws.String(region),
        Credentials: credentials.NewStaticCredentials(accKey, secKey, ""),
	})

	// Class Switch
	dClass := strings.ToLower(class)
	var identifierDiv string
	var identifierType string
	switch dClass {
	case "rds":
		identifierDiv = "RDS"
		identifierType = "DBInstanceIdentifier"
	case "redshift" :
		identifierDiv = "Redshift"
		identifierType = "ClusterIdentifier"
		
		// Metric Check
		var metricCheck bool = false
		for i:=0;i<len(rdMetring);i++ {
			if rdMetring[i] == metric {
				metricCheck = true
			}
		}
		if metricCheck == false {
			panic("Not Support Metric")
		}
	default:
		panic("Not Support Class.")
	}

	// Set Name Space
	vNameSpace := fmt.Sprintf("AWS/%s",identifierDiv)

	// Get Metric
	vl, err := cw.GetMetricStatistics(
		&cloudwatch.GetMetricStatisticsInput{
			Namespace: 		aws.String(vNameSpace),
			MetricName: 	aws.String(metric),
			Period: 		aws.Int64(60),
			StartTime:		aws.Time(start),
			EndTime:		aws.Time(end),
			Statistics:     []*string{
				aws.String(cloudwatch.StatisticAverage),
			},
			Dimensions: []*cloudwatch.Dimension{
				{
					Name:  aws.String(identifierType),
					Value: aws.String(istId),
				},
			},
		})
	if err != nil {
		panic(err)
	}
	
	// Last Time Data
	var lstIdx int = 0
	for i:=0; i<len(vl.Datapoints);i++ {
		if i != 0 {
			x := aws.TimeValue(vl.Datapoints[lstIdx].Timestamp)
			y := aws.TimeValue(vl.Datapoints[i].Timestamp)

			if x.After(y) {
				lstIdx = lstIdx
			} else {
				lstIdx = i
			}
		}
	}

	var result float64
	if metric == "FreeableMemory" || metric == "FreeLocalStorage"{
		 memByte := aws.Float64Value(vl.Datapoints[lstIdx].Average)
		 result = memByte / 1024 / 1024 / 1024
	} else {
		result = aws.Float64Value(vl.Datapoints[lstIdx].Average)
	}

	fmt.Println(result)
}