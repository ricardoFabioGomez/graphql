package cache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Get_From_cache(t *testing.T) {
	type (
		args struct {
			astDoc string
		}
		testcase struct {
			name string
			args args
		}
	)
	testscases := []testcase{
		{
			name: "document is empty",
			args: args{
				astDoc: "",
			},
		},
		{
			name: "document is nil",
			args: args{
				astDoc: "  ",
			},
		},
	}

	for _, tc := range testscases {
		t.Run(tc.name, func(t *testing.T) {
			cache := NewCache()
			cache.Set(tc.args.astDoc)
			if !cache.Get(tc.args.astDoc) {
				t.Error("Get() should return true")
			}
		})
	}
}

func Test_Get_From_cache_success(t *testing.T) {
	ass := assert.New(t)
	query := `query{
		getUser(id: %s){
			id
			name
		}
	}`
	cache := NewCache()
	cache.Set(fmt.Sprintf(query, "1"))
	ass.True(cache.Get(fmt.Sprintf(query, "1")))
	ass.True(cache.Get(fmt.Sprintf(query, "2")))
	ass.True(cache.Get(fmt.Sprintf(query, "3")))
	ass.True(cache.Get(fmt.Sprintf(query, "5")))
}

func Test_Get_From_cache_miss(t *testing.T) {
	ass := assert.New(t)
	query := `query{
		getUser(id: %s){
			id
			name
		}
	}`
	cache := NewCache()
	ass.False(cache.Get(fmt.Sprintf(query, "1")))
	ass.False(cache.Get(fmt.Sprintf(query, "2")))
	ass.False(cache.Get(fmt.Sprintf(query, "3")))
	ass.False(cache.Get(fmt.Sprintf(query, "5")))

}
