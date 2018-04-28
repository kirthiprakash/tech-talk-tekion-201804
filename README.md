# tech-talk-tekion-201804

**Summary of the talk**

The talk discusses database querying time for aggregation operations on a table with one million rows and 10 million
rows. (Avg. row size is ~100B)

We talked about running simple aggregation queries on an appointments table. An
appointment consists of appointment time, appointment type, appointment source etc.

```
mysql> desc appointments;
+--------------------+------------------+------+-----+---------+----------------+
| Field              | Type             | Null | Key | Default | Extra          |
+--------------------+------------------+------+-----+---------+----------------+
| id                 | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
| tenant_id          | int(11)          | YES  |     | NULL    |                |
| dealer_id          | int(11)          | YES  |     | NULL    |                |
| appointment_time   | timestamp        | YES  |     | NULL    |                |
| appointment_type   | varchar(255)     | YES  |     | NULL    |                |
| appointment_source | varchar(255)     | YES  |     | NULL    |                |
+--------------------+------------------+------+-----+---------+----------------+
```

**Running the script**

The script provided will help produce random appointments between year 2015 to 2018 and insert them into MysqlDB and MongoDB.
The database connection configuration can be found at the "/app/config.go" file. Change the values accordingly.
For MysqlDB, make sure to create the schema before running the scripts. The tables will be auto migrated.

```
$ go run main.go 
How may rows should I insert?: 10
Total execution time:  0.00187599835 mins
```
This will result in inserting 10 appointments in MysqlDB and MongoDB.

```
mysql> select * from appointments limit 5;
+-------+-----------+-----------+---------------------+------------------+--------------------+
| id    | tenant_id | dealer_id | appointment_time    | appointment_type | appointment_source |
+-------+-----------+-----------+---------------------+------------------+--------------------+
| 10642 |         2 |         1 | 2015-03-11 22:27:22 | PDI              | PHONE              |
| 10643 |         2 |         1 | 2015-04-09 03:47:44 | PDI              | WEB                |
| 10644 |         2 |         1 | 2016-03-03 09:38:30 | SCHEDULED        | PHONE              |
| 10645 |         2 |         1 | 2015-04-08 16:01:53 | WALKIN           | WEB                |
| 10646 |         2 |         1 | 2017-03-03 07:10:30 | PDI              | PHONE              |
+-------+-----------+-----------+---------------------+------------------+--------------------+
5 rows in set (0.00 sec)
```

```
> db.appointments.findOne()
{
	"_id" : ObjectId("5adf75e99f7cfbb4ff12c723"),
	"tenantID" : 2,
	"dealerID" : 1,
	"appointmentTime" : ISODate("2016-01-04T14:43:58Z"),
	"appointmentType" : "WALKIN",
	"appointmentSource" : "WEB"
}

```

The script takes approx. 30 mins to insert 1 million appointments and It takes approximately 5 hrs to insert 10 million appointments.
For the purpose of the talk, the appointments were pre-populated with 1 million appointment in the database name 'analyticsOne'.
Database name 'analyticsTen' had ten million appointments

*Execution Times*

```
mysql> use analyticsOne
mysql> select count(*) from appointments;
+----------+
| count(*) |
+----------+
|  1000008 |
+----------+
mysql> select appointment_type, count(*) from appointments group by appointment_type;
+------------------+----------+
| appointment_type | count(*) |
+------------------+----------+
| PDI              |   333316 |
| SCHEDULED        |   333490 |
| WALKIN           |   333202 |
+------------------+----------+
3 rows in set (0.52 sec)
```

```
mysql> user analyticsTen
mysql> select count(*) from appointments;
+----------+
| count(*) |
+----------+
| 10000008 |
+----------+
mysql> select appointment_type, count(*) from appointments group by appointment_type;
+------------------+----------+
| appointment_type | count(*) |
+------------------+----------+
| PDI              |  3334321 |
| SCHEDULED        |  3332551 |
| WALKIN           |  3333136 |
+------------------+----------+
3 rows in set (5.36 sec)
```

```
> use analyticsOne
> db.appointments.count()
1000008
> var before = new Date()
>  db.appointments.aggregate([{$group:{_id:"$appointmentType", count:{$sum:1}}}])
{ "_id" : "SCHEDULED", "count" : 333490 }
{ "_id" : "WALKIN", "count" : 333202 }
{ "_id" : "PDI", "count" : 333316 }
> var after = new Date()
> execution_mills = after - before
1020 
```

```
> use analyticsTen
> db.appointments.count()
10000008
> var before = new Date()
>  db.appointments.aggregate([{$group:{_id:"$appointmentType", count:{$sum:1}}}])
{ "_id" : "SCHEDULED", "count" : 3332551 }
{ "_id" : "PDI", "count" : 3334321 }
{ "_id" : "WALKIN", "count" : 3333136 }
> var after = new Date()
> execution_mills = after - before
30885 
//1st time hit reads from disk
```

```
> var before = new Date()
>  db.appointments.aggregate([{$group:{_id:"$appointmentType", count:{$sum:1}}}])
var after = new Date()
execution_mills = after - before
{ "_id" : "SCHEDULED", "count" : 3332551 }
{ "_id" : "PDI", "count" : 3334321 }
{ "_id" : "WALKIN", "count" : 3333136 }
> var after = new Date()
> execution_mills = after - before
10127 
//2nd time, the execution time reduced by 1/3rd as the data was cached in memory
```

**Table size**

```
mysql> use analyticsOne
mysql> SELECT   table_name AS 'Table',   ROUND((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024) AS 'Size (MB)' FROM   information_schema.TABLES WHERE   TABLE_SCHEMA = "analyticsOne" ORDER BY   (DATA_LENGTH + INDEX_LENGTH) DESC;
+------------------+-----------+
| Table            | Size (MB) |
+------------------+-----------+
| appointments     |        50 |
+------------------+-----------+
```
```
mysql use analyticsTen
mysql> SELECT   table_name AS 'Table',   ROUND((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024) AS 'Size (MB)' FROM   information_schema.TABLES WHERE   TABLE_SCHEMA = "analyticsTen" ORDER BY   (DATA_LENGTH + INDEX_LENGTH) DESC;
+--------------+-----------+
| Table        | Size (MB) |
+--------------+-----------+
| appointments |       468 |
+--------------+-----------+
```

```
> use analyticsOne
> db.appointments.totalSize() /1024 /1024
47.73046875 (MB)
```

```
> use analyticsTen
switched to db analyticsTen
> db.appointments.totalSize() /1024 /1024
479.1015625 (MB)
```

