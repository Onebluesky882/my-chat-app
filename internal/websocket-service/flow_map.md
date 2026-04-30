แผนภาพแบ่งเป็น 3 ช่วงหลัก:

1. เชื่อมต่อ — Client ส่ง GET /ws พร้อม query params → Fiber Router ดึง room_id/user_id ก่อน upgrade → FastHTTP Upgrader แปลง HTTP เป็น WebSocket → ส่ง conn เข้า HandleWS
2. wsSvc.Connect — ตรวจสิทธิ์ด้วย IsMember ก่อน ถ้าไม่ผ่านปิด connection ทันที → ลงทะเบียน userConns[userID] และ join(roomID)
3. readLoop — วนรับ message ไม่หยุด → parse JSON → แยก type: "message" บันทึกและ broadcast, "typing" broadcast typing event, "leave" ออกจาก loop → defer cleanup ทำงาน (leave room, ลบ userConn, ปิด conn)

[WebSocket Flow](../../assets/websocket-flow.png)
