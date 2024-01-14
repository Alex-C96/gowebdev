package database

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestNewDB(t *testing.T) {
	DB, err := NewDB(DatabaseString)
	if err != nil {
		t.Fail()
	}
	if fi, err := os.Stat(DatabaseString); err != nil {
		fmt.Println(fmt.Errorf("File path error %s", err))
	} else {
		fmt.Println(fi)
	}
	fmt.Println(DB.path)
}

func TestCreateChirp(t *testing.T) {
	DB, _ := NewDB(DatabaseString)
	chirp, err := DB.CreateChirp("This is a new chirp")
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	want := Chirp{Body: "This is a new chirp", Id: DB.ChirpId}
	if chirp.Body != want.Body {
		t.Errorf("Got %q, and wanted %q", chirp.Body, want.Body)
	}
	if chirp.Id != want.Id-1 {
		t.Errorf("Got %q, and wanted %q", chirp.Id, want.Id)
	}
}
