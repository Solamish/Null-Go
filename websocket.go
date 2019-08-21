package nullgo

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"net"
	"net/http"
)

type WebSocketConfig struct {
	OnOpen    func(*WebSocketContext)
	OnMessage func(*WebSocketContext, string)
	OnClose   func(*WebSocketContext)
	OnError   func(*WebSocketContext)
}

type WebSocketContext struct {
	Conn net.Conn
	buf  *bufio.ReadWriter
}


func enter(c *Context) {
	r := c.Request
	w := c.ResponseWriter
	config := c.config
	key := r.Header.Get("Sec-WebSocket-Key")
	s := sha1.New()

	// 将Sec-WebSocket-Key
	// 与258EAFA5-E914-47DA-95CA-C5AB0DC85B11拼接
	// 并使用sha1算法加密
	s.Write(QuickStringToBytes(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	b := s.Sum(nil)

	//将加密后的结果进行base64编码
	//得到Sec-WebSocket-Key-Accept
	secWebSocketAccept := base64.StdEncoding.EncodeToString(b)
	hijack := w.(http.Hijacker)
	con, buf, _ := hijack.Hijack()

	ctx := &WebSocketContext{
		Conn: con,
		buf:  buf,
	}

	//响应客户端的协议升级
	header := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: " + secWebSocketAccept + "\r\n\r\n"
	buf.Write(QuickStringToBytes(header))
	buf.Flush()

	config.OnOpen(ctx)
	go func() {
		for true {
			data := make([]byte, 2)
			_, err := buf.Read(data)
			if err != nil {
				wsClose(con, config.OnClose, ctx)
				break
			}
			bin1 := parseIntToBin(int(data[0]))
			bin2 := parseIntToBin(int(data[1]))

			if bin1[1] || bin1[2] || bin1[3] {
				wsClose(con, config.OnClose, ctx)
				break
			}

			if !bin2[0] {
				wsClose(con, config.OnClose, ctx)
				break
			}

			opcode := parseBinToInt(bin1[4:])
			payloadLen := parseBinToInt(bin2[1:])

			switch opcode {
			case 1:
				maskingKey := make([]byte, 4)
				buf.Read(maskingKey)

				payload := make([]byte, payloadLen)
				buf.Read(payload)

				data := make([]byte, payloadLen)

				//掩码运算
				for i := 0; i < payloadLen; i++ {
					data[i] = payload[i] ^ maskingKey[i%4]
				}

				config.OnMessage(ctx, QuickBytesToString(data))
			default:
				wsClose(con, config.OnClose, ctx)
				break
			}
			}
		}()
	}


func wsClose(conn net.Conn, onclose func(context *WebSocketContext), context *WebSocketContext) {
	onclose(context)
	conn.Close()
}



