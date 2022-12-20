package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/manifoldco/promptui"
	"github.com/rivo/tview"
	"github.com/schollz/progressbar"
	"github.com/spf13/cobra"
	"github.com/wgmouton/.dotfiles/manager/state"
	"github.com/wgmouton/.dotfiles/manager/types"
	"gopkg.in/yaml.v3"
)

var logColors []string = []string{
	"#f0f8ff", //aliceblue
	"#faebd7", //antiquewhite
	"#00ffff", //aqua
	"#7fffd4", //aquamarine
	"#f0ffff", //azure
	"#f5f5dc", //beige
	"#ffe4c4", //bisque
	"#000000", //black
	"#ffebcd", //blanchedalmond
	"#0000ff", //blue
	"#8a2be2", //blueviolet
	"#a52a2a", //brown
	"#deb887", //burlywood
	"#5f9ea0", //cadetblue
	"#7fff00", //chartreuse
	"#d2691e", //chocolate
	"#ff7f50", //coral
	"#6495ed", //cornflowerblue
	"#fff8dc", //cornsilk
	"#dc143c", //crimson
	"#00ffff", //cyan
	"#00008b", //darkblue
	"#008b8b", //darkcyan
	"#b8860b", //darkgoldenrod
	"#a9a9a9", //darkgray
	"#006400", //darkgreen
	"#a9a9a9", //darkgrey
	"#bdb76b", //darkkhaki
	"#8b008b", //darkmagenta
	"#556b2f", //darkolivegreen
	"#ff8c00", //darkorange
	"#9932cc", //darkorchid
	"#8b0000", //darkred
	"#e9967a", //darksalmon
	"#8fbc8f", //darkseagreen
	"#483d8b", //darkslateblue
	"#2f4f4f", //darkslategray
	"#2f4f4f", //darkslategrey
	"#00ced1", //darkturquoise
	"#9400d3", //darkviolet
	"#ff1493", //deeppink
	"#00bfff", //deepskyblue
	"#696969", //dimgray
	"#696969", //dimgrey
	"#1e90ff", //dodgerblue
	"#b22222", //firebrick
	"#fffaf0", //floralwhite
	"#228b22", //forestgreen
	"#ff00ff", //fuchsia
	"#dcdcdc", //gainsboro
	"#f8f8ff", //ghostwhite
	"#ffd700", //gold
	"#daa520", //goldenrod
	"#808080", //gray
	"#008000", //green
	"#adff2f", //greenyellow
	"#808080", //grey
	"#f0fff0", //honeydew
	"#ff69b4", //hotpink
	"#cd5c5c", //indianred
	"#4b0082", //indigo
	"#fffff0", //ivory
	"#f0e68c", //khaki
	"#e6e6fa", //lavender
	"#fff0f5", //lavenderblush
	"#7cfc00", //lawngreen
	"#fffacd", //lemonchiffon
	"#add8e6", //lightblue
	"#f08080", //lightcoral
	"#e0ffff", //lightcyan
	"#fafad2", //lightgoldenrodyellow
	"#d3d3d3", //lightgray
	"#90ee90", //lightgreen
	"#d3d3d3", //lightgrey
	"#ffb6c1", //lightpink
	"#ffa07a", //lightsalmon
	"#20b2aa", //lightseagreen
	"#87cefa", //lightskyblue
	"#778899", //lightslategray
	"#778899", //lightslategrey
	"#b0c4de", //lightsteelblue
	"#ffffe0", //lightyellow
	"#00ff00", //lime
	"#32cd32", //limegreen
	"#faf0e6", //linen
	"#ff00ff", //magenta
	"#800000", //maroon
	"#66cdaa", //mediumaquamarine
	"#0000cd", //mediumblue
	"#ba55d3", //mediumorchid
	"#9370db", //mediumpurple
	"#3cb371", //mediumseagreen
	"#7b68ee", //mediumslateblue
	"#00fa9a", //mediumspringgreen
	"#48d1cc", //mediumturquoise
	"#c71585", //mediumvioletred
	"#191970", //midnightblue
	"#f5fffa", //mintcream
	"#ffe4e1", //mistyrose
	"#ffe4b5", //moccasin
	"#ffdead", //navajowhite
	"#000080", //navy
	"#fdf5e6", //oldlace
	"#808000", //olive
	"#6b8e23", //olivedrab
	"#ffa500", //orange
	"#ff4500", //orangered
	"#da70d6", //orchid
	"#eee8aa", //palegoldenrod
	"#98fb98", //palegreen
	"#afeeee", //paleturquoise
	"#db7093", //palevioletred
	"#ffefd5", //papayawhip
	"#ffdab9", //peachpuff
	"#cd853f", //peru
	"#ffc0cb", //pink
	"#dda0dd", //plum
	"#b0e0e6", //powderblue
	"#800080", //purple
	"#ff0000", //red
	"#bc8f8f", //rosybrown
	"#4169e1", //royalblue
	"#8b4513", //saddlebrown
	"#fa8072", //salmon
	"#f4a460", //sandybrown
	"#2e8b57", //seagreen
	"#fff5ee", //seashell
	"#a0522d", //sienna
	"#c0c0c0", //silver
	"#87ceeb", //skyblue
	"#6a5acd", //slateblue
	"#708090", //slategray
	"#708090", //slategrey
	"#fffafa", //snow
	"#00ff7f", //springgreen
	"#4682b4", //steelblue
	"#d2b48c", //tan
	"#008080", //teal
	"#d8bfd8", //thistle
	"#ff6347", //tomato
	"#40e0d0", //turquoise
	"#ee82ee", //violet
	"#f5deb3", //wheat
	"#ffffff", //white
	"#f5f5f5", //whitesmoke
	"#ffff00", //yellow
	"#9acd32", //yellowgreen
}

type Progress struct {
	textView *tview.TextView
	full     int
	limit    int
	progress chan int
}

// full is the maximum amount of value can be sent to channel
// limit is the progress bar size
func (p *Progress) Init(full int, limit int) chan int {
	p.progress = make(chan int)
	p.full = full
	p.limit = limit

	go func() { // Simple channel status gauge (progress bar)
		progress := 0
		for {
			progress += <-p.progress

			if progress > full {
				break
			}

			x := progress * limit / full
			p.textView.Clear()
			_, _ = fmt.Fprintf(p.textView, "channel status:  %s%s %d/%d",
				strings.Repeat("■", x),
				strings.Repeat("□", limit-x),
				progress, full)

		}
	}()

	return p.progress
}

// TODO:: NEEDS COMPLETE REWRITE PROOF OF CONCEPT ONLY
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs all or specific tools",
	Long:  `This subcommand says hello`,
	Run: func(cmd *cobra.Command, args []string) {
		// ctx, cancelF := context.WithCancel(cmd.Context())
		// Set done status
		doneStatus := make(chan error, 1)
		defer close(doneStatus)
		defer func() {
			if err := <-doneStatus; err != nil {
				fmt.Println("error", err)
				os.Exit(1)
			}
			os.Exit(0)
		}()

		// Fetch tool list
		toolList := state.GetToolList()
		// toolListSize := len(toolList)

		runInitStage, err := cmd.Flags().GetBool("init")
		if err != nil {
			fmt.Println(err)
			return
		}

		// Set Table data
		toolCounter := 0
		doneCounter := 0

		// Get Stages
		stages := func() []int {
			stages := []int{}
			for stage := range toolList {
				if stage == 0 && !runInitStage {
					continue
				}
				stages = append(stages, stage)
			}
			sort.Ints(stages)
			return stages
		}()

		// Sort Stages

		testChan := make(chan string, 10000000)
		defer close(testChan)
		globalLogChan := make(chan string, 10000000)
		defer close(globalLogChan)

		progressChan := make(chan int, 10000)
		defer close(progressChan)

		reportChan := make(chan types.ExecutionReport, 10000)
		defer close(reportChan)

		finalReportChan := make(chan []types.ExecutionReport, 10000000)
		defer close(finalReportChan)
		// defer func() {
		// 	report := <-finalReportChan
		// 	yamlReport, _ := yaml.Marshal(&report)
		// 	os.WriteFile("./report.yaml", yamlReport, 0644)
		// }()

		// go func() {
		// 	reports := []types.ExecutionReport{}
		// 	for toolCounter != doneCounter {
		// 		reports = append(reports, <-reportChan)
		// 	}
		// 	finalReportChan <- reports
		// }()

		executeScript := func(script types.InstallScriptDefinition, colorCounter int) {
			logChan, statusChan := script.GetChannels()
			defer func() {
				// fmt.Println(toolCounter, (float64(doneCounter) / float64(toolCounter) * 100))
				doneCounter++
				progressChan <- int((float64(doneCounter) / float64(toolCounter) * 100))
			}()

			//TODO: Need to check if working
			// exec.Cmd("")

			// fmt.Printf("Installing %s...\n", script.Name)
			// fmt.Println(script.Description)

			for key, config := range script.Config {

				// Update Status
				if key == 0 {
					statusChan <- types.ExecutionStatusConfiguring
					globalLogChan <- fmt.Sprintf("[%s][ %s ] [white]%s", logColors[colorCounter], script.Name, "Configuaring tool before scripts...")
					logChan <- "Configuaring tool before scripts..."
				}

				if !config.GetBeforeScript() {
					continue
				}

				switch c := config.(type) {
				case types.ConfigLink:
					homePath, _ := os.UserHomeDir()
					targetDir := filepath.Dir(homePath + c.Target)

					_, err := os.Stat(targetDir)
					if errors.Is(err, &os.PathError{}) {
						os.MkdirAll(targetDir, fs.ModeDir)
					}

					if err := os.Symlink(script.ToolPath+c.Source, homePath+c.Target); err != nil {
						statusChan <- types.ExecutionStatusFailed
						globalLogChan <- fmt.Sprintf("[%s][ %s ] [white]%s", logColors[colorCounter], script.Name, err.Error())
						logChan <- err.Error()
					}
				case types.ConfigScript:
					configCmd := exec.Command("bash", "-c", c.Script)
					configStdout, _ := configCmd.StdoutPipe()
					configStderr, _ := configCmd.StderrPipe()
					configCmd.Start()

					go func() {
						scanner := bufio.NewScanner(configStderr)
						scanner.Split(bufio.ScanLines)
						for scanner.Scan() {
							m := scanner.Text()
							globalLogChan <- fmt.Sprintf("[%s][ %s ] [white]%s", logColors[colorCounter], script.Name, m)
							logChan <- m
						}
					}()

					go func() {
						scannerGood := bufio.NewScanner(configStdout)
						scannerGood.Split(bufio.ScanLines)
						for scannerGood.Scan() {
							x := scannerGood.Text()
							globalLogChan <- fmt.Sprintf("[%s][ %s ] [white]%s", logColors[colorCounter], script.Name, x)
							logChan <- x
						}
					}()

					if err := configCmd.Wait(); err != nil {
						globalLogChan <- fmt.Sprintf("[%s][ %s ] [white]%s", logColors[colorCounter], script.Name, err.Error())
						logChan <- err.Error()
					}
				}
			}

			cmd := exec.Command("bash", "-c", script.Scripts.MacosArm.Run)
			stdout, _ := cmd.StdoutPipe()
			stderr, _ := cmd.StderrPipe()
			statusChan <- types.ExecutionStatusInProgress
			cmd.Start()

			go func() {
				scanner := bufio.NewScanner(stderr)
				scanner.Split(bufio.ScanLines)
				for scanner.Scan() {
					m := scanner.Text()
					globalLogChan <- fmt.Sprintf("[%s][ %s ] [white]%s", logColors[colorCounter], script.Name, m)
					logChan <- m
				}
			}()

			go func() {
				scannerGood := bufio.NewScanner(stdout)
				scannerGood.Split(bufio.ScanLines)
				for scannerGood.Scan() {
					x := scannerGood.Text()
					globalLogChan <- fmt.Sprintf("[%s][ %s ] [white]%s", logColors[colorCounter], script.Name, x)
					logChan <- x
				}
			}()

			err := cmd.Wait()
			if err != nil {
				statusChan <- types.ExecutionStatusFailed
				return
			}

			for key, config := range script.Config {

				// Update Status
				if key == 0 {
					statusChan <- types.ExecutionStatusConfiguring
					globalLogChan <- fmt.Sprintf("[%s][ %s ] [white]%s", logColors[colorCounter], script.Name, "Configuaring tool after scripts...")
					logChan <- "Configuaring tool after scripts..."
				}

				if config.GetBeforeScript() {
					continue
				}

				switch c := config.(type) {
				case types.ConfigLink:
					homePath, _ := os.UserHomeDir()
					targetDir := filepath.Dir(homePath + c.Target)

					_, err := os.Stat(targetDir)
					if errors.Is(err, &os.PathError{}) {
						os.MkdirAll(targetDir, fs.ModeDir)
					}

					if err := os.Symlink(script.ToolPath+c.Source, homePath+c.Target); err != nil {
						statusChan <- types.ExecutionStatusFailed
						globalLogChan <- fmt.Sprintf("[%s][ %s ] [white]%s", logColors[colorCounter], script.Name, err.Error())
						logChan <- err.Error()
					}
				case types.ConfigScript:
					configCmd := exec.Command("bash", "-c", c.Script)
					configStdout, _ := configCmd.StdoutPipe()
					configStderr, _ := configCmd.StderrPipe()
					configCmd.Start()

					go func() {
						scanner := bufio.NewScanner(configStderr)
						scanner.Split(bufio.ScanLines)
						for scanner.Scan() {
							m := scanner.Text()
							globalLogChan <- fmt.Sprintf("[%s][ %s ] [white]%s", logColors[colorCounter], script.Name, m)
							logChan <- m
						}
					}()

					go func() {
						scannerGood := bufio.NewScanner(configStdout)
						scannerGood.Split(bufio.ScanLines)
						for scannerGood.Scan() {
							x := scannerGood.Text()
							globalLogChan <- fmt.Sprintf("[%s][ %s ] [white]%s", logColors[colorCounter], script.Name, x)
							logChan <- x
						}
					}()

					if err := configCmd.Wait(); err != nil {
						globalLogChan <- fmt.Sprintf("[%s][ %s ] [white]%s", logColors[colorCounter], script.Name, err.Error())
						logChan <- err.Error()
					}
				}
			}

			statusChan <- types.ExecutionStatusCompleted
			statusChan <- types.ExecutionStatusInstalled
			// reportChan <- types.ExecutionReport{
			// 	Name:             script.Name,
			// 	Status:           types.ExecutionStatusInstalled,
			// 	Version:          new(string),
			// 	InstallationPath: new(string),
			// }
		}

		// var concurantExecutionsWG sync.WaitGroup

		go func() {
			colorCounter := 0
			for _, stage := range stages {
				scripts := toolList[stage]

				var wgStage sync.WaitGroup
				var wgCommand sync.WaitGroup
				wgStage.Add(len(scripts))

				for _, script := range scripts {
					colorCounter++
					if script.Scripts.MacosArm.Async {
						wgStage.Add(1)
						go func(s types.InstallScriptDefinition) {
							defer wgStage.Done()
							defer wgCommand.Done()
							executeScript(s, colorCounter)
						}(script)
					} else {
						wgCommand.Wait()
						executeScript(script, colorCounter)
						wgStage.Done()
					}
				}
				wgStage.Wait()
			}
		}()

		// Setup the ui
		newPrimitive := func(text string) tview.Primitive {
			return tview.NewTextView().
				SetTextAlign(tview.AlignCenter).
				SetText(text)
		}

		app := tview.NewApplication()

		// Setup View
		table := tview.NewTable().
			SetBorders(false).
			SetWrapSelection(true, false).
			SetSelectable(true, false)
		// SetSelectedFunc(func(row, column int) {
		// 	app.
		// })

		modal := func(p tview.Primitive, width, height int) tview.Primitive {
			return tview.NewFlex().AddItem(p, 0, 1, true)
		}

		toolLogView := func(logChan <-chan string) tview.Primitive {
			flex := tview.NewFlex().SetDirection(tview.FlexRow)
			flex.SetBorder(true).SetTitle("Tool Logs")
			go func() {
				for {
					log := <-logChan
					if log == "" {
						break
					}

					app.QueueUpdateDraw(func() {
						flex.AddItem(tview.NewTextView().SetDynamicColors(true).SetText(log), 1, 0, false)
					})
				}
			}()
			return flex
		}

		toolListMap := map[int]types.InstallScriptDefinition{}
		toolListLogsMap := map[int]tview.Primitive{}

		// Set Table headings
		table.SetCell(0, 0, tview.NewTableCell("Stage").SetSelectable(false))
		table.SetCell(0, 1, tview.NewTableCell("Tool").SetSelectable(false))
		table.SetCell(0, 2, tview.NewTableCell("Status").SetSelectable(false))
		table.SetCell(0, 3, tview.NewTableCell("Platform").SetSelectable(false))
		table.SetCell(0, 4, tview.NewTableCell("Version").SetSelectable(false))

		for _, stage := range stages {
			for _, tool := range toolList[stage] {
				logChan, statusChan := tool.GetChannels()
				toolCounter++
				statusCell := tview.NewTableCell("Waiting")
				table.SetCell(toolCounter, 0, tview.NewTableCell(strconv.FormatInt(int64(stage), 32)))
				table.SetCell(toolCounter, 1, tview.NewTableCell(tool.Name))
				table.SetCell(toolCounter, 2, statusCell)

				toolListMap[toolCounter] = tool
				toolListLogsMap[toolCounter] = toolLogView(logChan)
				go func() {
					cont := true
					for cont {
						if !cont {
							return
						}
						switch <-statusChan {
						case types.ExecutionStatusInProgress:
							app.QueueUpdateDraw(func() {
								statusCell.SetText("In Progress")
								cont = true
							})
						case types.ExecutionStatusCompleted:
							app.QueueUpdateDraw(func() {
								statusCell.SetText("Completed")
								cont = true
							})
						case types.ExecutionStatusFailed:
							app.QueueUpdateDraw(func() {
								statusCell.SetText("Failed").SetTextColor(tcell.ColorRed)
								cont = false
							})
						case types.ExecutionStatusSkiped:
							app.QueueUpdateDraw(func() {
								statusCell.SetText("Skiped")
								cont = false
							})
						case types.ExecutionStatusInstalled:
							app.QueueUpdateDraw(func() {
								statusCell.SetText("Installed").SetTextColor(tcell.ColorGreen)
								cont = false
							})
						case types.ExecutionStatusConfiguring:
							app.QueueUpdateDraw(func() {
								statusCell.SetText("Configuring")
								cont = true
							})
						default:
							cont = true
						}
					}
				}()
			}
		}

		globalLogView := func() tview.Primitive {
			flex := tview.NewFlex().SetDirection(tview.FlexRow)
			flex.SetBorder(true).SetTitle("Global Logs")
			go func() {
				for {
					log := <-globalLogChan
					if log == "" {
						break
					}

					app.QueueUpdateDraw(func() {
						flex.AddItem(tview.NewTextView().SetDynamicColors(true).SetText(log), 1, 0, true)
					})
				}
			}()

			// flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			// 	old := c.focus
			// 	switch event.Key() {
			// 	case tcell.KeyUp:
			// 		if c.focus == 0 {
			// 			return event
			// 		}

			// 		c.focus--

			// 	case tcell.KeyDown:
			// 		if c.focus == len(c.videoElements)-1 {
			// 			return event
			// 		}

			// 		c.focus++

			// 	default:
			// 		// if c.customHandler != nil {
			// 		// 	return c.customHandler(ev)
			// 		// }

			// 		return event
			// 	}

			// 	// c.videoElements[old].SetBackgroundColor(tcell.ColorBlack)
			// 	// c.videoElements[old].SetTextColor(tcell.ColorWhite)

			// 	// c.videoElements[c.focus].SetBackgroundColor(tcell.ColorWhite)
			// 	// c.videoElements[c.focus].SetTextColor(tcell.ColorBlack)

			// 	c.SetOffset(c.focus-1, 0)
			// 	return event
			// })

			return flex
		}()

		table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'd':
				row, _ := table.GetSelection()

				// fmt.Println(row, col)
				tool := toolListMap[row]

				bytes, _ := yaml.Marshal(&tool)

				app.SetRoot(modal(tview.NewTextView().SetText(string(bytes)), 50, 50), true)
				return event
			case 'l':
				row, _ := table.GetSelection()
				tool := toolListMap[row]
				toolLogs := toolListLogsMap[row]

				logGrid := tview.NewFlex()
				command := tview.NewTextView().SetText(tool.Scripts.MacosArm.Run)
				command.SetTitle(" Command ").SetBorder(true)
				logGrid.SetBorder(true)
				logGrid.SetDirection(tview.FlexRow)
				logGrid.AddItem(command, 0, 1, false)
				logGrid.AddItem(toolLogs, 0, 9, true)

				app.SetRoot(logGrid, true)
				return event
			case 'L':
				app.SetRoot(globalLogView, true)
				return event
			}
			return event
		})

		footer := func() tview.Primitive {
			flexBox := tview.NewFlex()
			tv := tview.NewTextView().SetText(strings.Repeat("□", 100))
			flexBox.AddItem(tview.NewTextView().SetText("Progress"), 0, 1, false)
			flexBox.AddItem(tv, 0, 4, false)
			go func() {
				for {
					perc := <-progressChan
					app.QueueUpdateDraw(func() {
						tv.SetText(strings.Repeat("■", perc) + strings.Repeat("□", 100-perc))
					})
				}
			}()
			return flexBox
		}()

		grid := tview.NewGrid().
			SetRows(1, 0, 1).
			SetColumns(0).
			SetBorders(false).
			AddItem(newPrimitive("Header"), 0, 0, 1, 1, 0, 0, false).
			AddItem(table, 1, 0, 1, 1, 0, 0, false).
			AddItem(footer, 2, 0, 1, 1, 0, 0, false)

		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				app.SetRoot(grid, true).SetFocus(table)
				return nil
			}
			return event
		})
		// toolListSize := len(toolList)

		// executionReports := make(chan ExecutionReport, len(toolList))
		// go func() {
		// 	for r := 0; r < len(toolList); r++ {
		// 		select {}
		// 	}
		// }()

		// toolCounter := 0
		// for k, tool := range toolList {
		// 	// Build stage list
		// 	stages = append(stages, k)

		// 	// Build table ui
		// 	table.SetCell(toolCounter, 0, tview.NewTableCell(tool.Name))
		// 	// 	SetTextColor(color).
		// 	// 	SetAlign(tview.AlignCenter))
		// 	// word = (word + 1) % len(lorem)

		// 	toolCounter++

		// }

		// Sort stages

		if err := app.SetRoot(grid, true).SetFocus(table).Run(); err != nil {
			panic(err)
		}

		// Get all stages and short them
		// stages := func() []int {

		// 	return stages
		// }()

		// Setup Tool List Table

		// cols, rows := 5, toolListSize
		// for r := 0; r < rows; r++ {
		// 	for c := 0; c < cols; c++ {
		// 		color := tcell.ColorWhite
		// 		if c < 1 || r < 1 {
		// 			color = tcell.ColorYellow
		// 		}
		// 		table.SetCell(r, c,
		// 			tview.NewTableCell(lorem[word]).
		// 				SetTextColor(color).
		// 				SetAlign(tview.AlignCenter))
		// 		word = (word + 1) % len(lorem)
		// 	}
		// }
		// table.Select(0, 0).SetFixed(1, 1).SetDoneFunc

		// textView := tview.NewTextView().
		// 	SetChangedFunc(func() {
		// 		app.Draw()
		// 	})

		// progress := Progress{textView: textView}
		// progChan := progress.Init(360, 20)

		// go func() { // update progress bar

		// 	i := 0

		// 	for {
		// 		i++
		// 		progChan <- 1

		// 		if i > progress.full {
		// 			close(progChan)
		// 			break
		// 		}

		// 		time.Sleep(100 * time.Millisecond)
		// 	}
		// }()

		//
		// menu := newPrimitive("Menu")
		// main := newPrimitive("Main content")
		// sideBar := newPrimitive("Side Bar")

		// table := tview.NewTable().SetBorders(false)

		// lorem := strings.Split("Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.", " ")
		// cols, rows := 5, 40
		// word := 0
		// for r := 0; r < rows; r++ {
		// 	for c := 0; c < cols; c++ {
		// 		color := tcell.ColorWhite
		// 		if c < 1 || r < 1 {
		// 			color = tcell.ColorYellow
		// 		}
		// 		table.SetCell(r, c,
		// 			tview.NewTableCell(lorem[word]).
		// 				SetTextColor(color).
		// 				SetAlign(tview.AlignCenter))
		// 		word = (word + 1) % len(lorem)
		// 	}
		// }
		// // table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		// // 	if key == tcell.KeyEscape {
		// // 		app.Stop()
		// // 	}
		// // 	if key == tcell.KeyEnter {
		// // 		table.SetSelectable(true, true)
		// // 	}
		// // }).SetSelectedFunc(func(row int, column int) {
		// // 	table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		// // 	table.SetSelectable(false, false)
		// // })
		// // if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		// // 	panic(err)
		// // }

		// 	// AddItem(textView, 2, 0, 1, 1, 0, 0, false)

		// // // Layout for screens narrower than 100 cells (menu and side bar are hidden).
		// // grid.AddItem(menu, 0, 0, 0, 0, 0, 0, false).
		// // 	AddItem(main, 1, 0, 1, 3, 0, 0, false).
		// // 	AddItem(sideBar, 0, 0, 0, 0, 0, 0, false)

		// // Layout for screens wider than 100 cells.
		// // grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		// // AddItem(main, 1, 1, 1, 1, 0, 100, false).
		// // AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

		// if err := tview.NewApplication().SetRoot(grid, true).Run(); err != nil {
		// 	panic(err)
		// }

		// // frame.AddText()

		// app := tview.NewApplication()

		// textView.SetBorder(true)
		// textView.SetBackgroundColor(tcell.ColorDefault)

		// if err := app.SetRoot(textView, true).SetFocus(textView).Run(); err != nil {
		// 	panic(err)
		// }

		// app := tview.NewApplication()

		// box := tview.NewBox().SetBorder(false).SetTitle("Hello, world!")]

		// app := tview.NewApplication()
		// list := tview.NewList().
		// 	AddItem("List item 1", "Some explanatory text", 'a', nil).
		// 	AddItem("List item 2", "Some explanatory text", 'b', nil).
		// 	AddItem("List item 3", "Some explanatory text", 'c', nil).
		// 	AddItem("List item 4", "Some explanatory text", 'd', nil).
		// 	AddItem("Quit", "Press to exit", 'q', func() {
		// 		app.Stop()
		// 	})
		// if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		// 	doneStatus <- err
		// 	return
		// }

		// list := tview.NewList()

		// if err := tview.NewApplication().SetRoot(list, true).Run(); err != nil {
		// 	doneStatus <- err
		// 	return
		// }

		doneStatus <- nil
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates all or specific tools",
	Long:  `This subcommand says hello`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello called")

		// validate := func(input string) error {
		// 	_, err := strconv.ParseFloat(input, 64)
		// 	if err != nil {
		// 		return errors.New("Invalid number")
		// 	}
		// 	return nil
		// }

		prompt := promptui.Select{
			Label:        "Number",
			Items:        []string{"hi", "bue"},
			HideSelected: false,
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You choose %q\n", result)
	},
}

func runner() {
	cmd := exec.Command("tr", "a-z", "A-Z")

	cmd.Stdin = strings.NewReader("and old falcon")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("translated phrase: %q\n", out.String())

	bar := progressbar.New(100)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(40 * time.Millisecond)
	}
}

func init() {
	installCmd.Flags().Bool("init", false, "Run stage 0")
	RootCmd.AddCommand(installCmd, updateCmd)
}
