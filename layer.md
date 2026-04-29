Client/UI
↓
HTTP / WebSocket Layer
↓
Chat Service (business logic)
↓
Redis + Scylla (data layer)

```text
User กด send
    |
POST /chat/send
|
Chat Service.Send()
|
+--> Scylla
    | เก็บ message history ถาวร
    |
+--> Redis
    | cache latest messages
    | unread counter +1
    |
+--> WebSocket push
    realtime แจ้ง client
```
