package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"proj3/iooperation"
	"proj3/newstructs"
	"strconv"
	"strings"
	"sync"
)

//GenerateRandomMatrix generates we want
func GenerateRandomMatrix(n int, ranged float64) newstructs.FMatrix {
	var matrix newstructs.FMatrix
	matrix.Row = n
	a := make([][]float64, n)
	for i := range a {
		a[i] = make([]float64, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			a[i][j] = rand.Float64() * ranged
		}
	}
	matrix.Tables = a
	return matrix
}

func sequentialCreate(numberOfMatrix int, sizeOfMatrix int, path string) {
	var tempPath string
	for i := 0; i < numberOfMatrix; i++ {
		x := GenerateRandomMatrix(sizeOfMatrix, 100.0)
		tempPath = path + strconv.Itoa(i)
		iooperation.CreateFile(tempPath)
		iooperation.WriteFile(x.Tables, tempPath)
	}
}

func collectMatrixFromFile(s string, sizeOfMatrix int) newstructs.FMatrix {
	s = strings.TrimSpace(s)
	lst := strings.Split(s, "] [")

	var matrix newstructs.FMatrix
	matrix.Row = sizeOfMatrix
	a := make([][]float64, sizeOfMatrix)
	for i := range a {
		a[i] = make([]float64, sizeOfMatrix)
	}
	for i := 0; i < sizeOfMatrix; i++ {
		line := strings.Split(lst[i], " ")
		for j := 0; j < sizeOfMatrix; j++ {
			f, err := strconv.ParseFloat(line[j], 32)
			if err == nil {
				a[i][j] = float64(f)
			}
		}
	}
	matrix.Tables = a
	return matrix
}

func sequentialCalculate(matrixQueue *newstructs.Queue, path string, sizeOfMatrix int, numberOfMatrix int) {
	var tempPath string
	for i := 0; i < numberOfMatrix; i++ {
		tempPath = path
		f, err := os.Open(tempPath + strconv.Itoa(i))
		if iooperation.IsError(err) {
			return
		}

		b1 := make([]byte, 99999999)
		_, err = f.Read(b1)
		if iooperation.IsError(err) {
			return
		}

		x := bytes.Trim(b1, "\x00")
		res := string(x)
		res = strings.TrimSpace(res)
		res2 := res[2 : len(res)-2]

		newMatrix := collectMatrixFromFile(res2, sizeOfMatrix)
		matrixQueue.M = append(matrixQueue.M, newMatrix)
	}
}

func parallelGenerateMatrixs(done chan string, start int, end int, sizeOfMatrix int, path string) {
	var tempPath string
	flag := make(chan string, end-start)
	for i := start; i < end; i++ {
		x := GenerateRandomMatrix(sizeOfMatrix, 50.0)
		tempPath = path + strconv.Itoa(i)
		iooperation.CreateFile(tempPath)
		iooperation.WriteFile(x.Tables, tempPath)
		flag <- "a matrix is created"
	}
	for i := 0; i < len(flag); i++ {
		<-flag
	}
	done <- "a group of matrixes is created"
}

func parallelCalaulate(done chan string, start int, end int, matrixQueue *newstructs.Queue, path string, sizeOfMatrix int, mtx *sync.Mutex) {
	var tempPath string
	var queueOfThisGoroutine newstructs.Queue
	flag := make(chan string, end-start)

	for i := start; i < end; i++ {
		tempPath = path
		f, err := os.Open(tempPath + strconv.Itoa(i))
		if iooperation.IsError(err) {
			return
		}
		b1 := make([]byte, 99999999)
		_, err = f.Read(b1)
		if iooperation.IsError(err) {
			return
		}
		x := bytes.Trim(b1, "\x00")
		res := string(x)
		res = strings.TrimSpace(res)
		res2 := res[2 : len(res)-2]
		newMatrix := collectMatrixFromFile(res2, sizeOfMatrix)
		queueOfThisGoroutine.M = append(queueOfThisGoroutine.M, newMatrix)
		flag <- "a matrix is collected"
	}
	for i := 0; i < len(flag); i++ {
		<-flag
	}
	//fmt.Println(queueOfThisGoroutine)
	//So here we now have a queue owned by this goroutine,
	//and stored part of all the matrixs inside, we calculate them first
	if len(queueOfThisGoroutine.M) != 0 {
		sumOfThisThread := queueOfThisGoroutine.M[0]
		for i := 1; i < len(queueOfThisGoroutine.M); i++ {
			sumOfThisThread = sumOfThisThread.Multiply(queueOfThisGoroutine.M[i])
		}
		mtx.Lock()
		matrixQueue.M = append(matrixQueue.M, sumOfThisThread)
		mtx.Unlock()
	}

	done <- "a group of matrix is collected"
}

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return quotient, remainder
}

func main() {
	arg := os.Args
	sizeOfMatrix, _ := strconv.Atoi(arg[1])
	numberOfMatrix, _ := strconv.Atoi(arg[2])

	directory := "./matrixs/"
	path := directory

	//directory := "/Users/lisihan/Documents/52060 Parallel Programming/samli/proj3/matrixs/"

	var matrixQueue newstructs.Queue
	var mtx sync.Mutex
	if len(arg) == 3 {
		iooperation.ClearDir(directory)
		sequentialCreate(numberOfMatrix, sizeOfMatrix, path)
		sequentialCalculate(&matrixQueue, path, sizeOfMatrix, numberOfMatrix)
		finalResult := matrixQueue.M[0]
		for i := 1; i < len(matrixQueue.M); i++ {
			finalResult = finalResult.Multiply(matrixQueue.M[i])
		}
		fmt.Println(finalResult)
		return
	} else {
		numberOfThreads, _ := strconv.Atoi(arg[3])
		if numberOfThreads == 1 {
			iooperation.ClearDir(directory)
			sequentialCreate(numberOfMatrix, sizeOfMatrix, path)
			sequentialCalculate(&matrixQueue, path, sizeOfMatrix, numberOfMatrix)
			finalResult := matrixQueue.M[0]
			for i := 1; i < len(matrixQueue.M); i++ {
				finalResult = finalResult.Multiply(matrixQueue.M[i])
			}
			fmt.Println(finalResult)
			return
		} else {
			interval, _ := divmod(numberOfMatrix, numberOfThreads-1)

			//Before each round of operation, we delete all the original files
			iooperation.ClearDir(directory)

			//Create matrixs parallely
			createDone := make(chan string, numberOfMatrix)
			start := 0
			end := numberOfMatrix
			for i := 0; i < numberOfThreads; i++ {
				end = int(math.Min(float64(start+interval), float64(numberOfMatrix)))
				tempStart := start
				tempEnd := end
				go parallelGenerateMatrixs(createDone, tempStart, tempEnd, sizeOfMatrix, path)
				start = end
			}
			for i := 0; i < numberOfThreads; i++ {
				<-createDone
			}

			//Collect matrix from files parallely and add them into Queue
			anotherStart := 0
			anotherEnd := numberOfMatrix
			calculateDone := make(chan string, numberOfThreads)
			for i := 0; i < numberOfThreads; i++ {
				anotherEnd = int(math.Min(float64(anotherStart+interval), float64(numberOfMatrix)))
				tempStart := anotherStart
				tempEnd := anotherEnd
				go parallelCalaulate(calculateDone, tempStart, tempEnd, &matrixQueue, path, sizeOfMatrix, &mtx)
				anotherStart = anotherEnd
			}
			for i := 0; i < numberOfThreads; i++ {
				<-calculateDone
			}
			finalResult := matrixQueue.M[0]
			for i := 1; i < len(matrixQueue.M); i++ {
				finalResult = finalResult.Multiply(matrixQueue.M[i])
			}
			fmt.Println(finalResult)
			return
		}
	}
}
