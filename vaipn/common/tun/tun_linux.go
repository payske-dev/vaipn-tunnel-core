/*
 * Copyright (c) 2017, Vaipn Inc.
 * All rights reserved.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package tun

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/payske-dev/vaipn-tunnel-core/vaipn/common"
	"github.com/payske-dev/vaipn-tunnel-core/vaipn/common/errors"
	"golang.org/x/sys/unix"
)

const (
	DEFAULT_PUBLIC_INTERFACE_NAME = "eth0"
)

func IsSupported() bool {
	return true
}

func makeDeviceInboundBuffer(MTU int) []byte {
	return make([]byte, MTU)
}

func makeDeviceOutboundBuffer(MTU int) []byte {
	// On Linux, no outbound buffer is used
	return nil
}

// OpenTunDevice opens a file for performing device I/O with
// either a specified tun device, or a new tun device (when
// name is "").
func OpenTunDevice(name string) (*os.File, string, error) {

	// Prevent fork between creating fd and setting CLOEXEC
	// TODO: is this still necessary with unix.Open?
	syscall.ForkLock.RLock()
	defer syscall.ForkLock.RUnlock()

	// Requires process to run as root or have CAP_NET_ADMIN

	// As explained in https://github.com/golang/go/issues/30426, the fd must
	// not be added to the Go poller before the following TUNSETIFF ioctl
	// call. This is achieved by using unix.Open -- which opens a raw fd --
	// instead of os.FileOpen, followed by the ioctl and finally os.NewFile
	// to add the fd to the Go poller.
	//
	// Set CLOEXEC so file descriptor not leaked to network config command
	// subprocesses.

	fileName := "/dev/net/tun"

	fd, err := unix.Open(fileName, os.O_RDWR|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, "", errors.Trace(err)
	}

	// This code follows snippets in this thread:
	// https://groups.google.com/forum/#!msg/golang-nuts/x_c_pZ6p95c/8T0JBZLpTwAJ;

	// Definitions from <linux/if.h>, <linux/if_tun.h>

	// Note: using IFF_NO_PI, so packets have no size/flags header. This does mean
	// that if the MTU is changed after the tun device is initialized, packets could
	// be truncated when read.

	const (
		IFNAMSIZ        = 16
		IF_REQ_PAD_SIZE = 40 - 18
		IFF_TUN         = 0x0001
		IFF_NO_PI       = 0x1000
	)

	var ifName [IFNAMSIZ]byte
	if name == "" {
		copy(ifName[:], []byte("tun%d"))
	} else {
		copy(ifName[:], []byte(name))
	}

	ifReq := struct {
		name  [IFNAMSIZ]byte
		flags uint16
		pad   [IF_REQ_PAD_SIZE]byte
	}{
		ifName,
		uint16(IFF_TUN | IFF_NO_PI),
		[IF_REQ_PAD_SIZE]byte{},
	}

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TUNSETIFF),
		uintptr(unsafe.Pointer(&ifReq)))
	if errno != 0 {
		unix.Close(fd)
		return nil, "", errors.Trace(errno)
	}

	err = unix.SetNonblock(fd, true)
	if err != nil {
		unix.Close(fd)
		return nil, "", errors.Trace(err)
	}

	file := os.NewFile(uintptr(fd), fileName)

	deviceName := strings.Trim(string(ifReq.name[:]), "\x00")

	return file, deviceName, nil
}

func (device *Device) readTunPacket() (int, int, error) {

	// Assumes MTU passed to makeDeviceInboundBuffer is actual MTU and
	// so buffer is sufficiently large to always read a complete packet.

	n, err := device.deviceIO.Read(device.inboundBuffer)
	if err != nil {
		return 0, 0, errors.Trace(err)
	}
	return 0, n, nil
}

func (device *Device) writeTunPacket(packet []byte) error {

	// Doesn't need outboundBuffer since there's no header; write directly to device.

	_, err := device.deviceIO.Write(packet)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func resetNATTables(
	config *ServerConfig,
	IPAddress net.IP) error {

	// Uses the "conntrack" command, which is often not installed by default.

	// conntrack --delete -src-nat --orig-src <address> will clear NAT tables of existing
	// connections, making it less likely that traffic for a previous client using the
	// specified address will be forwarded to a new client using this address. This is in
	// the already unlikely event that there's still in-flight traffic when the address is
	// recycled.

	err := common.RunNetworkConfigCommand(
		config.Logger,
		config.SudoNetworkConfigCommands,
		"conntrack",
		"--delete",
		"--src-nat",
		"--orig-src",
		IPAddress.String())
	if err != nil {

		// conntrack exits with this error message when there are no flows
		// to delete, which is not a failure condition.
		if strings.Contains(err.Error(), "0 flow entries have been deleted") {
			return nil
		}

		return errors.Trace(err)
	}

	return nil
}

func configureServerInterface(
	config *ServerConfig,
	tunDeviceName string) error {

	// Set tun device network addresses and MTU

	IPv4Address, IPv4Netmask, err := splitIPMask(serverIPv4AddressCIDR)
	if err != nil {
		return errors.Trace(err)
	}

	err = common.RunNetworkConfigCommand(
		config.Logger,
		config.SudoNetworkConfigCommands,
		"ifconfig",
		tunDeviceName,
		IPv4Address, "netmask", IPv4Netmask,
		"mtu", strconv.Itoa(getMTU(config.MTU)),
		"up")
	if err != nil {
		return errors.Trace(err)
	}

	err = common.RunNetworkConfigCommand(
		config.Logger,
		config.SudoNetworkConfigCommands,
		"ifconfig",
		tunDeviceName,
		"add", serverIPv6AddressCIDR)
	if err != nil {
		if config.AllowNoIPv6NetworkConfiguration {
			config.Logger.WithTraceFields(
				common.LogFields{
					"error": err}).Warning(
				"assign IPv6 address failed")
		} else {
			return errors.Trace(err)
		}
	}

	egressInterface := config.EgressInterface
	if egressInterface == "" {
		egressInterface = DEFAULT_PUBLIC_INTERFACE_NAME
	}

	// NAT tun device to external interface

	// TODO: need only set forwarding for specific interfaces?

	err = common.RunNetworkConfigCommand(
		config.Logger,
		config.SudoNetworkConfigCommands,
		"sysctl",
		"net.ipv4.conf.all.forwarding=1")
	if err != nil {
		return errors.Trace(err)
	}

	err = common.RunNetworkConfigCommand(
		config.Logger,
		config.SudoNetworkConfigCommands,
		"sysctl",
		"net.ipv6.conf.all.forwarding=1")
	if err != nil {
		if config.AllowNoIPv6NetworkConfiguration {
			config.Logger.WithTraceFields(
				common.LogFields{
					"error": err}).Warning(
				"allow IPv6 forwarding failed")
		} else {
			return errors.Trace(err)
		}
	}

	// To avoid duplicates, first try to drop existing rule, then add

	for _, mode := range []string{"-D", "-A"} {

		err = common.RunNetworkConfigCommand(
			config.Logger,
			config.SudoNetworkConfigCommands,
			"iptables",
			"-t", "nat",
			mode, "POSTROUTING",
			"-s", privateSubnetIPv4.String(),
			"-o", egressInterface,
			"-j", "MASQUERADE")
		if mode != "-D" && err != nil {
			return errors.Trace(err)
		}

		err = common.RunNetworkConfigCommand(
			config.Logger,
			config.SudoNetworkConfigCommands,
			"ip6tables",
			"-t", "nat",
			mode, "POSTROUTING",
			"-s", privateSubnetIPv6.String(),
			"-o", egressInterface,
			"-j", "MASQUERADE")
		if mode != "-D" && err != nil {
			if config.AllowNoIPv6NetworkConfiguration {
				config.Logger.WithTraceFields(
					common.LogFields{
						"error": err}).Warning(
					"configure IPv6 masquerading failed")
			} else {
				return errors.Trace(err)
			}
		}
	}

	return nil
}

func configureClientInterface(
	config *ClientConfig,
	tunDeviceName string) error {

	// Set tun device network addresses and MTU

	IPv4Address, IPv4Netmask, err := splitIPMask(config.IPv4AddressCIDR)
	if err != nil {
		return errors.Trace(err)
	}

	err = common.RunNetworkConfigCommand(
		config.Logger,
		config.SudoNetworkConfigCommands,
		"ifconfig",
		tunDeviceName,
		IPv4Address,
		"netmask", IPv4Netmask,
		"mtu", strconv.Itoa(getMTU(config.MTU)),
		"up")
	if err != nil {
		return errors.Trace(err)
	}

	err = common.RunNetworkConfigCommand(
		config.Logger,
		config.SudoNetworkConfigCommands,
		"ifconfig",
		tunDeviceName,
		"add", config.IPv6AddressCIDR)
	if err != nil {
		if config.AllowNoIPv6NetworkConfiguration {
			config.Logger.WithTraceFields(
				common.LogFields{
					"error": err}).Warning(
				"assign IPv6 address failed")
		} else {
			return errors.Trace(err)
		}
	}

	// Set routing. Routes set here should automatically
	// drop when the tun device is removed.

	// TODO: appear to need explicit routing only for IPv6?

	for _, destination := range config.RouteDestinations {

		// Destination may be host (IP) or network (CIDR)

		IP := net.ParseIP(destination)
		if IP == nil {
			var err error
			IP, _, err = net.ParseCIDR(destination)
			if err != nil {
				return errors.Trace(err)
			}
		}
		if IP.To4() != nil {
			continue
		}

		// Note: use "replace" instead of "add" as route from
		// previous run (e.g., tun_test case) may not yet be cleared.

		err = common.RunNetworkConfigCommand(
			config.Logger,
			config.SudoNetworkConfigCommands,
			"ip",
			"-6",
			"route", "replace",
			destination,
			"dev", tunDeviceName)
		if err != nil {
			if config.AllowNoIPv6NetworkConfiguration {
				config.Logger.WithTraceFields(
					common.LogFields{
						"error": err}).Warning("add IPv6 route failed")
			} else {
				return errors.Trace(err)
			}
		}
	}

	return nil
}

// BindToDevice binds a socket to the specified interface.
func BindToDevice(fd int, deviceName string) error {
	err := syscall.BindToDevice(fd, deviceName)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func fixBindToDevice(logger common.Logger, useSudo bool, tunDeviceName string) error {

	// Fix the problem described here:
	// https://stackoverflow.com/questions/24011205/cant-perform-tcp-handshake-through-a-nat-between-two-nics-with-so-bindtodevice/
	//
	// > the linux kernel is configured on certain mainstream distributions
	// > (Ubuntu...) to act as a router and drop packets where the source
	// > address is suspect in order to prevent spoofing (search "rp_filter" on
	// > https://www.kernel.org/doc/Documentation/networking/ip-sysctl.txt and
	// > RFC3704)

	err := common.RunNetworkConfigCommand(
		logger,
		useSudo,
		"sysctl",
		"net.ipv4.conf.all.accept_local=1")
	if err != nil {
		return errors.Trace(err)
	}

	err = common.RunNetworkConfigCommand(
		logger,
		useSudo,
		"sysctl",
		"net.ipv4.conf.all.rp_filter=0")
	if err != nil {
		return errors.Trace(err)
	}

	err = common.RunNetworkConfigCommand(
		logger,
		useSudo,
		"sysctl",
		fmt.Sprintf("net.ipv4.conf.%s.rp_filter=0", tunDeviceName))
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}
