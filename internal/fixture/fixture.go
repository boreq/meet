package fixture

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/boreq/hydro/internal/config"
	"github.com/boreq/hydro/internal/wire"

	bolt "go.etcd.io/bbolt"
)

type CleanupFunc func()

func File(t *testing.T) (string, CleanupFunc) {
	file, err := ioutil.TempFile("", "eggplant_test")
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		err := os.Remove(file.Name())
		if err != nil {
			t.Fatal(err)
		}
	}

	return file.Name(), cleanup
}

func Bolt(t *testing.T) (*bolt.DB, CleanupFunc) {
	file, fileCleanup := File(t)

	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		defer fileCleanup()

		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}

	return db, cleanup
}

func Service(t *testing.T) (wire.ComponentTestService, CleanupFunc) {
	db, dbCleanup := Bolt(t)

	conf := config.Default()

	service, err := wire.BuildComponentTestService(db, conf)
	if err != nil {
		dbCleanup()
		t.Fatal(err)
	}

	if err := service.Service.Start(); err != nil {
		dbCleanup()
		t.Fatal(err)
	}

	cleanup := func() {
		defer dbCleanup()

		err := service.Service.Close()
		if err != nil {
			t.Fatal(err)
		}

		err = service.Service.Wait()
		if err != nil {
			t.Fatal(err)
		}
	}

	return service, cleanup
}
