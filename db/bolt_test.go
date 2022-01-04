package db

import (
	"os"
	"testing"
)

const (
	testFile   string = "./test.db"
	testBucket string = "test"
)

func TestBoltDB_Insert(t *testing.T) {
	// Create a test database
	d, err := NewBoltDB(testFile, 0600, testBucket)
	if err != nil {
		t.Fatalf("couldn't create database %s", err)
	}

	type args struct {
		shortURL string
		longURL  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Insert a key/value pair",
			args: args{
				shortURL: "/google",
				longURL:  "https://www.google.es/",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Insert(tt.args.shortURL, tt.args.longURL); (err != nil) != tt.wantErr {
				t.Errorf("BoltDB.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// Close the database
	d.Close()

	// Delete the test file
	os.Remove(testFile)
}

func TestBoltDB_GetFullURL(t *testing.T) {
	// Create a test database
	d, err := NewBoltDB(testFile, 0600, testBucket)
	if err != nil {
		t.Fatalf("couldn't create database %s", err)
	}

	// Add some records for testing
	d.Insert("/centos", "https://www.centos.org")
	d.Insert("/redhat", "https://www.redhat.com/es")

	type args struct {
		shortURL string
	}
	tests := []struct {
		name        string
		args        args
		wantLongURL string
		wantErr     bool
	}{
		{
			name: "Check that centos is retrieved accordingly",
			args: args{
				shortURL: "/centos",
			},
			wantLongURL: "https://www.centos.org",
			wantErr:     false,
		},
		{
			name: "Check that redhat is retrieved accordingly",
			args: args{
				shortURL: "/redhat",
			},
			wantLongURL: "https://www.redhat.com/es",
			wantErr:     false,
		},
		{
			name: "Try to retrieve some record that doesn't exist",
			args: args{
				shortURL: "/forocoches",
			},
			wantLongURL: "",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLongURL, err := d.GetFullURL(tt.args.shortURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("BoltDB.GetFullURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLongURL != tt.wantLongURL {
				t.Errorf("BoltDB.GetFullURL() = %v, want %v", gotLongURL, tt.wantLongURL)
			}
		})
	}

	// Close the database
	d.Close()

	// Delete the test file
	os.Remove(testFile)

}
