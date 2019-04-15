![Ack](https://raw.githubusercontent.com/txn2/ack/master/mast.jpg)

**Ack**: A TXN2 common API acknowledgement (ACK) struct.

## Example Implementation


## Run Example

```bash
AGENT=test SERVICE_ENV=dev SERVICE_NS=example go run ./example/server.go --test
```

## Example Ack Payload

```json
{
    "ack_version": 8,
    "agent": "test",
    "srv_env": "dev",
    "srv_ns": "example",
    "ack_uuid": "65143cad-8b06-4388-b618-ce3a364b7136",
    "req_uuid": "",
    "date_time": "2019-04-11T19:09:32-07:00",
    "success": true,
    "error_code": "",
    "error_message": "",
    "server_code": 200,
    "location": "/test",
    "payload_type": "Message",
    "payload": "A test message.",
    "duration": "514ns"
}
```

## Example Ack Headers

```plain
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< X-Ack-Agent: test
< X-Ack-Duration: 487ns
< X-Ack-Payload-Type: Message
< X-Ack-Req-Uuid: 2031e500-20d6-454a-a641-dc409334e8f3
< X-Ack-Srv-Env: dev
< X-Ack-Srv-Ns: example
< X-Ack-Uuid: 21fe6c71-6148-4ade-b64c-d9ccd71f8325
< X-Ack-Version: 8
< Date: Fri, 12 Apr 2019 02:58:52 GMT
< Content-Length: 366
```