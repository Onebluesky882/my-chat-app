## player active

— เก็บใน Redis แทน (เร็วกว่า)
presence data เปลี่ยนบ่อย เหมาะกับ Redis มากกว่า ScyllaDB:
go// online
s.redis.Set(ctx, "presence:"+userID, "online", 5\*time.Minute)

// offline
s.redis.Del(ctx, "presence:"+userID)

// เช็ค
val, err := s.redis.Get(ctx, "presence:"+userID).Result()
isOnline := err == nil && val == "online"
