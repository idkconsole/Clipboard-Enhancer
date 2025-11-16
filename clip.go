package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"syscall"
	"unsafe"

	"github.com/JamesHovious/w32"
	"github.com/atotto/clipboard"
)

var key = ""
var mdl = "gemini-2.5-flash-preview-05-20"

var pf = `
rewrite the text simply and naturally in human english.
fix grammar and clarity only.
keep the meaning the same.
do not add anything.
do not explain.
do not give options.
do not expand.
output only the rewritten text.
`

var pt = `
translate this text into natural simple english.
output only the translation.
`

var ps = `
answer this question in simple human english.
do not add anything extra.
`

var pd = `
explain this word or phrase in simple human english.
do not add anything extra.
`

func gen(x string) string {
	u := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", mdl, key)
	b := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": x}}},
		},
	}
	j, _ := json.Marshal(b)
	r, _ := http.NewRequest("POST", u, bytes.NewBuffer(j))
	r.Header.Set("content-type", "application/json")
	c := &http.Client{}
	res, err := c.Do(r)
	if err != nil {
		return "err"
	}
	defer res.Body.Close()
	d, _ := ioutil.ReadAll(res.Body)
	var o map[string]interface{}
	json.Unmarshal(d, &o)
	out := ""
	a, ok := o["candidates"].([]interface{})
	if ok && len(a) > 0 {
		x0 := a[0].(map[string]interface{})
		ct := x0["content"].(map[string]interface{})
		pp := ct["parts"].([]interface{})
		out = pp[0].(map[string]interface{})["text"].(string)
	}
	return strings.TrimSpace(out)
}

func fix(x string) string  { return gen(pf + "\n\n" + x) }
func trn(x string) string  { return gen(pt + "\n\n" + x) }
func ans(x string) string  { return gen(ps + "\n\n" + x) }
func defn(x string) string { return gen(pd + "\n\n" + x) }

func dofix() {
	t, _ := clipboard.ReadAll()
	t = strings.TrimSpace(t)
	if t == "" {
		fmt.Println("\nempty\n")
		return
	}
	o := fix(t)
	clipboard.WriteAll(o)
	fmt.Println("\nrewritten - {" + t + "} -> {" + o + "}\n")
}

func dotrn() {
	t, _ := clipboard.ReadAll()
	t = strings.TrimSpace(t)
	if t == "" {
		fmt.Println("\nempty\n")
		return
	}
	o := trn(t)
	clipboard.WriteAll(o)
	fmt.Println("\ntranslated - {" + t + "} -> {" + o + "}\n")
}

func doans() {
	t, _ := clipboard.ReadAll()
	t = strings.TrimSpace(t)
	if t == "" {
		fmt.Println("\nempty\n")
		return
	}
	o := ans(t)
	clipboard.WriteAll(o)
	fmt.Println("\nanswered - {" + t + "} -> {" + o + "}\n")
}

func dodef() {
	t, _ := clipboard.ReadAll()
	t = strings.TrimSpace(t)
	if t == "" {
		fmt.Println("\nempty\n")
		return
	}
	o := defn(t)
	clipboard.WriteAll(o)
	fmt.Println("\ndefined - {" + t + "} -> {" + o + "}\n")
}

var u = syscall.NewLazyDLL("user32.dll")
var hook = u.NewProc("SetWindowsHookExW")
var call = u.NewProc("CallNextHookEx")
var msgfn = u.NewProc("GetMessageW")
var used = false
var pressed = map[uint32]bool{}

const wh = 13
const kd = 0x0100
const ku = 0x0101

func kb(cb int, wp uintptr, lp uintptr) uintptr {
	k := (*w32.KBDLLHOOKSTRUCT)(unsafe.Pointer(lp))
	v := uint32(k.VkCode)
	if wp == kd {
		pressed[v] = true
		check()
	}
	if wp == ku {
		delete(pressed, v)
	}
	r, _, _ := call.Call(0, uintptr(cb), wp, lp)
	return r
}

func check() {
	if pressed[46] && (pressed[35] || pressed[34] || pressed[33] || pressed[36]) {
		if used {
			return
		}
		used = true
		if pressed[35] {
			go dofix()
		}
		if pressed[34] {
			go dotrn()
		}
		if pressed[33] {
			go doans()
		}
		if pressed[36] {
			go dodef()
		}
	}
	if !pressed[46] {
		used = false
	}
}

func main() {
	h, _, _ := hook.Call(uintptr(wh), syscall.NewCallback(kb), 0, 0)
	var m w32.MSG
	msgfn.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
	_ = h
}
