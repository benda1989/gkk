package compare

import (
	"github.com/benda1989/gkk/str"
)

func Map(n, o map[string]any) map[string][]string {
	return MapTrans(n, o, nil)
}

func MapTrans(n, o map[string]any, trans map[string]string) map[string][]string {
	re := map[string][]string{}
	for k, v := range n {
		if vv, ok := o[k]; ok {
			var ns, old string
			switch v.(type) {
			case string:
				ns = v.(string)
				old, _ = vv.(string)
			default:
				ns = str.String(v)
				old = str.String(vv)
			}
			if old != ns {
				if trans != nil {
					if tran, ok := trans[k]; ok {
						k = tran
					}
				}
				re[k] = []string{old, ns}
			}
		}
	}
	return re
}
