# fluentbit-plugin-out-detail

[![Build Status](https://travis-ci.org/nokute78/fluentbit-plugin-out-detail.svg?branch=master)](https://travis-ci.org/nokute78/fluentbit-plugin-out-detail)
[![Go Report Card](https://goreportcard.com/badge/github.com/nokute78/fluentbit-plugin-out-detail)](https://goreportcard.com/report/github.com/nokute78/fluentbit-plugin-out-detail)

Ouput plugin for [Fluent-Bit](https://fluentbit.io/) to show MessagePack in detail.

## Build

```
$ make
```

If success, `out_gdetail.so` will be generated.

## Try!

```
$ fluent-bit -e out_gdetail.so -i dummy -o gdetail
```

## Configuration Parameters

None

## Example Output
```
[taka@localhost build]$ ./bin/fluent-bit -e ~/git/oss/fluentbit-plugin-out-detail/out_gdetail.so -i dummy -o gdetail
Fluent Bit v1.2.0
Copyright (C) Treasure Data

[2019/05/15 09:25:23] [ info] [storage] initializing...
[2019/05/15 09:25:23] [ info] [storage] in-memory
[2019/05/15 09:25:23] [ info] [storage] normal synchronization mode, checksum disabled
[2019/05/15 09:25:23] [ info] [engine] started (pid=10486)
{"format":"fixarray", "header":"0x92", "length":2, "raw":"0x92d7005cdb5c74000df3f481a76d657373616765a564756d6d79", "value":
    [
        {"format":"event time", "header":"0xd7", "type":0, "raw":"0xd7005cdb5c74000df3f4", "value":"2019-05-15 09:25:24.00091442 +0900 JST"},
        {"format":"fixmap", "header":"0x81", "length":1, "raw":"0x81a76d657373616765a564756d6d79", "value":
            [
                {"key":
                    {"format":"fixstr", "header":"0xa7", "raw":"0xa76d657373616765", "value":"message"},
                 "value":
                    {"format":"fixstr", "header":"0xa5", "raw":"0xa564756d6d79", "value":"dummy"}
                }
            ]
        }
    ]
}
```

## License

[Apache License v2.0](https://www.apache.org/licenses/LICENSE-2.0)