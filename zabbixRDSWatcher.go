package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"flag"
	"time"
	"fmt"
)

/*
Metrics: Aurora 
@@ System
*FreeableMemory
*FreeLocalStorage
*CPUUtilization


@@ Network
CommitThroughput
DDLThroughput
DMLThroughput

*InsertThroughput
*SelectThroughput
*DeleteThroughput
*UpdateThroughput

@@ Connection and Database Info
*ActiveTransactions
*DatabaseConnections
*Deadlocks
*AbortedClients
*Queries
RowLockTime
EngineUptime

@@ Latency 
*InsertLatency
*SelectLatency
*UpdateLatency
*DeleteLatency
*DDLLatency
*DMLLatency
*CommitLatency

@@ Replica
*AuroraReplicaLag
ForwardingMasterDMLLatency
ForwardingReplicaDMLLatency
ForwardingReplicaReadWaitLatency
ForwardingReplicaSelectLatency


ForwardingMasterDMLThroughput
ForwardingReplicaDMLThroughput
ForwardingReplicaReadWaitThroughput
ForwardingReplicaReadWaitThroughput

*/

func main() {
	IninstanceId		:= flag.String("instance","","DBInstanceIdentifier")
	InRegion			:= flag.String("region","ap-northeast-2","AWS Region") 
	InMetric			:= flag.String("metric","","AWS CloudWatch Database Metric")
	flag.Parse()
	
	istId	:= *IninstanceId
	region	:= *InRegion
	metric	:= *InMetric

	// AWS KEY
	accKey 	:= ""
	secKey	:= ""

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

	// Get Metric
	vl, err := cw.GetMetricStatistics(
		&cloudwatch.GetMetricStatisticsInput{
			Namespace: 		aws.String("AWS/RDS"),
			MetricName: 	aws.String(metric),
			Period: 		aws.Int64(60),
			StartTime:		aws.Time(start),
			EndTime:		aws.Time(end),
			Statistics:     []*string{
				aws.String(cloudwatch.StatisticAverage),
			},
			Dimensions: []*cloudwatch.Dimension{
				{
					Name:  aws.String("DBInstanceIdentifier"),
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
	fmt.Println(vl.Datapoints)
	var result float64
	if metric == "FreeableMemory" || metric == "FreeLocalStorage"{
		 memByte := aws.Float64Value(vl.Datapoints[lstIdx].Average)
		 result = memByte / 1024 / 1024 / 1024
	} else {
		result = aws.Float64Value(vl.Datapoints[lstIdx].Average)
	}
  
	fmt.Println(vl.Datapoints)
	fmt.Println(result)
}
