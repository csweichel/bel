package bel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/go-test/deep"
)

func TestRoundtrip(t *testing.T) {
	if !isNpxAvailable() {
		t.Skip("npm is not available - skipping round trip")
		return
	}

	ws, err := ioutil.TempDir("", "")
	if !installTsNode(t, ws) {
		return
	}
	defer os.RemoveAll(ws)

	testdata := MyTestStruct{
		StringField: "hello world",
		// OptionalField      string `json:",omitempty"`
		NamedField: 42,
		// NamedOptionalField int32  `json:"thisIsOptional,omitempty"`
		SkipThisField: "has-a-value",
		Containment: AnotherTestStruct{
			Bar: true,
			Foo: "baz",
		},
		Referece: &AnotherTestStruct{
			Bar: false,
			Foo: "abc",
		},
	}

	if err != nil {
		t.Error(err)
		return
	}
	if !generateTypescript(t, ws, testdata) {
		return
	}

	fromTS := executeTypescript(t, ws)
	if fromTS == nil {
		return
	}

	testdata.SkipThisField = ""
	diff := deep.Equal(&testdata, fromTS)
	for _, d := range diff {
		t.Error(d)
	}
}

func installTsNode(t *testing.T, ws string) bool {
	cmd := exec.Command("/bin/sh", "-c", "npm install typescript ts-node")
	cmd.Dir = ws
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skip("cannot install ts-node: " + string(out))
		return false
	}
	return true
}

func generateTypescript(t *testing.T, ws string, testdata MyTestStruct) bool {
	handler, err := NewParsedSourceEnumHandler(".")
	if err != nil {
		t.Error(err)
		return false
	}

	extract, err := Extract(MyTestStruct{},
		WithEnumHandler(handler),
		FollowStructs,
	)
	if err != nil {
		t.Error(err)
		return false
	}

	f, err := os.Create(path.Join(ws, "struct.ts"))
	if err != nil {
		t.Error(err)
		return false
	}

	err = Render(extract, GenerateOutputTo(f))
	if err != nil {
		f.Close()
		t.Error(err)
		return false
	}
	f.Close()

	testjson, err := json.Marshal(testdata)
	if err != nil {
		t.Error(err)
		return false
	}

	indexts := []byte(fmt.Sprintf(`
import { MyTestStruct } from "./struct";

const data: MyTestStruct = %s;
console.log(JSON.stringify(data));
`, string(testjson)))
	err = ioutil.WriteFile(path.Join(ws, "index.ts"), indexts, 0744)
	if err != nil {
		t.Error(err)
		return false
	}

	return true
}

func executeTypescript(t *testing.T, ws string) *MyTestStruct {
	cmd := exec.Command("./node_modules/.bin/ts-node", "index.ts")
	cmd.Dir = ws

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(string(out))
		return nil
	}

	result := MyTestStruct{}
	err = json.Unmarshal(out, &result)
	if err != nil {
		t.Error(err)
		return nil
	}

	return &result
}

func isNpxAvailable() bool {
	cmd := exec.Command("/bin/sh", "-c", "which npx")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
