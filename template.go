package telegrabotlib

import (
    "bytes"
    "html/template"
)

func RenderTemplate(filePath string, data interface{}) (string, error) {
    tmpl := template.Must(template.ParseFiles(filePath))

    b := bytes.NewBufferString("")

    if err := tmpl.Execute(b, data); err != nil {
        return "", err
    }

    return b.String(), nil
}
