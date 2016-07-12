/**
 * Copyright 2016 Confluent Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"fmt"
	"os"
	"time"
)

/*
#include <librdkafka/rdkafka.h>
#cgo LDFLAGS: -lrdkafka

static const char *errdesc_to_string (const struct rd_kafka_err_desc *ed, int idx) {
   return ed[idx].name;
}

*/
import "C"

func main() {

	outfile := os.Args[1]

	f, err := os.Create(outfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("// Copyright 2016 Confluent Inc.\n")
	f.WriteString(fmt.Sprintf("// AUTOMATICALLY GENERATED BY %s ON %v USING librdkafka %s\n",
		os.Args[0], time.Now(), C.GoString(C.rd_kafka_version_str())))
	f.WriteString("package kafka\n")

	var errdescs *C.struct_rd_kafka_err_desc
	var csize C.size_t
	C.rd_kafka_get_err_descs(&errdescs, &csize)

	f.WriteString(`
/*
#include <librdkafka/rdkafka.h>
*/
import "C"

type KafkaErrorCode int

func (c KafkaErrorCode) String() string {
      return C.GoString(C.rd_kafka_err2str(C.rd_kafka_resp_err_t(c)))
}

const (
`)

	for i := 0; i < int(csize); i += 1 {
		errname := C.GoString(C.errdesc_to_string(errdescs, C.int(i)))
		if len(errname) == 0 {
			continue
		}
		f.WriteString(fmt.Sprintf("    ERR_%s KafkaErrorCode = KafkaErrorCode(C.RD_KAFKA_RESP_ERR_%s)\n",
			errname, errname))
	}

	f.WriteString(")\n")

}
