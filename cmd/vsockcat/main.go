package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	libvirt "github.com/libvirt/libvirt-go"
	"github.com/mdlayher/vsock"
)

type VirtDomain struct {
	XMLName xml.Name  `xml:"domain"`
	Name    string    `xml:"name"`
	Vsock   VirtVsock `xml:"devices>vsock"`
}

type VirtVsock struct {
	XMLName xml.Name     `xml:"vsock"`
	CID     VirtVsockCID `xml:"cid"`
}

type VirtVsockCID struct {
	XMLName xml.Name `xml:"cid"`
	Address uint32   `xml:"address,attr"`
}

var (
	libvirtUri string
	domainName string
	port       uint32
)

func init() {
	flag.StringVar(&libvirtUri, "c", "qemu:///system", "Libvirt URI")
	flag.Parse()

	domainName = flag.Arg(0)
	if len(domainName) <= 0 {
		fmt.Fprintf(os.Stderr, "A libvirt domain must be provided.\n")
		os.Exit(1)
	}

	parsedPort, err := strconv.ParseUint(flag.Arg(1), 10, 32)
	if err != nil {
		panic(err)
	}
	port = uint32(parsedPort)
}

func main() {
	conn, err := libvirt.NewConnectReadOnly(libvirtUri)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	dom, err := conn.LookupDomainByName(domainName)
	if err != nil {
		panic(err)
	}
	defer dom.Free()

	active, err := dom.IsActive()
	if err != nil {
		panic(err)
	} else if !active {
		fmt.Fprintf(os.Stderr, "Libvirt domain is not running.\n")
		os.Exit(2)
	}

	xmlDesc, err := dom.GetXMLDesc(0)
	if err != nil {
		panic(err)
	}

	domInfo := VirtDomain{}
	err = xml.Unmarshal([]byte(xmlDesc), &domInfo)
	if err != nil {
		panic(err)
	}

	c, err := vsock.Dial(domInfo.Vsock.CID.Address, port, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	go func() {
		io.Copy(c, os.Stdin)
	}()
	io.Copy(os.Stdout, c)
}
