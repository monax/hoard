package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/v8"
	"github.com/monax/hoard/v8/api"
	"github.com/monax/hoard/v8/client"
	"github.com/monax/hoard/v8/config"
	"github.com/monax/hoard/v8/encryption"
	"github.com/monax/hoard/v8/grant"
	"github.com/monax/hoard/v8/stores"
	"github.com/monax/hoard/v8/test/helpers"
	"google.golang.org/grpc"
)

const plaintextsDir = "plaintexts"
const storeDir = "store"
const grantsDir = "grants"
const grantExt = ".grant.json"

// We need to not use a random link nonce for this regression test to work (p.s. link nonces don't link nonces)
const linkNonce = "special-link-nonce-for-linking-nonces"

/** This program operates on a single directory <dir> (passed as first argument) with a structure:
 * <dir>/fixtures - original plaintext files for comparison and in order to reset regression tests
 * <dir>/snapshots/<snapshot>/plaintexts - containing some input data (not necessarily limited to text)
 * <dir>/snapshots/<snapshot>/store - the root directory of a Hoard filesystem store
 * <dir>/snapshots/<snapshot>/grants - a folder containing persisted grants for each file in the plaintexts directory
 *                (unless the plaintext is new for a particular run)
 *
 * Each time the program is run it walks the plaintexts directory and performs the following cycle:
 * - Saves the plaintext to Hoard obtaining a grant for that file (new grant)
 * - Tries to find a pre-saved (snapshot) grant from <dir>/output/grants corresponding to the same plaintext file (old grant)
 * - Tries to retrieve the plaintext using the old grant, falling back to the new grant if that does not exist
 * - Saves the retrieve plaintext to <dir>/output/plaintexts with the same file it was originally saved (overwriting any original)
 */
func main() {
	if len(os.Args) != 3 {
		fatalf("Pass a fixtures and snapshot directory as argument")
	}
	fixturesPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		fatalf("Could not get absolute file path for directory: %w", err)
	}
	snapshotPath, err := filepath.Abs(os.Args[2])
	if err != nil {
		fatalf("Could not get absolute file path for directory: %w", err)
	}
	secret, err := encryption.DeriveSecretKey([]byte("shhhh"), nil)
	secrets := config.SecretsManager{
		Provider: func(secretID string) (config.SymmetricSecret, error) {
			return config.SymmetricSecret{SecretKey: secret}, nil
		},
	}
	ctx := context.Background()
	store, err := stores.NewFileSystemStore(filepath.Join(snapshotPath, storeDir), base64.URLEncoding)
	if err != nil {
		fatalf("could not create FileSystemStore: %w", err)
	}
	logf("Running regression test cycle over '%s'", fixturesPath)
	err = helpers.RunWithTestServer(ctx,
		hoard.NewService(hoard.NewHoard(store,
			secrets, log.NewNopLogger()), 1024),
		func(server *grpc.Server, conn *grpc.ClientConn) error {
			cli := client.New(conn)
			// We place grants twice so that the second pass will store grants saved in the first and detect changes
			// in the plaintext or fixtures
			logf("\nPlacing grants (cycle 1/2)...")
			err := placeGrants(ctx, cli, fixturesPath, snapshotPath)
			if err != nil {
				return err
			}
			logf("\nPlacing grants (cycle 2/2)...")
			return placeGrants(ctx, cli, fixturesPath, snapshotPath)
		})
	if err != nil {
		fatalf("Error running regression tests: %w", err)
	}
}

func placeGrants(ctx context.Context, client *client.Client, fixturesPath, snapshotPath string) error {
	grantsPath := filepath.Join(snapshotPath, grantsDir)
	plaintextsPath := filepath.Join(snapshotPath, plaintextsDir)
	mkdirs(grantsPath, plaintextsPath)
	return filepath.Walk(fixturesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		name := info.Name()
		logf("\t%s:", name)
		grantFile := filepath.Join(grantsPath, name+grantExt)
		logf("\t\tuploading to hoard")
		grtToSave, err := putSealFile(ctx, client, filepath.Join(fixturesPath, name))
		if err != nil {
			return err
		}
		// Try to get the previously saved grant for this file to retrieve (from our grant snapshot dir)
		grtToRetrieve := new(grant.Grant)
		logf("\t\tchecking for existing grant")
		grtBytes, err := ioutil.ReadFile(grantFile)
		if err != nil {
			if os.IsNotExist(err) {
				logf("\t\tno grant found")
				// If the snapshot grant for this file does not exist we will just cycle this grant
				grtToRetrieve = grtToSave
			} else {
				return err
			}
		} else {
			logf("\t\tgrant found")
			err = json.Unmarshal(grtBytes, grtToRetrieve)
			if err != nil {
				return fmt.Errorf("could not unmarshal grant bytes: %w", err)
			}
		}

		logf("\t\tdownloading from hoard")
		err = unsealGetFile(ctx, client, grtToRetrieve, plaintextsPath)
		if err != nil {
			return err
		}

		grtBytes, err = json.Marshal(grtToSave)
		return ioutil.WriteFile(grantFile, grtBytes, 0600)
	})
}

func putSealFile(ctx context.Context, client *client.Client, srcPath string) (*grant.Grant, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	return client.PutSeal(ctx,
		&grant.Spec{
			LinkNonce: []byte(linkNonce),
			Symmetric: &grant.SymmetricSpec{
				PublicID: "DummySecretIsAlwaysUsed",
			},
		},
		&api.Header{
			Data: []byte(filepath.Base(srcPath)),
		},
		file)
}

func unsealGetFile(ctx context.Context, client *client.Client, grt *grant.Grant, destDir string) error {
	stream, err := client.UnsealGet(ctx, grt)
	if err != nil {
		return err
	}
	filename := stream.GetHead().GetData()
	if filename == nil {
		return fmt.Errorf("expected header metadata to contain filename but was nil")
	}
	// What could go wrong?!
	file, err := os.Create(filepath.Join(destDir, string(filename)))
	if err != nil {
		return err
	}
	_, err = stream.WriteTo(file)
	return err
}

func logf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Errorf(format, args...))
	os.Exit(1)
}

func mkdirs(dirs ...string) {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			fatalf("could not mkdir '%s': %w", dir, err)
		}
	}
}
