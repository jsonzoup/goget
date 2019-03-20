// Copyright 2019 JsonZou Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"strings"
	"fmt"
	"io/ioutil"
	"io"
	"bufio"
	"regexp"
	//"os/exec"
)

const  GITHUB = "github.com"
// 不需要扫描的dir,分号分割
const  NONE_SCAN_DIR = "views;static;vendor;docs;.idea;.mmwiki.session;conf;logs;data;Godeps;"
// 项目根目录
const PROJECT_ROOT_PATH = "D:/work/go/src/yourProject"

// get "gitahub.com/xxx/yyy"
func main() {
	fmt.Println(PROJECT_ROOT_PATH)
	remoteLibs := scanAll(PROJECT_ROOT_PATH)
	fmt.Println("\n--------所有依赖的远程库的 go get 命令--------")
	for _,rlib := range remoteLibs{
		goget(rlib)
	}
}
// 扫描所有的远程依赖包
func scanAll(path string) []string{
	remoteLibsMap := make(map[string]string)
	scanFile(path,remoteLibsMap)
	var remoteLibs []string
	for k,_ := range remoteLibsMap {
		remoteLibs = append(remoteLibs,k)
	}
	return remoteLibs
}

// 递归扫描指定工程下所有的go文件的代码
func scanFile(path string,remoteLibsMap map[string]string){
	files, _ := ioutil.ReadDir(path)
	for _, fi := range files {
		if fi.IsDir() {
			// 过滤不需要扫描的dir
			if !strings.Contains(NONE_SCAN_DIR,fi.Name()+";") {
				scanFile(path + "/" + fi.Name(),remoteLibsMap)
			}
		} else {
			// 只扫描go文件
			if strings.Contains(fi.Name(), ".go") {
				fmt.Println( path + "/" + fi.Name())
				file, _ := os.Open(path + "/" + fi.Name())
				buffer := bufio.NewReader(file)
				for {
					s, _, ok := buffer.ReadLine()
					if ok == io.EOF {
						break
					}
					result := matchRemoteLib(GITHUB,string(s));
					if result != "" {
					   remoteLibsMap[result]=""
					}

				}
				file.Close()
			}
		}
	}
}

// 逐行匹配获取远程lib
func matchRemoteLib(remoteRepo string,codeLine string) string{
	reg := regexp.MustCompile(`"(?P<remoteLib>`+remoteRepo+`/\S+/\S+){1}(/\S)?"`)
	match := reg.FindStringSubmatch(codeLine)
	groupNames := reg.SubexpNames()

	remoteLib :=""
	// 匹配远程仓库资源
	if len(match)>0{
		for i, name := range groupNames {
			if i != 0 && name == "remoteLib" {
				remoteLib = match[i]
				break
			}
		}
		if remoteRepo != "" {
			remoteLibSplit := strings.Split(remoteLib,"/")[0:3]
			return strings.Join(remoteLibSplit,"/");
		}
	}

	return ""
}
// 调用 go get [remoteLib] 命令
// 例如： go get github.com/snail007/go-activerecord
func goget(remoteLib string) {
	fmt.Println("go get " + remoteLib)
	// 一下代码执行有些问题，先按上面打出来的命令手动执行吧
	//cmd := exec.Command("go get " + remoteLib)
	//cmd.Stdout = os.Stdout //
	//err := cmd.Run()
	//if err != nil {
	//	fmt.Println("go get ",remoteLib," error.\n",err)
	//	return
	//}
	//fmt.Println("go get ",remoteLib," success.")
}
