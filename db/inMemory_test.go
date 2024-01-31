package db

import (
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	t.Parallel() // tests can run in parallel

	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "Test Set with valid key",
			key:   "validKey",
			value: "validValue",
		},
		{
			name:  "Test Set with empty key",
			key:   "",
			value: "emptyKey",
		},
		{
			name:  "Test Set with special characters in key",
			key:   "special@Key!",
			value: "specialValue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewDB()
			db.Set(tt.key, tt.value)
			gotValue, exists := db.Get(tt.key)

			if gotValue != tt.value {
				t.Errorf("Set() = %v, want %v", gotValue, tt.value)
			}

			if exists != true {
				t.Errorf("Set() key = %v doesn't exist", tt.key)
			}

			// Wait for more than 3 minutes to check if value is deleted
			//time.Sleep(3*time.Minute + 1*time.Second)
			//gotValue, exists = db.Get(tt.key)
			//if exists != false {
			//	t.Errorf("Set() key = %v should have been deleted after 3 minutes", tt.key)
			//}

		})
	}
}

func TestGet(t *testing.T) {
	db := NewDB()

	tests := []struct {
		name      string
		setup     func()
		input     string
		wantVal   string
		wantExist bool
	}{
		{
			name: "Exists",
			setup: func() {
				db.Set("foo", "bar")
			},
			input:     "foo",
			wantVal:   "bar",
			wantExist: true,
		},
		{
			name:      "Doesn't Exist",
			setup:     func() {},
			input:     "doesnt_exist",
			wantVal:   "",
			wantExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			gotVal, gotExist := db.Get(tt.input)
			if gotVal != tt.wantVal || gotExist != tt.wantExist {
				t.Errorf("inMemoryDb.Get() = value: %v, exist: %v, want value: %v, exist: %v", gotVal, gotExist, tt.wantVal, tt.wantExist)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name string
		set  map[string]string
		del  string
		want bool
	}{
		{
			name: "delete existing key",
			set:  map[string]string{"key1": "value1", "key2": "value2"},
			del:  "key1",
			want: false,
		},
		{
			name: "delete non-existing key",
			set:  map[string]string{"key1": "value1", "key2": "value2"},
			del:  "key3",
			want: false,
		},
		{
			name: "delete from empty inMemoryDb",
			set:  map[string]string{},
			del:  "key1",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			myDB := NewDB()
			for k, v := range tt.set {
				myDB.Set(k, v)
				time.Sleep(time.Millisecond)
			}
			myDB.Delete(tt.del)
			_, got := myDB.Get(tt.del)
			if got != tt.want {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}
