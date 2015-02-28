package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"unicode/utf8"
)

// Substitute character that will be used in case a Unicode
// character doesn't exist in the Windows 1252 Charset.
const substituteCharacter = byte(0x1A)

func getUnpackedWindows1252(char []byte) (code string) {
	// Here, we just used the standard library to encode
	// the received UTF-8 into Windows-1252.
	i := bytes.NewReader(char)
	o := transform.NewReader(i, charmap.Windows1252.NewEncoder())
	enc, error := ioutil.ReadAll(o)

	if error != nil {
		log.Println("Impossible to convert %s", char)
		return
	}

	// If the converted character is invalid we simply return an empty string.
	// Not that we check whether the substitute character was imputed (which
	// should never happen) or was the result of Unicode character
	// missing in the Windows-1252 charset.
	if len(enc) != 1 || (enc[0] == substituteCharacter && char[0] != enc[0]) {
		log.Printf("%s was invalid.", char)
		return
	}

	// Converts the byte into alt code.
	code = fmt.Sprintf("0%d", enc[0])
	return
}

func getUnpackedUtf8(character []byte) (hexCode string) {
	// UTF-8 code as hexadecimal.
	hexCode = hex.EncodeToString(character)
	return
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Get("/unicode", func(r render.Render, req *http.Request) {
		// Retrieve the character from URL parameters and
		// get its codepoint.
		char := req.URL.Query().Get("character")
		rune, _ := utf8.DecodeRuneInString(char)
		codePoint := fmt.Sprintf("%04X", rune)

		// Windows alt-codes
		utf8Code := getUnpackedUtf8([]byte(char))
		windows1252Code := getUnpackedWindows1252([]byte(char))

		renderStruct := struct {
			Unicode, Utf8Code, Windows1252Code, OldName, NewName string
		}{
			codePoint,
			utf8Code,
			windows1252Code,
			"",
			"",
		}

		r.HTML(200, "unicode", renderStruct)
	})

	m.Run()
}
