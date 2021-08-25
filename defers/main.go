package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	var function string

	fmt.Println("Which function-# to call?")
	fmt.Scanln(&function)

	switch function {
	case "0":
		deferSimple()
	case "1":
		deferFunc()
	case "2":
		deferGoroutine()
	case "3":
		deferExit()
	case "4":
		deferFatal()
	case "5":
		panicking()
	case "6":
		panicRecovery()
	default:
		fmt.Println("Give a number [0,6]")
	}

}

// Defers are pushed into a list/stack and will be processed Last In First Out

func deferSimple() {
	defer log.Println("I am last")
	defer log.Println("I am first")
}

// What are defers in functions doing??

func deferFunc() {
	defer log.Println("I am last")
	func() {
		defer log.Println("Am I second?")
	}()
	defer log.Println("I am first")
}

// Are concurrent go routines different?

func deferGoroutine() {
	done := make(chan bool)
	defer log.Println("I am last")
	go func() {
		defer log.Println("Am I second")
		done <- true
	}()
	<-done
	defer log.Println("I am first")
}

// What about os.Exit?

func deferExit() {
	done := make(chan bool)
	defer log.Println("I am last")
	go func() {
		defer log.Println("Am I second")
		os.Exit(1)
		done <- true
	}()
	<-done
	defer log.Println("I am first")
}

// What is happening here?

func deferFatal() {
	done := make(chan bool)
	defer log.Println("I am last")
	go func() {
		defer log.Println("Am I second")
		log.Fatalln("Something bad happend")
		done <- true
	}()
	<-done
	defer log.Println("I am first")
}

// What about error handling?

func panicking() {
	defer log.Println("I am last")
	done := make(chan bool)

	go func() {
		defer func() {
			log.Println("I am not second")
			if r := recover(); r != nil {
				log.Printf("Recovered a panic, %v", r)
			}
		}()
		err := errors.New("I'm gonna panic")
		if err != nil {
			log.Panic(err)
		}
		done <- true
	}()
	<-done
	defer log.Println("I am first")
}

// What about error handling?

func panicRecovery() {
	defer log.Println("I am last")
	done := make(chan bool)

	go func() {
		defer func() {
			log.Println("I am not second")
			if r := recover(); r != nil {
				log.Printf("%v: recovered a panic", r)
			}
			done <- true
		}()
		err := errors.New("I'm gonna panic")
		if err != nil {
			log.Panic(err)
		}
		// If we used another context which doesn't panic we'd have to end the
		// goroutine eventually by signaling through the channel
		// done <- true
		// otherwise we're deadlocking

	}()
	<-done
	defer log.Println("I am first")
}
