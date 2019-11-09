package internal

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
