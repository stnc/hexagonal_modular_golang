package main

import (
	"log"
	// "flag"

	"github.com/gofiber/fiber/v3"
	// api "hexagonalapp/internal/transport/api"
	all "hexagonalapp/internal/transport/common"

	"github.com/gofiber/template/django/v4"
)

func main() {

	engine := django.New("./views", ".j2")
	engine.Reload(true)

	app := fiber.New(fiber.Config{
		Views:          engine,
		ReadBufferSize: 16 * 1024, // 16KB buffer size for reading request bodies
		// ViewsLayout: "layouts/main", //TODO: neden var amaci ne ? 
	})

	if err := all.Runner(app); err != nil {
		log.Fatal(err)
	}

}

/*
   	fmt.Println(`
   |\     /|(  ____ \( \      \__   __/|\     /|
   | )   ( || (    \/| (         ) (   ( \   / )
   | (___) || (__    | |         | |    \ (_) /
   |  ___  ||  __)   | |         | |     ) _ (
   | (   ) || (      | |         | |    / ( ) \
   | )   ( || (____/\| (____/\___) (___( /   \ )
   |/     \|(_______/(_______/\_______/|/     \|       v3.1.0

   --------------------------------------------------------------------------------------
   `)

*/
/*

	fmt.Println(`
          _______           _______  _______  _______  _        _______  _
|\     /|(  ____ \|\     /|(  ___  )(  ____ \(  ___  )( (    /|(  ___  )( \
| )   ( || (    \/( \   / )| (   ) || (    \/| (   ) ||  \  ( || (   ) || (
| (___) || (__     \ (_) / | (___) || |      | |   | ||   \ | || (___) || |
|  ___  ||  __)     ) _ (  |  ___  || | ____ | |   | || (\ \) ||  ___  || |
| (   ) || (       / ( ) \ | (   ) || | \_  )| |   | || | \   || (   ) || |
| )   ( || (____/\( /   \ )| )   ( || (___) || (___) || )  \  || )   ( || (____/\
|/     \|(_______/|/     \||/     \|(_______)(_______)|/    )_)|/     \|(_______/

--------------------------------------------------------------------------------------
`)
*/
