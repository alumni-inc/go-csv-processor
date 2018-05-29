package util

import (
	"encoding/csv"
	"os"

	"io"

	"io/ioutil"
	"strings"

	"github.com/gocarina/gocsv"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"github.com/lunny/log"
	"errors"
)

var encodings = []encoding.Encoding{
	unicode.UTF8,
	japanese.ShiftJIS,
	japanese.EUCJP,
	japanese.ISO2022JP,
}

type Record struct {
	EmployeeNumber string `csv:"社員番号" validate:"required"`
	LastName       string `csv:"苗字" validate:"required"`
	FirstName      string `csv:"名前" validate:"required"`
	Email          string `csv:"メールアドレス" validate:"email"`
	Skills         string `csv:"スキル"`
	Job            string `csv:"職種"`
	Position       string `csv:"役職"`
	Department     string `csv:"最終所属"`
	RetireYear     uint   `csv:"退職年度"`
	RetireReason   string `csv:"退職理由"`

	NotUsed string `csv:"-"`
}

func GetRecords(f *os.File) ([]*Record, error) {

	var records []*Record

	err := replaceLineFeedCode(f)
	if err != nil {
		return records, err
	}

	enc, err := determineEncoding(f)
	if err != nil {
		return records, err
	}

	log.Infof("filename: %v, encoding: %v", f.Name(), enc)

	// Customize reader
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(transform.NewReader(f, enc.NewDecoder()))
		r.LazyQuotes = true
		return r
	})

	// Load records from file
	if err = gocsv.UnmarshalFile(f, &records); err != nil {
		return records, err
	}

	// Validate only first record
	//if len(records) > 0 {
	//	err = validate.Struct(records[0])
	//}

	return records, err
}

// replaceLineFeedCode replaces all line feed code to CR+LF
func replaceLineFeedCode(f *os.File) error {
	const nlcode = "\r\n" // CR+LF

	replacer := strings.NewReplacer(
		"\r\n", nlcode, // CR+LF (Windows, etc)
		"\r", nlcode, // CR (Mac OS9, etc)
		"\n", nlcode, // LF (Unix, Max OSX, etc)
	)

	// Read the all file content
	// FIXME: Watch out the file size...
	bytes, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return err
	}

	f, err = os.OpenFile(f.Name(), os.O_RDWR, os.ModePerm) // Will be overwritten
	if err != nil {
		return err
	}
	defer f.Close()

	// Overwrite whole file with the replaced string
	_, err = f.WriteString(replacer.Replace(string(bytes)))

	return err
}

func determineEncoding(f *os.File) (encoding.Encoding, error) {

	for _, enc := range encodings {
		// Go to the start of the file
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}

		r := csv.NewReader(transform.NewReader(f, enc.NewDecoder()))
		r.LazyQuotes = true

		// Read only first line
		record, err := r.Read()
		if err != nil {
			return nil, err
		}

		if record[0] == "id" { // TODO: fix hard coding
			if _, err := f.Seek(0, io.SeekStart); err != nil {
				return nil, err
			}
			return enc, nil
		}
	}

	return nil, errors.New("column rule is wrong or unsupported file encoding")
}
