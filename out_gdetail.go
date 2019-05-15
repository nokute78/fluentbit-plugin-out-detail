/*
   Copyright 2019 Takahiro Yamashita

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"C"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
	"github.com/nokute78/msgpack-microscope/pkg/msgpack"
)

func outputVerboseKV(obj *msgpack.MPObject, i uint32, out io.Writer, nest int) {
	spaces := strings.Repeat("    ", nest)

	fmt.Fprintf(out, "%s{\"key\":\n", spaces)
	outputVerboseJSON(obj.Child[i*2], out, nest+1)
	fmt.Fprint(out, ",\n")
	fmt.Fprintf(out, "%s \"value\":\n", spaces)
	outputVerboseJSON(obj.Child[i*2+1], out, nest+1)
	fmt.Fprintf(out, "\n%s}", spaces)
}

func outputVerboseJSON(obj *msgpack.MPObject, out io.Writer, nest int) {
	if obj == nil {
		return
	}
	spaces := strings.Repeat("    ", nest)

	switch {
	case msgpack.IsArray(obj.FirstByte):
		spaces2 := strings.Repeat("    ", nest+1)

		// array header info
		fmt.Fprintf(out, `%s{"format":"%s", "header":"0x%02x", "length":%d, "raw":"0x%0x", "value":`, spaces, obj.FormatName, obj.FirstByte, obj.Length, obj.Raw)

		if int(obj.Length) != len(obj.Child) {
			fmt.Fprintf(os.Stderr, "Error: size mismatch. length is %d, buf %d children.\n", obj.Length, len(obj.Child))
			return
		}

		// array body info
		fmt.Fprintf(out, "\n%s[\n", spaces2)
		if obj.Length > 0 {
			var i uint32
			for i = 0; i < obj.Length-1; i++ {
				outputVerboseJSON(obj.Child[i], out, nest+2)
				fmt.Fprintf(out, ",\n")
			}
			outputVerboseJSON(obj.Child[obj.Length-1], out, nest+2)
		}
		fmt.Fprintf(out, "\n%s]\n%s}\n", spaces2, spaces)
	case msgpack.IsMap(obj.FirstByte):
		spaces2 := strings.Repeat("    ", nest+1)
		// map header info
		fmt.Fprintf(out, `%s{"format":"%s", "header":"0x%02x", "length":%d, "raw":"0x%0x", "value":`, spaces, obj.FormatName, obj.FirstByte, obj.Length, obj.Raw)

		if int(obj.Length*2) != len(obj.Child) {
			fmt.Fprintf(os.Stderr, "Error: size mismatch. length is %d, buf %d(!=length*2) children.\n", obj.Length, len(obj.Child))
			return
		}

		// map body info
		fmt.Fprintf(out, "\n%s[\n", spaces2)
		var i uint32
		if obj.Length > 0 {
			for i = 0; i < obj.Length-1; i++ {
				outputVerboseKV(obj, i, out, nest+2)
				fmt.Fprint(out, ",\n")
			}
			outputVerboseKV(obj, obj.Length-1, out, nest+2)
		}
		fmt.Fprintf(out, "\n%s]\n%s}", spaces2, spaces)

	case msgpack.IsString(obj.FirstByte) || msgpack.IsBin(obj.FirstByte):
		fmt.Fprintf(out, `%s{"format":"%s", "header":"0x%02x", "raw":"0x%0x", "value":"%s"}`, spaces, obj.FormatName, obj.FirstByte, obj.Raw, obj.DataStr)
	case msgpack.IsExt(obj.FirstByte):
		fmt.Fprintf(out, `%s{"format":"%s", "header":"0x%02x", "type":%d, "raw":"0x%0x", "value":"%s"}`, spaces, obj.FormatName, obj.FirstByte, obj.ExtType, obj.Raw, obj.DataStr)
	case msgpack.NilFormat == obj.FirstByte:
		fmt.Fprintf(out, `%s{"format":"%s", "header":"0x%02x", "raw":"0x%0x", "value":null}`, spaces, obj.FormatName, obj.FirstByte, obj.Raw)
	case msgpack.NeverUsedFormat == obj.FirstByte:
		fmt.Fprintf(out, `%s{"format":"%s", "header":"0x%02x", "raw":"0x%0x", "value":%s}`, spaces, obj.FormatName, obj.FirstByte, obj.Raw, obj.DataStr)
		fmt.Fprintf(os.Stderr, "Error: Never Used Format detected\n")
		return
	default:
		fmt.Fprintf(out, `%s{"format":"%s", "header":"0x%02x", "raw":"0x%0x", "value":%s}`, spaces, obj.FormatName, obj.FirstByte, obj.Raw, obj.DataStr)
	}
}

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	return output.FLBPluginRegister(def, "gdetail", "Show MessagePack in detail")
}

//export FLBPluginInit
// (fluentbit will call this)
// plugin (context) pointer to fluentbit context (state/ c code)
func FLBPluginInit(plugin unsafe.Pointer) int {
	msgpack.RegisterFluentdEventTime()

	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	b := C.GoBytes(data, C.int(length))

	buf := bytes.NewBuffer(b)
	out := os.Stdout

	for buf.Len() > 0 {
		ret, err := msgpack.Decode(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error(%s) detected. Incoming data may be broken.\n", err)
			if ret == nil {
				return output.FLB_ERROR
			}
			/* ret is broken, but try to output as much as possible. */
		}
		outputVerboseJSON(ret, out, 0)
	}

	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	return output.FLB_OK
}

func main() {
}
