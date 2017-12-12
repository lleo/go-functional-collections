package sorted_map

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidPos(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, black,
					mknod(110, red, nil, nil),
					mknod(130, red, nil, nil)))))

	var err = m.valid()

	if err != nil {
		t.Fatal("valid map is shown as not valid")
	}
}

func TestValidNegBlackCount(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				nil)))

	var err = m.valid()

	if err == nil {
		t.Fatal("invalid map is shown as valid")
	}

	var errStr = fmt.Sprintf("%s", err)
	if !strings.Contains(errStr, "count") {
		t.Fatal("invalid map did not show a black count violation")
	}
}

func TestValidNegRedRed(t *testing.T) {
	var m = mkmap(
		mknod(60, black,
			mknod(20, black,
				mknod(10, black, nil, nil),
				mknod(40, black,
					mknod(30, red, nil, nil),
					mknod(50, red, nil, nil))),
			mknod(100, black,
				mknod(80, black,
					mknod(70, red, nil, nil),
					mknod(90, red, nil, nil)),
				mknod(120, red,
					mknod(110, red,
						mknod(105, black, nil, nil),
						mknod(115, black, nil, nil)),
					mknod(130, red,
						mknod(125, black, nil, nil),
						mknod(135, black, nil, nil))))))

	var err = m.valid()

	if err == nil {
		t.Fatal("invalid map is shown as valid")
	}

	var errStr = fmt.Sprintf("%s", err)
	if !strings.Contains(errStr, "red") {
		t.Fatal("invalid map did not show red-red violation")
	}
}

func TestValidNegNumEntries(t *testing.T) {
	var m = New()
	m.root = newNode(IntKey(10), 10)
	//m.numEnts = 0

	var err = m.valid()

	var errStr = fmt.Sprintf("%s", err)
	if !strings.Contains(errStr, "NumEntries") {
		t.Fatal("invalid map did not show incorrect NumEntries() value")
	}
}
