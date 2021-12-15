/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    qrcode_test
 * @Date:    2021/12/15 4:52 下午
 * @package: qrcode
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package qrcode

import (
	"image/png"
	"os"
	"testing"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func TestCode(t *testing.T) {

		// Create the barcode
		qrCode, _ := qr.Encode("http://blog.hawtech.cn", qr.M, qr.Auto)

		// Scale the barcode to 200x200 pixels
		qrCode, _ = barcode.Scale(qrCode, 200, 200)

		// create the output file
		file, _ := os.Create("qrcode.png")
		defer file.Close()

		// encode the barcode as png
		png.Encode(file, qrCode)
}
