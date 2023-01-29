package maildir

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeFolderName(t *testing.T) {

	{
		result, err := decodeFolderName("A")
		assert.NoError(t, err)
		assert.Equal(t, "A", result)
	}
	{
		result, err := decodeFolderName("a.b.c")
		assert.NoError(t, err)
		assert.Equal(t, "a.b.c", result)
	}
	{
		result, err := decodeFolderName("INBOX.&MMYwuTDI-.A-&MEI-&-1")
		assert.NoError(t, err)
		assert.Equal(t, "INBOX.テスト.A-あ&1", result)
	}
}

func TestDecodeFolderName_DecodeError(t *testing.T) {

	_, err := decodeFolderName("&AAA")
	assert.EqualError(t, err, "&AAA is invalid folder name: utf7: invalid UTF-7")
}
