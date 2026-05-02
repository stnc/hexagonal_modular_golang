package stnchelper

import (
	"math/rand"

	"time"

	"github.com/gosimple/slug"
)

/*
usage
	rand.Seed(time.Now().UnixNano())

	fmt.Println(RandSlugV1(5))
*/
//RandSlugV1 random slug
func RandSlugV1(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Slugify slug genereate
func Slugify(title string, size int) string {
	slug.MaxLength = size
	return slug.MakeLang(title, "tr")
}

// GenericName for uplaod generic name
func GenericName(title string, size int) string {
	var name string
	name = Slugify(title, size)
	currentTime := time.Now()
	dateadd := currentTime.Format("15_04_05")
	return name + "_" + dateadd
}
