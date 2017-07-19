package client

import (
	"testing"
)

func TestGet(t *testing.T) {
	client := &HTTPClient{Addr: "http://localhost:5678"}
	k := "hello.world"
	v, err := client.Get(k)
	if err != nil {
		t.Error("Get错误：", err)
	}
	if v != "1" {
		t.Errorf("Get错误：期望 1，得到 %s", v)
	}
}

func TestSet(t *testing.T) {
	client := &HTTPClient{Addr: "http://localhost:5678"}
	k := "wo.ai.ni"
	v := "1"
	if err := client.Set(k, v); err != nil {
		t.Error("Set错误：", err)
		return
	}

	v2, err := client.Get(k)
	if err != nil {
		t.Error("Get error after set: ", err)
		return
	}

	if v2 != v {
		t.Errorf("Get error after set: want %s, got %s", v, v2)
		return
	}

	if err := client.Del(k); err != nil {
		t.Error("Del error after set: ", err)
		return
	}
}
