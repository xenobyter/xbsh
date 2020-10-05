package view

import (
	"unicode/utf8"
	"time"
	"fmt"
	"testing"
)

func TestGetStatusBarContent( t *testing.T) {
	fmt.Println("getStatusBarContent should return a string with a given length")
	wLen := 100
	gLen := utf8.RuneCountInString(getStatusBarContent(time.Now(),wLen))
	if wLen != gLen {
		t.Errorf("got: %v, want: %v", gLen, wLen)
	}

	fmt.Println("getStatusBarContent should cut time first if len is too small")
	wDate := time.Now().Format("2006-01-02") + " <<<"
	c := getStatusBarContent(time.Now(),41)
	gDate := c[utf8.RuneCountInString(c)-len(wDate):]
	if wDate != gDate {
		t.Errorf("got: %v, want: %v", gDate, wDate)
	}

	fmt.Println("getStatusBarContent should cut date next if len is too small")
	wDate = " <<<"
	c = getStatusBarContent(time.Now(),14)
	gDate = c[len(c)-len(wDate):]
	if wDate != gDate {
		t.Errorf("got: %v, want: %v", gDate, wDate)
	}
	
	fmt.Println("getStatusBarContent should behave with very short length")
	wString := " <<<"
	c = getStatusBarContent(time.Now(),17)
	gString := c[len(c)-len(wDate):]
	if wString != gString {
		t.Errorf("got: %v, want: %v", gString, wString)
	}

}