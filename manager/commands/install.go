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
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

var logColors []string = []string{
	"aquamarine",
	"blue",
	"coral",
	"cyan",
	"goldenrod",
	"green",
	"khaki",
	"magenta",
	"orange",
	"orchid",
	"pink",
	"purple",
	"red",
	"salmon",
	"turquoise",
	"violet",
	"yellow",
	"aquamarine::d",
	"blue::d",
	"coral::d",
	"cyan::d",
	"goldenrod::d",
	"green::d",
	"khaki::d",
	"magenta::d",
	"orange::d",
	"orchid::d",
	"pink::d",
	"purple::d",
	"red::d",
	"salmon::d",
	"turquoise::d",
	"violet::d",
	"yellow::d",
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

		// Set Table data
		toolCounter := 0
		doneCounter := 0

		// Get Stages
		stages := maps.Keys(toolList)

		// Sort Stages
		sort.Ints(stages)

		testChan := make(chan string)
		defer close(testChan)
		globalLogChan := make(chan string)
		defer close(globalLogChan)

		progressChan := make(chan int)
		defer close(progressChan)

		reportChan := make(chan types.ExecutionReport)
		defer close(reportChan)

		finalReportChan := make(chan []types.ExecutionReport)
		defer close(finalReportChan)
		defer func() {
			report := <-finalReportChan
			yamlReport, _ := yaml.Marshal(&report)
			os.WriteFile("./report.yaml", yamlReport, 0644)
		}()

		go func() {
			reports := []types.ExecutionReport{}
			for toolCounter != doneCounter {
				reports = append(reports, <-reportChan)
			}
			finalReportChan <- reports
		}()

		executeScript := func(wg *sync.WaitGroup, script types.InstallScriptDefinition, colorCounter int) {
			logChan, statusChan := script.GetChannels()
			defer func() {
				wg.Done()
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
						return
					}
				case types.ConfigScript:
					// return
				}
			}

			cmd := exec.Command("bash", "-c", *script.Scripts.MacosArm)
			stdin, _ := cmd.StdoutPipe()
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

			func() {
				scannerGood := bufio.NewScanner(stdin)
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
						return
					}
				case types.ConfigScript:
					// return
				}
			}

			statusChan <- types.ExecutionStatusCompleted
			statusChan <- types.ExecutionStatusInstalled
			reportChan <- types.ExecutionReport{
				Name:             script.Name,
				Status:           types.ExecutionStatusInstalled,
				Version:          new(string),
				InstallationPath: new(string),
			}
		}

		// var concurantExecutionsWG sync.WaitGroup

		go func() {
			colorCounter := 0
			for _, stage := range stages {
				scripts := toolList[stage]

				var wg sync.WaitGroup
				wg.Add(len(scripts))

				for _, script := range scripts {
					colorCounter++
					go executeScript(&wg, script, colorCounter)
				}
				wg.Wait()
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

		toolListMap := map[int]types.InstallScriptDefinition{}
		// Set Table headings
		table.SetCell(0, 0, tview.NewTableCell("Stage").SetSelectable(false))
		table.SetCell(0, 1, tview.NewTableCell("Tool").SetSelectable(false))
		table.SetCell(0, 2, tview.NewTableCell("Status").SetSelectable(false))
		table.SetCell(0, 3, tview.NewTableCell("Platform").SetSelectable(false))
		table.SetCell(0, 4, tview.NewTableCell("Version").SetSelectable(false))

		for _, stage := range stages {
			for _, tool := range toolList[stage] {
				toolCounter++
				statusCell := tview.NewTableCell("Waiting")
				table.SetCell(toolCounter, 0, tview.NewTableCell(strconv.FormatInt(int64(stage), 32)))
				table.SetCell(toolCounter, 1, tview.NewTableCell(tool.Name))
				table.SetCell(toolCounter, 2, statusCell)

				toolListMap[toolCounter] = tool
				_, statusChan := tool.GetChannels()
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
						flex.AddItem(tview.NewTextView().SetDynamicColors(true).SetText(log), 1, 0, false)
					})
				}
			}()
			return flex
		}()

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
				logChan, _ := tool.GetChannels()

				logGrid := tview.NewFlex()
				command := tview.NewTextView().SetText(*tool.Scripts.MacosArm)
				command.SetTitle(" Command ").SetBorder(true)
				logGrid.SetBorder(true)
				logGrid.SetDirection(tview.FlexRow)
				logGrid.AddItem(command, 0, 1, false)
				logGrid.AddItem(toolLogView(logChan), 0, 9, true)

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
		sort.Ints(stages)

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
	RootCmd.AddCommand(installCmd, updateCmd)
}
