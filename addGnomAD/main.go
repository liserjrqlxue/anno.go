package main

import (
	"flag"
	"fmt"
	"github.com/liserjrqlxue/annogo/GnomAD"
	"github.com/liserjrqlxue/simple-util"
	"os"
	"strings"
)

var (
	gnomAD = flag.String(
		"gnomAD",
		"",
		"GnomAD vcf.gz file path",
	)
	input = flag.String(
		"input",
		"",
		"input tsv",
	)
	output = flag.String(
		"output",
		"",
		"output file path",
	)
)

func main() {
	flag.Parse()
	if *gnomAD == "" || *input == "" || *output == "" {
		flag.Usage()
		os.Exit(1)
	}
	tbx, err := GnomAD.New(*gnomAD)
	simple_util.CheckErr(err)
	defer simple_util.DeferClose(tbx)

	outF, err := os.Create(*output)
	simple_util.CheckErr(err)
	defer simple_util.DeferClose(outF)

	mapArray, title := simple_util.File2MapArray(*input, "\t", nil)
	title = append(title, "GnomAD AF", "GnomAD EAS AF", "GnomAD HomoAlt Count", "GnomAD EAS HomoAlt Count")
	_, err = fmt.Fprintln(outF, strings.Join(title, "\t"))
	simple_util.CheckErr(err)
	for _, item := range mapArray {
		GnomAD.AddGnomAD(tbx, item)
		var array []string
		for _, key := range title {
			array = append(array, item[key])
		}
		_, err = fmt.Fprintln(outF, strings.Join(array, "\t"))
		simple_util.CheckErr(err)
	}
}
