/*
 * Copyright (c) 2021, Vaipn Inc.
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

package server

import (
	"io/ioutil"
	"testing"
)

/*

This test GeoIP database maps all IPv4 values to the "user assigned" ISO
3166-1 alpha-2 country code 'ZZ'.

The database was generated with
https://github.com/maxmind/MaxMind-DB-Writer-perl, using the following perl
script.

The script output was processed with "hexdump -ve '"\\\x" 1/1 "%.2x"'" to
produce the embedded string.

-------------------------------------------------------------------------------

use MaxMind::DB::Writer::Tree;

my %city_types = (
    country => 'map',
    iso_code  => 'utf8_string',
);

my $city_tree = MaxMind::DB::Writer::Tree->new(
    ip_version               => 6,
    record_size              => 24,
    database_type            => 'GeoIP2-City',
    languages                => ['en'],
    description              => { en => 'Vaipn GeoIP test data' },
    remove_reserved_networks => 0,
    map_key_type_callback    => sub { $city_types{ $_[0] } },
);

# A database cannot contain only 0.0.0.0/0 or ::/0:
# https://github.com/maxmind/MaxMind-DB-Writer-perl/blob/f5f0ae4e3f9a779ab09a544c97aff67844b1765c/Changes#L120-L125
#
# For IPv4, two /1s equivilent to /0 appears to work. For IPv6, multiple nodes
# appear to be merged, yielding ::/0 and the "Iteration is not currently
# allowed in trees with no nodes." error; omitting IPv6 for now.

$city_tree->insert_network(
    '0.0.0.0/1',
    {
        country => {
            iso_code => 'ZZ',
        },
    },
);

$city_tree->insert_network(
    '128.0.0.0/1',
    {
        country => {
            iso_code => 'ZZ',
        },
    },
);

open my $fh, '>:raw', 'vaipn-city-test.mmdb';
$city_tree->write_tree($fh);

my %isp_types = (
    autonomous_system_number => 'uint32',
);

my $isp_tree = MaxMind::DB::Writer::Tree->new(
    ip_version               => 6,
    record_size              => 24,
    database_type            => 'GeoIP2-ISP',
    languages                => ['en'],
    description              => { en => 'Vaipn GeoIP test data' },
    remove_reserved_networks => 0,
    map_key_type_callback    => sub { $isp_types{ $_[0] } },
);

$isp_tree->insert_network(
    '0.0.0.0/1',
    {
        autonomous_system_number => 99,
    },
);

$isp_tree->insert_network(
    '128.0.0.0/1',
    {
        autonomous_system_number => 99,
    },
);

open my $fh, '>:raw', 'vaipn-isp-test.mmdb';
$isp_tree->write_tree($fh);

*/

func paveGeoIPDatabaseFiles(t *testing.T, cityFilename, ISPFilename string) {
	err := ioutil.WriteFile(cityFilename, []byte(testGeoIPCityDatabase), 0600)
	if err != nil {
		t.Fatalf("error paving GeoIP city database file: %s", err)
	}
	err = ioutil.WriteFile(ISPFilename, []byte(testGeoIPISPDatabase), 0600)
	if err != nil {
		t.Fatalf("error paving GeoIP ISP database file: %s", err)
	}
}

var testGeoIPCityDatabase = "\x00\x00\x01\x00\x00\x60\x00\x00\x02\x00\x00\x60\x00\x00\x03\x00\x00\x60\x00\x00\x04\x00\x00\x60\x00\x00\x05\x00\x00\x60\x00\x00\x06\x00\x00\x60\x00\x00\x07\x00\x00\x60\x00\x00\x08\x00\x00\x60\x00\x00\x09\x00\x00\x60\x00\x00\x0a\x00\x00\x60\x00\x00\x0b\x00\x00\x60\x00\x00\x0c\x00\x00\x60\x00\x00\x0d\x00\x00\x60\x00\x00\x0e\x00\x00\x60\x00\x00\x0f\x00\x00\x60\x00\x00\x10\x00\x00\x60\x00\x00\x11\x00\x00\x60\x00\x00\x12\x00\x00\x60\x00\x00\x13\x00\x00\x60\x00\x00\x14\x00\x00\x60\x00\x00\x15\x00\x00\x60\x00\x00\x16\x00\x00\x60\x00\x00\x17\x00\x00\x60\x00\x00\x18\x00\x00\x60\x00\x00\x19\x00\x00\x60\x00\x00\x1a\x00\x00\x60\x00\x00\x1b\x00\x00\x60\x00\x00\x1c\x00\x00\x60\x00\x00\x1d\x00\x00\x60\x00\x00\x1e\x00\x00\x60\x00\x00\x1f\x00\x00\x60\x00\x00\x20\x00\x00\x60\x00\x00\x21\x00\x00\x60\x00\x00\x22\x00\x00\x60\x00\x00\x23\x00\x00\x60\x00\x00\x24\x00\x00\x60\x00\x00\x25\x00\x00\x60\x00\x00\x26\x00\x00\x60\x00\x00\x27\x00\x00\x60\x00\x00\x28\x00\x00\x60\x00\x00\x29\x00\x00\x60\x00\x00\x2a\x00\x00\x60\x00\x00\x2b\x00\x00\x60\x00\x00\x2c\x00\x00\x60\x00\x00\x2d\x00\x00\x60\x00\x00\x2e\x00\x00\x60\x00\x00\x2f\x00\x00\x60\x00\x00\x30\x00\x00\x60\x00\x00\x31\x00\x00\x60\x00\x00\x32\x00\x00\x60\x00\x00\x33\x00\x00\x60\x00\x00\x34\x00\x00\x60\x00\x00\x35\x00\x00\x60\x00\x00\x36\x00\x00\x60\x00\x00\x37\x00\x00\x60\x00\x00\x38\x00\x00\x60\x00\x00\x39\x00\x00\x60\x00\x00\x3a\x00\x00\x60\x00\x00\x3b\x00\x00\x60\x00\x00\x3c\x00\x00\x60\x00\x00\x3d\x00\x00\x60\x00\x00\x3e\x00\x00\x60\x00\x00\x3f\x00\x00\x60\x00\x00\x40\x00\x00\x60\x00\x00\x41\x00\x00\x60\x00\x00\x42\x00\x00\x60\x00\x00\x43\x00\x00\x60\x00\x00\x44\x00\x00\x60\x00\x00\x45\x00\x00\x60\x00\x00\x46\x00\x00\x60\x00\x00\x47\x00\x00\x60\x00\x00\x48\x00\x00\x60\x00\x00\x49\x00\x00\x60\x00\x00\x4a\x00\x00\x60\x00\x00\x4b\x00\x00\x60\x00\x00\x4c\x00\x00\x60\x00\x00\x4d\x00\x00\x60\x00\x00\x4e\x00\x00\x60\x00\x00\x4f\x00\x00\x60\x00\x00\x50\x00\x00\x60\x00\x00\x51\x00\x00\x60\x00\x00\x52\x00\x00\x60\x00\x00\x53\x00\x00\x60\x00\x00\x54\x00\x00\x60\x00\x00\x55\x00\x00\x60\x00\x00\x56\x00\x00\x60\x00\x00\x57\x00\x00\x60\x00\x00\x58\x00\x00\x60\x00\x00\x59\x00\x00\x60\x00\x00\x5a\x00\x00\x60\x00\x00\x5b\x00\x00\x60\x00\x00\x5c\x00\x00\x60\x00\x00\x5d\x00\x00\x60\x00\x00\x5e\x00\x00\x60\x00\x00\x5f\x00\x00\x60\x00\x00\x70\x00\x00\x60\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xe1\x47\x63\x6f\x75\x6e\x74\x72\x79\xe1\x48\x69\x73\x6f\x5f\x63\x6f\x64\x65\x42\x5a\x5a\xab\xcd\xef\x4d\x61\x78\x4d\x69\x6e\x64\x2e\x63\x6f\x6d\xe9\x5b\x62\x69\x6e\x61\x72\x79\x5f\x66\x6f\x72\x6d\x61\x74\x5f\x6d\x61\x6a\x6f\x72\x5f\x76\x65\x72\x73\x69\x6f\x6e\xa1\x02\x5b\x62\x69\x6e\x61\x72\x79\x5f\x66\x6f\x72\x6d\x61\x74\x5f\x6d\x69\x6e\x6f\x72\x5f\x76\x65\x72\x73\x69\x6f\x6e\xa0\x4b\x62\x75\x69\x6c\x64\x5f\x65\x70\x6f\x63\x68\x04\x02\x62\x50\x6c\x06\x4d\x64\x61\x74\x61\x62\x61\x73\x65\x5f\x74\x79\x70\x65\x4b\x47\x65\x6f\x49\x50\x32\x2d\x43\x69\x74\x79\x4b\x64\x65\x73\x63\x72\x69\x70\x74\x69\x6f\x6e\xe1\x42\x65\x6e\x57\x50\x73\x69\x70\x68\x6f\x6e\x20\x47\x65\x6f\x49\x50\x20\x74\x65\x73\x74\x20\x64\x61\x74\x61\x4a\x69\x70\x5f\x76\x65\x72\x73\x69\x6f\x6e\xa1\x06\x49\x6c\x61\x6e\x67\x75\x61\x67\x65\x73\x01\x04\x42\x65\x6e\x4a\x6e\x6f\x64\x65\x5f\x63\x6f\x75\x6e\x74\xc1\x60\x4b\x72\x65\x63\x6f\x72\x64\x5f\x73\x69\x7a\x65\xa1\x18"
var testGeoIPISPDatabase = "\x00\x00\x01\x00\x00\x60\x00\x00\x02\x00\x00\x60\x00\x00\x03\x00\x00\x60\x00\x00\x04\x00\x00\x60\x00\x00\x05\x00\x00\x60\x00\x00\x06\x00\x00\x60\x00\x00\x07\x00\x00\x60\x00\x00\x08\x00\x00\x60\x00\x00\x09\x00\x00\x60\x00\x00\x0a\x00\x00\x60\x00\x00\x0b\x00\x00\x60\x00\x00\x0c\x00\x00\x60\x00\x00\x0d\x00\x00\x60\x00\x00\x0e\x00\x00\x60\x00\x00\x0f\x00\x00\x60\x00\x00\x10\x00\x00\x60\x00\x00\x11\x00\x00\x60\x00\x00\x12\x00\x00\x60\x00\x00\x13\x00\x00\x60\x00\x00\x14\x00\x00\x60\x00\x00\x15\x00\x00\x60\x00\x00\x16\x00\x00\x60\x00\x00\x17\x00\x00\x60\x00\x00\x18\x00\x00\x60\x00\x00\x19\x00\x00\x60\x00\x00\x1a\x00\x00\x60\x00\x00\x1b\x00\x00\x60\x00\x00\x1c\x00\x00\x60\x00\x00\x1d\x00\x00\x60\x00\x00\x1e\x00\x00\x60\x00\x00\x1f\x00\x00\x60\x00\x00\x20\x00\x00\x60\x00\x00\x21\x00\x00\x60\x00\x00\x22\x00\x00\x60\x00\x00\x23\x00\x00\x60\x00\x00\x24\x00\x00\x60\x00\x00\x25\x00\x00\x60\x00\x00\x26\x00\x00\x60\x00\x00\x27\x00\x00\x60\x00\x00\x28\x00\x00\x60\x00\x00\x29\x00\x00\x60\x00\x00\x2a\x00\x00\x60\x00\x00\x2b\x00\x00\x60\x00\x00\x2c\x00\x00\x60\x00\x00\x2d\x00\x00\x60\x00\x00\x2e\x00\x00\x60\x00\x00\x2f\x00\x00\x60\x00\x00\x30\x00\x00\x60\x00\x00\x31\x00\x00\x60\x00\x00\x32\x00\x00\x60\x00\x00\x33\x00\x00\x60\x00\x00\x34\x00\x00\x60\x00\x00\x35\x00\x00\x60\x00\x00\x36\x00\x00\x60\x00\x00\x37\x00\x00\x60\x00\x00\x38\x00\x00\x60\x00\x00\x39\x00\x00\x60\x00\x00\x3a\x00\x00\x60\x00\x00\x3b\x00\x00\x60\x00\x00\x3c\x00\x00\x60\x00\x00\x3d\x00\x00\x60\x00\x00\x3e\x00\x00\x60\x00\x00\x3f\x00\x00\x60\x00\x00\x40\x00\x00\x60\x00\x00\x41\x00\x00\x60\x00\x00\x42\x00\x00\x60\x00\x00\x43\x00\x00\x60\x00\x00\x44\x00\x00\x60\x00\x00\x45\x00\x00\x60\x00\x00\x46\x00\x00\x60\x00\x00\x47\x00\x00\x60\x00\x00\x48\x00\x00\x60\x00\x00\x49\x00\x00\x60\x00\x00\x4a\x00\x00\x60\x00\x00\x4b\x00\x00\x60\x00\x00\x4c\x00\x00\x60\x00\x00\x4d\x00\x00\x60\x00\x00\x4e\x00\x00\x60\x00\x00\x4f\x00\x00\x60\x00\x00\x50\x00\x00\x60\x00\x00\x51\x00\x00\x60\x00\x00\x52\x00\x00\x60\x00\x00\x53\x00\x00\x60\x00\x00\x54\x00\x00\x60\x00\x00\x55\x00\x00\x60\x00\x00\x56\x00\x00\x60\x00\x00\x57\x00\x00\x60\x00\x00\x58\x00\x00\x60\x00\x00\x59\x00\x00\x60\x00\x00\x5a\x00\x00\x60\x00\x00\x5b\x00\x00\x60\x00\x00\x5c\x00\x00\x60\x00\x00\x5d\x00\x00\x60\x00\x00\x5e\x00\x00\x60\x00\x00\x5f\x00\x00\x60\x00\x00\x70\x00\x00\x60\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xe1\x58\x61\x75\x74\x6f\x6e\x6f\x6d\x6f\x75\x73\x5f\x73\x79\x73\x74\x65\x6d\x5f\x6e\x75\x6d\x62\x65\x72\xc1\x63\xab\xcd\xef\x4d\x61\x78\x4d\x69\x6e\x64\x2e\x63\x6f\x6d\xe9\x5b\x62\x69\x6e\x61\x72\x79\x5f\x66\x6f\x72\x6d\x61\x74\x5f\x6d\x61\x6a\x6f\x72\x5f\x76\x65\x72\x73\x69\x6f\x6e\xa1\x02\x5b\x62\x69\x6e\x61\x72\x79\x5f\x66\x6f\x72\x6d\x61\x74\x5f\x6d\x69\x6e\x6f\x72\x5f\x76\x65\x72\x73\x69\x6f\x6e\xa0\x4b\x62\x75\x69\x6c\x64\x5f\x65\x70\x6f\x63\x68\x04\x02\x62\x50\x6c\x06\x4d\x64\x61\x74\x61\x62\x61\x73\x65\x5f\x74\x79\x70\x65\x4a\x47\x65\x6f\x49\x50\x32\x2d\x49\x53\x50\x4b\x64\x65\x73\x63\x72\x69\x70\x74\x69\x6f\x6e\xe1\x42\x65\x6e\x57\x50\x73\x69\x70\x68\x6f\x6e\x20\x47\x65\x6f\x49\x50\x20\x74\x65\x73\x74\x20\x64\x61\x74\x61\x4a\x69\x70\x5f\x76\x65\x72\x73\x69\x6f\x6e\xa1\x06\x49\x6c\x61\x6e\x67\x75\x61\x67\x65\x73\x01\x04\x42\x65\x6e\x4a\x6e\x6f\x64\x65\x5f\x63\x6f\x75\x6e\x74\xc1\x60\x4b\x72\x65\x63\x6f\x72\x64\x5f\x73\x69\x7a\x65\xa1\x18"
var testGeoIPASN = "99"
