package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/DusanKasan/parsemail"
)

// EXPORT_PATH=`defaults read com.freron.MailMate MmExportPath 2>/dev/null`
// if [ -z "${EXPORT_PATH}" ]; then
// 	EXPORT_PATH=`"${MM_BUNDLE_SUPPORT}/bin/select_export_folder"`
// fi
// if [ -z "${EXPORT_PATH}" ]; then
//    EXPORT_PATH="${HOME}/Desktop/MailMate Export"
// fi

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func main() {
	folder, err := exec.Command("defaults", "read", "com.freron.MailMate", "MmAttachmentsExportPath").Output()

	folderstr := string(folder)

	if err != nil || (folderstr == "") {
		prefix := os.Getenv("MM_BUNDLE_SUPPORT")
		folder, err = exec.Command(fmt.Sprintf("%s/bin/select_export_folder", prefix)).Output()
	}

	folderstr = string(folder)

	if err != nil || folderstr == "" {
		home := os.Getenv("HOME")
		folderstr = fmt.Sprintf("%s/Desktop/MailMate Export", home)
	}

	reader := bufio.NewReader(os.Stdin)
	email, err := parsemail.Parse(reader)
	check(err)

	for _, a := range email.Attachments {
		err := os.MkdirAll(folderstr, 0755)
		check(err)
		err = ioutil.WriteFile(path.Join(folderstr, a.Filename), StreamToByte(a.Data), 0644)
		check(err)
	}

}
