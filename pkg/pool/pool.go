package pool

import "sync"

type worker[D any] func(item D)

func Process[D any](data *[]D, w worker[D], size int) {
	wg := sync.WaitGroup{}

	length := len(*data)

	var workerPoolSize int

	if length < size {
		workerPoolSize = length
	} else {
		workerPoolSize = size
	}

	dataCh := make(chan D, workerPoolSize)

	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for data := range dataCh {
				w(data)
			}
		}()
	}

	for _, v := range *data {
		dataCh <- v
	}

	close(dataCh)

	wg.Wait()
}
