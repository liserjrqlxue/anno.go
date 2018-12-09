package main

import (
	"bytes"
	"fmt"
	"github.com/brentp/bix"
	"github.com/brentp/irelate/interfaces"
	"github.com/brentp/vcfgo"
	"github.com/liserjrqlxue/simple-util"
	"os"
	"strings"
)

func main() {
	db := "clinvar.vcf.gz"
	//f,_:=os.Open(db)
	//rdr, err := vcfgo.NewReader(f,false)
	//simple_util.CheckErr(err)
	var rdr = vcfgo.Reader{}

	tbx, err := bix.New(db)
	simple_util.CheckErr(err)
	defer simple_util.DeferClose(tbx)
	vals, err := tbx.Query(interfaces.AsIPosition("1", 976124, 976124))
	simple_util.CheckErr(err)
	for {
		v, err := vals.Next()
		if err != nil {
			break
		}
		var fields = makeFields([]byte(v.(interfaces.IVariant).String()))
		var variant = rdr.Parse(fields)
		fmt.Println(variant.Alternate)
		fmt.Println(variant.Alt())
		fmt.Println(variant.Alternate)
		info := strings.Split(v.(interfaces.IVariant).String(), "\t")[7]
		var dataHash = parseInfo(info)
		for k, v := range dataHash {
			fmt.Println(k + "\t" + v)
		}
	}
}

func parseInfo(info string) map[string]string {
	var Info = make(map[string]string)
	for _, str := range strings.Split(info, ";") {
		kv := strings.SplitN(str, "=", 2)
		if len(kv) == 1 {
			Info[kv[0]] = "TRUE"
		} else if len(kv) == 2 {
			Info[kv[0]] = kv[1]
		}
	}
	return Info
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

// from github.com/brentp/vcfgo/reader.go
func makeFields(line []byte) [][]byte {
	fields := bytes.SplitN(line, []byte{'\t'}, 9)
	s := 0
	for i, f := range fields {
		if i == 7 {
			break
		}
		s += len(f) + 1
	}
	if s >= len(line) {
		fmt.Fprintf(os.Stderr, "XXXXX: bad VCF line '%s'", line)
		return fields
	}
	e := bytes.IndexByte(line[s:], '\t')
	if e == -1 {
		e = len(line)
	} else {
		e += s
	}

	fields[7] = line[s:e]
	return fields
}
