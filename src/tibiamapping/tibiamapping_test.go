package tibiamapping

import (
	"errors"
	"testing"
)

func TestFakeTest(t *testing.T) {
	if tibiaAssetsDataMinJsonURL != "https://assets.tibiadata.com/data.min.json" {
		t.Error(errors.New("tibiaAssetsDataMinJsonURL is wrong"))
	}

	if tibiaAssetsSha256SumURL != "https://assets.tibiadata.com/sha256sum.txt" {
		t.Error(errors.New("tibiaAssetsSha256SumURL is wrong"))
	}

	if tibiaAssetsSha512SumURL != "https://assets.tibiadata.com/sha512sum.txt" {
		t.Error(errors.New("tibiaAssetsSha512SumURL is wrong"))
	}
}

func TestInitiator(t *testing.T) {
	mapping, err := Run("TibiaData-API/v3/testing")
	if err != nil {
		t.Fatal(err)
	}

	if len(mapping.RawData) == 0 {
		t.Fatal(errors.New("RawData is empty"))
	}

	if mapping.Sha256Sum == "" {
		t.Fatal(errors.New("Sha256Sum is empty"))
	}

	if mapping.Sha512Sum == "" {
		t.Fatal(errors.New("Sha512Sum is empty"))
	}
}
