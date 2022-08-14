package specific

import (
	"bytes"
	"fmt"
)

func CheckSliceShare1()  {
	// 1.a and b share the same memory
	a := make([]int, 32)
	b := a[1:16]
	fmt.Printf("%T\t %v\n", a, a)
	fmt.Printf("%T\t %v\n", b, b)
	fmt.Printf("a: %p\t b: %p\n", &a, &b)
	fmt.Printf("a: %p\t b: %p\n", &a[1], &b[0])

	// 2. 1)a cap is not enough -> a reallocates the memory
	// 2) pointer a and b has not changed
	a = append(a, 1)
	a[2] = 43

	fmt.Printf("%T\t %v\n", a, a)
	fmt.Printf("%T\t %v\n", b, b)
	fmt.Printf("a: %p\t b: %p\n", &a, &b)
	fmt.Printf("a: %p\t b: %p\n", &a[1], &b[0]) // if a := make([]int, 32, 33), then here is the same
}

func CheckSliceShare2()  {
	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/')

	// 1.a b c share the same memory
	a := path[:sepIndex]
	b := path[sepIndex + 1:]
	c := path[:sepIndex:sepIndex] // Full Slice Expression
	fmt.Println("dir1 =>",string(a)) //prints: dir1 => AAAA
	fmt.Println("dir2 =>",string(b)) //prints: dir2 => BBBBBBBBB
	fmt.Println("dir3 =>",string(c)) //prints: dir3 => AAAA
	fmt.Printf("%T\t %v\t %v\n", a, len(a), cap(a))
	fmt.Printf("%T\t %v\t %v\n", b, len(b), cap(b))
	fmt.Printf("%T\t %v\t %v\n", c, len(c), cap(c))
	fmt.Printf("a: %p\t b: %p\t c: %p\n", &a, &b, c)
	fmt.Printf("a: %p\t b: %p\t c: %p\n", &a[3], &b[0], &c[3]) // 3 is man len a, can't get 5

	// 2. 1)a cap is enough -> a won't reallocates the memory -> effect b's memory
	// 2) pointer a and b has not changed forever
	// 3) c cap is not enough -> c reallocate -> won't effect a and b
	a = append(a, "suffix"...)
	// c = append(c, "suffix"...) when use c to append, a and b won't change, then a can not get a[5]

	fmt.Println("dir1 =>",string(a)) //prints: dir1 => AAAAsuffix
	fmt.Println("dir2 =>",string(b)) //prints: dir2 => uffixBBBB
	fmt.Println("dir3 =>",string(c)) //prints: dir3 => AAAA/AAAAsuffix
	fmt.Printf("%T\t %v\t %v\n", a, len(a), cap(a))
	fmt.Printf("%T\t %v\t %v\n", b, len(b), cap(b))
	fmt.Printf("%T\t %v\t %v\n", c, len(c), cap(c))
	fmt.Printf("a: %p\t b: %p\t c: %p\n", &a, &b, c)
	fmt.Printf("a: %p\t b: %p\t c: %p\n", &a[5], &b[0], &c[3]) // because they share the same memory, so can get 5
}

