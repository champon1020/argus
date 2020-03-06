package service

import "reflect"

func GenFlg(st interface{}, fieldNames ...string) (flg uint32) {
	v := reflect.Indirect(reflect.ValueOf(st))
	t := v.Type()
	for _, fn := range fieldNames {
		for j := 0; j < t.NumField(); j++ {
			if fn == t.Field(j).Name {
				flg |= 1 << j
			}
		}
	}
	return
}
