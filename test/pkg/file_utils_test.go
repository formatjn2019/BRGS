package test

import (
	"BRGS/pkg/tools"
	"BRGS/pkg/utils"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/go-ini/ini"
)

func TestWalk(t *testing.T) {
	walkedDic := tools.WalkDir("D:\\testf")
	fmt.Println(len(walkedDic))
	for k, v := range walkedDic {
		println(k, v)
		break
	}
	//fmt.Println(walkedDic2)
}

func TestWrite(t *testing.T) {
	context := []map[string]string{{"标题1": "332", "标题2": "4422"}, {"标题1": "33", "标题2": "44"}}
	utils.WriteCsvWithDict("config.csv", context)
}

func TestWriteUid(t *testing.T) {
	pathUidDic, _ := tools.CalculateAllUid("D:\\testDir")
	itemList := make([]map[string]string, 0, len(pathUidDic))
	for path, uid := range pathUidDic {
		itemList = append(itemList, map[string]string{"path": path, "uid": uid})
	}
	println(len(pathUidDic))
	println(utils.WriteCsvWithDict("uid.csv", itemList))
}

func TestSaveConfig(t *testing.T) {
	temp := os.TempDir()
	fmt.Println(temp)

	time := time.Now()
	fmt.Println(time.GoString())
	fmt.Println(time.Format("2006-01-02 15:04:05"))
	fmt.Println(time.Format("20060102_150405"))
}

func TestZip(t *testing.T) {
	inputPath := "D:\\testf"
	outputPath := "D:\\testf.zip"
	start := time.Now().UnixNano()
	dict := tools.WalkDir(inputPath)
	println(tools.WriteZip(outputPath, dict))
	end := time.Now().UnixNano()
	fmt.Println(start)
	fmt.Println(end)
	fmt.Println(end - start)
}

func TestIni(t *testing.T) {
	cfg, err := ini.Load("../test.ini")
	println()
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// // 典型读取操作，默认分区可以使用空字符串表示
	// fmt.Println("App Mode:", cfg.Section("").Key("app_mode").String())
	// fmt.Println("Data Path:", cfg.Section("paths").Key("data").String())

	// // 我们可以做一些候选值限制的操作
	// fmt.Println("Server Protocol:",
	// 	cfg.Section("server").Key("protocol").In("http", []string{"http", "https"}))
	// // 如果读取的值不在候选列表内，则会回退使用提供的默认值
	// fmt.Println("Email Protocol:",
	// 	cfg.Section("server").Key("protocol").In("smtp", []string{"imap", "smtp"}))

	// // 试一试自动类型转换
	// fmt.Printf("Port Number: (%[1]T) %[1]d\n", cfg.Section("server").Key("http_port").MustInt(9999))
	// fmt.Printf("Enforce Domain: (%[1]T) %[1]v\n", cfg.Section("server").Key("enforce_domain").MustBool(false))

	// // 差不多了，修改某个值然后进行保存
	// cfg.Section("").Key("app_mode").SetValue("production")
	fmt.Println(cfg.Section("excel_default_tip").Key("name").String())
	fmt.Println(cfg.Section("excel_default_tip").Key("watchDir").String())
	fmt.Println(cfg.Section("excel_default_tip").Key("tempDir").String())
	fmt.Println(cfg.Section("excel_default_tip").Key("archiveDir").String())
	fmt.Println(cfg.Section("excel_default_tip").Key("syncInterval").String())
	fmt.Println(cfg.Section("excel_default_tip").Key("archiveInterval").String())
	// cfg.SaveTo("my.ini.local")
}

func TestCrc32(t *testing.T) {
	// fmt.Println(CRC32("123456"))
	fmt.Println(CRC32("1234567890123456789012345678901221098765432109876543210987654321"))
	fmt.Println(CRC32("2109876543210987654321098765432112345678901234567890123456789012"))
	bytes := []byte{0, 0, 0, 0}
	fmt.Printf("crc32.ChecksumIEEE(bytes): %v\n", crc32.ChecksumIEEE(bytes))
	bytes2 := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	// fmt.Println(MD5("123456"))

	fmt.Printf("crc32.ChecksumIEEE(bytes2): %v\n", crc32.ChecksumIEEE(bytes2))
	// fmt.Println(SHA1("123456"))
}

func MD5(file string) (value [md5.Size]byte, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	value = md5.Sum(data)
	return

}

// 生成sha1
func SHA1(str string) string {
	c := sha1.New()
	c.Write([]byte(str))
	return hex.EncodeToString(c.Sum(nil))
}

func CRC32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

func TestSameCrc32(t *testing.T) {
	var num uint64
	crc32Map := map[uint32]uint64{}
outter:
	for num = 0; ; num++ {
		bytesArr := md5.Sum([]byte(strconv.FormatUint(num, 10)))
		crc32Num := crc32.ChecksumIEEE(bytesArr[:])
		if oldNum, ok := crc32Map[crc32Num]; ok {
			println(oldNum)
			println(num)

			break outter
		} else {
			crc32Map[crc32Num] = num
		}
	}

	num2 := num
	bytesArr := md5.Sum([]byte(strconv.FormatUint(num, 10)))
	num1 := crc32Map[crc32.ChecksumIEEE(bytesArr[:])]
	fmt.Println("crc32")
	fmt.Println(crc32.ChecksumIEEE(bytesArr[:]))
	fmt.Println(num1, num2)
	fmt.Println("md5:")
	fmt.Println(md5.Sum([]byte(strconv.FormatUint(num1, 10))))
	fmt.Println(md5.Sum([]byte(strconv.FormatUint(num2, 10))))
}

func TestSameMd5(t *testing.T) {
	fmt.Println(MD5("plane.jpg"))
	fmt.Println(MD5("ship.jpg"))

	arr, _ := os.ReadFile("plane.jpg")
	arr2, _ := os.ReadFile("ship.jpg")
	fmt.Printf("crc32.ChecksumIEEE(arr): %v\n", crc32.ChecksumIEEE(arr))
	fmt.Printf("crc32.ChecksumIEEE(arr): %v\n", crc32.ChecksumIEEE(arr2))

}

func TestSameMd5String(t *testing.T) {
	data1, _ := hex.DecodeString("4dc968ff0ee35c209572d4777b721587d36fa7b21bdc56b74a3dc0783e7b9518afbfa200a8284bf36e8e4b55b35f427593d849676da0d1555d8360fb5f07fea2")
	data2, _ := hex.DecodeString("4dc968ff0ee35c209572d4777b721587d36fa7b21bdc56b74a3dc0783e7b9518afbfa202a8284bf36e8e4b55b35f427593d849676da0d1d55d8360fb5f07fea2")
	fmt.Println(string(data1) == string(data2))
	//fmt.Println(string(data1))
	//fmt.Println(string(data2))
	strMd51, strMd52 := fmt.Sprintf("%x", md5.Sum(data1)), fmt.Sprintf("%x", md5.Sum(data2))
	fmt.Println(strMd51 + "\n" + strMd52)
	fmt.Println(strMd51 == strMd52)
}

func TestCrc32Same(t *testing.T) {
	n1, n2 := 4003, 20671
	md51, md52 := md5.Sum([]byte(strconv.Itoa(n1))), md5.Sum([]byte(strconv.Itoa(n2)))
	crc321, crc322 := crc32.ChecksumIEEE(md51[:]), crc32.ChecksumIEEE(md52[:])
	fmt.Println(n1, n2)
	fmt.Println(md51, md52)
	fmt.Println(crc321, crc322)
}

func TestSpeedCompare(t *testing.T) {
	testDir := "D:\\testDir\\input"

	result := map[string]string{}
	var walkDir func(string, string)
	walkDir = func(root string, prefix string) {
		fl, err := ioutil.ReadDir(root)
		if err != nil {
			panic(err)
		} else {
			for _, file := range fl {
				if file.IsDir() {
					walkDir(filepath.Join(root, file.Name()), prefix+"/"+file.Name())
				} else {
					result[filepath.Join(root, file.Name())] = (prefix + "/" + file.Name())[1:]
				}
			}
		}
	}
	walkDir(testDir, "")
	// for key, value := range result {
	// 	fmt.Printf("%-70s\t--------->\t%-50s", key, value)
	// }

	crc32Map := map[string]string{}
	startTime := time.Now().UnixMicro()
	for k, v := range result {
		ctx, _ := ioutil.ReadFile(k)
		u32 := crc32.ChecksumIEEE(ctx)
		crc32Map[v] = strconv.Itoa(int(u32))
	}
	endTime := time.Now().UnixMicro()

	fmt.Println("经历时间", endTime-startTime)
	// for k, v := range crc32Map {
	// 	fmt.Println(k, v)
	// }

	md5Map := map[string]string{}
	startTime = time.Now().UnixMicro()
	for k, v := range result {
		ctx, _ := ioutil.ReadFile(k)
		m5 := sha1.Sum(ctx)
		md5Map[v] = string(m5[:])
	}
	endTime = time.Now().UnixMicro()

	fmt.Println("经历时间", endTime-startTime)
	// for k, v := range md5Map {
	// 	fmt.Println(k, v)
	// }

}

func TestFilecompare(t *testing.T) {
	fmt.Printf("util.CompareFile(\"plane.jpg\", \"ship.jpg\"): %v\n", utils.CompareFile("plane.jpg", "ship.jpg"))
}

func TestCompareFileCopy(t *testing.T) {
	file1, file2 := "input/01.前言.md", "output/01.前言.md"
	fmt.Printf("util.CompareFile(file1, file2): %v\n", utils.CompareFile(file1, file2))
	ctx1, _ := ioutil.ReadFile(file1)
	ctx2, _ := ioutil.ReadFile(file2)
	fmt.Printf("md5.Sum(ctx1): %v\n", md5.Sum(ctx1))
	fmt.Printf("md5.Sum(ctx2): %v\n", md5.Sum(ctx2))
	fmt.Printf("md5.Sum(ctx2): %v\n", sha1.Sum(ctx1))
	fmt.Printf("md5.Sum(ctx2): %v\n", sha1.Sum(ctx2))
}

func TestCalculateAllFile(t *testing.T) {
	path := `D:\testf`
	pathUidDic, err := tools.CalculateAllUid(path)
	fmt.Println(pathUidDic)
	fmt.Println(err)
	result, err := tools.CheckUid(pathUidDic)
	fmt.Println(result, err)
}
