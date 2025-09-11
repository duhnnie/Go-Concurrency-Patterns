package main

import "fmt"

func generator(n int) <-chan int {
	output := make(chan int)

	go func() {
		for i := 2; i < n; i++ {
			output <- i
		}

		close(output)
	}()

	return output
}

func filterEven(input <-chan int) <-chan int {
	output := make(chan int)

	go func() {
		for number := range input {
			if number == 2 || number%2 != 0 {
				output <- number
			}
		}

		close(output)
	}()

	return output
}

func primeChecker(input <-chan int) <-chan int {
	output := make(chan int)

	go func() {
		for number := range input {
			isPrime := true

			for i := 2; i*i <= number; i++ {
				if number%i == 0 {
					isPrime = false
					break
				}
			}

			if isPrime {
				output <- number
			}
		}

		close(output)
	}()

	return output
}

func collector(input <-chan int) <-chan []int {
	output := make(chan []int)

	go func() {
		primeNumbers := []int{}

		for number := range input {
			primeNumbers = append(primeNumbers, number)
		}

		output <- primeNumbers
		close(output)
	}()

	return output
}

func main() {
	upperLimit := 50
	generatorChannel := generator(upperLimit)
	filterChannel := filterEven(generatorChannel)
	primeChannel := primeChecker(filterChannel)
	collectorChannel := collector(primeChannel)

	fmt.Printf("Primes up to %d: %v", upperLimit, <-collectorChannel)
}
