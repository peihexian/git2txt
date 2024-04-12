package main

import (
	"github.com/go-git/go-git/v5"
	"os"
	"path/filepath"
	"strings"
)

const sourceCodeFileExtension = ".asp|.aspx|.bat|.c|.cc|.conf|.cpp|.cs|.css|.go|.gradle|.h|.hpp|.html|.ini|.java|.js|.json|.jsp|.jsx|.kt|.less|.log|.lua|.md|.php|.pl|.properties|.proto|.py|.rb|.scss|.sh|.sql|.swift|.ts|.tsx|.txt|.vue|.xml|.yaml|.yml"

func main() {
	/**
	  本程序用于克隆gitlub 或者gitlab 仓库的代码到客户端本地以后并将其合并为一个单独的文本文件
	*/
	//检测命令行有无仓库地址输入参数,如果没有则打印程序用法
	if len(os.Args) < 4 {
		println("Usage: git2txt [gitlab url] [directory] [txt file name]")
		os.Exit(1)
	}
	//获取命令行输入的仓库地址
	url := os.Args[1]
	directory := os.Args[2]
	//获取命令行输入的文本文件名
	txtFileName := os.Args[3]

	//检查仓库地址是否合法
	if url == "" {
		println("url is empty")
		os.Exit(1)
	}
	if url[:4] != "http" && url[:3] != "git" {
		println("url is not a valid url")
		os.Exit(1)
	}

	//检查本地地址是否合法
	if directory == "" {
		println("local directory is empty")
		os.Exit(1)
	}
	//判断本地路径是否存在
	_, err := os.Stat(directory)
	if err != nil {
		//如果本地路径不存在则创建
		err = os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			println("create directory error:", err)
			os.Exit(1)
		}
	}

	//判断本地路径是否是一个文件夹
	if !filepath.IsAbs(directory) {
		println("local directory is not a directory")
		os.Exit(1)
	}

	//判断本地文件夹是否有内容
	files, err := os.ReadDir(directory)
	if err != nil {
		println("read directory error:", err)
		os.Exit(1)
	}
	if len(files) > 0 {
		//如果本地文件夹有内容则清空
		err = os.RemoveAll(directory)
		if err != nil {
			println("remove directory error:", err)
			os.Exit(1)
		}
	}

	//克隆仓库代码到本地
	err = clone(url, directory)
	if err != nil {
		println("clone error:", err)
		os.Exit(1)
	}
	println("clone success")
	//合并仓库代码为一个单独的文本文件
	err = merge(directory, txtFileName)
	if err != nil {
		println("merge error:", err)
		os.Exit(1)
	}
	println("merge success,output file is:", txtFileName)
}

func clone(url string, directory string) error {
	_, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		println("clone error:", err)
		return err
	}
	return nil
}

func merge(directory string, txtFileName string) error {
	//创建一个新的文本文件
	file, err := os.Create(txtFileName)
	if err != nil {
		println("create file error:", err)
		return err
	}
	defer file.Close()

	//分割源码文件后缀
	sourceCodeFileExtensionList := strings.Split(sourceCodeFileExtension, "|")

	//遍历本地仓库文件夹下的所有文件
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		//判断文件是否是源码文件
		isSourceCodeFile := false
		for _, extension := range sourceCodeFileExtensionList {
			if strings.HasSuffix(path, extension) {
				isSourceCodeFile = true
				break
			}
		}
		if !isSourceCodeFile {
			return nil
		}
		//读取文件内容
		content, err := os.ReadFile(path)
		if err != nil {
			println("read file error:", err)
			return err
		}
		//先写入文件名
		_, err = file.Write([]byte("-------------------------\nThe following content from file: " + path + "\r\n-------------------------\r\n\r\n"))
		if err != nil {
			println("write file error:", err)
			return err
		}
		//再写入文件内容
		_, err = file.Write(content)
		if err != nil {
			println("write file error:", err)
			return err
		}
		return nil
	})
	return nil
}
