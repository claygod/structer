package structer

// Structer
// Tests and benchmarks
// Copyright © 2017 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "fmt"
import "testing"
import "os"
import "encoding/gob"
import "strconv"

//import "encoding/gob"
//import "bytes"
import "log"
import "time"

//import "reflect"

//import "unsafe"
/*
func TestInt64ToInt32(t *testing.T) {
	var x int64 = -5
	log.Print("Сконвертировали `-5` из int64 b int32 и получили: ", int32(x))
}
*/
func TestEngine(t *testing.T) {
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  244234201,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"News", "April"},
	}

	s, _ := New(a, "Id", []string{"Title", "Date", "Tags"})

	s.Add(a)

	a.Tags = []string{"News", "May"}
	a.Id = "a222"
	a.Date = 244234209
	s.Add(a)

	a.Tags = []string{"News", "May"}
	a.Id = "a333"
	a.Date = 244234203
	s.Add(a)
	a.Tags = []string{"News", "May"}
	a.Id = "a444"
	a.Date = 244234204
	s.Add(a)

	a.Tags = []string{"Docs", "May"}
	a.Id = "a555"
	a.Date = 244234205
	s.Add(a)

	for i := 100; i < 1600; i++ {
		a.Id = "b" + strconv.Itoa(i)
		a.Tags = []string{"Fotos", "May"}
		//s.Add(a)
	}
	log.Print("Start SELECT", time.Now().UnixNano())
	tStart := time.Now().UnixNano()

	rez2, _ := s.Find().
		ByFields([]string{"TagsMay", "TagsNews"}).
		OrderBy("Date", DESC).
		Limit(0, 5).
		Do()

	tStart2 := time.Now().UnixNano()
	log.Print("Время проведения теста(1): ", tStart2-tStart)
	log.Print("Получилось 2): ", len(rez2))
	log.Print("Получилось 2): ", rez2)
	//for k, v := range s.tags.subTags {
	//	log.Print(" -- ", k, ":", v.count)
	//}
}

func TestSelect(t *testing.T) {
	a := Article{
		Id:    "a",
		Pub:   true,
		Date:  244234201,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"News", "April"},
	}
	for i := 0; i < 100; i++ {
		a.Text += " bla-bla-bla"
	}
	s, _ := New(a, "Id", []string{"Date", "Tags"})

	for i := 0; i < 50; i++ {
		a.Id = "a" + strconv.Itoa(i)
		s.Add(a)
	}

	for i := 500; i < 1000; i++ {
		a.Id = "a" + strconv.Itoa(i)
		a.Tags = []string{"News", "May"}
		s.Add(a)
	}

	for i := 1000; i < 3000; i++ {
		a.Id = "a" + strconv.Itoa(i)
		a.Tags = []string{"Docs", "May"}
		s.Add(a)
	}

	for i := 3000; i < 5000; i++ {
		a.Id = "a" + strconv.Itoa(i)
		a.Tags = []string{"Scans", "May"}
		s.Add(a)
	}
	for i := 5000; i < 15000; i++ {
		a.Id = "a" + strconv.Itoa(i)
		a.Tags = []string{"Scans", "April"}
		s.Add(a)
	}

	for i := 15000; i < 66000; i++ {
		a.Id = "a" + strconv.Itoa(i)
		a.Tags = []string{"Fotos", "July"}
		//s.Add(a)
	}

	for i := 66000; i < 497000; i++ {
		a.Id = "a" + strconv.Itoa(i)
		a.Tags = []string{"Fotos", "May"}
		//s.Add(a)
	}

	for i := 497000; i < 497500; i++ {
		a.Id = "a" + strconv.Itoa(i)
		a.Tags = []string{"May"}
		s.Add(a)
	}

	for i := 497500; i < 499000; i++ {
		a.Id = "a" + strconv.Itoa(i)
		a.Tags = []string{"News"}
		s.Add(a)
	}

	//log.Print("Start TEST! ", time.Now().UnixNano())
	//tStart := time.Now().UnixNano()

	rez2, _ := s.Find().
		ByFields([]string{"TagsNews", "TagsMay"}).
		OrderBy("Date", ASC).
		Limit(0, 600).
		Do()

		//tStart2 := time.Now().UnixNano()
		//log.Print("Время проведения теста(2) ", tStart2-tStart)
		//log.Print("Результат теста(2)", len(rez2))

	if len(rez2) != 500 {
		t.Error(fmt.Sprintf("Ошибка SELECT. Ожидаемое значение 500 а полученное: %i", len(rez2)))
	}
}

func TestDel(t *testing.T) {
	a := Article{
		Id:    "a",
		Pub:   true,
		Date:  244234201,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"News", "April"},
	}

	s, _ := New(a, "Id", []string{"Date", "Tags"})

	for i := 0; i < 5; i++ {
		a.Id = "a" + strconv.Itoa(i)
		s.Add(a)
	}

	// log.Print("Start DEL! ", time.Now().UnixNano())
	log.Print("Фрагментированно(0): ", s.Fragmentation(), " ", s.Count())
	s.Del("a1")
	s.Del("a2")

	rez2, _ := s.Find().
		ByFields([]string{"TagsNews"}).
		OrderBy("Date", ASC).
		Limit(0, 600).
		Do()
	if len(rez2) != 3 {
		t.Error(fmt.Sprintf("Ошибка удаления. Ожидаемое значение 3 а полученное: %i", len(rez2)))
	}
	// log.Print("Результат теста по удалению: ", len(rez2), " ", rez2)
	log.Print("Фрагментированно(1): ", s.Fragmentation(), " ", s.Count())
	s.Revision()
	log.Print("Фрагментированно(2): ", s.Fragmentation(), " ", s.Count())
}

func TestUpdate(t *testing.T) {
	a := Article{
		Id:    "a",
		Pub:   true,
		Date:  244234201,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"News", "April"},
	}

	s, _ := New(a, "Id", []string{"Date", "Rate", "Tags"})

	for i := 0; i < 6; i++ {
		a.Id = "a" + strconv.Itoa(i)
		s.Add(a)
	}
	a.Id = "a7"
	a.Rate = 772
	s.Add(a)
	a.Id = "a8"
	a.Rate = 771
	s.Add(a)

	// log.Print("Start UPDATE! ", time.Now().UnixNano())
	a.Id = "a4"
	a.Date = 244777201
	a.Rate = 999
	s.Update(a)

	rez2, _ := s.Find().
		ByFields([]string{"TagsNews"}).
		OrderBy("Rate", ASC).
		Limit(0, 600).
		Do()
	// log.Print("Результат теста по изменению: ", len(rez2), " ", rez2)

	for _, art := range rez2 {
		art2 := art.(Article)
		if art2.Id == "a7" && art2.Rate != 772 {
			t.Error(fmt.Sprintf("Ошибка UPDATE. Ожидаемое значение 772 а полученное: %i", art2.Rate))
		}
		if art2.Id == "a4" && art2.Rate != 999 {
			t.Error(fmt.Sprintf("Ошибка UPDATE. Ожидаемое значение 999 а полученное: %i", art2.Rate))
		}
	}
}

func TestGet(t *testing.T) {
	a := Article{
		Id:    "a",
		Pub:   true,
		Date:  244234201,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"News", "April"},
	}

	s, _ := New(a, "Id", []string{"Date", "Tags"})

	for i := 0; i < 5; i++ {
		a.Id = "a" + strconv.Itoa(i)
		s.Add(a)
	}

	log.Print("Start GET! ", time.Now().UnixNano())

	rez2 := s.Get("a1")
	log.Print("Результат теста GET: ", rez2)
}

func TestSave(t *testing.T) {
	a := Article{
		Id:    "a",
		Pub:   true,
		Date:  244234201,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"News", "April"},
	}

	s, _ := New(a, "Id", []string{"Date", "Tags"})

	for i := 0; i < 70000; i++ {
		a.Id = "a" + strconv.Itoa(i)
		s.Add(a)
	}
	s.Del("a2")

	log.Print("Start Save! ", time.Now().UnixNano())

	rez2 := s.Save("test.gob")
	log.Print("Результат теста Save: ", rez2)
}

func TestLoad(t *testing.T) {
	a := Article{
		Id:    "a",
		Pub:   true,
		Date:  244234201,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"News", "April"},
	}

	s, _ := New(a, "Id", []string{"Date", "Tags"})

	log.Print("Start Load! ", time.Now().UnixNano())

	//s.Load("test.gob")

	file, err := os.Open("test.gob")
	//log.Print("=начало2= ", err)
	if err != nil {
		log.Print("========ERROR LOAD!! +++++= ")
	}
	//log.Print("=начало3= ", path)
	decoder := gob.NewDecoder(file)
	for {
		err := decoder.Decode(&a)
		if err != nil {
			break
		}
		//log.Print("=====================обратно= ", err, a)
		s.AddUnsafe(a)
	}
	//log.Print("Результат теста Load: ", rez2)
	rez2, _ := s.Find().
		ByFields([]string{"TagsNews"}).
		OrderBy("Date", ASC).
		Limit(0, 600).
		Do()
	log.Print("Результат теста по LOAD ", len(rez2), " ")
}

func TestRevision(t *testing.T) {
	a := Article{
		Id:    "a",
		Pub:   true,
		Date:  244234201,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"News", "April"},
	}

	s, _ := New(a, "Id", []string{"Date", "Tags"})

	for i := 0; i < 70000; i++ {
		a.Id = "a" + strconv.Itoa(i)
		s.Add(a)
	}
	s.Del("a2")
	s.Revision()

	//log.Print("Start Save! ", time.Now().UnixNano())

	//rez2 := s.Save("test.gob")
	//log.Print("Результат теста Save: ", rez2)
}

/*
func Test201(t *testing.T) {
	p := GetTestStruct()
	//t1 := reflect.TypeOf(p)
	//log.Print(t1)
	z := (*iface)(unsafe.Pointer(&p))
	log.Print("----------------")
	//log.Print(z.data)
	log.Print(*(*string)(unsafe.Pointer(uintptr(unsafe.Pointer(z.data)) + 16)))
	//log.Print(z.tab)
}
*/
/*
func Test200(t *testing.T) {
	//a := Article{Id: "a1", Title: "abc"}
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  244234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}

	edb, err := New(a)
	if err != nil {
		t.Error("Failed to create a new database:", err)
	}
	edb.Add(a)
}
*/
/*
func Test100(t *testing.T) {
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  244234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	//dec := gob.NewDecoder(&network) // Will read from network.

	err := enc.Encode(a)
	if err != nil {
		t.Error("encode error:", err)
	}
}

func Test101(t *testing.T) {
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  244234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}
	b := Article{
		Id:    "a1",
		Pub:   true,
		Date:  244234231,
		Title: "Two article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}
	fmt.Print(unsafe.Offsetof(a.Desc), "\r\n")
	fmt.Print(unsafe.Offsetof(b.Desc), "\r\n")

}
*/

/*


func BenchmarkStorageAdd(b *testing.B) {
	b.StopTimer()
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  244234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}
	s := NewStorage()

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		s.Add(uint64(n), a)
	}
}

func BenchmarkStorageAddParallel(b *testing.B) {
	b.StopTimer()
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  244234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}
	s := NewStorage()
	n := uint64(0)
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Add(uint64(n), a)
			n++
		}
	})
}

func BenchmarkFieldUnsafe(b *testing.B) {
	b.StopTimer()
	p := GetTestStruct()
	//z := (*iface)(unsafe.Pointer(&p))
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		//z := (*iface)(unsafe.Pointer(&p))
		_ = *(*string)(unsafe.Pointer(uintptr(unsafe.Pointer((*iface)(unsafe.Pointer(&p)).data)) + 16))
	}
}

func BenchmarkFieldReflect(b *testing.B) {
	b.StopTimer()
	p := GetTestStruct()

	//z := (*iface)(unsafe.Pointer(&p))
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		t1 := reflect.TypeOf(p)
		t1.Field(1)
	}
}

func BenchmarkGobEn(b *testing.B) {
	b.StopTimer()
	//arr = make([]Article, 0)
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  244234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}
	//for i := 0; i < 10000; i++ {
	//	arr = append(arr, a)
	//}

	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	//dec := gob.NewDecoder(&network) // Will read from network.

	err := enc.Encode(a)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	//var q Article
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		//	a.Id = strconv.Itoa(n)
		//	d.Add(a, fu)
		enc.Encode(a)
		//dec.Decode(&q)
	}
}
func BenchmarkGobEnDe(b *testing.B) {
	b.StopTimer()
	//arr = make([]Article, 0)
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  244234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}
	//for i := 0; i < 10000; i++ {
	//	arr = append(arr, a)
	//}

	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	dec := gob.NewDecoder(&network) // Will read from network.

	err := enc.Encode(a)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	var q Article
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		//	a.Id = strconv.Itoa(n)
		//	d.Add(a, fu)
		enc.Encode(a)
		dec.Decode(&q)
	}
}

func BenchmarkGobEnParallel(b *testing.B) {
	b.StopTimer()
	//arr = make([]Article, 0)
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  244234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}

	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	//dec := gob.NewDecoder(&network) // Will read from network.

	err := enc.Encode(a)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	//for i := 0; i < 10000; i++ {
	//	enc.Encode(a)
	//}

	//var q Article
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			enc.Encode(a)
			//dec.Decode(&q)
		}
	})
}
*/
/*
func BenchmarkMapCreate(b *testing.B) {
	b.StopTimer()
	q := make(map[int]bool)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		q = make(map[int]bool)
	}
	q[0] = true
}

func BenchmarkMapAdd(b *testing.B) {
	b.StopTimer()
	q := make(map[int]bool)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		k := int(uint16(n))
		q[k] = true
	}
	q[0] = true
}

func BenchmarkMapAddDel(b *testing.B) {
	b.StopTimer()
	q := make(map[int]bool)
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		k := int(uint16(n))
		q[k] = true
		delete(q, k)
	}
	q[0] = true
}

func BenchmarkMapIsset(b *testing.B) {
	b.StopTimer()
	q := make(map[int]bool)

	for n := 0; n < 2000; n++ {
		k := int(uint16(n))
		q[k] = true

	}

	b.StartTimer()

	for n := 0; n < b.N; n++ {
		k := int(uint16(n))
		if _, ok := q[k]; ok {

		}
	}
	q[0] = true
}

func BenchmarkMapGetNil(b *testing.B) {
	b.StopTimer()
	q := make(map[int]bool)

	for n := 0; n < 2000; n++ {
		k := int(uint16(n))
		q[k] = true

	}
	var x bool
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		k := int(uint16(n))
		x = q[k]
	}
	q[0] = x
}

func BenchmarkMapFor20000(b *testing.B) {
	b.StopTimer()
	q := make(map[int]bool)

	for n := 0; n < 20000; n++ {
		q[n] = true

	}
	var x int
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		for k, _ := range q {
			x = k
		}
	}
	q[x] = true
}

func BenchmarkMapMake(b *testing.B) {
	b.StopTimer()
	q := make(map[string]int)

	b.StartTimer()

	for n := 0; n < b.N; n++ {
		q = make(map[string]int)
	}
	q["o"] = 0
}

func BenchmarkSliceIsset(b *testing.B) {
	b.StopTimer()
	q := make([]int8, 70000)

	var x int8
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		k := int(uint16(n))
		if q[k] == 1 {
			x = q[k]
		}
	}
	q[0] = x
}

func BenchmarkSliceAdd(b *testing.B) {
	b.StopTimer()
	q := make([]int, 0, 100)

	b.StartTimer()

	for n := 0; n < b.N; n++ {
		q = append(q, n)
	}
	q[0] = 0
}

func BenchmarkSliceMake(b *testing.B) {
	b.StopTimer()
	q := make([]int, 0, 100)

	b.StartTimer()

	for n := 0; n < b.N; n++ {
		q = make([]int, 0, 100)
	}
	q = append(q, 0)
}

func BenchmarkSliceFor20000(b *testing.B) {
	b.StopTimer()
	q := make([]int, 20000)

	var x int
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		for _, k := range q {
			x = k
		}
	}
	q = append(q, x)
}
*/

func BenchmarkEngineAdd(b *testing.B) {
	b.StopTimer()
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  244234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}
	s, _ := New(a, "Id", []string{"Date", "Tags"})

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		a.Id = "a" + strconv.Itoa(n)
		tags := make([]string, 2)
		tags[0] = strconv.Itoa(int((byte(n) << 4) >> 4)) // только их тут 16
		tags[1] = strconv.Itoa(int(byte(n)))             // категория - тут их 256

		a.Tags = tags
		a.Date++ // = int(byte(n))
		s.Add(a)
	}
}

func BenchmarkEngineAddParallel(b *testing.B) {
	b.StopTimer()
	a := Article{
		Id:    "a1",
		Pub:   true,
		Date:  244234231,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"news", "april"},
	}
	s, _ := New(a, "Id", []string{"Date", "Tags"})
	n := 0
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			a.Id = "a" + strconv.Itoa(n)
			tags := make([]string, 2)
			tags[0] = strconv.Itoa(int((byte(n) << 4) >> 4)) // только их тут 16
			tags[1] = strconv.Itoa(int(byte(n)))             // категория - тут их 256

			a.Tags = tags
			a.Date++ // = int(byte(n))
			s.Add(a)
			n++
		}
	})
}

/*
 */
func BenchmarkEngineSelect256(b *testing.B) {
	b.StopTimer()
	a := Article{
		Id:    "a",
		Pub:   true,
		Date:  244234201,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"News", "April"},
	}
	for i := 0; i < 100; i++ {
		a.Text += " bla-bla-bla"
	}
	s, _ := New(a, "Id", []string{"Date", "Tags"})

	for i := 0; i < 60000; i++ {
		a.Id = "a" + strconv.Itoa(i)
		tags := make([]string, 2)
		tags[0] = strconv.Itoa(int((byte(i) << 4) >> 4)) // только их тут 16
		tags[1] = strconv.Itoa(int(byte(i)))             // категория - тут их 256

		a.Tags = tags
		a.Date++
		s.Add(a)
	}
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		s.Find().
			ByFields([]string{"Tags" + strconv.Itoa(int((byte(n)<<4)>>4)),
				"Tags" + strconv.Itoa(int(byte(n)>>1))}).
			OrderBy("Date", ASC).
			Limit(0, 10).
			Do()
	}
}

func BenchmarkEngineSelect256Parallel(b *testing.B) {
	b.StopTimer()
	a := Article{
		Id:    "a",
		Pub:   true,
		Date:  244234201,
		Title: "First article",
		Desc:  "Description",
		Text:  "Big text body",
		Tags:  []string{"News", "April"},
	}
	for i := 0; i < 100; i++ {
		a.Text += " bla-bla-bla"
	}
	s, _ := New(a, "Id", []string{"Date", "Tags"})

	for i := 0; i < 60000; i++ {
		a.Id = "a" + strconv.Itoa(i)
		tags := make([]string, 2)
		tags[0] = strconv.Itoa(int((byte(i) << 4) >> 4)) // только их тут 16
		tags[1] = strconv.Itoa(int(byte(i)))             // категория - тут их 256

		a.Tags = tags
		s.Add(a)
	}
	b.StartTimer()
	n := 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Find().
				ByFields([]string{"Tags" + strconv.Itoa(int((byte(n)<<4)>>4)),
					"Tags" + strconv.Itoa(int(byte(n)>>1))}).
				OrderBy("Date", ASC).
				Limit(0, 10).
				Do()
			n++
		}
	})
}

/*
=== laptop
BenchmarkEngineAdd-4                 	  300000	      6327 ns/op
BenchmarkEngineAddParallel-4         	  500000	      4294 ns/op
BenchmarkEngineSelect256-4           	  200000	      6680 ns/op
BenchmarkEngineSelect256Parallel-4   	  500000	      4044 ns/op
===
*/
/*
const POOL_DEGREE uint64 = 16
BenchmarkEngineAdd-4              	  500000	      4562 ns/op
BenchmarkEngineAddParallel-4      	 1000000	      3481 ns/op
BenchmarkEngineSelect-4           	    2000	    883588 ns/op
BenchmarkEngineSelectParallel-4   	    5000	    542854 ns/op

const POOL_DEGREE uint64 = 8
BenchmarkEngineAdd-4              	  500000	      4948 ns/op
BenchmarkEngineAddParallel-4      	  300000	      3440 ns/op
BenchmarkEngineSelect-4           	    2000	    891500 ns/op
BenchmarkEngineSelectParallel-4   	    5000	    460046 ns/op

const POOL_DEGREE uint64 = 4
BenchmarkEngineAdd-4              	  500000	      5860 ns/op
BenchmarkEngineAddParallel-4      	  300000	      3737 ns/op
BenchmarkEngineSelect-4           	    3000	    760409 ns/op
BenchmarkEngineSelectParallel-4   	    5000	    349634 ns/op

const POOL_DEGREE uint64 = 2
BenchmarkEngineAdd-4              	  500000	      5524 ns/op
BenchmarkEngineAddParallel-4      	  300000	      3490 ns/op
BenchmarkEngineSelect-4           	    2000	    657000 ns/op
BenchmarkEngineSelectParallel-4   	   10000	    361000 ns/op




BenchmarkMapCreate-4  				10000000	       173 ns/op
BenchmarkMapAdd-4                 	20000000	        75.0 ns/op
BenchmarkMapAddDel-4              	10000000	       118 ns/op
BenchmarkMapIsset-4               	50000000	        40.7 ns/op
BenchmarkMapGetNil-4              	50000000	        38.4 ns/op
BenchmarkMapFor20000-4            	    3000	    552364 ns/op
BenchmarkMapMake-4                	10000000	       138 ns/op
BenchmarkSliceIsset-4             	2000000000	         1.75 ns/op
BenchmarkSliceAdd-4               	100000000	        21.4 ns/op
BenchmarkSliceMake-4              	 5000000	       405 ns/op
BenchmarkSliceFor20000-4          	  200000	     10570 ns/op
BenchmarkEngineAdd-4              	  500000	      3782 ns/op
BenchmarkEngineAddParallel-4      	  500000	      2736 ns/op
BenchmarkEngineSelect-4           	   10000	    177410 ns/op
BenchmarkEngineSelectParallel-4   	   20000	     88955 ns/op
*/

type Article struct {
	Id    string
	Pub   bool
	Date  int
	Rate  int
	Title string
	Desc  string
	Text  string
	Tags  []string
}
