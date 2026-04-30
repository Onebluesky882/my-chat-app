1-1:
POST /rooms/direct { members: [user-1, user-2] }
│
▼
สร้าง room + เพิ่ม member ทั้ง 2 ในครั้งเดียว
│
▼
ได้ room_id กลับมา → connect WebSocket

Group:
POST /rooms/group { members: [user-1, user-2, user-3] }
│
▼
สร้าง room + เพิ่ม member ทุกคน
│
▼
ได้ room_id กลับมา → connect WebSocket
