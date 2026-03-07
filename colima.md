Colima
├─ ScyllaDB
├─ Redis
└─ Backend Go

scylladb

```sql

เข้าใปใน scylla-node1 (main node)
docker exec -it scylla-node1 cqlsh


DROP KEYSPACE chat;

CREATE KEYSPACE chat
WITH replication = {
'class': 'NetworkTopologyStrategy',
'datacenter1': 2
}
AND TABLETS = {'enabled': false};

--check
DESCRIBE KEYSPACE chat;
```

// ตรวจสอบ  ข้างนอก 
docker exec -it scylla-node1 nodetool status
