package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Config 结构体用于存储配置信息
type Config struct {
	CodeQL struct {
		Path             string            `json:"path"`
		RulesPath        string            `json:"rules_path"`
		RulesPathEnabled bool              `json:"rules_path_enabled"`
		Queries          map[string]string `json:"queries"`
	} `json:"codeql"`
	Maven struct {
		Path string `json:"path"`
	} `json:"maven"`
	Workspace struct {
		TempDir string `json:"temp_dir"`
	} `json:"workspace"`
}

// 加载配置文件
func loadConfig() (*Config, error) {
	configFile := "config.json"
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

func main() {
	// 获取当前程序运行目录
	execDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		return
	}

	// 加载配置文件
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		return
	}

	// 如果配置文件中的临时目录是相对路径，则基于当前运行目录
	if !filepath.IsAbs(config.Workspace.TempDir) {
		config.Workspace.TempDir = filepath.Join(execDir, config.Workspace.TempDir)
	}

	myApp := app.New()
	window := myApp.NewWindow("CodeQL扫描工具")

	// 界面元素
	resultArea := widget.NewMultiLineEntry()
	resultArea.SetPlaceHolder("扫描结果将显示在这里...")

	// 配置界面元素
	codeqlPathEntry := widget.NewEntry()
	codeqlPathEntry.SetText(config.CodeQL.Path)
	codeqlPathEntry.SetMinRowsVisible(1)
	codeqlPathEntry.Resize(fyne.NewSize(600, 30))
	codeqlPathButton := widget.NewButton("选择CodeQL路径", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err == nil && uri != nil {
				codeqlPathEntry.SetText(filepath.ToSlash(uri.URI().Path()))
			}
		}, window)
	})

	rulesPathEntry := widget.NewEntry()
	rulesPathEntry.SetText(config.CodeQL.RulesPath)
	rulesPathEntry.SetMinRowsVisible(1)
	rulesPathEntry.Resize(fyne.NewSize(600, 30))
	rulesPathButton := widget.NewButton("选择规则库路径", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				rulesPathEntry.SetText(filepath.ToSlash(uri.Path()))
			}
		}, window)
	})

	rulesPathEnabled := widget.NewCheck("启用自定义规则库", func(enabled bool) {
		config.CodeQL.RulesPathEnabled = enabled
		if enabled {
			rulesPathEntry.Enable()
			rulesPathButton.Enable()
		} else {
			rulesPathEntry.Disable()
			rulesPathButton.Disable()
		}
	})
	rulesPathEnabled.SetChecked(config.CodeQL.RulesPathEnabled)
	if config.CodeQL.RulesPathEnabled {
		rulesPathEntry.Enable()
		rulesPathButton.Enable()
	} else {
		rulesPathEntry.Disable()
		rulesPathButton.Disable()
	}

	mavenPathEntry := widget.NewEntry()
	mavenPathEntry.SetText(config.Maven.Path)
	mavenPathEntry.SetMinRowsVisible(1)
	mavenPathEntry.Resize(fyne.NewSize(600, 30))
	mavenPathButton := widget.NewButton("选择Maven路径", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err == nil && uri != nil {
				mavenPathEntry.SetText(filepath.ToSlash(uri.URI().Path()))
			}
		}, window)
	})

	tempDirEntry := widget.NewEntry()
	tempDirEntry.SetText(config.Workspace.TempDir)
	tempDirEntry.SetMinRowsVisible(1)
	tempDirEntry.Resize(fyne.NewSize(600, 30))
	tempDirButton := widget.NewButton("选择临时目录", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				tempDirEntry.SetText(filepath.ToSlash(uri.Path()))
			}
		}, window)
	})

	// 保存配置按钮
	saveConfigButton := widget.NewButton("保存配置", func() {
		config.CodeQL.Path = codeqlPathEntry.Text
		config.CodeQL.RulesPath = rulesPathEntry.Text
		config.Maven.Path = mavenPathEntry.Text
		config.Workspace.TempDir = tempDirEntry.Text

		configData, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			dialog.ShowError(fmt.Errorf("生成配置文件失败: %v", err), window)
			return
		}

		if err := os.WriteFile("config.json", configData, 0644); err != nil {
			dialog.ShowError(fmt.Errorf("保存配置文件失败: %v", err), window)
			return
		}

		dialog.ShowInformation("成功", "配置已保存", window)
	})

	dirEntry := widget.NewEntry()
	dirEntry.SetMinRowsVisible(1)
	dirEntry.Resize(fyne.NewSize(600, 30))
	dirButton := widget.NewButton("选择目录", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				dirEntry.SetText(filepath.ToSlash(uri.Path()))
			}
		}, window)
	})

	langSelect := widget.NewSelect([]string{"Java", "Python", "JavaScript", "Go"}, nil)
	scanButton := widget.NewButton("开始扫描", func() {
		if dirEntry.Text == "" || langSelect.Selected == "" {
			resultArea.SetText("请先选择目录和语言")
			return
		}

		// 使用配置的临时目录或项目目录
		tmpDir := config.Workspace.TempDir
		if tmpDir == "" {
			tmpDir = dirEntry.Text
		}

		// 创建CodeQL数据库
		projectName := filepath.Base(dirEntry.Text)
		dbPath := filepath.Join(tmpDir, projectName+"-codeql_db")
		lang := strings.ToLower(langSelect.Selected)
		cmdArgs := []string{"database", "create", dbPath,
			"--language=" + lang,
			"--source-root=" + dirEntry.Text}

		// 如果是Java项目，设置MAVEN_HOME环境变量
		env := os.Environ()
		if lang == "java" && config.Maven.Path != "" {
			mavenHome := filepath.Dir(filepath.Dir(config.Maven.Path))
			env = append(env, "MAVEN_HOME="+mavenHome)
			env = append(env, "PATH="+filepath.Dir(config.Maven.Path)+string(os.PathListSeparator)+os.Getenv("PATH"))
		}

		cmdCreate := exec.Command(config.CodeQL.Path, cmdArgs...)
		cmdCreate.Env = env

		if out, err := cmdCreate.CombinedOutput(); err != nil {
			resultArea.SetText(fmt.Sprintf("创建数据库失败:\n%s\n%v", out, err))
			return
		}

		// 执行默认查询
		reportFile := filepath.Join(tmpDir, projectName+"-results.csv")

		queryPath := config.CodeQL.Queries[lang]
		if config.CodeQL.RulesPathEnabled && config.CodeQL.RulesPath != "" {
			queryPath = config.CodeQL.RulesPath
		}

		cmdAnalyze := exec.Command(config.CodeQL.Path, "database", "analyze", dbPath,
			queryPath,
			"--format=csv",
			"--output="+reportFile)

		if out, err := cmdAnalyze.CombinedOutput(); err != nil {
			resultArea.SetText(fmt.Sprintf("分析失败:\n%s\n%v", out, err))
			return
		}

		// 读取结果文件
		resultData, err := os.ReadFile(reportFile)
		if err != nil {
			resultArea.SetText(fmt.Sprintf("读取结果文件失败: %v", err))
			return
		}
		resultArea.SetText(string(resultData))
		dialog.ShowInformation("扫描完成", "报告已生成到："+reportFile, window)
	})

	// 布局
	form := container.NewVBox(
		widget.NewLabel("配置设置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("CodeQL路径:"),
			container.NewHBox(codeqlPathEntry, codeqlPathButton),
			widget.NewLabel("规则库路径:"),
			container.NewVBox(
				rulesPathEnabled,
				container.NewHBox(rulesPathEntry, rulesPathButton),
			),
			widget.NewLabel("Maven路径:"),
			container.NewHBox(mavenPathEntry, mavenPathButton),
			widget.NewLabel("临时工作目录:"),
			container.NewHBox(tempDirEntry, tempDirButton),
		),
		saveConfigButton,
		widget.NewSeparator(),
		widget.NewLabel("代码项目目录:"),
		container.NewHBox(dirEntry, dirButton),
		widget.NewLabel("选择语言:"),
		langSelect,
		scanButton,
		widget.NewLabel("扫描结果:"),
		resultArea,
	)

	window.SetContent(form)
	window.Resize(fyne.NewSize(1000, 700))
	window.ShowAndRun()
}
