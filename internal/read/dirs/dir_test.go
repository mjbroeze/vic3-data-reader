package dirs

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"vic3-data-reader/internal/env"
	"vic3-data-reader/internal/testframework/tempenv"
)

// expectedDir appends the relPath to the user's home directory
func expectedDir(relPath string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("could not get user's home dir: ", err)
	}
	return filepath.Join(home, relPath)
}

func Test_dataRootDir_default(t *testing.T) {
	expected := expectedDir(".local/share/Steam/steamapps/common/Victoria 3/game/common")
	actual, err := dataRootPath()
	if err != nil {
		t.Error("unexpected error: ", err)
	} else if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func Test_dataRootDir_custom(t *testing.T) {
	mockVic3Dir := "/my/custom/path"
	expected := "/my/custom/path/game/common"
	test := func() {
		actual, err := dataRootPath()
		if err != nil {
			t.Error("unexpected error: ", err)
		} else if expected != actual {
			t.Errorf("expected: %s, actual: %s", expected, actual)
		}
	}
	err := tempenv.Mock(t, string(env.Vic3Dir), mockVic3Dir, test)
	if err != nil {
		t.Fatal("Mock error: ", err)
	}
}

func TestDirPath_BuildingGroups(t *testing.T) {
	expected := expectedDir(".local/share/Steam/steamapps/common/Victoria 3/game/common/building_groups")
	actual, err := BuildingGroups.DirPath()
	if err != nil {
		t.Error("unexpected error: ", err)
	} else if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestDirPath_Buildings(t *testing.T) {
	expected := expectedDir(".local/share/Steam/steamapps/common/Victoria 3/game/common/buildings")
	actual, err := Buildings.DirPath()
	if err != nil {
		t.Error("unexpected error: ", err)
	} else if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestDirPath_Goods(t *testing.T) {
	expected := expectedDir(".local/share/Steam/steamapps/common/Victoria 3/game/common/goods")
	actual, err := Goods.DirPath()
	if err != nil {
		t.Error("unexpected error: ", err)
	} else if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestDirPath_ProductionMethodGroups(t *testing.T) {
	expected := expectedDir(".local/share/Steam/steamapps/common/Victoria 3/game/common/production_method_groups")
	actual, err := ProductionMethodGroups.DirPath()
	if err != nil {
		t.Error("unexpected error: ", err)
	} else if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestDirPath_ProductionMethods(t *testing.T) {
	expected := expectedDir(".local/share/Steam/steamapps/common/Victoria 3/game/common/production_methods")
	actual, err := ProductionMethods.DirPath()
	if err != nil {
		t.Error("unexpected error: ", err)
	} else if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

func TestDirPath_Technologies(t *testing.T) {
	expected := expectedDir(".local/share/Steam/steamapps/common/Victoria 3/game/common/technology/technologies")
	actual, err := Technologies.DirPath()
	if err != nil {
		t.Error("unexpected error: ", err)
	} else if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}

const (
	TestEmptyDir      DataDir = "empty"
	TestDummyOnlyDir  DataDir = "dummy-only"
	TestReadmeOnlyDir DataDir = "readme-only"
	TestSingleFileDir DataDir = "single-file"
	TestSmokeDir      DataDir = "smoke"
)

// filesTestHelper mocks the vic3 dir env variable to the full path of the testdata/mockVic3Dir dir
func filesTestHelper(t *testing.T, test func()) {
	mockPath, err := filepath.Abs("testdata/mockVic3Dir")
	if err != nil {
		t.Fatalf("could not get absolute path of mock Vic3Dir from rel path: %s", err)
	}
	_, err = os.Stat(mockPath)
	if err != nil && os.IsNotExist(err) {
		t.Fatalf("could not find mock Vic3Dir at %s", mockPath)
	} else if err != nil {
		t.Fatalf("invalid mock Vic3Dir at %s: %s", mockPath, err)
	}

	err = tempenv.Mock(t, string(env.Vic3Dir), mockPath, test)
	if err != nil {
		t.Fatalf("error mocking env variable: %s", err)
	}
}

func TestFiles_emptyDirHasNoFiles(t *testing.T) {
	mockPath, err := filepath.Abs("testdata/mockVic3Dir")
	if err != nil {
		t.Fatalf("could not get absolute path of mock Vic3Dir from rel path: %s", err)
	}

	// git doesn't allow committing a truly empty dir, so we will make one
	emptyDirPath := filepath.Join(mockPath, "game", "common", string(TestEmptyDir))
	err = os.Mkdir(emptyDirPath, 0700)
	if err != nil {
		t.Fatalf("error creating empty directory: %s", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("error removing empty directory: %s", err)
		}
	}(emptyDirPath)

	test := func() {
		files, err := TestEmptyDir.Files()
		if err != nil {
			t.Error("unexpected error: ", err)
		} else if len(files) != 0 {
			t.Errorf("expected to find no files in empty dir, found %d", len(files))
		}
	}
	filesTestHelper(t, test)
}

func TestFiles_dummyOnlyDirHasNoFiles(t *testing.T) {
	test := func() {
		files, err := TestDummyOnlyDir.Files()
		if err != nil {
			t.Error("unexpected error: ", err)
		} else if len(files) != 0 {
			t.Errorf("expected to find no valid files in dummy only dir, found %d", len(files))
		}
	}
	filesTestHelper(t, test)
}

func TestFiles_readmeOnlyDirHasNoFiles(t *testing.T) {
	test := func() {
		files, err := TestReadmeOnlyDir.Files()
		if err != nil {
			t.Error("unexpected error: ", err)
		} else if len(files) != 0 {
			t.Errorf("expected to find no valid files in readme only dir, found %d", len(files))
		}
	}
	filesTestHelper(t, test)
}

func TestFiles_singleFileDirHasOneFile(t *testing.T) {
	test := func() {
		files, err := TestSingleFileDir.Files()
		if err != nil {
			t.Error("unexpected error: ", err)
		} else if len(files) != 1 {
			t.Errorf("expected to find one file in single file dir, found %d", len(files))
		}

		fp := files[0]
		_, name := filepath.Split(string(fp))
		expectedName := "00_actual.txt"
		if expectedName != name {
			t.Errorf("expected: %s, actual: %s", name, "00_actual.txt")
		}

	}
	filesTestHelper(t, test)
}

func TestFiles_SmokeDirHasOneFile(t *testing.T) {
	test := func() {
		files, err := TestSmokeDir.Files()
		if err != nil {
			t.Error("unexpected error: ", err)
		} else if len(files) != 1 {
			t.Errorf("expected to find one valid file in smoke dir, found %d", len(files))
		}

		fp := files[0]
		_, name := filepath.Split(string(fp))
		expectedName := "01_actual.txt"
		if expectedName != name {
			t.Errorf("expected: %s, actual: %s", name, "00_actual.txt")
		}
	}
	filesTestHelper(t, test)
}
