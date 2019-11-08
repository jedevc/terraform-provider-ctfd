package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"reflect"
)

func MultipartMarshal(writer *multipart.Writer, content interface{}) error {
	v := reflect.ValueOf(content)
	t := reflect.TypeOf(content)

	for i := 0; i < t.NumField(); i++ {
		// extract field names and tags
		name := t.Field(i).Tag.Get("multipart")
		value := v.Field(i).Interface()

		// ignore explicitly marked blank fields
		if name == "-" {
			continue
		}

		// define default
		if len(name) == 0 {
			name = t.Field(i).Name
		}

		file, ok := value.(*os.File)
		if ok {
			// create form field
			field, err := writer.CreateFormFile(name, file.Name())
			if err != nil {
				return err
			}

			_, err = io.Copy(field, file)
			if err != nil {
				return err
			}
		} else {
			// create normal field
			field, err := writer.CreateFormField(name)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(field, "%v", value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
