package filemaker

import (
    "bufio"
    "bytes"
    "github.com/whencome/gotil/fileutil"
    "io"
    "io/ioutil"
    "os"
    "path"
)

// 获取目标文件
func (fm *FileMaker) getTargetFile() string {
    return path.Join(fm.OutputDir, fm.FileName)
}

// 生成文件，如果文件已经存在，则更新model部分
func (fm *FileMaker) MakeFile() (bool, error) {
    targetFile := fm.getTargetFile()
    var rs bool = false
    var err error = nil
    // 检查目标文件是否已经存在
    if fileutil.Exists(targetFile) {
        // 刷新文件
        rs, err = fm.refreshFile(targetFile)
    } else {
        // 生成新的文件
        rs, err = fm.createFile(targetFile)
    }
    // 格式化go文件
    if rs {
        FmtGoFile(targetFile)
    }
    return rs, err
}

// 创建文件
func (fm *FileMaker) createFile(f string) (bool, error) {
    file, err := os.OpenFile(f, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
    if err != nil {
        return false, err
    }
    defer file.Close()
    // 开始写文件
    writer := bufio.NewWriter(file)
    // 写入头信息
    header := fm.buildHeader()
    // 写入头信息
    _, err = writer.WriteString(header)
    if err != nil {
        return false, err
    }
    // 写入model内容
    _, err = writer.WriteString(fm.getModelContent())
    if err != nil {
        return false, err
    }
    // 写入文件
    err = writer.Flush()
    if err != nil {
        return false, err
    }
    // 返回结果
    return true, nil
}

// 刷新文件
func (fm *FileMaker) refreshFile(f string) (bool, error) {
    file, err := os.OpenFile(f, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
    if err != nil {
        return false, err
    }
    defer file.Close()
    // 更新文件内容，不再修改package以及import中的内容
    reader := bufio.NewReader(file)
    buf := bytes.Buffer{}
    autoGeneratedCode := false
    for {
        line, _, err := reader.ReadLine()
        if err == io.EOF {
            break
        }
        if string(line[:3]) == "//:" && string(line[len(line)-3:]) == "://" {
            if !autoGeneratedCode {
                autoGeneratedCode = true
                buf.WriteString(fm.getModelContent())
            } else {
                autoGeneratedCode = false
            }
            continue
        }
        if autoGeneratedCode {
            continue
        }
        if !autoGeneratedCode {
            buf.Write(line)
            buf.WriteString(NEW_LINE)
        }
    }
    // 直接覆盖整个文件
    err = ioutil.WriteFile(f, buf.Bytes(), os.ModePerm)
    if err != nil {
        return false, err
    }
    return true, nil
}

// 写入model内容
func (fm *FileMaker) getModelContent() string {
    buf := bytes.Buffer{}
    // 写入注释信息
    buf.WriteString(MODEL_COMMENT_START)
    buf.WriteString(NEW_LINE)
    buf.WriteString("// the following model was generated by go model creator")
    buf.WriteString(NEW_LINE)
    buf.WriteString("// do not modify the code between the comment model-start and model-end")
    buf.WriteString(NEW_LINE)
    buf.WriteString("// or the code will be overwritten when update model by go model creator")
    buf.WriteString(NEW_LINE)
    // model 定义内容
    body := fm.buildModel()
    buf.WriteString(body)
    buf.WriteString(NEW_LINE)
    // 添加注释结尾
    buf.WriteString(NEW_LINE)
    buf.WriteString("// your own code should write below (after the mode-end comment) or place in a single file")
    buf.WriteString(NEW_LINE)
    buf.WriteString(MODEL_COMMENT_END)
    buf.WriteString(NEW_LINE)
    return buf.String()
}

