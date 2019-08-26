package GnomAD

import (
	"github.com/brentp/bix"
	"github.com/brentp/irelate/interfaces"
	"github.com/liserjrqlxue/simple-util"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Tbx struct {
	Bix *bix.Bix
}

type Variant struct {
	Chrom      string
	Start, End uint32
	Ref        string
	Alt        []string
	Info       map[string]interface{}
}

var FieldFloat32 = []string{
	"AF",
	"AF_eas",
}

var FieldInt = []string{
	"AC",
	"AN",
	"nhomalt",
	"AN_eas",
	"AC_eas",
	"nhomalt_eas",
}

func New(path string) (*Tbx, error) {
	tbx := new(Tbx)
	Bix, err := bix.New(path)
	tbx.Bix = Bix
	return tbx, err
}

func (tbx Tbx) Close() error {
	return tbx.Bix.Close()
}

func (tbx Tbx) Query(chrom string, start, end int) []Variant {
	vals, err := tbx.Bix.Query(interfaces.AsIPosition(chrom, start, end))
	simple_util.CheckErr(err)
	var items []Variant
	for {
		v, err := vals.Next()
		if err != nil {
			break
		}
		vt := v.(interfaces.IVariant)
		/*fmt.Println("Id:\t",vt.Id())
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
		*/
		Chrom := vt.Chrom()
		Start := vt.Start()
		End := vt.End()
		Ref := vt.Ref()
		Alt := vt.Alt()

		var variant Variant
		variant.Chrom = Chrom
		variant.Start = Start
		variant.End = End
		variant.Ref = Ref
		variant.Alt = Alt
		var info = make(map[string]interface{})
		for _, k := range FieldInt {
			info[k], _ = vt.Info().Get(k)
		}
		for _, k := range FieldFloat32 {
			info[k], _ = vt.Info().Get(k)
		}
		variant.Info = info
		items = append(items, variant)
	}
	return items
}

func chopVariant(ref, alt string) (refn, altn string, offset int) {
	refs := strings.Split(ref, "")
	alts := strings.Split(alt, "")
	offset = 0
	for j := range refs {
		if j < len(alts) && refs[j] == alts[j] {
			offset++
		} else {
			break
		}
	}
	refn = strings.Join(refs[offset:], "")
	altn = strings.Join(alts[offset:], "")
	return //refn,altn,offset
}

func (tbx Tbx) Hit(chrom string, start, end int, ref, alt string, vals []Variant) Variant {
	var hit Variant
	ref, alt, offset := chopVariant(ref, alt)
	start += offset
	hit.Chrom = chrom
	hit.Start = uint32(start)
	hit.End = uint32(end)
	hit.Ref = ref
	hit.Alt = append(hit.Alt, alt)

	//fmt.Println(hit)

	if len(vals) == 0 {
		return hit
	}
	for _, val := range vals {
		//fmt.Println(val)

		for i, Alt := range val.Alt {
			Ref, Alt, offset := chopVariant(val.Ref, Alt)
			//fmt.Println(Ref,Alt,offset,val)
			if ref != Ref || alt != Alt {
				continue
			}
			val.Start += uint32(offset)
			if val.Start != uint32(start) || val.End != uint32(end) {
				continue
			}

			var info = make(map[string]interface{})
			for _, k := range FieldInt {
				if val.Info[k] == nil {
					continue
				}
				t, ok := val.Info[k].([]int)
				if ok {
					info[k] = t[i]
				} else {
					log.Fatal("[", chrom, " ", start, " ", end, " ", ref, " ", alt, "]\tkey:{"+k+"} can not parse to []int:value:{", val.Info[k], "}")
				}
				//info[k]=t[i]
			}
			for _, k := range FieldFloat32 {
				if val.Info[k] == nil {
					continue
				}
				t, ok := val.Info[k].([]float32)
				if ok {
					info[k] = t[i]
				} else {
					log.Fatal("[", chrom, " ", start, " ", end, " ", ref, " ", alt, "]\tkey:{", k, "} can not parse to []fload32:value:{", val.Info[k], "}")
				}
				//info[k]=t[i]
			}
			hit.Info = info
			break
		}
	}
	return hit
}

var (
	chrPrefix = regexp.MustCompile("^chr")
)

func AddGnomAD(tbx *Tbx, item map[string]string) {
	chr := item["#Chr"]
	chr = chrPrefix.ReplaceAllString(chr, "")
	start, err := strconv.Atoi(item["Start"])
	simple_util.CheckErr(err)
	stop, err := strconv.Atoi(item["Stop"])
	qStart := start
	if start == stop {
		qStart -= 1
	}
	vals := tbx.Query(chr, start-1, stop)
	if vals == nil {
		return
	}

	ref := item["Ref"]
	if ref == "." {
		ref = ""
	}
	alt := item["Call"]
	if alt == "." {
		alt = ""
	}

	hit := tbx.Hit(chr, start, stop, ref, alt, vals)
	if hit.Info == nil {
		return
	}
	if hit.Info["AF"] == nil {
		item["GnomAD AF"] = ""
	} else {
		item["GnomAD AF"] = strconv.FormatFloat(float64(hit.Info["AF"].(float32)), 'f', -1, 32)
	}
	if hit.Info["AF_eas"] == nil {
		item["GnomAD EAS AF"] = ""
	} else {
		item["GnomAD EAS AF"] = strconv.FormatFloat(float64(hit.Info["AF_eas"].(float32)), 'f', -1, 32)
	}
	item["GnomAD HomoAlt Count"] = strconv.Itoa(hit.Info["nhomalt"].(int))
	item["GnomAD EAS HomoAlt Count"] = strconv.Itoa(hit.Info["nhomalt"].(int))
	return
}
