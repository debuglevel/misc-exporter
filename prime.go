package main

// Borrowed from ChatGPT
func getPrimes(maximum int) []int {
	//zap.S().Debugf("Getting primes up to %v...\n", maximum)

	// Create a slice to store the primes
	var primes []int

	// Iterate through the numbers from 2 to 1000
	for i := 2; i <= maximum; i++ {
		// Assume the number is prime
		isPrime := true

		// Check if the number is prime by checking if it is evenly divisible by any of the primes we have found so far
		for _, prime := range primes {
			if i%prime == 0 {
				// If the number is evenly divisible by a prime, it is not prime
				isPrime = false
				break
			}
		}

		// If the number is prime, add it to the slice of primes
		if isPrime {
			//zap.S().Debugf("Found prime %v\n", i)
			primes = append(primes, i)
		}
	}

	return primes
}
