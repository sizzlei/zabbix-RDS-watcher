# zabbix-RDS-watcher
AWS-GO-SDK를 활용하여 Aurora RDS에 대한 CloudWatch Metric을 Zabbix로 불러올수 있는 프로그램 입니다. 

## Import
```
go get github.com/aws/aws-sdk-go/aws
```

## Configure
```
// AWS KEY
accKey 	:= ""
secKey	:= ""
```
Flag로 Key를 처리하는 경우 Zabbix 서버 프로세스상 Key가 노출되기 때문에 코드 내부에 작성.


## Complie
```
GOOS=linux go build main.go
```

## View Metric
해당 프로그램은 Aurora for MySQL에 대해 아래의 Metric을 수집
+ FreeableMemory
+ FreeLocalStorage
+ CPUUtilization
+ CommitThroughput
+ DDLThroughput
+ DMLThroughput
+ InsertThroughput
+ SelectThroughput
+ DeleteThroughput
+ UpdateThroughput
+ ActiveTransactions
+ DatabaseConnections
+ Deadlocks
+ AbortedClients
+ Queries
+ RowLockTime
+ EngineUptime
+ InsertLatency
+ SelectLatency
+ UpdateLatency
+ DeleteLatency
+ DDLLatency
+ DMLLatency
+ CommitLatency
+ AuroraReplicaLag
+ ForwardingMasterDMLLatency
+ ForwardingReplicaDMLLatency
+ ForwardingReplicaReadWaitLatency
+ ForwardingReplicaSelectLatency
+ ForwardingMasterDMLThroughput
+ ForwardingReplicaDMLThroughput
+ ForwardingReplicaReadWaitThroughput
+ ForwardingReplicaReadWaitThroughput

## Zabbix Usage
Zabbix에 Item 등록시 Key 부분에 아래와 같이 작성.
![ZabbixItem](./img/Zabbix-item.png)
```
zabbixRDSWatcher["-metric","{Metric Name}","-instance","{HOST.HOST}"]
```
+ {Metric Name} : 위에서 수집하려는 Metric Name
+ {HOST.HOST}   : Zabbix Host에 추가된 서버의 Name
> 해당 프로그램은 Zabbix Host에 RDS 등록 시  Host name 부분에 DBInstanceIndentier를 넣어야 합니다.
![ZabbixHostAdd](./img/Zabbix-host-add.png)