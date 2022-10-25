package main

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/textproto"
)

// https://gist.github.com/mattetti/5914158/f4d1393d83ebedc682a3c8e7bdc6b49670083b84
// https://gist.github.com/SLonger/34030e57fc49ff74f62c1b1b2f591ff0
// https://freshman.tech/snippets/go/multipart-upload-google-drive/
//
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-mime-multi
/*
Content-Type: multipart/mixed; boundary="//"
MIME-Version: 1.0

--//
Content-Type: text/cloud-config; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit
Content-Disposition: attachment; filename="cloud-config.txt"

#cloud-config
runcmd:
 - [ mkdir, /test-cloudinit ]
write_files:
 - path: /test-cloudinit/cloud-init.txt
   content: Created by cloud-init

--//
Content-Type: text/x-shellscript; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit
Content-Disposition: attachment; filename="userdata.txt"

#!/bin/bash
  mkdir test-userscript
  touch /test-userscript/userscript.txt
  echo "Created by bash shell script" >> /test-userscript/userscript.txt
--//--

*/
func upload() error {
	// New empty buffer
	body := &bytes.Buffer{}
	// Creates a new multipart Writer with a random boundary
	// writing to the empty buffer
	writer := multipart.NewWriter(body)
	err := writer.SetBoundary(`//`)
	if err != nil {
		return err
	}

	// Metadata part
	cloudInitHeader := textproto.MIMEHeader{}
	// Set the Content-Type header
	cloudInitHeader.Set("Content-Type", `text/cloud-config; charset="us-ascii"`)
	cloudInitHeader.Set("MIME-Version", `1.0`)
	cloudInitHeader.Set("Content-Transfer-Encoding", `7bit`)
	cloudInitHeader.Set("Content-Disposition", `attachment; filename="cloud-config.txt"`)

	// Create new multipart part
	cloudInitPart, err := writer.CreatePart(cloudInitHeader)
	if err != nil {
		return err
	}
	// Write the part body
	cloudInitPart.Write([]byte(`#cloud-config
runcmd:
 - [ mkdir, /test-cloudinit ]
write_files:
 - path: /test-cloudinit/cloud-init.txt
   content: Created by cloud-init
`))

	// Metadata part
	scriptHeader := textproto.MIMEHeader{}
	// Set the Content-Type header
	scriptHeader.Set("Content-Type", `text/x-shellscript; charset="us-ascii"`)
	scriptHeader.Set("MIME-Version", `1.0`)
	scriptHeader.Set("Content-Transfer-Encoding", `7bit`)
	scriptHeader.Set("Content-Disposition", `attachment; filename="userdata.txt"`)

	// Create new multipart part
	scriptPart, err := writer.CreatePart(scriptHeader)
	if err != nil {
		return err
	}
	// Write the part body
	scriptPart.Write([]byte(`#!/bin/bash
mkdir test-userscript
touch /test-userscript/userscript.txt
echo "Created by bash shell script" >> /test-userscript/userscript.txt
`))

	// Finish constructing the multipart request body
	writer.Close()

	fmt.Println(string(body.Bytes()))

	return nil
}

func main() {
	if err := upload(); err != nil {
		panic(err)
	}
}
