package main

import (
	"fmt"
	"github.com/gtygo/artmap"
	"sync"
	"time"
)

func main() {
	m := artmap.New()

	benchSize := 10000000
	goroutineSize := 4
	dataSet1 := make([][]byte, benchSize)
	dataSet2 := make([]string, benchSize)

	println("prepare data set ...")
	for i := 0; i < benchSize; i++ {
		dataSet1 = append(dataSet1, []byte(fmt.Sprintf("test:%d", i)))
	}
	for i := 0; i < benchSize; i++ {
		dataSet2 = append(dataSet2, fmt.Sprintf("test:%d", i))
	}
	wg := sync.WaitGroup{}
	println("start set art tree ...")
	start := time.Now()
	for i := 0; i < goroutineSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < benchSize; i++ {
				m.Set(dataSet1[i], "1")
			}

		}()
	}
	wg.Wait()
	end := time.Since(start)

	fmt.Printf("cost : %v \n", end)
	m1 := sync.Map{}
	wg2 := sync.WaitGroup{}
	println("start set map ...")
	start1 := time.Now()
	for i := 0; i < goroutineSize; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			for i := 0; i < benchSize; i++ {
				m1.Store(dataSet2[i], 1)
			}

		}()
	}
	wg2.Wait()
	end1 := time.Since(start1)
	fmt.Printf("cost : %v \n", end1)
}
