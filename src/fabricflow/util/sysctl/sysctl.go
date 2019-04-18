// -*- coding: utf-8 -*-

package sysctl

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type SysctlEntry struct {
	isComment bool
	isDeleted bool
	key       string
	value     string
	rawValue  string
}

func (e *SysctlEntry) String() string {
	return e.rawValue
}

func (c *SysctlEntry) Key() string {
	return c.key
}

func (c *SysctlEntry) Value() string {
	return c.value
}

func (e *SysctlEntry) SetValue(value string) {
	if !e.isComment {
		e.value = value
		e.rawValue = fmt.Sprintf("%s=%s", e.key, value)
		e.isDeleted = false
	}
}

func NewSysctlEntryComment(text string) *SysctlEntry {
	return &SysctlEntry{
		isComment: true,
		rawValue:  text,
	}
}

func NewSysctlEntryConfig(key, val, raw string) *SysctlEntry {
	if len(raw) == 0 {
		raw = fmt.Sprintf("%s=%s", key, val)
	}

	return &SysctlEntry{
		isComment: false,
		key:       key,
		value:     val,
		rawValue:  raw,
	}
}

func ParseSysctlEntry(s string) *SysctlEntry {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return NewSysctlEntryComment("")
	}

	if s[0] == '#' {
		return NewSysctlEntryComment(s)
	}

	items := strings.SplitN(s, "=", 2)
	if len(items) != 2 {
		return NewSysctlEntryComment(s)
	}

	return NewSysctlEntryConfig(strings.TrimSpace(items[0]), strings.TrimSpace(items[1]), s)
}

type SysctlConfig struct {
	entries []*SysctlEntry
	index   map[string]int
}

func NewSysctlConfig() *SysctlConfig {
	return &SysctlConfig{
		entries: []*SysctlEntry{},
		index:   map[string]int{},
	}
}

func (c *SysctlConfig) mkIndex() {
	for index, entry := range c.entries {
		if entry.isDeleted {
			continue
		}

		if entry.isComment {
			continue
		}

		c.index[entry.key] = index
	}
}

func (c *SysctlConfig) find(key string) *SysctlEntry {
	if index, ok := c.index[key]; ok {
		return c.entries[index]
	}

	return nil

}

func (c *SysctlConfig) AddEntry(e *SysctlEntry) {
	c.entries = append(c.entries, e)
	c.mkIndex()
}

func (c *SysctlConfig) Set(key, val string) {
	e := c.find(key)
	if e == nil {
		c.AddEntry(NewSysctlEntryConfig(key, val, ""))
		return
	}

	e.SetValue(val)
}

func (c *SysctlConfig) Del(key string) {
	if e := c.find(key); e != nil {
		e.isDeleted = true
	}
}

func (c *SysctlConfig) Entry(key string) *SysctlEntry {
	e := c.find(key)
	if e == nil {
		return nil
	}

	if e.isDeleted {
		return nil
	}

	return e
}

func (c *SysctlConfig) Shrink(f func(*SysctlEntry) bool) {
	c.Range(func(e *SysctlEntry) error {
		e.isDeleted = !f(e)
		return nil
	})
}

func (c *SysctlConfig) Range(f func(*SysctlEntry) error) error {
	for _, entry := range c.entries {
		if entry.isDeleted {
			continue
		}

		if err := f(entry); err != nil {
			return err
		}
	}

	return nil
}

func (c *SysctlConfig) WriteTo(w io.Writer) (int64, error) {
	var sum int64
	writer := bufio.NewWriter(w)
	err := c.Range(func(e *SysctlEntry) error {
		n, err := writer.WriteString(e.String())
		if err != nil {
			return err
		}

		sum += int64(n)

		n, err = writer.WriteString("\n")
		if err != nil {
			return err
		}

		sum += int64(n)

		return nil
	})

	if err != nil {
		return sum, err
	}

	return sum, writer.Flush()
}

func ReadConfig(r io.Reader) *SysctlConfig {
	c := NewSysctlConfig()
	s := bufio.NewScanner(r)
	for s.Scan() {
		e := ParseSysctlEntry(s.Text())
		c.AddEntry(e)
	}

	return c
}
