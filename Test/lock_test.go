package Test

import (
	"sync"
	"testing"
)

var wg sync.WaitGroup

func fibonacci(num int) int{
	if num<2{
		return 1
	}
	return fibonacci(num-1) + fibonacci(num-2)
}

func BenchmarkLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var mu sync.Mutex
		cnt:=0
		for j:=0;j<1000;j++{
			go func() {
				for k:=0;k<100;k++{
					mu.Lock()
					cnt++
					mu.Unlock()
				}
			}()
		}
	}
}

func BenchmarkSpinLock(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		sp:=NewSpinLock()
		cnt:=0
		for j:=0;j<1000;j++{
			go func() {
				for k:=0;k<100;k++{
					sp.Lock()
					cnt++
					sp.Unlock()
				}
			}()
		}
	}
}

func BenchmarkLongTaskMutex(b *testing.B) {
	wg.Add(100)
	for i := 0; i < b.N; i++ {
		var mu sync.Mutex
		for j:=0;j<1000;j++{
			for k:=0;k<100;k++{
				mu.Lock()
				fibonacci(20)
				mu.Unlock()
			}
		}
	}
}


func BenchmarkLongTaskSpinLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sp:= NewSpinLock()
		for j:=0;j<1000;j++{
			go func() {
				for k:=0;k<100;k++{
					sp.Lock()
					fibonacci(20)
					sp.Unlock()
				}
			}()
		}
	}
}