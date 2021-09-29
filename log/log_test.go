package log

import (
	"os"
	"sync"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	Println("hello")
	Printf("hello %s", "name")
}

func TestDefaultLoggerPanic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error(err)
		}
	}()
	Panicf("hello %s", "name")
}

func TestNew(t *testing.T) {
	file := "./test.log"
	datefile := FileNameToTodayName(file)
	logger := New(file, 0, "[test-log] ")
	logger.Println("hello")
	if _, err := os.Stat(datefile); err != nil {
		t.Errorf("file: %s not found", datefile)
	} else {
		os.Remove(datefile)
	}
}

func TestDeleteFile(t *testing.T) {
	var files = []string{
		"./test-del-2021-09-01.log",
		"./test-del-2021-09-02.log",
		"./test-del-2021-09-03.log",
		"./test-del-2021-09-04.log",
	}

	//create new file
	for _, f := range files {
		if file, err := os.Create(f); err == nil {
			defer file.Close()
		}
	}

	file := "./test-del.log"
	datefile := FileNameToTodayName(file)
	logger := New(file, 2, "[test-log] ")

	if len(logger.DeleteFiles()) != 3 {
		t.Error("delete file count:", len(logger.DeleteFiles()))
	}

	logger.Println("hello")

	if len(logger.DeleteFiles()) != 0 {
		t.Error("delete file error:", len(logger.DeleteFiles()))
	}

	//delete all files
	for _, f := range files {
		os.Remove(f)
	}
	os.Remove(datefile)
}

func TestDeleteFile2(t *testing.T) {
	var files = []string{
		"./test-del2-2021-09-01.log",
		"./test-del2-2021-09-02.log",
		"./test-del2-2021-09-03.log",
		"./test-del2-2021-09-04.log",
	}

	//create new file
	for _, f := range files {
		if file, err := os.Create(f); err == nil {
			defer file.Close()
		}
	}

	file := "./test-del2.log"
	datefile := FileNameToTodayName(file)
	logger := New(file, 2, "[test-log] ")

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			logger.Printf("test remove file2 msg: %d", i)
			wg.Done()
		}(i)
	}
	wg.Wait()

	if len(logger.DeleteFiles()) != 0 {
		t.Error("delete file error:", len(logger.DeleteFiles()))
	}

	//delete all files
	for _, f := range files {
		os.Remove(f)
	}
	os.Remove(datefile)
}

func TestRotateFile(t *testing.T) {
	file := "./test-rotate.log"
	datefile := FileNameToTodayName(file)
	logger := New(file, 2, "[test-log] ")

	logger.Println("test msg")

	os.Remove(datefile)

	logger.Println("test msg2")

	if _, err := os.Stat(datefile); err != nil {
		t.Errorf("file: %s not found", datefile)
	} else {
		os.Remove(datefile)
	}
}
