## จะสร้าง group struture อย่างไร ถ้าไม่มี join

คำถามนี้สำคัญมาก เพราะมันคือ “หัวใจ” ของการทำระบบด้วย Cassandra เลย:

❗ Cassandra ไม่มี JOIN → คุณต้อง “เก็บข้อมูลซ้ำ (denormalize)” ให้ตรงกับ query

ดังนั้นการทำ group / chat room ต้องออกแบบเป็นหลาย table ตามการใช้งานจริง

## แนวคิดก่อนออกแบบ

ในระบบ group chat ปกติคุณจะมี query แบบนี้:

ดูว่า user อยู่ group อะไรบ้าง
ดูสมาชิกใน group
ดูข้อความใน group
ดึงข้อมูล group (name, owner ฯลฯ)

👉 แต่ละ query = 1 table

### โครงสร้างที่ควรมี

1. ข้อมูล group หลัก

```sql
CREATE TABLE groups_by_id (
  group_id uuid PRIMARY KEY,
  group_name text,
  owner_id uuid,
  created_at timestamp
);
```

2. 👥 group → members

```sql
   CREATE TABLE members_by_group (
   group_id uuid,
   user_id uuid,
   joined_at timestamp,
   PRIMARY KEY (group_id, user_id)
   );
```

👉 ใช้เมื่อ:

```sql
SELECT * FROM members_by_group WHERE group_id = ?;
```

3. 👤 user → groups (สำคัญมาก)

```sql
CREATE TABLE groups_by_user (
  user_id uuid,
  updated_at timeuuid,
  group_id uuid,
  group_name text,
  last_message text,
  PRIMARY KEY (user_id, updated_at, group_id)
) WITH CLUSTERING ORDER BY (updated_at DESC, group_id ASC);
```

👉 ใช้เมื่อ:

```sql
SELECT * FROM groups_by_user WHERE user_id = ? LIMIT 20;
```

4. 💬 messages ใน group

```sql
CREATE TABLE messages_by_group (
  group_id uuid,
  message_id timeuuid,
  sender_id uuid,
  content text,
  PRIMARY KEY (group_id, message_id)
) WITH CLUSTERING ORDER BY (message_id DESC);
```

5. 🔎 (optional) message by id \*reply chat

```sql
CREATE TABLE message_by_id (
  message_id timeuuid PRIMARY KEY,
  group_id uuid,
  sender_id uuid,
  content text
);
```

### เวลาเขียนข้อมูล (สำคัญ)

สมมติ user ส่ง message:

คุณต้องเขียนหลาย table พร้อมกัน เช่น:

messages_by_group
message_by_id
update groups_by_user (last_message, updated_at)

#### Cassandra design = write-heavy + duplicate data

Cassandra ไม่มี join

do not

```sql
SELECT *
FROM groups g
JOIN members m ON g.group_id = m.group_id
WHERE user_id = ?;
```

## สรุป mindset

ไม่มี JOIN → ต้อง pre-join ตอนเขียน
1 query → 1 table
ยอม duplicate data เพื่อ speed
