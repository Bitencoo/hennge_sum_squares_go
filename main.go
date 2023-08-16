// I programmed this code based in concurrency and parallelism as go is intended to. It doesn't necessarilly runs faster then
// without using, but I wanted to study more and apply some of my knowledge 
package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	//"time"
)

type nList struct {
	sum     int
	numbers []int
}

func main() {
	numCores := runtime.NumCPU() // Get the number of available CPU cores
	runtime.GOMAXPROCS(numCores) // Utilize all available CPU cores
	var nLines int               // Number of lines of input
	var aux *int                 // Pointer to track remaining lines
	var pos *int                 // Pointer to track current position
	var wg sync.WaitGroup        // WaitGroup for synchronization
	var mu sync.Mutex            // Mutex for synchronization

	numbers := []nList{} // Slice to store line information
	fmt.Scan(&nLines)    // Read the number of lines from input

	aux = &nLines                          // Set aux to point to nLines
	pos = new(int)                         // Create a new integer variable pointed by pos
	*pos = 0                               // Initialize pos to 0
	numbers = readInput(aux, pos, numbers) // Read input lines

	// Starting a timer at the beginning of the compute
	// start := time.Now()

	wg.Add(1) // Adding the goroutine to the WaitGroup
	go calculateSumRecursive(numbers, 0, &wg, &mu) // Start parallel calculation
	wg.Wait() // Waits for all goroutines to be done before executing the rest of the code

	// Starting the timer at the end of the compute
	// elapsed := time.Since(start)

	// Only to check how much time my code used to finish the computation
	// fmt.Printf("Time: %s", elapsed)

	// Prints the sum values for each line using recursion
	printSliceRecursive(numbers, 0)

}

// Recursive function to read quantity of elements in the respective line
func readInput(aux *int, pos *int, numbers []nList) []nList {
	if *aux == 0 {
		return numbers
	}

	var qtyNumbers int
	var newLine nList

	fmt.Scan(&qtyNumbers) // Read the quantity of numbers in the line

	newLine.numbers = make([]int, qtyNumbers) // Initialize slice for numbers
	newLine.sum = 0                           // Initialize sum for the line

	numbers = append(numbers, newLine)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	num := scanner.Text() // Reads the line of numbers
	numStrings := strings.Fields(num) // Transform the line of numbers into a slice of Strings
	numbers[*pos].numbers = readLineNumbers(qtyNumbers, *pos, numbers, numStrings, 0) // Recursively transform all string numbers to int

	*pos++
	*aux--
	return readInput(aux, pos, numbers) // Recursively read input for other lines
}

// Recursive function to read every element of the current line
func readLineNumbers(qtyNumbers int, linePos int, numbers []nList, numbersString []string, numberPosition int) []int {
	if qtyNumbers == 0 {
		return numbers[linePos].numbers
	}

	num, err := strconv.Atoi(numbersString[numberPosition]) // Converts String input to int
	if err != nil {
		fmt.Println("Error converting input to int:", err)
		return numbers[linePos].numbers
	}

	numbers[linePos].numbers[numberPosition] = num
	qtyNumbers--
	numberPosition++

	return readLineNumbers(qtyNumbers, linePos, numbers, numbersString, numberPosition) // Recursively read numbers
}

// Iterate through each line saying
func calculateSumRecursive(numbers []nList, indexLine int, wg *sync.WaitGroup, mu *sync.Mutex) {
	if indexLine >= len(numbers) {
		wg.Done()
		return
	}

	line := &numbers[indexLine]
	line.sum = 0

	calculateLineSumRecursive(line.numbers, 0, &line.sum, mu) // Calculate sum for the line

	go calculateSumRecursive(numbers, indexLine+1, wg, mu) // Recursively calculate for the next line
}

// Iterate through each element of the slice and realizes the square of that number
func calculateLineSumRecursive(numbers []int, currentIndex int, sumPointer *int, mu *sync.Mutex) {
	if currentIndex >= len(numbers) {
		return
	}

	value := numbers[currentIndex]
	if value > 0 {
		mu.Lock() //Locking the mutex, we need to guarantee that only one goroutine can access the variable at a time
		*sumPointer += value * value
		mu.Unlock() //Unlocking the mutex after making all the procedures needed
	}

	calculateLineSumRecursive(numbers, currentIndex+1, sumPointer, mu)
}

// Recursive function to print the result of the sum of each line of the slice
func printSliceRecursive(numbers []nList, currentIndex int) {
	if len(numbers) == currentIndex {
		return
	}

	fmt.Println(numbers[currentIndex].sum)

	printSliceRecursive(numbers, currentIndex+1)
}
