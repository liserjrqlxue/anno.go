package main

import (
	"bytes"
	"fmt"
	"github.com/liserjrqlxue/annogo/GnomAD"
	"github.com/liserjrqlxue/simple-util"
	"os"
	"strings"
)

func main() {

	db := "D:\\data\\gnomad-public\\release\\2.1\\vcf\\exomes\\gnomad.exomes.r2.1.sites.vcf.gz"
	//db:="D:\\data\\ftp.ncbi.nlm.nih.gov\\pub\\clinvar\\vcf_GRCh37\\clinvar.vcf.gz"
	//f,_:=os.Open(db)
	//rdr, err := vcfgo.NewReader(f,false)
	//simple_util.CheckErr(err)
	//var rdr= vcfgo.Reader{}

	tbx, err := GnomAD.New(db)
	simple_util.CheckErr(err)
	defer simple_util.DeferClose(tbx)
	//15:98513086-98513088
	vals := tbx.Query("15", 98513086, 98513088)
	hit := tbx.Hit("15", 98513086, 98513088, "T", "C", vals)
	//fmt.Println(tbx.GetHeaderDescription("AC_eas_female"))
	simple_util.CheckErr(err)
	fmt.Println(hit.Ref)
	fmt.Println(hit.Alt)

	for _, k := range GnomAD.FieldInt {
		fmt.Println(k, hit.Info[k])
	}
	for _, k := range GnomAD.FieldFloat32 {
		fmt.Println(k, hit.Info[k])
	}
}

/*
	for {
		v, err := vals.Next()
		if err != nil {
			break
		}
		vt:=v.(interfaces.IVariant)
		fmt.Println("Id:\t",vt.Id())
		fmt.Println("Chrom:\t",vt.Chrom())
		fmt.Println("Start:\t",vt.Start())
		fmt.Println("End:\t",vt.End())
		s,e,b:=vt.CIPos()
		fmt.Println("CIPos:\t",s,e,b)
		s,e,b=vt.CIEnd()
		fmt.Println("CIEnd:\t",s,e,b)
		fmt.Println("Ref:\t",vt.Ref())
		fmt.Println("Alt:\t",vt.Alt())
		af,err:=vt.Info().Get("AF")
		fmt.Println("Info:\t",af,err)

		var fields = makeFields([]byte(v.(interfaces.IVariant).String()))
		var variant = rdr.Parse(fields)
		fmt.Println(variant.Alternate)
		fmt.Println(variant.Alt())
		fmt.Println(variant.Alternate)
		info := strings.Split(v.(interfaces.IVariant).String(), "\t")[7]
		var dataHash = parseInfo(info)
		fmt.Println("AC_female:\t"+dataHash["AC_female"])
		//for k, v := range dataHash {
			//fmt.Println(k + "\t" + v)
		//}
	}
}*/

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
