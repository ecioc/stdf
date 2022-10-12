package stdf

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func TestAnalysis(t *testing.T) {
	//读取stdf内容
	f, err := os.OpenFile("./mock/HF0062B_V1_F4DUT_TEST_P11464__09_06282021_101511.stdf", os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	b1 := make([]byte, 4)
	index := 0
	mI := 0
	// var sg sync.WaitGroup
	for {
		f.ReadAt(b1, int64(index))
		// 通过4个字节串生成一个stdf结构对象
		t.Log("\n 结构头:", b1)
		o1 := NewStdfRecord(b1)
		if o1 == nil {
			return
		}
		// 通过反射获取对象
		v1 := reflect.ValueOf(o1)
		// t1 := reflect.TypeOf(o1)
		// 创建一个content的byte切片
		ii := v1.Elem().FieldByName("Rec_Len").Uint()
		t.Log("\n 对象对应的内容长度字节数:", ii)
		b2 := make([]byte, ii)
		index += 4
		f.ReadAt(b2, int64(index))
		t.Log("\n 对象对应的字节内容:", b2)
		TransB2S(b2, o1)
		// var params []reflect.Value
		// params = append(params, reflect.ValueOf(b2))
		// v1.MethodByName("ToRecord").Call(params)
		t.Log("\n 对象对应的字符串内容:", v1.MethodByName("ToString").Call(nil))

		// t.Log("Number of fields", v1.Elem().NumField())
		// for i := 0; i < v1.Elem().NumField(); i++ {
		// 	field := t1.Elem().Field(i)
		// 	value := v1.Elem().Field(i).Interface()
		// 	fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
		// 	t.Logf("Field:%d \t type:%T \t value:%v\n",
		// 		i, v1.Elem().Field(i), v1.Elem().Field(i))
		// }

		index += int(ii)
		mI++
		if mI == 3 {
			break
		}

		// f.ReadAt()
	}
	// sg.Wait()
	// fmt.Println(far)
}
