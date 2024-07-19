package migres

import (
	"errors"
	"testing"

	"github.com/annybs/go-version"
)

func TestModule_Upgrade(t *testing.T) {
	type TestCase struct {
		Input Module
		Up    bool
		From  string
		To    string

		Expected []string
		Err      error
	}

	testOutput := []string{}

	testMigration := func(v string, ok bool) Migration {
		f := func() error {
			if ok {
				testOutput = append(testOutput, v)
				return nil
			}
			return errors.New("test")
		}
		return Func(f, f)
	}

	testCases := []TestCase{
		{
			Input: Module{
				"1.0.0": testMigration("1.0.0", true),
				"2.0.0": testMigration("2.0.0", true),
				"3.0.0": testMigration("3.0.0", true),
			},
			Up:   true,
			From: "0",
			To:   "3",

			Expected: []string{"1.0.0", "2.0.0", "3.0.0"},
		},
		{
			Input: Module{
				"1.0.0": testMigration("1.0.0", true),
				"2.0.0": testMigration("2.0.0", true),
				"3.0.0": testMigration("3.0.0", true),
			},
			Up:   true,
			From: "1",
			To:   "3",

			Expected: []string{"2.0.0", "3.0.0"},
		},
		{
			Input: Module{
				"1.0.0": testMigration("1.0.0", true),
				"2.0.0": testMigration("2.0.0", true),
				"3.0.0": testMigration("3.0.0", true),
			},
			From: "3",
			To:   "0",

			Expected: []string{"3.0.0", "2.0.0", "1.0.0"},
		},
		{
			Input: Module{
				"1.0.0": testMigration("1.0.0", true),
				"2.0.0": testMigration("2.0.0", true),
				"3.0.0": testMigration("3.0.0", true),
			},
			From: "3",
			To:   "2",

			Expected: []string{"3.0.0"},
		},
		{
			Input: Module{
				"1.0.0": testMigration("1.0.0", true),
				"2.0.0": testMigration("2.0.0", false),
				"3.0.0": testMigration("3.0.0", true),
			},
			Up:   true,
			From: "0",
			To:   "3",

			Err: failMigration(errors.New("test"), version.MustParse("2.0.0"), version.MustParse("1.0.0")),
		},
		{
			Input: Module{
				"1.0.0": testMigration("1.0.0", true),
				"2.0.0": testMigration("2.0.0", false),
				"3.0.0": testMigration("3.0.0", true),
			},
			From: "4",
			To:   "1",

			Err: failMigration(errors.New("test"), version.MustParse("2.0.0"), version.MustParse("3.0.0")),
		},
		{
			Input: Module{
				"1.0.0": testMigration("1.0.0", true),
				"2.0.0": testMigration("2.0.0", true),
				"3.0.0": testMigration("3.0.0", false),
			},
			From: "4",
			To:   "1",

			Err: failMigration(errors.New("test"), version.MustParse("3.0.0"), nil),
		},
	}

	for i, testCase := range testCases {
		testOutput = []string{}

		var err error
		if testCase.Up {
			err = testCase.Input.Upgrade(testCase.From, testCase.To)
		} else {
			err = testCase.Input.Downgrade(testCase.From, testCase.To)
		}

		if err != nil {
			if testCase.Err == nil {
				t.Errorf("test %d failed (expected nil error, got error %v)", i, err)
			} else if !errors.Is(err, testCase.Err) {
				t.Errorf("test %d failed (expected error %v, got error %v)", i, testCase.Err, err)
			} else {
				a := err.(*Error)
				e := testCase.Err.(*Error)
				if !a.Version.Equal(e.Version) {
					t.Errorf("test %d failed (expected error.Version %s, got error.Version %s)", i, e.Version, a.Version)
				} else if !a.LastVersion.Equal(e.LastVersion) {
					t.Errorf("test %d failed (expected error.LastVersion %s, got error.LastVersion %s)", i, e.LastVersion, a.LastVersion)
				} else {
					t.Logf("test %d passed (expected error %v, got error %v)", i, testCase.Err, err)
				}
			}
			continue
		} else if testCase.Err != nil {
			t.Errorf("test %d failed (expected error %v, got nil error)", i, testCase.Err)
			continue
		}

		if len(testOutput) != len(testCase.Expected) {
			t.Errorf("test %d failed (expected %v, got %v)", i, testCase.Expected, testOutput)
			continue
		}
		ok := true
		for j, v := range testOutput {
			if v != testCase.Expected[j] {
				t.Errorf("test %d failed (expected version %q at index %d, got version %q)", i, testCase.Expected[j], j, v)
				ok = false
			}
		}
		if ok {
			t.Logf("test %d passed (expected %v, got %v)", i, testCase.Expected, testOutput)
		}
	}
}
