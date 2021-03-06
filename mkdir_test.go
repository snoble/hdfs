package hdfs

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var mode os.FileMode = 0777 | os.ModeDir

func TestMkdir(t *testing.T) {
	client := getClient(t)

	baleet(t, "/_test/dir2")

	err := client.Mkdir("/_test/dir2", mode)
	assert.Nil(t, err)

	fi, err := client.Stat("/_test/dir2")
	require.Nil(t, err)
	assert.True(t, fi.IsDir())
}

func TestMkdirExists(t *testing.T) {
	client := getClient(t)

	mkdirp(t, "/_test/existingdir")

	err := client.Mkdir("/_test/existingdir", mode)
	assertPathError(t, err, "mkdir", "/_test/existingdir", os.ErrExist)
}

func TestMkdirWithoutParent(t *testing.T) {
	client := getClient(t)

	baleet(t, "/_test/nonexistent")

	err := client.Mkdir("/_test/nonexistent/foo", mode)
	assertPathError(t, err, "mkdir", "/_test/nonexistent/foo", os.ErrNotExist)

	_, err = client.Stat("/_test/nonexistent/foo")
	assertPathError(t, err, "stat", "/_test/nonexistent/foo", os.ErrNotExist)

	_, err = client.Stat("/_test/nonexistent")
	assertPathError(t, err, "stat", "/_test/nonexistent", os.ErrNotExist)
}

func TestMkdirAll(t *testing.T) {
	client := getClient(t)

	baleet(t, "/_test/dir3")

	err := client.MkdirAll("/_test/dir3/foo", mode)
	assert.Nil(t, err)

	fi, err := client.Stat("/_test/dir3/foo")
	require.Nil(t, err)
	assert.True(t, fi.IsDir())

	fi, err = client.Stat("/_test/dir3")
	require.Nil(t, err)
	assert.True(t, fi.IsDir())
	assert.Equal(t, 0, fi.Size())
}

func TestMkdirWIthoutPermission(t *testing.T) {
	client := getClient(t)
	otherClient := getClientForUser(t, "other")

	mkdirp(t, "/_test/accessdenied")

	err := otherClient.Mkdir("/_test/accessdenied/dir", mode)
	assertPathError(t, err, "mkdir", "/_test/accessdenied/dir", os.ErrPermission)

	_, err = client.Stat("/_test/accessdenied/dir")
	assertPathError(t, err, "stat", "/_test/accessdenied/dir", os.ErrNotExist)
}

func TestMkdirAllWIthoutPermission(t *testing.T) {
	client := getClient(t)
	otherClient := getClientForUser(t, "other")

	mkdirp(t, "/_test/accessdenied")

	err := otherClient.Mkdir("/_test/accessdenied/dir2/foo", mode)
	assertPathError(t, err, "mkdir", "/_test/accessdenied/dir2/foo", os.ErrPermission)

	_, err = client.Stat("/_test/accessdenied/dir2/foo")
	assertPathError(t, err, "stat", "/_test/accessdenied/dir2/foo", os.ErrNotExist)

	_, err = client.Stat("/_test/accessdenied/dir2")
	assertPathError(t, err, "stat", "/_test/accessdenied/dir2", os.ErrNotExist)
}
