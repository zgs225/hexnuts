package server

import (
	"testing"
)

func TestConfigurer(t *testing.T) {
	c := Configurer{
		Items: make(map[string]interface{}),
	}

	if err := c.Set("hello.world", "nihao"); err != nil {
		t.Errorf("设置配置错误：%v", err)
	}

	if err := c.Set("hello.world", "123"); err == nil {
		t.Error("设置配置错误，重复的键不能被设置")
	}

	if v, err := c.Get("hello.world"); err != nil {
		t.Errorf("获取配置错误：%v", err)
	} else {
		if v != "nihao" {
			t.Errorf("获取配置错误：\n\t期望：nihao\n\t获得：%s", v)
		}
	}

	if err := c.Update("hello.world", "123"); err != nil {
		t.Error("更新配置错误：", err)
	} else {
		if v, err := c.Get("hello.world"); err != nil {
			t.Errorf("获取配置错误：%v", err)
		} else {
			if v != "123" {
				t.Errorf("获取配置错误：\n\t期望：123\n\t获得：%s", v)
			}
		}
	}

	if _, err := c.Get("not.exists"); err == nil {
		t.Error("获取配置错误：应该返回不存在此配置的错误")
	}

	if err := c.Del("hello"); err == nil {
		t.Error("删除配置错误：父配置不能被删除")
	}

	if err := c.Del("hello.world"); err != nil {
		t.Error("删除配置错误：", err)
	}

	if _, err := c.Get("hello.world"); err == nil {
		t.Error("获取配置错误：已经删除的配置不应该存在")
	}
}
