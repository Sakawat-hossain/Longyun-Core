package updater

// FlClash/Longyun compatibility shim.
//
// Upstream (chen08209 v0.8.94) refactored geo-database updating to the path-less
// UpdateMMDB()/UpdateASN()/UpdateGeoIp()/UpdateGeoSite() functions that write to
// the globally configured C.Path.* locations. The Longyun-Core Flutter wrapper
// (core/hub.go -> handleUpdateGeoData) still targets an explicit destination
// path resolved from the Flutter side, so we keep the original *WithPath helpers
// here. They reuse the same downloadForBytes/saveFile primitives the package
// already provides, so behaviour is identical to the pre-upgrade core and the
// wrapper/FFI surface does not change.

import (
	"fmt"

	"github.com/metacubex/mihomo/component/geodata"
	"github.com/metacubex/mihomo/component/mmdb"
	"github.com/oschwald/maxminddb-golang"
)

func UpdateMMDBWithPath(path string) (err error) {
	defer mmdb.ReloadIP()
	data, err := downloadForBytes(geodata.MmdbUrl())
	if err != nil {
		return fmt.Errorf("can't download MMDB database file: %w", err)
	}
	instance, err := maxminddb.FromBytes(data)
	if err != nil {
		return fmt.Errorf("invalid MMDB database file: %s", err)
	}
	_ = instance.Close()

	mmdb.IPInstance().Reader.Close()
	if err = saveFile(data, path); err != nil {
		return fmt.Errorf("can't save MMDB database file: %w", err)
	}
	return nil
}

func UpdateASNWithPath(path string) (err error) {
	defer mmdb.ReloadASN()
	data, err := downloadForBytes(geodata.ASNUrl())
	if err != nil {
		return fmt.Errorf("can't download ASN database file: %w", err)
	}

	instance, err := maxminddb.FromBytes(data)
	if err != nil {
		return fmt.Errorf("invalid ASN database file: %s", err)
	}
	_ = instance.Close()

	mmdb.ASNInstance().Reader.Close()
	if err = saveFile(data, path); err != nil {
		return fmt.Errorf("can't save ASN database file: %w", err)
	}
	return nil
}

func UpdateGeoIpWithPath(path string) (err error) {
	geoLoader, err := geodata.GetGeoDataLoader("standard")
	if err != nil {
		return fmt.Errorf("can't get geo data loader: %w", err)
	}
	data, err := downloadForBytes(geodata.GeoIpUrl())
	if err != nil {
		return fmt.Errorf("can't download GeoIP database file: %w", err)
	}
	if _, err = geoLoader.LoadIPByBytes(data, "cn"); err != nil {
		return fmt.Errorf("invalid GeoIP database file: %s", err)
	}
	if err = saveFile(data, path); err != nil {
		return fmt.Errorf("can't save GeoIP database file: %w", err)
	}
	return nil
}

func UpdateGeoSiteWithPath(path string) (err error) {
	geoLoader, err := geodata.GetGeoDataLoader("standard")
	if err != nil {
		return fmt.Errorf("can't get geo data loader: %w", err)
	}
	data, err := downloadForBytes(geodata.GeoSiteUrl())
	if err != nil {
		return fmt.Errorf("can't download GeoSite database file: %w", err)
	}

	if _, err = geoLoader.LoadSiteByBytes(data, "cn"); err != nil {
		return fmt.Errorf("invalid GeoSite database file: %s", err)
	}

	if err = saveFile(data, path); err != nil {
		return fmt.Errorf("can't save GeoSite database file: %w", err)
	}
	return nil
}
