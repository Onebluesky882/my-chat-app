## setup docker

❯ docker-compose up -d
[+] up 4/4
✔ Network my-chat-app_scylla-net Created 0.0s
✔ Container scylla-node1 Started 0.1s
✔ Container scylla-node3 Started 0.2s
✔ Container scylla-node2 Started

## step

checking : เช็ค cluster (สำคัญสุด)
docker exec -it scylla-node1 nodetool status

1. สร้าง keyspace (ให้ cluster ใช้จริง)

content db with app TablePlus
db : Cassandra

```sql
CREATE KEYSPACE chat_app
WITH replication = { 'class': 'SimpleStrategy', 'replication_factor': 3};
```

2. เลือกใช้ database USE chat_app;

3. สร้าง table

```sql
CREATE TABLE messages (
room_id text,
message_id timeuuid,
user_id text,
content text,
created_at timestamp,
PRIMARY KEY (room_id, message_id  )
) WITH CLUSTERING ORDER BY (message_id DESC);
```

4. insert ข้อมูล

```sql
INSERT INTO messages (room_id, message_id, user_id, content, created_at)
VALUES ('room1', now(), 'user1', 'hello world', toTimestamp(now()));
```

5. query ดูข้อมูล

```sql
SELECT * FROM messages WHERE room_id = 'room1';
```

## คำสั่งพื้นฐานที่ต้องรู้

```sql
ดู keyspace ทั้งหมด
DESCRIBE KEYSPACES;
ดู tables ใน keyspace ปัจจุบัน
DESCRIBE TABLES;
ดูโครงสร้าง table
DESCRIBE TABLE messages;
```

ใน Cassandra คำว่า keyspace คือระดับโครงสร้างที่ใหญ่ที่สุดของฐานข้อมูล — เปรียบง่าย ๆ ได้ว่าเหมือน “database” ในระบบอื่น ๆ เช่น MySQL

Keyspace = container ของ tables + การตั้งค่า replication มันไม่ได้แค่เก็บ table แต่ยังกำหนดว่า:

ข้อมูลจะถูก copy ไปกี่เครื่อง (replication factor)
ใช้กลยุทธ์การกระจายข้อมูลแบบไหน

ตัวอย่าง

```sql
CREATE KEYSPACE chat_app
WITH replication = {
  'class': 'SimpleStrategy',
  'replication_factor': 3
};
```

### จุดสำคัญที่ต้องเข้าใจ

ความหมาย: chat_app = ชื่อ keyspace
replication_factor: 3 = ข้อมูลถูกเก็บ 3 โหนด (กันพัง)

Keyspace ไม่ใช่แค่ “folder”
มันควบคุม data distribution + fault tolerance
เป็นหัวใจของ scalability ใน Cassandra

## ⚠️ ข้อสำคัญ

ห้าม query แบบนี้

```sql
SELECT * FROM messages;
```

## 🔥 หลักคิดก่อนออกแบบ (สำคัญมาก)

ScyllaDB ≠ SQL
มันเป็น query-first database
ต้องออกแบบ “ตาม query ที่จะใช้” ไม่ใช่ตาม structure

### Step ต่อไปที่คุณควรทำ (จริงจัง)

1. Define use-case ก่อน (chat app ของคุณ)

ตัวอย่าง:
• user ส่ง message
• เปิดห้อง chat
• โหลด message ล่าสุด
• โหลด message ย้อนหลัง
• online / last seen

2. แตกเป็น Query จริง

Q1: get messages in room ล่าสุด 50 ข้อความ
Q2: load older messages (pagination)
Q3: insert message
Q4: get user chats list

3. ออกแบบ table (สำคัญสุด)

❌ ห้ามคิดแบบ SQL

✅ ต้องคิดแบบ “1 table ต่อ 1 query”
