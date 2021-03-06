package nullgo

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

//定义一些可能会用到的变量
var (
	Comment = []byte{'#'}
	Empty   = []byte{}
	Equal   = []byte{'='}
	Quote  = []byte{'"'}
)

type Config struct {
	filename string


	//[]{comment, key...}
	//用来存放注释和参数的key...
	comment  map[int][]string

	// key:value
	// 用来存放参数键值对的map
	data     map[string]string

	//key:offset
	//map的key就是参数的key，value是偏移量
	offset   map[string]int64
	sync.RWMutex
}

func LoadConfig(name string) (*Config, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		filename: file.Name(),
		comment:  make(map[int][]string),
		data:     make(map[string]string),
		offset:   make(map[string]int64),
		RWMutex:  sync.RWMutex{},
	}
	cfg.Lock()
	defer cfg.Unlock()
	defer file.Close()

	var comment bytes.Buffer
	buf := bufio.NewReader(file)

	//一行一行的读取数据
	for nComment, off := 0, int64(1); ; {
		line, _ ,err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if bytes.Equal(line, Empty) {
			continue
		}
		off += int64(len(line))

		if bytes.HasPrefix(line, Comment) {
			line = bytes.TrimLeft(line,"#")
			line = bytes.TrimLeftFunc(line, unicode.IsSpace)
			comment.Write(line)
			comment.WriteByte('\n')
			continue
		}
		//将注释写入comment
		if comment.Len() !=0 {
			cfg.comment[nComment] = []string{comment.String()}
			comment.Reset()
			nComment++
		}

		//取出键值对
		value := bytes.SplitN(line, Equal, 2)
		if bytes.HasPrefix(value[1], Quote) {
			value[1] = bytes.Trim(value[1], `"`)
		}

		key := strings.TrimSpace(string(value[0]))
		cfg.comment[nComment-1] = append(cfg.comment[nComment-1], key)
		cfg.data[key] = strings.TrimSpace(string(value[1]))
		cfg.offset[key] = off
	}


	return cfg, nil
}

//将参数以布尔值返回
func (c *Config) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.data[key])
}

//将参数以int值返回
func (c *Config) Int(key string) (int, error) {
	return strconv.Atoi(c.data[key])
}

//将参数以浮点数返回
func (c *Config) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.data[key], 64)
}

//将参数以string值返回
func (c *Config) String(key string) string {
	return c.data[key]
}

