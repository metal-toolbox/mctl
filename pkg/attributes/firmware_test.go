package attributes

import (
	"testing"

	"github.com/stretchr/testify/require"
	ss "go.hollow.sh/serverservice/pkg/api/v1"
)

func TestIsFirmwareInventory(t *testing.T) {
	t.Parallel()
	t.Run("bad namespace", func(t *testing.T) {
		t.Parallel()
		attr := ss.VersionedAttributes{
			Namespace: "not.alloy.ns",
		}
		ok, fw := isFirmwareAttribute(&attr)
		require.False(t, ok)
		require.Nil(t, fw)
	})
	t.Run("bad json data", func(t *testing.T) {
		t.Parallel()
		attr := ss.VersionedAttributes{
			Namespace: alloyNamespace,
			Data:      []byte("this is invalid json"),
		}
		ok, fw := isFirmwareAttribute(&attr)
		require.False(t, ok)
		require.Nil(t, fw)
	})
	t.Run("data had no firmware", func(t *testing.T) {
		t.Parallel()
		attr := ss.VersionedAttributes{
			Namespace: alloyNamespace,
			Data:      []byte(`{"msg":"unrelated data"}`),
		}
		ok, fw := isFirmwareAttribute(&attr)
		require.False(t, ok)
		require.Nil(t, fw)
	})
	t.Run("good firmware data", func(t *testing.T) {
		t.Parallel()
		attr := ss.VersionedAttributes{
			Namespace: alloyNamespace,
			Data:      []byte(`{"firmware":{"installed":"fw-version-string"}}`),
		}
		ok, fw := isFirmwareAttribute(&attr)
		require.True(t, ok)
		require.NotNil(t, fw)
		require.Equal(t, "fw-version-string", fw.Installed)
	})
}

func TestFirmwareFromComponents(t *testing.T) {
	t.Parallel()
	cmps := []ss.ServerComponent{
		{
			Name:   "NoNameComponent",
			Vendor: "Lowest Bidder, Gmbh",
			Model:  "Excellence",
			VersionedAttributes: []ss.VersionedAttributes{
				{
					Namespace: "irrelevantNamespace",
				},
				{
					Namespace: alloyNamespace,
					Data:      []byte(`{"firmware":{"installed":"best-version"}}`),
				},
			},
		},
		{
			Name:   "NoFirmware",
			Vendor: "NoCode LLC",
			Model:  "meh",
		},
	}
	set := FirmwareFromComponents(cmps)
	require.Equal(t, 1, len(set))
	require.Equal(t, "Lowest Bidder, Gmbh", set[0].Vendor)
	require.Equal(t, "best-version", set[0].Firmware.Installed)

}
