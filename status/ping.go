package status

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

var (
	emptySample  = make([]Player, 0)
	emptyModList = make([]Mod, 0)
)

func GetStatus(server string, port int) Status {
	var status Status

	data, err := getServerData(server, port)
	if err != nil {
		status.Type = "error"
		status.Error = err
	} else {
		status.Type = "success"
		status.Data = data
	}

	return status
}

func getServerData(server string, port int) (*Data, error) {
	con, err := net.Dial("tcp", fmt.Sprintf("%s:%v", server, port))
	if err != nil {
		return nil, err
	}

	defer con.Close()
	deadline := time.Now().Add(5 * time.Second)
	con.SetDeadline(deadline)
	con.SetReadDeadline(deadline)
	con.SetWriteDeadline(deadline)

	// handshake
	_, err = con.Write(handshake(server, port).Bytes())
	if err != nil {
		return nil, err
	}

	// requestStatus request
	_, err = con.Write(requestStatus().Bytes())
	if err != nil {
		return nil, err
	}

	// read length
	r := bufio.NewReader(con)
	l, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}

	// read data
	data := make([]byte, l)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}

	// trim to json string
	_, i0 := binary.Uvarint(data)
	_, i1 := binary.Uvarint(data[i0:])

	return parseData(data[i0+i1:])
}

func parseData(b []byte) (*Data, error) {
	var s Data
	err := json.Unmarshal(b, &s)
	tidyData(s)
	return &s, err
}

// handshake payload
func handshake(server string, port int) *bytes.Buffer {
	var buf bytes.Buffer
	buf.WriteByte(0x00)                               // id
	buf.WriteByte(0x47)                               // proto
	buf.Write(varInt(len(server)))                    // length
	buf.WriteString(server)                           // ip
	binary.Write(&buf, binary.BigEndian, int16(port)) // port
	buf.WriteByte(0x01)
	return wrap(&buf)
}

// requestStatus request payload
func requestStatus() *bytes.Buffer {
	var buf bytes.Buffer
	buf.WriteByte(0x00) // id
	return wrap(&buf)
}

// wraps the buffer into a packet
func wrap(b *bytes.Buffer) *bytes.Buffer {
	var buf bytes.Buffer
	buf.Write(varInt(len(b.Bytes())))
	buf.Write(b.Bytes())
	return &buf
}

// encode int
func varInt(i int) []byte {
	buf := make([]byte, 10)
	l := binary.PutUvarint(buf, uint64(i))
	return buf[:l]
}

func tidyData(s Data) {
	if s.Players.Sample == nil {
		s.Players.Sample = emptySample
	}
	if s.ModInfo != nil && s.ModInfo.ModList == nil {
		s.ModInfo.ModList = emptyModList
	}
}
