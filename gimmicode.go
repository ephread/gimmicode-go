package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/fzzy/radix/redis"
	"github.com/go-martini/martini"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"encoding/csv"
	"io"
	"io/ioutil"
	"net/http"
	"unicode/utf8"
	"log"
	"os"
)

// Substitute character that will be used in case a Unicode
// character doesn't exist in the Windows 1252 Charset.
const substituteCharacter = byte(0x1A)

func ensureDatabasePresent(client *redis.Client) {

	seeded, err := client.Cmd("EXISTS", "seeded").Bool()

	if err != nil {
		log.Fatal(err)
	}

	if seeded {
		return
	}

	// Load up the official Unicode CSV-like data.
	res, err := http.Get("http://unicode.org/Public/UNIDATA/UnicodeData.txt")
	defer res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	// Read each the names for each code point.
	reader := csv.NewReader(res.Body)

	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	for {
		record, err := reader.Read()

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		client.Cmd("HSET", record[0], "unicode_new_name", record[1])
		client.Cmd("HSET", record[0], "unicode_old_name", record[10])
	}

	client.Cmd("SET", "seeded", "true")
}

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

// Retrieve Unicode character names from pre-seeded Redis store.
func getNamesFromRedis(client *redis.Client, characterHexCode string) (oldName, newName string) {
	newName, errNewName := client.Cmd("HGET", characterHexCode, "unicode_new_name").Str()
	oldName, errOldName := client.Cmd("HGET", characterHexCode, "unicode_old_name").Str()

	if errNewName != nil || errOldName != nil {
		// Error doesn't make any difference, if it fails or doesn't exist,
		// we'll simply display nothing.
	}
	return
}

func main() {

	client, err := redis.Dial("tcp", fmt.Sprintf("%s:6379", os.Getenv("REDIS_PORT_6379_TCP_ADDR")))
	if err != nil {
		log.Fatal(err)
	}
	client.Cmd("SELECT", "2")
	defer client.Close()

	ensureDatabasePresent(client)

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

		// Retrieve unicode names from the redis store.
		oldName, newName := getNamesFromRedis(client, codePoint)

		// Windows alt-codes
		utf8Code := getUnpackedUtf8([]byte(char))
		windows1252Code := getUnpackedWindows1252([]byte(char))

		renderStruct := struct {
			Unicode, Utf8Code, Windows1252Code, OldName, NewName string
		}{
			codePoint,
			utf8Code,
			windows1252Code,
			oldName,
			newName,
		}

		r.HTML(200, "unicode", renderStruct)
	})

	m.Run()
}
